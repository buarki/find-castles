package fanin_test

import (
	"context"
	"testing"
	"time"

	"github.com/buarki/find-castles/fanin"
)

func TestMerge(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		defer close(ch1)
		for i := 0; i < 5; i++ {
			ch1 <- i
		}
	}()

	go func() {
		defer close(ch2)
		for i := 5; i < 10; i++ {
			ch2 <- i
		}
	}()

	mergedChannel := fanin.Merge(ctx, ch1, ch2)

	var results []int
	for item := range mergedChannel {
		results = append(results, item)
	}

	expectedSet := map[int]bool{
		0: true,
		1: true,
		2: true,
		3: true,
		4: true,
		5: true,
		6: true,
		7: true,
		8: true,
		9: true,
	}
	if len(results) != len(expectedSet) {
		t.Fatalf("expected %v elements, got %v elements", len(expectedSet), len(results))
	}

	for r := range results {
		if !expectedSet[r] {
			t.Errorf("expected element %d not found in results", r)
		}
	}
}

func TestMergeWithContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		defer close(ch1)
		for i := 0; i < 5; i++ {
			time.Sleep(100 * time.Millisecond)
			ch1 <- i
		}
	}()

	go func() {
		defer close(ch2)
		for i := 5; i < 10; i++ {
			time.Sleep(100 * time.Millisecond)
			ch2 <- i
		}
	}()

	mergedChannel := fanin.Merge(ctx, ch1, ch2)
	go func() {
		time.Sleep(200 * time.Millisecond)
		cancel()
	}()

	var results []int
	for item := range mergedChannel {
		results = append(results, item)
	}

	if len(results) >= 10 {
		t.Fatalf("expected fewer than 10 elements, got %v elements", len(results))
	}
}
