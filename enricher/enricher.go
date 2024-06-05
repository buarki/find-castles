package enricher

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"github.com/buarki/find-castles/castle"
	"github.com/buarki/find-castles/fanin"
)

type Enricher interface {
	CollectCastlesToEnrich(ctx context.Context) ([]castle.Model, error)

	EnrichCastle(ctx context.Context, c castle.Model) (castle.Model, error)
}

type EnchimentExecutor struct {
	enrichers map[castle.Country]Enricher
	cpus      int
}

func New(cpusToUse int, httpClient *http.Client, fetchHTML func(ctx context.Context, link string, httpClient *http.Client) ([]byte, error)) *EnchimentExecutor {
	cpus := cpusToUse
	availableCPUs := runtime.NumCPU()
	if cpusToUse > availableCPUs {
		cpus = availableCPUs
	}
	return &EnchimentExecutor{
		cpus: cpus,
		enrichers: map[castle.Country]Enricher{
			castle.Ireland:  NewIrishEnricher(httpClient, fetchHTML),
			castle.Portugal: NewPortugueseEnricher(httpClient, fetchHTML),
			castle.UK:       NewBritishEnricher(httpClient, fetchHTML),
		},
	}
}

func (ex *EnchimentExecutor) Enrich(ctx context.Context) (<-chan castle.Model, <-chan error) {
	castlesToEnrich, errChan := ex.collectCastles(ctx)
	enrichedCastlesBuf := []<-chan castle.Model{}
	castlesEnrichmentErr := []<-chan error{errChan}
	for i := 0; i < ex.cpus; i++ {
		receivedEnrichedCastlesChan, enrichErrs := ex.extractData(ctx, castlesToEnrich)
		enrichedCastlesBuf = append(enrichedCastlesBuf, receivedEnrichedCastlesChan)
		castlesEnrichmentErr = append(castlesEnrichmentErr, enrichErrs)
	}

	enrichedCastles := fanin.Merge(ctx, enrichedCastlesBuf...)
	enrichmentErrs := fanin.Merge(ctx, castlesEnrichmentErr...)

	return enrichedCastles, enrichmentErrs
}

func (ex *EnchimentExecutor) toChanel(ctx context.Context, e Enricher) (<-chan castle.Model, <-chan error) {
	castlesToEnrich := make(chan castle.Model)
	errChan := make(chan error)
	go func() {
		defer close(castlesToEnrich)
		defer close(errChan)

		englandCastles, err := e.CollectCastlesToEnrich(ctx)
		if err != nil {
			errChan <- err
		}
		for _, c := range englandCastles {
			castlesToEnrich <- c
		}
	}()
	return castlesToEnrich, errChan
}

func (ex *EnchimentExecutor) collectCastles(ctx context.Context) (<-chan castle.Model, <-chan error) {
	var collectingChan []<-chan castle.Model
	var errChan []<-chan error
	for _, enricher := range ex.enrichers {
		castlesChan, castlesErrChan := ex.toChanel(ctx, enricher)
		collectingChan = append(collectingChan, castlesChan)
		errChan = append(errChan, castlesErrChan)
	}
	return fanin.Merge(ctx, collectingChan...), fanin.Merge(ctx, errChan...)
}

func (ex *EnchimentExecutor) extractData(ctx context.Context, castlesToEnrich <-chan castle.Model) (chan castle.Model, chan error) {
	enrichedCastles := make(chan castle.Model)
	errChan := make(chan error)

	go func() {
		defer close(enrichedCastles)
		defer close(errChan)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("[enrichCastle] done")
				return
			case castleToEnrich, ok := <-castlesToEnrich:
				if ok {
					fmt.Println(castleToEnrich)
					enricher := ex.enrichers[castleToEnrich.Country]
					enrichedCastle, err := enricher.EnrichCastle(ctx, castleToEnrich)
					if err != nil {
						errChan <- err
					} else {
						fmt.Println("CPU", enrichedCastle)
						enrichedCastles <- enrichedCastle
					}
				} else {
					fmt.Println("received zero value from chanel")
					return
				}
			}
		}
	}()

	return enrichedCastles, errChan
}
