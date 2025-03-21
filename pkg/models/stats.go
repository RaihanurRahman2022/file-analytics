package models

import (
	"fmt"
	"sync"
)

// Numeric is a constraint that permits any numeric type
// Demonstrates generic constraints
type Numeric interface {
	~int | ~int64 | ~float64 | ~uint64
}

// StatValue represents a statistical value with a name
// Demonstrates struct with generics
type StatValue[T Numeric] struct {
	Name  string
	Value T
}

// StatsCollector is a generic statistics collector
// Demonstrates generic type with constraints
type StatsCollector[T Numeric] struct {
	mu    sync.RWMutex
	stats map[string]T
}

// NewStatsCollector creates a new statistics collector
// Demonstrates generic constructor
func NewStatsCollector[T Numeric]() *StatsCollector[T] {
	return &StatsCollector[T]{
		stats: make(map[string]T),
	}
}

// Add adds a value to a named statistic
// Demonstrates generic method
func (s *StatsCollector[T]) Add(name string, value T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stats[name] += value
}

// Get retrieves a named statistic
// Demonstrates multiple return values with generics
func (s *StatsCollector[T]) Get(name string) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.stats[name]
	return val, ok
}

// Iterator represents an iterator over statistics
// Demonstrates iterator interface
type Iterator[T Numeric] interface {
	// Next advances the iterator
	Next() bool
	// Value returns the current value
	Value() StatValue[T]
}

// statsIterator implements Iterator
// Demonstrates struct embedding and iterator pattern
type statsIterator[T Numeric] struct {
	collector *StatsCollector[T]
	keys      []string
	current   int
}

// Iterate creates an iterator for the statistics
// Demonstrates factory method for iterator
func (s *StatsCollector[T]) Iterate() Iterator[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Get all keys
	keys := make([]string, 0, len(s.stats))
	for k := range s.stats {
		keys = append(keys, k)
	}

	return &statsIterator[T]{
		collector: s,
		keys:      keys,
		current:   -1,
	}
}

// Next implements Iterator interface
func (it *statsIterator[T]) Next() bool {
	it.current++
	return it.current < len(it.keys)
}

// Value implements Iterator interface
func (it *statsIterator[T]) Value() StatValue[T] {
	if it.current >= len(it.keys) {
		panic("iterator out of bounds")
	}

	key := it.keys[it.current]
	value, _ := it.collector.Get(key)
	return StatValue[T]{
		Name:  key,
		Value: value,
	}
}

// String provides a string representation of the statistics
// Demonstrates stringer interface implementation
func (s *StatsCollector[T]) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := "Statistics:\n"
	for name, value := range s.stats {
		result += fmt.Sprintf("%s: %v\n", name, value)
	}
	return result
}

// Reset clears all statistics
// Demonstrates method with pointer receiver
func (s *StatsCollector[T]) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stats = make(map[string]T)
}
