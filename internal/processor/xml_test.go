package processor

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestXMLProcessor(t *testing.T) {
	// Create test XML file
	testXML := `<?xml version="1.0" encoding="UTF-8"?>
<root>
    <item>
        <name>test1</name>
        <value>1</value>
    </item>
    <item>
        <name>test2</name>
        <value>2</value>
    </item>
</root>`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.xml")

	if err := os.WriteFile(testFile, []byte(testXML), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test processor
	processor := NewXMLProcessor(4096)
	if !processor.CanHandle(testFile) {
		t.Error("Processor should handle XML files")
	}

	result, err := processor.Process(context.Background(), testFile)
	if err != nil {
		t.Fatalf("Failed to process file: %v", err)
	}

	// XML has 2 items with 2 elements each (name and value) plus root element
	expectedElements := 5 // root + 2 items + 2 names + 2 values
	if result.Lines < expectedElements {
		t.Errorf("Expected at least %d elements, got %d", expectedElements, result.Lines)
	}
} 