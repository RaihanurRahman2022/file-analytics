package processor

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/RaihanurRahman2022/file-analytics/pkg/models"
)

// XMLProcessor implements the Processor interface for XML files
type XMLProcessor struct {
	*models.BaseProcessor
}

// NewXMLProcessor creates a new XML processor
func NewXMLProcessor(bufferSize int) *XMLProcessor {
	return &XMLProcessor{
		BaseProcessor: models.NewBaseProcessor("xml", bufferSize),
	}
}

// CanHandle implements the Processor interface
func (p *XMLProcessor) CanHandle(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".xml"
}

// Process implements the Processor interface
func (p *XMLProcessor) Process(ctx context.Context, path string) (models.ProcessResult, error) {
	result := models.ProcessResult{
		FileInfo: models.FileInfo{
			Path:      path,
			Type:      "xml",
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

	// Process the XML file
	start := time.Now()
	decoder := xml.NewDecoder(file)

	// Count elements and calculate size
	var elements, textNodes int
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.Error = fmt.Errorf("failed to decode XML: %w", err)
			return result, result.Error
		}

		switch t := token.(type) {
		case xml.StartElement:
			elements++
		case xml.CharData:
			if len(strings.TrimSpace(string(t))) > 0 {
				textNodes++
			}
		}
	}

	result.Duration = time.Since(start)
	result.Lines = elements + textNodes // Count both elements and text nodes
	result.Words = textNodes            // Use text nodes as word count
	result.Bytes = int(info.Size())

	return result, nil
}
