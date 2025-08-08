package game

import (
	"testing"
)

func TestSecureRandom(t *testing.T) {
	sr := NewSecureRandom()

	t.Run("Intn returns values in range", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			val := sr.Intn(10)
			if val < 0 || val >= 10 {
				t.Errorf("Intn(10) returned %d, expected value in [0, 10)", val)
			}
		}
	})

	t.Run("Intn with zero returns zero", func(t *testing.T) {
		val := sr.Intn(0)
		if val != 0 {
			t.Errorf("Intn(0) returned %d, expected 0", val)
		}
	})

	t.Run("Intn with negative returns zero", func(t *testing.T) {
		val := sr.Intn(-5)
		if val != 0 {
			t.Errorf("Intn(-5) returned %d, expected 0", val)
		}
	})

	t.Run("Float64 returns values in range", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			val := sr.Float64()
			if val < 0.0 || val >= 1.0 {
				t.Errorf("Float64() returned %f, expected value in [0.0, 1.0)", val)
			}
		}
	})

	t.Run("Global functions work", func(t *testing.T) {
		val1 := SecureIntn(10)
		if val1 < 0 || val1 >= 10 {
			t.Errorf("SecureIntn(10) returned %d, expected value in [0, 10)", val1)
		}

		val2 := SecureFloat64()
		if val2 < 0.0 || val2 >= 1.0 {
			t.Errorf("SecureFloat64() returned %f, expected value in [0.0, 1.0)", val2)
		}
	})

	t.Run("Values are reasonably distributed", func(t *testing.T) {
		// Test that we get different values (not all the same)
		values := make(map[int]bool)
		for i := 0; i < 50; i++ {
			val := sr.Intn(10)
			values[val] = true
		}

		// We should get at least 3 different values in 50 tries
		if len(values) < 3 {
			t.Errorf("Got only %d different values in 50 tries, expected at least 3", len(values))
		}
	})
}

func TestSecureRandomFallback(t *testing.T) {
	// This test ensures the fallback mechanism works
	// We can't easily test crypto/rand failure, but we can test the fallback RNG
	sr := NewSecureRandom()

	// Test that fallback RNG is initialized
	if sr.fallbackRNG == nil {
		t.Error("Fallback RNG not initialized")
	}

	// Test multiple calls to ensure consistency
	for i := 0; i < 10; i++ {
		val := sr.Intn(100)
		if val < 0 || val >= 100 {
			t.Errorf("Intn(100) returned %d, expected value in [0, 100)", val)
		}
	}
}

func BenchmarkSecureIntn(b *testing.B) {
	sr := NewSecureRandom()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sr.Intn(100)
	}
}

func BenchmarkSecureFloat64(b *testing.B) {
	sr := NewSecureRandom()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sr.Float64()
	}
}
