package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/RaihanurRahman2022/file-analytics/internal/monitor"
)

// Server represents the HTTP API server
type Server struct {
	addr     string
	server   *http.Server
	handlers map[string]http.HandlerFunc
}

// NewServer creates a new HTTP API server
func NewServer(addr string) *Server {
	s := &Server{
		addr:     addr,
		handlers: make(map[string]http.HandlerFunc),
	}

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.withMiddleware(s.handleHealth))
	mux.HandleFunc("/metrics", s.withMiddleware(s.handleMetrics))

	s.server = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s
}

// Start begins listening for HTTP requests
func (s *Server) Start() error {
	log.Printf("Starting server on %s", s.addr)
	return s.server.ListenAndServe()
}

// Stop gracefully shuts down the server
func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// Middleware function type
type Middleware func(http.HandlerFunc) http.HandlerFunc

// withMiddleware applies common middleware to handlers
func (s *Server) withMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	// Apply middleware in order
	return s.logRequest(
		s.timeRequest(
			s.recoverPanic(handler),
		),
	)
}

// logRequest logs incoming HTTP requests
func (s *Server) logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next(w, r)
	}
}

// timeRequest measures request duration
func (s *Server) timeRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		duration := time.Since(start)
		log.Printf("Request processed in %v", duration)
	}
}

// recoverPanic recovers from panics in handlers
func (s *Server) recoverPanic(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

// handleHealth handles health check requests
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleMetrics handles metrics requests
func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := map[string]interface{}{
		"uptime": time.Since(time.Now()),
		"requests": map[string]int{
			"total":   100, // Example values
			"success": 95,
			"error":   5,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// Handlers represents the API handlers
type Handlers struct {
	metrics *monitor.MetricsCollector
	mux     *http.ServeMux
}

// NewHandlers creates new API handlers
func NewHandlers(metrics *monitor.MetricsCollector) *Handlers {
	h := &Handlers{
		metrics: metrics,
		mux:     http.NewServeMux(),
	}
	h.setupRoutes()
	return h
}

// Router returns the HTTP router
func (h *Handlers) Router() http.Handler {
	return h.mux
}

// setupRoutes configures API routes
func (h *Handlers) setupRoutes() {
	h.mux.HandleFunc("/api/v1/analyze", h.handleAnalyze)
	h.mux.HandleFunc("/api/v1/hash", h.handleHash)
	h.mux.HandleFunc("/api/v1/metrics", h.handleMetrics)
}

// handleAnalyze handles file analysis requests
func (h *Handlers) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// handleHash handles file hash requests
func (h *Handlers) handleHash(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// handleMetrics handles metrics requests
func (h *Handlers) handleMetrics(w http.ResponseWriter, r *http.Request) {
	processed, errors, avgDuration := h.metrics.GetMetrics()
	metrics := map[string]interface{}{
		"processed": processed,
		"errors":    errors,
		"duration":  avgDuration.String(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
