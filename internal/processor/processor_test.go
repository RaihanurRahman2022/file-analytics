package processor

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/RaihanurRahman2022/file-analytics/pkg/utils"
)

func TestJSONProcessor(t *testing.T) {
	// Create test JSON file
	testData := []map[string]interface{}{
		{"name": "test1", "value": 1},
		{"name": "test2", "value": 2},
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.json")

	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	encoder := json.NewEncoder(file)
	for _, data := range testData {
		if err := encoder.Encode(data); err != nil {
			t.Fatalf("Failed to write test data: %v", err)
		}
	}
	file.Close()

	// Test processor
	processor := NewJSONProcessor(4096)
	if !processor.CanHandle(testFile) {
		t.Error("Processor should handle JSON files")
	}

	result, err := processor.Process(context.Background(), testFile)
	if err != nil {
		t.Fatalf("Failed to process file: %v", err)
	}

	if result.Lines != len(testData) {
		t.Errorf("Expected %d lines, got %d", len(testData), result.Lines)
	}
}

func TestCSVProcessor(t *testing.T) {
	// Create test CSV file
	testData := [][]string{
		{"name", "value"},
		{"test1", "1"},
		{"test2", "2"},
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.csv")

	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	writer := csv.NewWriter(file)
	for _, row := range testData {
		if err := writer.Write(row); err != nil {
			t.Fatalf("Failed to write test data: %v", err)
		}
	}
	writer.Flush()
	file.Close()

	// Test processor
	processor := NewCSVProcessor(4096)
	if !processor.CanHandle(testFile) {
		t.Error("Processor should handle CSV files")
	}

	result, err := processor.Process(context.Background(), testFile)
	if err != nil {
		t.Fatalf("Failed to process file: %v", err)
	}

	if result.Lines != len(testData) {
		t.Errorf("Expected %d lines, got %d", len(testData), result.Lines)
	}
}

func TestHashAndBase64(t *testing.T) {
	// Create test file
	testContent := "Hello, World!"
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test hash
	hash, err := utils.HashFile(testFile)
	if err != nil {
		t.Fatalf("Failed to hash file: %v", err)
	}
	if hash == "" {
		t.Error("Hash should not be empty")
	}

	// Test base64 encoding
	encoded, err := utils.Base64EncodeFile(testFile)
	if err != nil {
		t.Fatalf("Failed to encode file: %v", err)
	}
	if encoded == "" {
		t.Error("Encoded content should not be empty")
	}

	// Test base64 decoding
	outputFile := filepath.Join(tmpDir, "decoded.txt")
	if err := utils.Base64DecodeFile(encoded, outputFile); err != nil {
		t.Fatalf("Failed to decode file: %v", err)
	}

	decodedContent, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read decoded file: %v", err)
	}

	if string(decodedContent) != testContent {
		t.Error("Decoded content should match original content")
	}
}
