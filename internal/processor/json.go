package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/RaihanurRahman2022/file-analytics/pkg/models"
)

// JSONProcessor implements the Processor interface for JSON files
type JSONProcessor struct {
	*models.BaseProcessor
}

// NewJSONProcessor creates a new JSON processor
func NewJSONProcessor(bufferSize int) *JSONProcessor {
	return &JSONProcessor{
		BaseProcessor: models.NewBaseProcessor("json", bufferSize),
	}
}

// CanHandle implements the Processor interface
func (p *JSONProcessor) CanHandle(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), ".json")
}

// Process implements the Processor interface
func (p *JSONProcessor) Process(ctx context.Context, path string) (models.ProcessResult, error) {
	result := models.ProcessResult{
		FileInfo: models.FileInfo{
			Path:      path,
			Type:      "json",
			Processed: time.Now(),
		},
	}

	// Get file info
	info, err := os.Stat(path)
	if err != nil {
		result.Error = fmt.Errorf("failed to get file info: %w", err)
		return result, result.Error
	}

	result.Size = info.Size()
	result.Modified = info.ModTime()

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		result.Error = fmt.Errorf("failed to open file: %w", err)
		return result, result.Error
	}
	defer file.Close()

	// Process the JSON file
	start := time.Now()
	decoder := json.NewDecoder(file)

	// Count objects and calculate size
	var count int
	for {
		var json interface{}
		if err := decoder.Decode(&json); err != nil {
			if err == io.EOF {
				break
			}
			result.Error = fmt.Errorf("failed to decode JSON: %w", err)
			return result, result.Error
		}
		count++
	}

	result.Duration = time.Since(start)
	result.Lines = count // In JSON, each object is counted as a line
	result.Bytes = int(info.Size())

	return result, nil
}
