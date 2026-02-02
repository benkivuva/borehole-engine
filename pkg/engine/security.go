package engine

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// CertificatePayload represents the data to be signed.
type CertificatePayload struct {
	Score     float64 `json:"score"`
	Timestamp int64   `json:"iat"` // Issued At (Unix)
	Expires   int64   `json:"exp"` // Expiry (Unix)
	UserID    string  `json:"uid"` // Anonymous ID (e.g., Device ID hash)
	Tampered  bool    `json:"tampered"`
}

// SecurityModule handles cryptographic operations.
type SecurityModule struct {
	publicKey  ed25519.PublicKey
	privateKey ed25519.PrivateKey
	mu         sync.RWMutex
}

var (
	secInstance *SecurityModule
	secOnce     sync.Once
)

// GetSecurityModule returns the singleton security module.
// In a real app, keys would be loaded from a secure vault.
// Here, we generate a fresh pair on startup for demonstration.
func GetSecurityModule() *SecurityModule {
	secOnce.Do(func() {
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			// simplified panic for critical security failure in init
			panic(fmt.Sprintf("failed to generate ed25519 keys: %v", err))
		}
		secInstance = &SecurityModule{
			publicKey:  pub,
			privateKey: priv,
		}
	})
	return secInstance
}

// IssueCertificate creates a signed payload for a credit score.
// Returns two strings: formatted payload (JSON) and the Base64 signature.
func (s *SecurityModule) IssueCertificate(score float64, uid string) (string, string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 1. Create Payload
	payload := CertificatePayload{
		Score:     score,
		Timestamp: time.Now().Unix(),
		Expires:   time.Now().Add(24 * time.Hour).Unix(),
		UserID:    uid,
		Tampered:  false, // Hardcoded engine is immutable by design
	}

	// 2. Serialize
	data, err := json.Marshal(payload)
	if err != nil {
		return "", "", fmt.Errorf("marshal error: %v", err)
	}

	// 3. Sign
	signature := ed25519.Sign(s.privateKey, data)

	// 4. Encode
	// We return the raw JSON string (so the verifier knows what was signed)
	// and the Base64 signature.
	return string(data), base64.StdEncoding.EncodeToString(signature), nil
}

// VerifyCertificate checks if a score claim is valid and signed by this engine.
// Returns true if valid.
func (s *SecurityModule) VerifyCertificate(payloadJSON string, signatureB64 string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 1. Decode Signature
	sig, err := base64.StdEncoding.DecodeString(signatureB64)
	if err != nil {
		return false, fmt.Errorf("invalid base64 signature: %v", err)
	}

	// 2. Verify
	isValid := ed25519.Verify(s.publicKey, []byte(payloadJSON), sig)
	return isValid, nil
}

// GetPublicKeyBase64 returns the public key to display or share.
func (s *SecurityModule) GetPublicKeyBase64() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return base64.StdEncoding.EncodeToString(s.publicKey)
}
