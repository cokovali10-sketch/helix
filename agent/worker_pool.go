// Package agent provides worker pool
package agent

import (
	"sync"
)

// WorkerPool manages workers
type WorkerPool struct {
	workers int
	jobs    chan func()
	wg      sync.WaitGroup
	stop    chan struct{}
}

// NewWorkerPool creates pool
func NewWorkerPool(workers int) *WorkerPool {
	wp := &WorkerPool{
		workers: workers,
		jobs:    make(chan func(), workers*4),
		stop:    make(chan struct{}),
	}
	for i := 0; i < workers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
	return wp
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	for {
		select {
		case job := <-wp.jobs:
			job()
		case <-wp.stop:
			return
		}
	}
}

// Submit submits job
func (wp *WorkerPool) Submit(job func()) {
	wp.jobs <- job
}

// Stop stops pool
func (wp *WorkerPool) Stop() {
	close(wp.stop)
	wp.wg.Wait()
}

// Workers returns count
func (wp *WorkerPool) Workers() int {
	return wp.workers
}
