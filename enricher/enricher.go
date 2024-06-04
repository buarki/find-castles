package enricher

import (
	"context"
	"fmt"
	"net/http"

	"github.com/buarki/find-castles/castle"
)

type Enricher func(
	ctx context.Context,
	httpClient *http.Client,
	castle castle.Model,
) (castle.Model, error)

var (
	enrichers = map[castle.Country]Enricher{
		castle.Portugal: enrichCastleFromPortugal,
		castle.UK:       enrichCastleFromUK,
	}
)

func Enrich(ctx context.Context, httpClient *http.Client, castlesToEnrich <-chan castle.Model, cpuId int) (chan castle.Model, chan error) {
	enrichedCastles := make(chan castle.Model)
	errChan := make(chan error)

	go func() {
		defer close(enrichedCastles)
		defer close(errChan)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("[enrichCastle] done", cpuId)
				return
			case castleToEnrich, ok := <-castlesToEnrich:
				if ok {
					enricher := enrichers[castleToEnrich.Country]
					enrichedCastle, err := enricher(ctx, httpClient, castleToEnrich)
					if err != nil {
						errChan <- err
					} else {
						fmt.Println("CPU", cpuId, enrichedCastle)
						enrichedCastles <- enrichedCastle
					}
				} else {
					fmt.Println("received zero value from chanel", cpuId)
					return
				}
			}
		}
	}()

	return enrichedCastles, errChan
}
