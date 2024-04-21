package main

import (
	"context"
	"net/http"
	"runtime"

	"github.com/buarki/find-castles/castle"
	"github.com/buarki/find-castles/collector"
	"github.com/buarki/find-castles/enricher"
	"github.com/buarki/find-castles/fanin"
)

func findCastles(ctx context.Context, httpClient *http.Client) (<-chan castle.Model, <-chan error) {
	castlesToEnrich, err1 := collector.CollectCastlesToEnrich(ctx, httpClient)

	enrichedCastlesBuf := []<-chan castle.Model{}
	castlesEnrichmentErr := []<-chan error{err1}
	for i := 0; i < runtime.NumCPU(); i++ {
		receivedEnrichedCastlesChan, enrichErrs := enricher.Enrich(ctx, httpClient, castlesToEnrich, i)
		enrichedCastlesBuf = append(enrichedCastlesBuf, receivedEnrichedCastlesChan)
		castlesEnrichmentErr = append(castlesEnrichmentErr, enrichErrs)
	}

	enrichedCastles := fanin.Merge(ctx, enrichedCastlesBuf...)
	enrichmentErrs := fanin.Merge(ctx, castlesEnrichmentErr...)

	return enrichedCastles, enrichmentErrs
}
