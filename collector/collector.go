package collector

import (
	"context"
	"net/http"

	"github.com/buarki/find-castles/castle"
	"github.com/buarki/find-castles/fanin"
)

func CollectCastlesToEnrich(ctx context.Context, httpClient *http.Client) (<-chan castle.Model, <-chan error) {
	portugueseCastlesToEnrich, errCollectingPortugueseCastles := collectForPotugal(ctx, httpClient)
	englishCastlesToEnrich, errCollectingEnglishCastles := collectForUK(ctx, httpClient)
	return fanin.Merge(ctx, portugueseCastlesToEnrich, englishCastlesToEnrich), fanin.Merge(ctx, errCollectingEnglishCastles, errCollectingPortugueseCastles)
}
