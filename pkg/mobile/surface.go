package mobile

import (
	"context"
	"encoding/json"
	"fmt"

	"borehole/core/pkg/engine"
	"borehole/core/pkg/parser"
)

// MobileEngine is the JNI-compatible bridge for Android integration.
type MobileEngine struct {
	parser parser.Parser
}

// NewMobileEngine initializes the bridge. Engine is managed as a singleton.
func NewMobileEngine() *MobileEngine {
	return &MobileEngine{
		parser: parser.NewParser(),
	}
}

// CalculateBoreholeScore orchestrates the full ETL and Inference pipeline.
// Parser (ETL) -> Mapper (Transform) -> Engine (Inference) -> Result (Output).
func (m *MobileEngine) CalculateBoreholeScore(jsonLogs string) string {
	var logs []string

	if err := json.Unmarshal([]byte(jsonLogs), &logs); err != nil {
		return `{"error": "invalid_json_input"}`
	}

	// 1. ETL: Parse raw SMS logs into structured Transaction objects
	txns, err := m.parser.ParseLogs(context.Background(), logs)
	if err != nil {
		return fmt.Sprintf(`{"error": "parsing_failed", "details": "%v"}`, err)
	}

	// 2. Transform: Map transactions to 20-dimension feature vector
	features := engine.MapFeatures(txns)

	// 3. Inference: Get prediction from singleton ML engine
	mlEngine, err := engine.GetEngine()
	if err != nil {
		// Fallback or error reporting
		return fmt.Sprintf(`{"error": "engine_initialization_failed", "details": "%v"}`, err)
	}

	score := mlEngine.Predict(features)

	// 4. Output: Package results for React Native
	result := parser.ScoreResult{
		Score:    score,
		Features: features,
		TxnCount: len(txns),
	}

	resBytes, _ := json.Marshal(result)
	return string(resBytes)
}

// GenerateSignedScore creates a verifiable certificate for a given score.
// Returns a JSON string containing {payload, signature, public_key}.
func (m *MobileEngine) GenerateSignedScore(score float64) string {
	sec := engine.GetSecurityModule()

	// For MVP, we use a random Anonymous ID.
	// In production, this would be a hash of the device ID or user ID.
	uid := "anon_user_xyz"

	payloadStr, signature, err := sec.IssueCertificate(score, uid)
	if err != nil {
		return fmt.Sprintf(`{"error": "signing_failed", "details": "%v"}`, err)
	}

	response := map[string]string{
		"payload":    payloadStr,
		"signature":  signature,
		"public_key": sec.GetPublicKeyBase64(),
	}

	bytes, _ := json.Marshal(response)
	return string(bytes)
}
