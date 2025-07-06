package workerpool

import (
	"context"
	"errors"
	"log"
	"sync"
)

const workerPoolKey = "workerPool"

type WorkerPool interface {
	Submit(task Task) error
	Resize(n int)
	Shutdown()
}

type Task func(ctx context.Context)

type Pool struct {
	tasks   chan Task
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	workers int
}

func NewWorkerPool(parent context.Context, n int) *Pool {
	ctx, cancel := context.WithCancel(parent)
	wp := &Pool{
		tasks:   make(chan Task),
		ctx:     ctx,
		cancel:  cancel,
		workers: 0,
	}
	wp.Resize(n)
	return wp
}

func (wp *Pool) worker(id int) {
	defer wp.wg.Done()

	log.Printf("worker %d started", id)

	for {
		select {
		case <-wp.ctx.Done():
			log.Printf("worker %d shutting down", id)
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

func (wp *Pool) Submit(task Task) error {
	if task == nil {
		return errors.New("task is nil")
	}

	select {
	case <-wp.ctx.Done():
		return errors.New("worker pool is shutting down")
	case wp.tasks <- task:
		return nil
	}
}

func (wp *Pool) Resize(n int) {
	select {
	case <-wp.ctx.Done():
		log.Printf("worker pool is shutting down")
		return
	default:
	}

	if n < 0 {
		log.Printf("mw invalid resize: %d", n)
		return
	}

	diff := n - wp.workers
	if diff == 0 {
		log.Printf("mw resize already %d", n)
		return
	}

	log.Printf("mw resizing from %d to %d", wp.workers, n)

	if diff > 0 {
		for i := 0; i < diff; i++ {
			wp.wg.Add(1)
			go wp.worker(wp.workers + i + 1)
		}
	}
	wp.workers = n
}

func (wp *Pool) Shutdown() {
	log.Println("Graceful shutdown starts")
	wp.cancel()

	wp.wg.Wait()
	log.Println("Graceful shutdown done")
}

func WorkerPoolValue(ctx context.Context, wp WorkerPool) context.Context {
	return context.WithValue(ctx, workerPoolKey, wp)
}

func WorkerPoolFromContext(ctx context.Context) (WorkerPool, bool) {
	wp, ok := ctx.Value(workerPoolKey).(WorkerPool)
	return wp, ok
}
