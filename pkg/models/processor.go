package models

import (
	"context"
	"io"
	"time"
)

// FileInfo represents metadata about a processed file
// Demonstrates struct usage
type FileInfo struct {
	Path      string
	Size      int64
	Modified  time.Time
	Processed time.Time
	Type      string
}

// ProcessResult represents the result of file processing
// Demonstrates struct composition
type ProcessResult struct {
	FileInfo
	Lines    int
	Words    int
	Bytes    int
	Error    error
	Duration time.Duration
}

// Processor defines the interface for file processors
// Demonstrates interface definition
type Processor interface {
	// Process handles a single file
	// Demonstrates multiple return values
	Process(ctx context.Context, path string) (ProcessResult, error)

	// CanHandle checks if this processor can handle the given file type
	// Demonstrates simple return values
	CanHandle(path string) bool

	// Name returns the processor name
	// Demonstrates method definition
	Name() string
}

// BaseProcessor provides common functionality for processors
// Demonstrates struct embedding and composition
type BaseProcessor struct {
	name       string
	bufferSize int
}

// NewBaseProcessor demonstrates a constructor function
func NewBaseProcessor(name string, bufferSize int) *BaseProcessor {
	// Demonstrates if/else with single line
	if bufferSize <= 0 {
		bufferSize = 4096
	}

	return &BaseProcessor{
		name:       name,
		bufferSize: bufferSize,
	}
}

// Name implements the Processor interface
func (p *BaseProcessor) Name() string {
	return p.name
}

// readLines demonstrates working with io.Reader and error handling
func (p *BaseProcessor) readLines(reader io.Reader) (lines, words, bytes int, err error) {
	// Create a buffer for reading
	// Demonstrates array usage
	buf := make([]byte, p.bufferSize)

	// Variables to track state
	var (
		inWord bool
		count  int
	)

	// Read the file in chunks
	// Demonstrates for loop with multiple conditions
	for {
		count, err = reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}

		bytes += count

		// Process the buffer
		// Demonstrates range loop over slice
		for _, b := range buf[:count] {
			// Count lines
			if b == '\n' {
				lines++
			}

			// Count words
			// Demonstrates switch statement
			switch {
			case b == ' ' || b == '\n' || b == '\t':
				inWord = false
			case !inWord:
				words++
				inWord = true
			}
		}
	}

	// Adjust final counts
	if bytes > 0 && !inWord {
		lines++
	}

	return
}
