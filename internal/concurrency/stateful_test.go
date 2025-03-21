package concurrency

import (
	"testing"
	"time"
)

func TestStatefulPool(t *testing.T) {
	// Create a pool with 2 workers
	pool := NewStatefulPool(2, 10, 100*time.Millisecond)
	pool.Start()
	defer pool.Stop()

	// Submit some tasks
	tasks := []int{1, 2, 3, 4, 5}
	for _, task := range tasks {
		if err := pool.Submit(task); err != nil {
			t.Fatalf("Failed to submit task: %v", err)
		}
	}

	// Collect results
	results := make([]interface{}, 0)
	for i := 0; i < len(tasks); i++ {
		select {
		case result := <-pool.Results():
			results = append(results, result)
		case <-time.After(2 * time.Second):
			t.Fatal("Timeout waiting for results")
		}
	}

	// Verify results
	if len(results) != len(tasks) {
		t.Errorf("Expected %d results, got %d", len(tasks), len(results))
	}

	// Check worker stats
	stats := pool.GetWorkerStats()
	if len(stats) != 2 {
		t.Errorf("Expected 2 workers, got %d", len(stats))
	}

	// Verify work distribution
	totalWork := int64(0)
	for _, stat := range stats {
		totalWork += stat.WorkCount
	}
	if totalWork != int64(len(tasks)) {
		t.Errorf("Expected total work count of %d, got %d", len(tasks), totalWork)
	}
}

func TestStatefulPoolGracefulShutdown(t *testing.T) {
	pool := NewStatefulPool(2, 10, 100*time.Millisecond)
	pool.Start()

	// Submit some tasks
	for i := 0; i < 5; i++ {
		pool.Submit(i)
	}

	// Stop the pool
	pool.Stop()

	// Try to submit more tasks (should fail)
	err := pool.Submit(6)
	if err == nil {
		t.Error("Expected error when submitting to stopped pool")
	}

	// Collect remaining results
	results := make([]interface{}, 0)
	for {
		select {
		case result := <-pool.Results():
			results = append(results, result)
		case <-time.After(100 * time.Millisecond):
			goto done
		}
	}
done:
	if len(results) == 0 {
		t.Error("Expected to receive some results before shutdown")
	}
}

func TestStatefulPoolRateLimiting(t *testing.T) {
	pool := NewStatefulPool(2, 10, 200*time.Millisecond)
	pool.Start()
	defer pool.Stop()

	// Submit tasks rapidly
	start := time.Now()
	for i := 0; i < 4; i++ {
		pool.Submit(i)
	}

	// Collect results
	for i := 0; i < 4; i++ {
		<-pool.Results()
	}
	duration := time.Since(start)

	// With 2 workers and 200ms rate limit, processing 4 tasks should take at least 400ms
	if duration < 400*time.Millisecond {
		t.Errorf("Expected duration >= 400ms, got %v", duration)
	}
}
