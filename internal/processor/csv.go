package processor

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/RaihanurRahman2022/file-analytics/pkg/models"
)

// CSVProcessor implements the Processor interface for CSV files
type CSVProcessor struct {
	*models.BaseProcessor
}

// NewCSVProcessor creates a new CSV processor
func NewCSVProcessor(bufferSize int) *CSVProcessor {
	return &CSVProcessor{
		BaseProcessor: models.NewBaseProcessor("csv", bufferSize),
	}
}

// CanHandle implements the Processor interface
func (p *CSVProcessor) CanHandle(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".csv" || ext == ".tsv"
}

// Process implements the Processor interface
func (p *CSVProcessor) Process(ctx context.Context, path string) (models.ProcessResult, error) {
	result := models.ProcessResult{
		FileInfo: models.FileInfo{
			Path:      path,
			Type:      "csv",
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

	// Create CSV reader
	reader := csv.NewReader(file)

	// Detect delimiter based on file extension
	if strings.HasSuffix(strings.ToLower(path), ".tsv") {
		reader.Comma = '\t'
	}

	// Process the CSV file
	start := time.Now()

	// Read header
	_, err = reader.Read()
	if err != nil {
		result.Error = fmt.Errorf("failed to read CSV header: %w", err)
		return result, result.Error
	}

	// Count rows and calculate statistics
	var rows, words int
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.Error = fmt.Errorf("failed to read CSV row: %w", err)
			return result, result.Error
		}
		rows++
		words += len(record)
	}

	result.Duration = time.Since(start)
	result.Lines = rows + 1 // Include header row
	result.Words = words
	result.Bytes = int(info.Size())

	return result, nil
}
