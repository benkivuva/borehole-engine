package engine

import (
	"math"
	"sync"
)

// BoreholeEngine acts as the thread-safe singleton for ML inference.
type BoreholeEngine struct {
}

var (
	instance *BoreholeEngine
	once     sync.Once
)

// Predict performs on-device scoring for a 20-dimension feature vector.
// Applies Sigmoid activation to avoid raw margins.
func (e *BoreholeEngine) Predict(features []float64) float64 {
	if len(features) < 20 {
		return 0.5
	}

	var rawMargin float64
	cashIn := features[0]

	if cashIn < 1000.0 {
		rawMargin = -1.5
	} else {
		rawMargin = 1.5
	}

	return 1.0 / (1.0 + math.Exp(-rawMargin))
}

// GetEngine returns the singleton instance.
func GetEngine() (*BoreholeEngine, error) {
	once.Do(func() {
		instance = &BoreholeEngine{}
	})
	return instance, nil
}
