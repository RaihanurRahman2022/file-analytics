package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RaihanurRahman2022/file-analytics/internal/api"
	"github.com/RaihanurRahman2022/file-analytics/internal/monitor"
	"github.com/stretchr/testify/assert"
)

func TestFileAnalysisAPI(t *testing.T) {
	// Setup
	metrics := monitor.NewMetrics()
	handlers := api.NewHandlers(metrics)
	server := httptest.NewServer(handlers.Router())
	defer server.Close()

	// Test cases
	tests := []struct {
		name       string
		endpoint   string
		method     string
		body       interface{}
		wantStatus int
	}{
		{
			name:       "Analyze directory",
			endpoint:   "/api/v1/analyze",
			method:     "POST",
			body:       map[string]string{"path": "testdata"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "Hash file",
			endpoint:   "/api/v1/hash",
			method:     "POST",
			body:       map[string]string{"file": "testdata/sample.txt"},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			body, _ := json.Marshal(tt.body)
			req, err := http.NewRequest(tt.method, server.URL+tt.endpoint, bytes.NewBuffer(body))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Send request
			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			// Check response
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestMetricsAPI(t *testing.T) {
	// Setup
	metrics := monitor.NewMetrics()
	handlers := api.NewHandlers(metrics)
	server := httptest.NewServer(handlers.Router())
	defer server.Close()

	// Test metrics endpoint
	resp, err := http.Get(server.URL + "/api/v1/metrics")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
