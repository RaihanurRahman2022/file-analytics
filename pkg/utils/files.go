package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// FileFilter is a function type that determines if a file should be processed
// Demonstrates function type definition
type FileFilter func(path string) bool

// WalkFunc is a function type that processes a file
// Demonstrates function type for callbacks
type WalkFunc func(path string) error

// CreateExtensionFilter demonstrates closure creation
// Returns a FileFilter that checks file extensions
func CreateExtensionFilter(extensions ...string) FileFilter {
	// Convert extensions to lowercase for comparison
	// Demonstrates slice manipulation
	lowerExt := make([]string, len(extensions))
	for i, ext := range extensions {
		lowerExt[i] = strings.ToLower(ext)
	}

	// Return a closure that captures lowerExt
	return func(path string) bool {
		ext := strings.ToLower(filepath.Ext(path))
		// Demonstrates slice searching
		for _, validExt := range lowerExt {
			if ext == validExt {
				return true
			}
		}
		return false
	}
}

// CreateSizeFilter demonstrates closure with multiple parameters
// Returns a FileFilter that checks file size
func CreateSizeFilter(minSize, maxSize int64) FileFilter {
	return func(path string) bool {
		info, err := os.Stat(path)
		if err != nil {
			return false
		}

		size := info.Size()
		// Demonstrates logical operators
		return (minSize <= 0 || size >= minSize) &&
			(maxSize <= 0 || size <= maxSize)
	}
}

// CombineFilters demonstrates variadic functions
// Returns a FileFilter that combines multiple filters with AND logic
func CombineFilters(filters ...FileFilter) FileFilter {
	return func(path string) bool {
		// Demonstrates short-circuit evaluation
		for _, filter := range filters {
			if !filter(path) {
				return false
			}
		}
		return true
	}
}

// WalkFiles demonstrates recursive directory traversal
// Processes files in a directory tree that match the filter
func WalkFiles(root string, filter FileFilter, fn WalkFunc) error {
	// Demonstrates recursive function
	var walkFn filepath.WalkFunc = func(path string, info os.FileInfo, err error) error {
		// Error handling
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Apply filter
		if filter != nil && !filter(path) {
			return nil
		}

		// Process file
		return fn(path)
	}

	// Start recursive walk
	return filepath.Walk(root, walkFn)
}

// CountFiles demonstrates a simple use of WalkFiles
// Returns the number of files matching the filter
func CountFiles(root string, filter FileFilter) (count int, err error) {
	// Demonstrates closure capturing a variable
	err = WalkFiles(root, filter, func(path string) error {
		count++
		return nil
	})
	return
}
