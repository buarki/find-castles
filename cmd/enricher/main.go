package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/buarki/find-castles/castle"
	"github.com/buarki/find-castles/db"
	"github.com/buarki/find-castles/enricher"
	"github.com/buarki/find-castles/executor"
	"github.com/buarki/find-castles/htmlfetcher"
	"github.com/buarki/find-castles/httpclient"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	databaseName   = "find-castles"
	collectionName = "castles"
	bufferSize     = 10
)

func main() {
	mongoURI := os.Getenv("DB_URI")
	if mongoURI == "" {
		log.Fatal("missing env var DB_URI")
	}
	operationTimeoutInSeconds := os.Getenv("ENRICHMENT_TIMEOUT_IN_SECONDS")
	if mongoURI == "" {
		log.Fatal("missing env var ENRICHMENT_TIMEOUT_IN_SECONDS")
	}
	timeoutAsNumber, err := strconv.ParseInt(operationTimeoutInSeconds, 10, 0)
	if err != nil {
		log.Fatalf("failed to parse given ENRICHMENT_TIMEOUT_IN_SECONDS [%s] into number, got %v", operationTimeoutInSeconds, err)
	}
	connectionTimeout := time.Duration(timeoutAsNumber) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	dbClient, err := db.NewClient(ctx, mongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer dbClient.Disconnect(ctx)

	collection := dbClient.Database(databaseName).Collection(collectionName)

	if err := db.AddIndexes(ctx, collection); err != nil {
		log.Fatal(err)
	}

	httpClient := httpclient.New()
	enrichers := map[castle.Country]enricher.Enricher{
		castle.Portugal: enricher.NewCastelosDePortugalEnricher(httpClient, htmlfetcher.Fetch),
		castle.Ireland:  enricher.NewHeritageIreland(httpClient, htmlfetcher.Fetch),
		castle.UK:       enricher.NewMedievalBritainEnricher(httpClient, htmlfetcher.Fetch),
		castle.Slovakia: enricher.NewEbidatEnricher(httpClient, htmlfetcher.Fetch),
	}
	cpus := runtime.NumCPU()
	castlesEnricher := executor.New(int(float64(cpus)*0.3), int(float64(cpus)*0.7), httpClient, enrichers)
	castlesChan, errChan := castlesEnricher.Enrich(ctx)

	checkingCastlesBuffer := make([]castle.Model, 0, bufferSize)

	for {
		select {
		case <-ctx.Done():
			return
		case castle, ok := <-castlesChan:
			if !ok {
				if len(checkingCastlesBuffer) > 0 {
					if err := processBuffer(ctx, collection, checkingCastlesBuffer); err != nil {
						log.Fatal(err)
					}
				}
				return
			}
			checkingCastlesBuffer = append(checkingCastlesBuffer, castle)
			if len(checkingCastlesBuffer) == bufferSize {
				if err := processBuffer(ctx, collection, checkingCastlesBuffer); err != nil {
					log.Fatal(err)
				}
				checkingCastlesBuffer = checkingCastlesBuffer[:0]
			}
		case err := <-errChan:
			if err != nil {
				log.Printf("error enriching castles: %v", err)
			}
		}
	}
}

func processBuffer(ctx context.Context, collection *mongo.Collection, buffer []castle.Model) error {
	if len(buffer) == 0 {
		return errors.New("cannot process empty buffer")
	}

	similarCastlesFound, err := db.TryToFindCastles(ctx, collection, buffer)
	if err != nil {
		return err
	}

	castlesToSave, err := reconcileCastles(buffer, similarCastlesFound)
	if err != nil {
		return err
	}

	if err := db.SaveCastles(ctx, collection, castlesToSave); err != nil {
		return err
	}

	return nil
}

func reconcileCastles(newCastles, similarCastles []castle.Model) ([]castle.Model, error) {
	var result []castle.Model

	// O(nˆ2), can we improve it?
	for _, newCastle := range newCastles {
		found := false
		for _, existingCastle := range similarCastles {
			if newCastle.IsProbably(existingCastle) {
				slog.Info("found similar castle", "current castle name", newCastle.Name, "found castle name", existingCastle.Name)
				reconciliatedCastle, err := newCastle.ReconcileWith(existingCastle)
				if err != nil {
					return nil, err
				}
				result = append(result, reconciliatedCastle)
				found = true
				break
			}
		}
		if !found {
			result = append(result, newCastle)
		}
	}

	return result, nil
}
