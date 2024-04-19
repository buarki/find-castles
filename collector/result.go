package collector

import "github.com/buarki/find-castles/castle"

type CollectResult struct {
	Castle castle.Model
	Err    error
}
