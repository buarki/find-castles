package fanin

import (
	"context"
	"sync"
)

func Merge[T any](ctx context.Context, channels ...<-chan T) <-chan T {
	mergedChannel := make(chan T)
	var wg sync.WaitGroup
	wg.Add(len(channels))
	for _, channel := range channels {
		go func(c <-chan T) {
			defer wg.Done()
			for castle := range c {
				select {
				case <-ctx.Done():
					return
				case mergedChannel <- castle:
				}
			}
		}(channel)
	}
	go func() {
		wg.Wait()
		close(mergedChannel)
	}()
	return mergedChannel
}
