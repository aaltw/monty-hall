package config

import (
	"path/filepath"
	"testing"

	"github.com/westhuis/monty-hall/pkg/stats"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	// Test that default config is valid
	if err := config.Validate(); err != nil {
		t.Errorf("Default config should be valid: %v", err)
	}

	// Test specific defaults
	if config.UI.ColorScheme != "default" {
		t.Errorf("Expected default color scheme 'default', got '%s'", config.UI.ColorScheme)
	}

	if config.UI.AnimationSpeed != 2 {
		t.Errorf("Expected default animation speed 2, got %d", config.UI.AnimationSpeed)
	}

	if config.Game.DefaultStrategy != "ask" {
		t.Errorf("Expected default strategy 'ask', got '%s'", config.Game.DefaultStrategy)
	}

	if config.Stats.ExportFormat != stats.ExportJSON {
		t.Errorf("Expected default export format JSON, got %v", config.Stats.ExportFormat)
	}

	if config.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", config.Version)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		modifyFunc  func(*Config)
		expectError bool
	}{
		{
			name: "Valid config",
			modifyFunc: func(c *Config) {
				// No modifications - should be valid
			},
			expectError: false,
		},
		{
			name: "Invalid color scheme",
			modifyFunc: func(c *Config) {
				c.UI.ColorScheme = "invalid"
			},
			expectError: true,
		},
		{
			name: "Invalid animation speed - negative",
			modifyFunc: func(c *Config) {
				c.UI.AnimationSpeed = -1
			},
			expectError: true,
		},
		{
			name: "Invalid animation speed - too high",
			modifyFunc: func(c *Config) {
				c.UI.AnimationSpeed = 4
			},
			expectError: true,
		},
		{
			name: "Invalid terminal dimensions",
			modifyFunc: func(c *Config) {
				c.UI.TerminalWidth = -1
			},
			expectError: true,
		},
		{
			name: "Invalid default strategy",
			modifyFunc: func(c *Config) {
				c.Game.DefaultStrategy = "invalid"
			},
			expectError: true,
		},
		{
			name: "Invalid max history size",
			modifyFunc: func(c *Config) {
				c.Stats.MaxHistorySize = -1
			},
			expectError: true,
		},
		{
			name: "Valid edge cases",
			modifyFunc: func(c *Config) {
				c.UI.AnimationSpeed = 0    // Disabled
				c.UI.TerminalWidth = 0     // Auto
				c.Stats.MaxHistorySize = 0 // Will be set to default
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultConfig()
			tt.modifyFunc(config)

			err := config.Validate()
			if tt.expectError && err == nil {
				t.Error("Expected validation error, but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no validation error, but got: %v", err)
			}
		})
	}
}

func TestConfigApplyDefaults(t *testing.T) {
	config := &Config{
		UI: UIConfig{
			// Leave some fields empty
			AnimationSpeed: 0,
		},
		Game: GameConfig{
			// Leave strategy empty
		},
		Stats: StatsConfig{
			// Leave max history size as 0
		},
		// Leave version empty
	}

	config.ApplyDefaults()

	// Check that defaults were applied
	if config.UI.ColorScheme != "default" {
		t.Errorf("Expected color scheme to be set to default, got '%s'", config.UI.ColorScheme)
	}

	if config.Game.DefaultStrategy != "ask" {
		t.Errorf("Expected default strategy to be set to 'ask', got '%s'", config.Game.DefaultStrategy)
	}

	if config.Stats.MaxHistorySize != 10000 {
		t.Errorf("Expected max history size to be set to 10000, got %d", config.Stats.MaxHistorySize)
	}

	if config.Version != "1.0.0" {
		t.Errorf("Expected version to be set to '1.0.0', got '%s'", config.Version)
	}
}

func TestConfigClone(t *testing.T) {
	original := DefaultConfig()
	original.UI.ColorScheme = "high-contrast"
	original.Game.ShowHints = false

	clone := original.Clone()

	// Verify clone has same values
	if clone.UI.ColorScheme != original.UI.ColorScheme {
		t.Errorf("Clone color scheme doesn't match: expected '%s', got '%s'",
			original.UI.ColorScheme, clone.UI.ColorScheme)
	}

	if clone.Game.ShowHints != original.Game.ShowHints {
		t.Errorf("Clone show hints doesn't match: expected %t, got %t",
			original.Game.ShowHints, clone.Game.ShowHints)
	}

	// Verify they are separate objects
	clone.UI.ColorScheme = "colorblind-safe"
	if original.UI.ColorScheme == clone.UI.ColorScheme {
		t.Error("Clone is not independent of original")
	}
}

func TestGetConfigDir(t *testing.T) {
	configDir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("Failed to get config directory: %v", err)
	}

	if configDir == "" {
		t.Error("Config directory should not be empty")
	}

	// Should contain platform-appropriate path
	if !filepath.IsAbs(configDir) {
		t.Error("Config directory should be an absolute path")
	}
}

func TestGetConfigPath(t *testing.T) {
	configPath, err := GetConfigPath()
	if err != nil {
		t.Fatalf("Failed to get config path: %v", err)
	}

	if configPath == "" {
		t.Error("Config path should not be empty")
	}

	if !filepath.IsAbs(configPath) {
		t.Error("Config path should be an absolute path")
	}

	if filepath.Ext(configPath) != ".json" {
		t.Error("Config path should have .json extension")
	}
}

func TestConfigString(t *testing.T) {
	config := DefaultConfig()
	str := config.String()

	if str == "" {
		t.Error("Config string representation should not be empty")
	}

	// Should be valid JSON
	if str[0] != '{' || str[len(str)-1] != '}' {
		t.Error("Config string should be valid JSON")
	}
}

func TestGetColorSchemes(t *testing.T) {
	schemes := GetColorSchemes()

	expectedSchemes := []string{"default", "high-contrast", "colorblind-safe"}
	if len(schemes) != len(expectedSchemes) {
		t.Errorf("Expected %d color schemes, got %d", len(expectedSchemes), len(schemes))
	}

	for _, expected := range expectedSchemes {
		found := false
		for _, scheme := range schemes {
			if scheme == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected color scheme '%s' not found", expected)
		}
	}
}

func TestGetAnimationSpeeds(t *testing.T) {
	speeds := GetAnimationSpeeds()

	expectedSpeeds := []int{0, 1, 2, 3}
	if len(speeds) != len(expectedSpeeds) {
		t.Errorf("Expected %d animation speeds, got %d", len(expectedSpeeds), len(speeds))
	}

	for _, expected := range expectedSpeeds {
		if description, exists := speeds[expected]; !exists {
			t.Errorf("Expected animation speed %d not found", expected)
		} else if description == "" {
			t.Errorf("Animation speed %d should have a description", expected)
		}
	}
}

func TestGetDefaultStrategies(t *testing.T) {
	strategies := GetDefaultStrategies()

	expectedStrategies := []string{"switch", "stay", "ask"}
	if len(strategies) != len(expectedStrategies) {
		t.Errorf("Expected %d strategies, got %d", len(expectedStrategies), len(strategies))
	}

	for _, expected := range expectedStrategies {
		if description, exists := strategies[expected]; !exists {
			t.Errorf("Expected strategy '%s' not found", expected)
		} else if description == "" {
			t.Errorf("Strategy '%s' should have a description", expected)
		}
	}
}

// Benchmark tests
func BenchmarkDefaultConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DefaultConfig()
	}
}

func BenchmarkConfigValidate(b *testing.B) {
	config := DefaultConfig()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		config.Validate()
	}
}

func BenchmarkConfigClone(b *testing.B) {
	config := DefaultConfig()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		config.Clone()
	}
}
