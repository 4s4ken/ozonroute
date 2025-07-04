package app

import (
	"context"
	"testing"
	"time"
)

func TestWorkerPool_Resize(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wp := NewWorkerPool(ctx, 2)
	if wp.workers != 2 {
		t.Fatalf("expected 2 workers, got %d", wp.workers)
	}

	wp.Resize(4)
	if wp.workers != 4 {
		t.Fatalf("expected 4 workers after resize up, got %d", wp.workers)
	}

	wp.Resize(1)
	if wp.workers != 1 {
		t.Fatalf("expected 1 worker after resize down, got %d", wp.workers)
	}

	taskExecuted := make(chan struct{}, 1)

	err := wp.Submit(func(ctx context.Context) {
		select {
		case <-ctx.Done():
			return
		case taskExecuted <- struct{}{}:
		}
	})
	if err != nil {
		t.Fatalf("unexpected error from Submit: %v", err)
	}

	select {
	case <-taskExecuted:
	case <-time.After(time.Second):
		t.Fatal("task was not executed after resize")
	}
	
	wp.Shutdown()
	if wp.workers != 1 {
		t.Fatalf("expected workers to remain at 1 after shutdown, got %d", wp.workers)
	}
}
