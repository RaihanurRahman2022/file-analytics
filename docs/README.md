# File Analytics Tool Documentation

## Overview
File Analytics is a comprehensive file analysis tool that provides both CLI and web interfaces for:
- File content analysis (lines, words, bytes)
- File format processing (Text, JSON, CSV, XML)
- File hashing (SHA256)
- Base64 encoding/decoding
- Real-time monitoring and metrics
- Report generation (HTML, Markdown)

## Project Structure
```
file-analytics/
├── cmd/
│   ├── analyzer/         # CLI application entry point
│   └── server/          # HTTP API server entry point
├── internal/
│   ├── processor/       # File processing implementations
│   ├── worker/         # Worker pool implementation
│   ├── monitor/        # Monitoring and metrics
│   └── api/           # HTTP API handlers
├── pkg/
│   ├── models/         # Data models
│   ├── errors/         # Error handling
│   ├── utils/          # Shared utilities
│   └── templates/      # Report generation templates
├── configs/            # Configuration files
├── test/              # Test files and fixtures
│   └── integration/   # Integration tests
├── web/              # Web interface assets
│   ├── css/         # Stylesheets
│   └── js/          # JavaScript files
└── docs/            # Documentation
    ├── testdata/    # Sample test files
    ├── README.md    # This file
    ├── WORKFLOW.md  # Detailed workflow
    ├── TESTING.md   # Testing guide
    └── API.md      # API documentation
```

## Quick Start

### CLI Tool
1. Build the CLI tool:
   ```bash
   go build -o analyzer cmd/analyzer/main.go
   ```

2. Run the analyzer:
   ```bash
   ./analyzer analyze [path]
   ./analyzer hash [file]
   ./analyzer encode [file]
   ./analyzer decode [base64] [output]
   ```

### Web Interface
1. Build the server:
   ```bash
   go build -o server cmd/server/main.go
   ```

2. Start the server:
   ```bash
   ./server -port 8080
   ```

3. Open web browser:
   ```
   http://localhost:8080
   ```

## Features

### CLI Features
- File analysis with multiple processors
- File hashing
- Base64 encoding/decoding
- Report generation
- Progress tracking

### Web Interface Features
- Interactive file analysis
- Real-time results display
- File hash calculation
- Responsive design
- Error handling and feedback

### API Features
- RESTful endpoints
- JSON responses
- Metrics monitoring
- Rate limiting
- Graceful shutdown

## Documentation Sections
- [Workflow Guide](WORKFLOW.md) - Detailed explanation of how the tool works
- [Testing Guide](TESTING.md) - How to test the tool with sample data
- [API Documentation](API.md) - Detailed API documentation

## Requirements
- Go 1.21 or later
- Modern web browser (for web interface)
- Dependencies listed in go.mod

## Development
1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Run tests: `go test ./...`
4. Build: `go build ./...`

## Contributing
1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License
This project is licensed under the MIT License - see the LICENSE file for details. 