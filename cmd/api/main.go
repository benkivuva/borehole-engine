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
	e := engine.NewEngine()

	// Setup router using Go 1.22+ ServeMux
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /health", healthHandler)

	// Main scoring endpoint
	mux.HandleFunc("POST /v1/score", scoreHandler(p, e, logger))

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
func scoreHandler(p parser.Parser, e engine.Vectorizer, logger *log.Logger) http.HandlerFunc {
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
		features := e.Vectorize(txns)

		// Calculate score (simple weighted sum for demo)
		score := calculateScore(features)

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

// calculateScore computes a credit score from the feature vector.
// This is a simplified scoring function. In production, this would
// use an XGBoost model loaded via go:embed or file.
func calculateScore(features []float64) float64 {
	if len(features) < 15 {
		return 0
	}

	// Feature weights (simplified model)
	weights := []float64{
		0.10,  // total_income (positive)
		-0.05, // total_expenses (negative)
		0.15,  // net_flow (positive)
		0.05,  // avg_txn_amount
		0.02,  // txn_count
		-0.10, // income_regularity (lower is better)
		-0.25, // gambling_index (strongly negative)
		0.05,  // utility_ratio (positive - responsible spending)
		-0.15, // fuliza_usage (negative)
		0.10,  // fuliza_repay_rate (positive)
		-0.02, // p2p_ratio
		0.05,  // max_single_txn
		-0.05, // balance_volatility
		0.05,  // days_active
		0.02,  // avg_daily_volume
	}

	var score float64
	for i, weight := range weights {
		if i < len(features) {
			// Normalize feature contribution
			contribution := weight * normalizeFeature(features[i], i)
			score += contribution
		}
	}

	// Scale to 0-1 range using sigmoid-like function
	score = 1 / (1 + sigmoid(-score))

	// Clamp to valid range
	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}

	return score
}

// normalizeFeature scales features to comparable ranges.
func normalizeFeature(value float64, index int) float64 {
	// Scale factors based on expected ranges for Kenyan transactions
	scales := []float64{
		100000, // total_income (up to 100k KES)
		100000, // total_expenses
		50000,  // net_flow
		5000,   // avg_txn_amount
		100,    // txn_count
		1,      // income_regularity (already 0-1 scale)
		1,      // gambling_index (already 0-1 scale)
		1,      // utility_ratio (already 0-1 scale)
		1,      // fuliza_usage (already 0-1 scale)
		1,      // fuliza_repay_rate (already 0-1 scale)
		1,      // p2p_ratio (already 0-1 scale)
		50000,  // max_single_txn
		10000,  // balance_volatility
		30,     // days_active
		10000,  // avg_daily_volume
	}

	if index >= len(scales) || scales[index] == 0 {
		return value
	}

	return value / scales[index]
}

// sigmoid helper function.
func sigmoid(x float64) float64 {
	if x > 500 {
		return 1
	}
	if x < -500 {
		return 0
	}
	return 1 / (1 + exp(-x))
}

// exp is a simple exponential approximation.
func exp(x float64) float64 {
	// Use math.Exp via type assertion to avoid import cycle concerns
	// In production, use math.Exp directly
	const e = 2.718281828459045
	result := 1.0
	term := 1.0
	for i := 1; i < 20; i++ {
		term *= x / float64(i)
		result += term
	}
	return result
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
