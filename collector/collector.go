package collector

import (
	"context"
	"net/http"

	"github.com/buarki/find-castles/castle"
	"github.com/buarki/find-castles/fanin"
)

func CollectCastlesToEnrich(ctx context.Context, httpClient *http.Client) (<-chan castle.Model, <-chan error) {
	portugueseCastlesToEnrich, errCollectingPortugueseCastles := collectCastlesFromPortugal(ctx, httpClient)
	ukCastlesToEnrich, errCollectingEnglishCastles := collectCastlesFromUK(ctx, httpClient)
	return fanin.Merge(ctx, portugueseCastlesToEnrich, ukCastlesToEnrich), fanin.Merge(ctx, errCollectingEnglishCastles, errCollectingPortugueseCastles)
}
