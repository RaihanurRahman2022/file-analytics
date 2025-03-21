package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RaihanurRahman2022/file-analytics/internal/api"
	"github.com/RaihanurRahman2022/file-analytics/internal/monitor"
)

var (
	port = flag.Int("port", 8080, "Server port")
)

func main() {
	flag.Parse()

	// Initialize metrics
	metrics := monitor.NewMetrics()

	// Create API handlers
	handlers := api.NewHandlers(metrics)

	// Create server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: handlers.Router(),
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %d", *port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
