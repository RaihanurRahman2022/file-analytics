package worker

import (
	"context"
	"sync"
	"time"
)

// Task represents a unit of work to be processed
type Task interface {
	Process() error
	ID() string
}

// Pool manages a pool of workers with rate limiting
type Pool struct {
	workers     int
	rateLimiter chan struct{}
	tasks       chan Task
	results     chan error
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewPool creates a new worker pool with specified parameters
func NewPool(workers int, queueSize int, rateLimit time.Duration) *Pool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &Pool{
		workers:     workers,
		rateLimiter: make(chan struct{}, workers),
		tasks:       make(chan Task, queueSize),
		results:     make(chan error, queueSize),
		ctx:         ctx,
		cancel:      cancel,
	}

	// Initialize rate limiter tokens
	for i := 0; i < workers; i++ {
		pool.rateLimiter <- struct{}{}
	}

	return pool
}

// Start launches the worker pool
func (p *Pool) Start() {
	// Launch workers
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

// Submit adds a task to the pool
func (p *Pool) Submit(task Task) error {
	select {
	case p.tasks <- task:
		return nil
	case <-p.ctx.Done():
		return p.ctx.Err()
	}
}

// Stop gracefully shuts down the worker pool
func (p *Pool) Stop() {
	p.cancel()
	close(p.tasks)
	p.wg.Wait()
	close(p.results)
}

// Results returns the channel for receiving task results
func (p *Pool) Results() <-chan error {
	return p.results
}

// worker processes tasks with rate limiting
func (p *Pool) worker(id int) {
	defer p.wg.Done()

	for task := range p.tasks {
		select {
		case <-p.ctx.Done():
			return
		case <-p.rateLimiter:
			// Process task with rate limiting
			err := task.Process()
			p.results <- err

			// Return token to rate limiter
			p.rateLimiter <- struct{}{}
		}
	}
}

// Stats represents pool statistics
type Stats struct {
	ActiveWorkers  int
	QueuedTasks    int
	CompletedTasks int
}

// GetStats returns current pool statistics
func (p *Pool) GetStats() Stats {
	return Stats{
		ActiveWorkers:  p.workers - len(p.rateLimiter),
		QueuedTasks:    len(p.tasks),
		CompletedTasks: cap(p.tasks) - len(p.tasks),
	}
}
