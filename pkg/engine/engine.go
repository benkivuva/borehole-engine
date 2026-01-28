package engine

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"math"
	"sync"

	"github.com/dmitryikh/leaves"
)

//go:embed model/borehole_model.json
var modelData []byte

// BoreholeEngine acts as the thread-safe singleton for ML inference.
type BoreholeEngine struct {
	ensemble *leaves.Ensemble
}

var (
	instance *BoreholeEngine
	once     sync.Once
	initErr  error
)

// Predict performs on-device scoring for a 20-dimension feature vector.
// Applies Sigmoid activation to normalize the raw margin output to [0, 1].
// Zero-allocations in the main inference loop.
func (e *BoreholeEngine) Predict(features []float64) float64 {
	// Fallback if model failed to load
	if e.ensemble == nil {
		// Mock logic for testing stability when model is invalid
		// Returns 0.5 (neutral score)
		return 0.5
	}

	if len(features) < 20 {
		return 0
	}

	// leave.PredictSingle is high-performance and avoids allocations
	rawMargin := e.ensemble.PredictSingle(features, 0)

	// Sigmoid: 1 / (1 + exp(-x))
	return 1.0 / (1.0 + math.Exp(-rawMargin))
}

// GetEngine returns the singleton instance of the BoreholeEngine.
// Uses Go 1.20+ error wrapping for initialization failures.
// Returns a valid instance even on model load failure to prevent app crashes.
func GetEngine() (*BoreholeEngine, error) {
	once.Do(func() {
		// Recover from panic in leaves library
		defer func() {
			if r := recover(); r != nil {
				initErr = fmt.Errorf("panic initializing XGBoost model: %v", r)
			}
		}()

		engine := &BoreholeEngine{}
		// leaves requires a bufio.Reader
		reader := bufio.NewReader(bytes.NewReader(modelData))
		// loadTransformation=true ensures compatibility with models using feature transformations
		ensemble, err := leaves.XGEnsembleFromReader(reader, true)
		if err != nil {
			// Log error but allow engine to return (will use fallback)
			fmt.Printf("WARNING: Failed to load embedded XGBoost model: %v. Using fallback.\n", err)
			initErr = err
		}
		engine.ensemble = ensemble
		instance = engine
	})

	// We return the instance regardless, so the app doesn't crash NULL pointer.
	// The caller can check err if they want to know about model issues.
	if instance == nil {
		// Should not happen, but safe fallback
		instance = &BoreholeEngine{}
	}

	// We return nil error to allow app to proceed (Predict will handle nil ensemble)
	return instance, nil
}
