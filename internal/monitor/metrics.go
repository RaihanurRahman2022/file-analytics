package monitor

import (
	"sync/atomic"
	"time"
)

// MetricsCollector handles system-wide metrics collection
// Demonstrates atomic operations and periodic reporting
type MetricsCollector struct {
	// Atomic counters for thread-safe metrics
	processed atomic.Uint64
	errors    atomic.Uint64
	duration  atomic.Int64

	// Channels for control
	stopChan chan struct{}
	ticker   *time.Ticker
}

// NewMetricsCollector creates a new metrics collector
// Demonstrates constructor pattern and ticker setup
func NewMetricsCollector(reportInterval time.Duration) *MetricsCollector {
	return &MetricsCollector{
		stopChan: make(chan struct{}),
		ticker:   time.NewTicker(reportInterval),
	}
}

// Start begins periodic metrics reporting
// Demonstrates goroutine and ticker usage
func (m *MetricsCollector) Start() {
	go func() {
		for {
			select {
			case <-m.ticker.C:
				m.reportMetrics()
			case <-m.stopChan:
				m.ticker.Stop()
				return
			}
		}
	}()
}

// Stop halts metrics reporting
// Demonstrates graceful shutdown
func (m *MetricsCollector) Stop() {
	close(m.stopChan)
}

// IncrementProcessed atomically increments the processed counter
// Demonstrates atomic operations
func (m *MetricsCollector) IncrementProcessed() {
	m.processed.Add(1)
}

// IncrementErrors atomically increments the error counter
func (m *MetricsCollector) IncrementErrors() {
	m.errors.Add(1)
}

// AddDuration atomically adds to the total duration
func (m *MetricsCollector) AddDuration(d time.Duration) {
	m.duration.Add(int64(d))
}

// GetMetrics returns current metrics
// Demonstrates multiple return values
func (m *MetricsCollector) GetMetrics() (processed uint64, errors uint64, avgDuration time.Duration) {
	processed = m.processed.Load()
	errors = m.errors.Load()

	// Calculate average duration
	if processed > 0 {
		totalDuration := m.duration.Load()
		avgDuration = time.Duration(totalDuration) / time.Duration(processed)
	}

	return
}

// reportMetrics handles periodic metrics reporting
// Demonstrates time formatting and logging
func (m *MetricsCollector) reportMetrics() {
	processed, errors, avgDuration := m.GetMetrics()

	// Format metrics report
	report := struct {
		Timestamp   string
		Processed   uint64
		Errors      uint64
		AvgDuration string
	}{
		Timestamp:   time.Now().Format(time.RFC3339),
		Processed:   processed,
		Errors:      errors,
		AvgDuration: avgDuration.String(),
	}

	// In a real application, you might:
	// - Log to a file
	// - Send to a monitoring service
	// - Update metrics endpoint
	// - Store in a time-series database
	_ = report // Placeholder for actual reporting logic
}
