package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const Prefix = "prompt_live_"

// GenerateAPIKey generates a Stripe-like API Key:
// prefix (12 chars) + 64 random hex chars = 76 characters total.
// Returns raw key, sha256 hash of raw key, and a masked key representation.
func GenerateAPIKey() (string, string, string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", "", "", err
	}

	randomHex := hex.EncodeToString(b)
	rawKey := Prefix + randomHex

	hash := HashAPIKey(rawKey)

	// Masked key: e.g. prompt_live_xxxx...1234 (last 4 chars)
	masked := fmt.Sprintf("%sxxxx...%s", Prefix, randomHex[len(randomHex)-4:])

	return rawKey, hash, masked, nil
}

// HashAPIKey returns the SHA-256 hex string of the raw API key.
func HashAPIKey(rawKey string) string {
	h := sha256.New()
	h.Write([]byte(rawKey))
	return hex.EncodeToString(h.Sum(nil))
}
