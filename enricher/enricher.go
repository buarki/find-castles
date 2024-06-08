package enricher

import (
	"context"

	"github.com/buarki/find-castles/castle"
)

type Enricher interface {
	CollectCastlesToEnrich(ctx context.Context) ([]castle.Model, error)

	EnrichCastle(ctx context.Context, c castle.Model) (castle.Model, error)
}
