package processor

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yourusername/file-analytics/pkg/models"
)

// TextProcessor implements the Processor interface for text files
// Demonstrates struct embedding
type TextProcessor struct {
	*models.BaseProcessor
	// Supported extensions
	extensions []string
}

// NewTextProcessor demonstrates a constructor function with variadic parameters
func NewTextProcessor(bufferSize int, extensions ...string) *TextProcessor {
	// If no extensions provided, use defaults
	// Demonstrates slice operations
	if len(extensions) == 0 {
		extensions = []string{".txt", ".log", ".md"}
	}

	return &TextProcessor{
		BaseProcessor: models.NewBaseProcessor("text", bufferSize),
		extensions:    extensions,
	}
}

// CanHandle implements the Processor interface
// Demonstrates string operations and loops
func (p *TextProcessor) CanHandle(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	// Demonstrates range loop over slice
	for _, supported := range p.extensions {
		if ext == supported {
			return true
		}
	}
	return false
}

// Process implements the Processor interface
// Demonstrates error handling and multiple return values
func (p *TextProcessor) Process(ctx context.Context, path string) (models.ProcessResult, error) {
	// Initialize result with embedded struct
	result := models.ProcessResult{
		FileInfo: models.FileInfo{
			Path:      path,
			Type:      "text",
			Processed: time.Now(),
		},
	}

	// Get file info
	info, err := os.Stat(path)
	if err != nil {
		result.Error = fmt.Errorf("failed to get file info: %w", err)
		return result, result.Error
	}

	// Update file info
	result.Size = info.Size()
	result.Modified = info.ModTime()

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		result.Error = fmt.Errorf("failed to open file: %w", err)
		return result, result.Error
	}
	defer file.Close()

	// Process the file content
	// Demonstrates multiple assignment from function return
	start := time.Now()
	result.Lines, result.Words, result.Bytes, err = p.readLines(file)
	result.Duration = time.Since(start)

	if err != nil {
		result.Error = fmt.Errorf("failed to process file: %w", err)
		return result, result.Error
	}

	return result, nil
}

// readLines counts lines, words, and bytes in a reader
// Demonstrates working with io.Reader and multiple return values
func (p *TextProcessor) readLines(reader io.Reader) (lines, words, bytes int, err error) {
	// Create a buffer for reading
	// Demonstrates array usage
	buf := make([]byte, 4096)

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

// SupportedExtensions demonstrates a method returning a slice
func (p *TextProcessor) SupportedExtensions() []string {
	// Demonstrates creating a new slice
	result := make([]string, len(p.extensions))
	// Demonstrates copy
	copy(result, p.extensions)
	return result
}

// AddExtension demonstrates method with pointer receiver
func (p *TextProcessor) AddExtension(ext string) {
	// Demonstrates string manipulation
	ext = strings.ToLower(strings.TrimSpace(ext))
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	// Demonstrates slice append
	p.extensions = append(p.extensions, ext)
}
