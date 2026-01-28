// Package main provides the local test API for the Borehole Edge-Scoring Engine.
// This API is for local development and testing of the SMS scoring pipeline.
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"borehole/core/pkg/engine"
	"borehole/core/pkg/parser"
)

const (
	defaultAddr     = ":8080"
	readTimeout     = 10 * time.Second
	writeTimeout    = 10 * time.Second
	shutdownTimeout = 5 * time.Second
)

func main() {
	// Logger setup
	logger := log.New(os.Stdout, "[borehole] ", log.LstdFlags|log.Lshortfile)

	// Initialize dependencies
	p := parser.NewParser()
	// Engine is now a singleton, initialized on first use

	// Setup router using Go 1.22+ ServeMux
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /health", healthHandler)

	// Main scoring endpoint
	mux.HandleFunc("POST /v1/score", scoreHandler(p, logger))

	// Create server
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = defaultAddr
	}

	server := &http.Server{
		Addr:         addr,
		Handler:      loggingMiddleware(logger, mux),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	// Graceful shutdown setup
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Printf("Starting server on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server error: %v", err)
		}
	}()

	<-done
	logger.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server shutdown failed: %v", err)
	}

	logger.Println("Server stopped gracefully")
}

// ScoreRequest is the JSON input for the scoring endpoint.
type ScoreRequest struct {
	Logs []string `json:"logs"`
}

// ScoreResponse is the JSON output for the scoring endpoint.
type ScoreResponse struct {
	Score    float64   `json:"score"`
	Features []float64 `json:"features"`
	TxnCount int       `json:"txn_count"`
	Message  string    `json:"message,omitempty"`
}

// healthHandler returns a simple health check response.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

// scoreHandler processes SMS logs and returns a credit score.
func scoreHandler(p parser.Parser, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request
		var req ScoreRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, "invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Validate input
		if len(req.Logs) == 0 {
			writeError(w, "logs array cannot be empty", http.StatusBadRequest)
			return
		}

		// Parse SMS logs
		txns, err := p.ParseLogs(r.Context(), req.Logs)
		if err != nil {
			logger.Printf("Parse error: %v", err)
			writeError(w, "failed to parse logs", http.StatusInternalServerError)
			return
		}

		// Generate feature vector
		features := engine.MapFeatures(txns)

		// Calculate score using the ML Engine
		mlEngine, err := engine.GetEngine()
		var score float64
		if err != nil {
			logger.Printf("Engine init error: %v", err)
			// Fallback to 0 or handle error appropriately.
			// For this test API, we'll return 0 and log the error.
		} else {
			score = mlEngine.Predict(features)
		}

		// Build response
		resp := ScoreResponse{
			Score:    score,
			Features: features,
			TxnCount: len(txns),
		}

		if len(txns) == 0 {
			resp.Message = "no transactions could be parsed from provided logs"
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

// writeError sends a JSON error response.
func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// loggingMiddleware logs HTTP requests.
func loggingMiddleware(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create response writer wrapper to capture status
		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		logger.Printf("%s %s %d %v",
			r.Method,
			r.URL.Path,
			wrapped.status,
			time.Since(start),
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
