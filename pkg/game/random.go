package game

import (
	"crypto/rand"
	"math"
	"math/big"
	mathrand "math/rand"
	"time"
)

// SecureRandom provides cryptographically secure random number generation
// with fallback to math/rand if crypto/rand fails
type SecureRandom struct {
	fallbackRNG *mathrand.Rand
}

// NewSecureRandom creates a new SecureRandom instance
func NewSecureRandom() *SecureRandom {
	// Initialize fallback with time-based seed
	source := mathrand.NewSource(time.Now().UnixNano())
	fallbackRNG := mathrand.New(source)

	return &SecureRandom{
		fallbackRNG: fallbackRNG,
	}
}

// Intn returns a secure random integer in [0, n)
// Falls back to math/rand if crypto/rand fails
func (sr *SecureRandom) Intn(n int) int {
	if n <= 0 {
		return 0
	}

	// Try crypto/rand first
	bigN := big.NewInt(int64(n))
	result, err := rand.Int(rand.Reader, bigN)
	if err == nil {
		return int(result.Int64())
	}

	// Fallback to math/rand if crypto/rand fails
	return sr.fallbackRNG.Intn(n)
}

// Float64 returns a secure random float64 in [0.0, 1.0)
// Falls back to math/rand if crypto/rand fails
func (sr *SecureRandom) Float64() float64 {
	// Try to generate a secure random float64
	// Generate 8 random bytes and convert to float64
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err == nil {
		// Convert bytes to uint64, then to float64 in [0, 1)
		var value uint64
		for i, b := range bytes {
			value |= uint64(b) << (8 * i)
		}
		// Convert to float64 in range [0, 1)
		return float64(value) / float64(math.MaxUint64)
	}

	// Fallback to math/rand if crypto/rand fails
	return sr.fallbackRNG.Float64()
}

// Global secure random instance for convenience
var globalSecureRandom = NewSecureRandom()

// SecureIntn returns a secure random integer in [0, n) using the global instance
func SecureIntn(n int) int {
	return globalSecureRandom.Intn(n)
}

// SecureFloat64 returns a secure random float64 in [0.0, 1.0) using the global instance
func SecureFloat64() float64 {
	return globalSecureRandom.Float64()
}
