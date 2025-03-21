package processor

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/file-analytics/pkg/models"
)

// WorkRequest represents a file processing request
// Demonstrates struct usage
type WorkRequest struct {
	FilePath string
	// Demonstrates channel directions with responses
	ResponseChan chan<- models.ProcessResult
}

// WorkerPool manages a pool of worker goroutines
// Demonstrates struct with channels
type WorkerPool struct {
	size      int
	processor models.Processor
	// Demonstrates buffered channels
	requests chan WorkRequest
	// Demonstrates channel for worker pool control
	done chan struct{}
	// Demonstrates error handling with channels
	errors chan error
}

// NewWorkerPool creates a new worker pool
// Demonstrates constructor pattern
func NewWorkerPool(size int, processor models.Processor) *WorkerPool {
	if size <= 0 {
		size = 1
	}

	return &WorkerPool{
		size:      size,
		processor: processor,
		// Buffered channel demonstration
		requests: make(chan WorkRequest, size*2),
		done:     make(chan struct{}),
		errors:   make(chan error, size),
	}
}

// Start launches the worker pool
// Demonstrates goroutine management
func (p *WorkerPool) Start(ctx context.Context) {
	// Launch workers
	for i := 0; i < p.size; i++ {
		// Demonstrates goroutine launch
		go p.worker(ctx, i)
	}
}

// Stop gracefully shuts down the worker pool
// Demonstrates channel closing
func (p *WorkerPool) Stop() {
	close(p.requests)
	// Wait for workers to finish
	<-p.done
}

// Submit adds a file to be processed
// Demonstrates non-blocking channel operations
func (p *WorkerPool) Submit(path string) (chan models.ProcessResult, error) {
	// Create response channel
	responseChan := make(chan models.ProcessResult, 1)

	// Demonstrates select with timeout
	select {
	case p.requests <- WorkRequest{FilePath: path, ResponseChan: responseChan}:
		return responseChan, nil
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("submission timeout: worker pool is full")
	}
}

// worker processes files from the request channel
// Demonstrates goroutine worker pattern
func (p *WorkerPool) worker(ctx context.Context, id int) {
	// Demonstrates defer for cleanup
	defer func() {
		if id == 0 { // Only one worker needs to close the done channel
			close(p.done)
		}
	}()

	// Demonstrates range over channels
	for req := range p.requests {
		// Demonstrates select for cancellation
		select {
		case <-ctx.Done():
			return
		default:
			// Process the file
			result, err := p.processor.Process(ctx, req.FilePath)
			if err != nil {
				// Demonstrates error channel
				select {
				case p.errors <- err:
				default: // Don't block if error channel is full
				}
			}

			// Send result back through response channel
			// Demonstrates channel direction usage
			select {
			case req.ResponseChan <- result:
			default: // Don't block if receiver is gone
			}

			// Close the response channel
			close(req.ResponseChan)
		}
	}
}

// Errors returns a channel that receives processing errors
// Demonstrates channel as return value
func (p *WorkerPool) Errors() <-chan error {
	return p.errors
}
