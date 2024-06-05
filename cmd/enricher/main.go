package main

import (
	"context"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/buarki/find-castles/castle"
	"github.com/buarki/find-castles/db"
	"github.com/buarki/find-castles/enricher"
	"github.com/buarki/find-castles/htmlfetcher"
	"github.com/buarki/find-castles/httpclient"
)

const (
	databaseName   = "castles"
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
	castlesEnricher := enricher.New(runtime.NumCPU(), httpClient, htmlfetcher.Fetch)
	castlesChan, errChan := castlesEnricher.Enrich(ctx)

	var buffer []castle.Model

	for {
		select {
		case castle, ok := <-castlesChan:
			if !ok {
				if len(buffer) > 0 {
					if err := db.SaveCastles(ctx, collection, buffer); err != nil {
						log.Fatal(err)
					}
				}
				return
			}
			buffer = append(buffer, castle)
			if len(buffer) >= bufferSize {
				if err := db.SaveCastles(ctx, collection, buffer); err != nil {
					log.Fatal(err)
				}
				buffer = buffer[:0]
			}
		case err := <-errChan:
			if err != nil {
				log.Printf("error enriching castles: %v", err)
			}
		}
	}
}
