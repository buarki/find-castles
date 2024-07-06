package executor

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/buarki/find-castles/castle"
	"github.com/buarki/find-castles/enricher"
)

type EnchimentExecutor struct {
	enrichers      map[enricher.Source]enricher.Enricher
	collectingCPUs int
	extractingCPUs int
}

func New(
	collectingCPUs,
	extractingCPUs int,
	httpClient *http.Client,
	enrichers map[enricher.Source]enricher.Enricher) *EnchimentExecutor {
	return &EnchimentExecutor{
		enrichers:      enrichers,
		collectingCPUs: collectingCPUs,
		extractingCPUs: extractingCPUs,
	}
}

func (ex *EnchimentExecutor) Enrich(ctx context.Context) (<-chan castle.Model, <-chan error) {
	enrichedCastles := make(chan castle.Model)
	errChan := make(chan error)

	enrichersChan := make(chan enricher.Enricher, len(ex.enrichers))

	for _, enricher := range ex.enrichers {
		enrichersChan <- enricher
	}
	close(enrichersChan)

	castlesToEnrichChan := make(chan castle.Model)
	var castlesToEnrichChanWg sync.WaitGroup
	castlesToEnrichChanWg.Add(ex.collectingCPUs)

	for i := 0; i < ex.collectingCPUs; i++ {
		go func() {
			defer castlesToEnrichChanWg.Done()
			for {
				select {
				case <-ctx.Done():
					fmt.Println("done getting castle channels")
					return
				case enricher, ok := <-enrichersChan:
					if !ok {
						return
					}
					ex.collectCastlesToEnrich(ctx, enricher, castlesToEnrichChan, errChan)
				}
			}
		}()
	}

	go func() {
		castlesToEnrichChanWg.Wait()

		close(castlesToEnrichChan)
	}()

	var enrichedCastlesWg sync.WaitGroup
	enrichedCastlesWg.Add(ex.extractingCPUs)

	for i := 0; i < ex.extractingCPUs; i++ {
		go func() {
			defer enrichedCastlesWg.Done()

			for {
				select {
				case <-ctx.Done():
					fmt.Println("Main done")
					return
				case c, ok := <-castlesToEnrichChan:
					if !ok {
						return
					}
					enrichedCastle, err := ex.enrichers[enricher.Source(c.CurrentEnrichmentSource)].EnrichCastle(ctx, c)
					if err != nil {
						errChan <- err
					} else {
						enrichedCastles <- enrichedCastle
					}
				}
			}
		}()
	}

	go func() {
		enrichedCastlesWg.Wait()

		close(enrichedCastles)
		close(errChan)
	}()

	return enrichedCastles, errChan
}

func (ex *EnchimentExecutor) collectCastlesToEnrich(
	ctx context.Context,
	enricher enricher.Enricher,
	castlesToEnrichChan chan castle.Model,
	errChan chan error,
) {
	castlesChan, eChan := enricher.CollectCastlesToEnrich(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case c, ok := <-castlesChan:
			if !ok {
				return
			}
			castlesToEnrichChan <- c
		case e, ok := <-eChan:
			if !ok {
				return
			}
			errChan <- e
		}
	}
}
