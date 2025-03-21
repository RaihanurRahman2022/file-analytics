package concurrency

import (
	"context"
	"sync"
	"time"
)

// StatefulWorker represents a worker that maintains state
type StatefulWorker struct {
	ID        int
	State     interface{}
	LastWork  time.Time
	WorkCount int64
	mu        sync.RWMutex
}

// StatefulPool manages a pool of stateful workers
type StatefulPool struct {
	workers     []*StatefulWorker
	tasks       chan interface{}
	results     chan interface{}
	done        chan struct{}
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	rateLimiter chan struct{}
}

// NewStatefulPool creates a new pool of stateful workers
func NewStatefulPool(workers int, queueSize int, rateLimit time.Duration) *StatefulPool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &StatefulPool{
		workers:     make([]*StatefulWorker, workers),
		tasks:       make(chan interface{}, queueSize),
		results:     make(chan interface{}, queueSize),
		done:        make(chan struct{}),
		ctx:         ctx,
		cancel:      cancel,
		rateLimiter: make(chan struct{}, workers),
	}

	// Initialize workers
	for i := 0; i < workers; i++ {
		pool.workers[i] = &StatefulWorker{
			ID:        i,
			LastWork:  time.Now(),
			WorkCount: 0,
		}
		pool.rateLimiter <- struct{}{}
	}

	return pool
}

// Start launches the worker pool
func (p *StatefulPool) Start() {
	for i, worker := range p.workers {
		p.wg.Add(1)
		go p.runWorker(i, worker)
	}
}

// Submit adds a task to the pool
func (p *StatefulPool) Submit(task interface{}) error {
	select {
	case p.tasks <- task:
		return nil
	case <-p.ctx.Done():
		return p.ctx.Err()
	}
}

// Stop gracefully shuts down the pool
func (p *StatefulPool) Stop() {
	p.cancel()
	close(p.tasks)
	p.wg.Wait()
	close(p.results)
}

// Results returns the channel for receiving task results
func (p *StatefulPool) Results() <-chan interface{} {
	return p.results
}

// runWorker runs a single stateful worker
func (p *StatefulPool) runWorker(id int, worker *StatefulWorker) {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-p.rateLimiter:
			// Process task with rate limiting
			select {
			case task := <-p.tasks:
				// Update worker state
				worker.mu.Lock()
				worker.LastWork = time.Now()
				worker.WorkCount++
				worker.mu.Unlock()

				// Process task
				result := p.processTask(worker, task)
				p.results <- result

				// Return token to rate limiter
				p.rateLimiter <- struct{}{}
			case <-p.ctx.Done():
				return
			}
		}
	}
}

// processTask processes a single task and updates worker state
func (p *StatefulPool) processTask(worker *StatefulWorker, task interface{}) interface{} {
	// Example task processing
	// In a real application, this would be customized based on the task type
	time.Sleep(100 * time.Millisecond) // Simulate work
	return task
}

// GetWorkerStats returns statistics for all workers
func (p *StatefulPool) GetWorkerStats() []WorkerStats {
	stats := make([]WorkerStats, len(p.workers))
	for i, worker := range p.workers {
		worker.mu.RLock()
		stats[i] = WorkerStats{
			ID:        worker.ID,
			LastWork:  worker.LastWork,
			WorkCount: worker.WorkCount,
		}
		worker.mu.RUnlock()
	}
	return stats
}

// WorkerStats represents statistics for a single worker
type WorkerStats struct {
	ID        int
	LastWork  time.Time
	WorkCount int64
} 