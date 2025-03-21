# API Documentation

## Core Interfaces

### Processor Interface
```go
type Processor interface {
    CanHandle(path string) bool
    Process(ctx context.Context, path string) (models.ProcessResult, error)
}
```

### ProcessResult Structure
```go
type ProcessResult struct {
    Lines     int64
    Words     int64
    Bytes     int64
    Duration  time.Duration
}
```

## Available Processors

### TextProcessor
```go
func NewTextProcessor(bufferSize int) *TextProcessor
```
- Handles plain text files
- Counts lines, words, and bytes
- Uses buffered reading for efficiency

### JSONProcessor
```go
func NewJSONProcessor(bufferSize int) *JSONProcessor
```
- Processes JSON files
- Validates JSON structure
- Counts tokens and bytes

### CSVProcessor
```go
func NewCSVProcessor(bufferSize int) *CSVProcessor
```
- Handles CSV files
- Supports custom delimiters
- Counts rows and fields

### XMLProcessor
```go
func NewXMLProcessor(bufferSize int) *XMLProcessor
```
- Processes XML files
- Validates XML structure
- Counts elements and attributes

## Utility Functions

### File Operations
```go
func HashFile(path string) (string, error)
func Base64EncodeFile(path string) (string, error)
func Base64DecodeFile(content string, outputPath string) error
```

### File Walking
```go
func WalkFiles(path string, filter func(string) bool, handler func(string) error) error
```

## Configuration

### Processor Configuration
```go
type Config struct {
    BufferSize    int
    MaxWorkers    int
    RateLimit     time.Duration
    Extensions    []string
}
```

## Error Types
```go
type ProcessingError struct {
    File    string
    Message string
    Err     error
}

type ValidationError struct {
    File    string
    Message string
}
```

## Usage Examples

### Basic File Processing
```go
processor := NewTextProcessor(4096)
result, err := processor.Process(context.Background(), "file.txt")
```

### Concurrent Processing
```go
pool := NewWorkerPool(4, 100)
err := pool.ProcessFiles("directory", processors)
```

### File Hashing
```go
hash, err := HashFile("file.txt")
```

### Base64 Operations
```go
encoded, err := Base64EncodeFile("file.txt")
err = Base64DecodeFile(encoded, "output.txt")
```

## Best Practices
1. Always use context for cancellation
2. Handle errors appropriately
3. Use buffered operations for large files
4. Implement rate limiting for file operations
5. Clean up resources after processing
6. Use appropriate buffer sizes
7. Implement proper error handling
8. Follow Go idioms and patterns 