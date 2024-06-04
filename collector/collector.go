package collector

import (
	"context"
	"net/http"

	"github.com/buarki/find-castles/castle"
	"github.com/buarki/find-castles/fanin"
)

func CollectCastlesToEnrich(ctx context.Context, httpClient *http.Client) (<-chan castle.Model, <-chan error) {
	// IoC please
	ukCastlesToEnrich, errCollectingEnglishCastles := collectCastlesFromUK(ctx, httpClient)
	irishCastlesToEnrich, errCollectingIrishCastles := collectCastlesFromIreland(ctx, httpClient)
	portugueseCastlesToEnrich, errCollectingPortugueseCastles := collectCastlesFromPortugal(ctx, httpClient)
	return fanin.Merge(ctx, portugueseCastlesToEnrich, ukCastlesToEnrich, irishCastlesToEnrich), fanin.Merge(ctx, errCollectingEnglishCastles, errCollectingPortugueseCastles, errCollectingIrishCastles)
}
