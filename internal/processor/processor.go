package processor

import (
	"context"

	"github.com/RaihanurRahman2022/file-analytics/pkg/models"
)

// Processor defines the interface for file processors
type Processor interface {
	// CanHandle determines if the processor can handle the given file
	CanHandle(path string) bool
	// Process handles the file and returns processing results
	Process(ctx context.Context, path string) (models.ProcessResult, error)
}
