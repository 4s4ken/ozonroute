package app

import (
	"context"
	"sync"
)

type Task func(ctx context.Context)

type WorkerPool struct {
	tasks  chan Task
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewWorkerPool(parent context.Context, workers int) *WorkerPool {
	ctx, cancel := context.WithCancel(parent)
	wp := &WorkerPool{
		tasks:  make(chan Task),
		ctx:    ctx,
		cancel: cancel,
	}
	wp.start(workers)
	return wp
}

func (wp *WorkerPool) start(workers int) {
	wp.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wp.wg.Done()
			for {
				select {
				case <-wp.ctx.Done():
					return
				case task, ok := <-wp.tasks:
					if !ok {
						return
					}
					task(wp.ctx)
				}
			}
		}()
	}
}

func (wp *WorkerPool) Submit(task Task) {
	select {
	case <-wp.ctx.Done():
	case wp.tasks <- task:
	}
}

func (wp *WorkerPool) Shutdown() {
	wp.cancel()
	close(wp.tasks)
	wp.wg.Wait()
}
