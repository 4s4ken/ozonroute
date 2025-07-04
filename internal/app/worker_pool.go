package app

import (
	"context"
	"errors"
	"log"
	"sync"
)

type Task func(ctx context.Context)

type WorkerPool struct {
	tasks     chan Task
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	mu        sync.Mutex
	workers   int
	stopChans []chan struct{}
}

func NewWorkerPool(parent context.Context, n int) *WorkerPool {
	ctx, cancel := context.WithCancel(parent)
	wp := &WorkerPool{
		tasks:     make(chan Task),
		ctx:       ctx,
		cancel:    cancel,
		workers:   0,
		stopChans: make([]chan struct{}, 0),
	}
	wp.Resize(n)
	return wp
}

func (wp *WorkerPool) worker(id int, stopCh chan struct{}) {
	wp.wg.Add(1)
	defer wp.wg.Done()

	log.Printf("worker %d started", id)

	for {
		select {
		case <-wp.ctx.Done():
			log.Printf("worker %d shutting down", id)
			return
		case <-stopCh:
			log.Printf("worker %d stopping", id)
			return
		case task, ok := <-wp.tasks:
			if !ok {
				log.Printf("worker %d stopping", id)
				return
			}
			task(wp.ctx)
		}
	}
}

func (wp *WorkerPool) Submit(task Task) error {
	select {
	case <-wp.ctx.Done():
		return errors.New("worker pool is shutting down")
	case wp.tasks <- task:
		return nil
	}
}

func (wp *WorkerPool) Resize(n int) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if n < 0 {
		log.Printf("wp invalid resize: %d", n)
		return
	}

	diff := n - wp.workers
	if diff == 0 {
		log.Printf("wp resize already %d", n)
		return
	}

	log.Printf("wp resizing from %d to %d", wp.workers, n)

	if diff > 0 {
		for i := 0; i < diff; i++ {
			stopCh := make(chan struct{})
			wp.stopChans = append(wp.stopChans, stopCh)
			go wp.worker(wp.workers+i+1, stopCh)
		}
	} else {
		for i := 0; i < -diff; i++ {
			idx := wp.workers - 1 - i
			close(wp.stopChans[idx])
		}
		wp.stopChans = wp.stopChans[:n]
	}

	wp.workers = n
}

func (wp *WorkerPool) Shutdown() {
	log.Println("Graceful shutdown starts")
	wp.cancel()

	wp.mu.Lock()
	defer wp.mu.Unlock()

	for _, stop := range wp.stopChans {
		close(stop)
	}
	close(wp.tasks)

	wp.wg.Wait()
	log.Println("Graceful shutdown done")
}

func (s *Service) ResizeWorkerPool(ctx context.Context, n int) error {
	s.wp.Resize(n)
	return nil
}
