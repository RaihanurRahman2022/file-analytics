package templates

import (
	"bytes"
	"html/template"
	"time"
)

// ReportData represents the data structure for report generation
type ReportData struct {
	Title          string
	Timestamp      time.Time
	Files          []FileInfo
	Statistics     Statistics
	Errors         []string
	ProcessingTime time.Duration
}

// FileInfo represents information about a processed file
type FileInfo struct {
	Name           string
	Size           int64
	Type           string
	WordCount      int
	LineCount      int
	Hash           string
	ProcessingTime time.Duration
}

// Statistics represents overall processing statistics
type Statistics struct {
	TotalFiles   int
	TotalSize    int64
	TotalWords   int
	TotalLines   int
	SuccessCount int
	ErrorCount   int
	AverageTime  time.Duration
}

// HTMLTemplate is the template for HTML reports
const HTMLTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background: #f5f5f5; padding: 20px; border-radius: 5px; }
        .stats { margin: 20px 0; }
        .file-list { margin: 20px 0; }
        .error-list { color: red; }
        table { width: 100%; border-collapse: collapse; }
        th, td { padding: 8px; border: 1px solid #ddd; text-align: left; }
        th { background: #f5f5f5; }
    </style>
</head>
<body>
    <div class="header">
        <h1>{{.Title}}</h1>
        <p>Generated at: {{.Timestamp.Format "2006-01-02 15:04:05"}}</p>
    </div>

    <div class="stats">
        <h2>Statistics</h2>
        <table>
            <tr><th>Total Files</th><td>{{.Statistics.TotalFiles}}</td></tr>
            <tr><th>Total Size</th><td>{{.Statistics.TotalSize}} bytes</td></tr>
            <tr><th>Total Words</th><td>{{.Statistics.TotalWords}}</td></tr>
            <tr><th>Total Lines</th><td>{{.Statistics.TotalLines}}</td></tr>
            <tr><th>Success Count</th><td>{{.Statistics.SuccessCount}}</td></tr>
            <tr><th>Error Count</th><td>{{.Statistics.ErrorCount}}</td></tr>
            <tr><th>Average Processing Time</th><td>{{.Statistics.AverageTime}}</td></tr>
        </table>
    </div>

    <div class="file-list">
        <h2>Processed Files</h2>
        <table>
            <tr>
                <th>Name</th>
                <th>Size</th>
                <th>Type</th>
                <th>Words</th>
                <th>Lines</th>
                <th>Hash</th>
                <th>Processing Time</th>
            </tr>
            {{range .Files}}
            <tr>
                <td>{{.Name}}</td>
                <td>{{.Size}}</td>
                <td>{{.Type}}</td>
                <td>{{.WordCount}}</td>
                <td>{{.LineCount}}</td>
                <td>{{.Hash}}</td>
                <td>{{.ProcessingTime}}</td>
            </tr>
            {{end}}
        </table>
    </div>

    {{if .Errors}}
    <div class="error-list">
        <h2>Errors</h2>
        <ul>
            {{range .Errors}}
            <li>{{.}}</li>
            {{end}}
        </ul>
    </div>
    {{end}}

    <p>Total Processing Time: {{.ProcessingTime}}</p>
</body>
</html>
`

// MarkdownTemplate is the template for Markdown reports
const MarkdownTemplate = `# {{.Title}}

Generated at: {{.Timestamp.Format "2006-01-02 15:04:05"}}

## Statistics

| Metric | Value |
|--------|-------|
| Total Files | {{.Statistics.TotalFiles}} |
| Total Size | {{.Statistics.TotalSize}} bytes |
| Total Words | {{.Statistics.TotalWords}} |
| Total Lines | {{.Statistics.TotalLines}} |
| Success Count | {{.Statistics.SuccessCount}} |
| Error Count | {{.Statistics.ErrorCount}} |
| Average Processing Time | {{.Statistics.AverageTime}} |

## Processed Files

| Name | Size | Type | Words | Lines | Hash | Processing Time |
|------|------|------|-------|-------|------|-----------------|
{{range .Files}}| {{.Name}} | {{.Size}} | {{.Type}} | {{.WordCount}} | {{.LineCount}} | {{.Hash}} | {{.ProcessingTime}} |
{{end}}

{{if .Errors}}
## Errors

{{range .Errors}}- {{.}}
{{end}}
{{end}}

Total Processing Time: {{.ProcessingTime}}
`

// GenerateHTMLReport generates an HTML report from the provided data
func GenerateHTMLReport(data ReportData) (string, error) {
	tmpl, err := template.New("html").Parse(HTMLTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// GenerateMarkdownReport generates a Markdown report from the provided data
func GenerateMarkdownReport(data ReportData) (string, error) {
	tmpl, err := template.New("markdown").Parse(MarkdownTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
