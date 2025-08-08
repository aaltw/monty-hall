package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/westhuis/monty-hall/pkg/stats"
)

// Config represents the application configuration
type Config struct {
	UI        UIConfig        `json:"ui"`
	Game      GameConfig      `json:"game"`
	Stats     StatsConfig     `json:"stats"`
	Education EducationConfig `json:"education"`
	Version   string          `json:"version"`
}

// UIConfig contains user interface configuration options
type UIConfig struct {
	ColorScheme    string `json:"color_scheme"`    // "default", "high-contrast", "colorblind-safe"
	AnimationSpeed int    `json:"animation_speed"` // 0=disabled, 1=slow, 2=normal, 3=fast
	ShowTutorial   bool   `json:"show_tutorial"`   // Show tutorial on first run
	AutoSave       bool   `json:"auto_save"`       // Auto-save statistics
	TerminalWidth  int    `json:"terminal_width"`  // Preferred terminal width (0=auto)
	TerminalHeight int    `json:"terminal_height"` // Preferred terminal height (0=auto)
	ShowAnimations bool   `json:"show_animations"` // Enable/disable animations
	ReducedMotion  bool   `json:"reduced_motion"`  // Accessibility: reduce motion
	HighContrast   bool   `json:"high_contrast"`   // Accessibility: high contrast mode
	LargeText      bool   `json:"large_text"`      // Accessibility: larger text
}

// GameConfig contains game-specific configuration options
type GameConfig struct {
	AutoAdvance     bool   `json:"auto_advance"`     // Auto-advance through game phases
	ConfirmChoices  bool   `json:"confirm_choices"`  // Require confirmation for choices
	ShowProbability bool   `json:"show_probability"` // Show probability information
	DefaultStrategy string `json:"default_strategy"` // "switch", "stay", or "ask"
	ShowHints       bool   `json:"show_hints"`       // Show strategy hints
	PlaySounds      bool   `json:"play_sounds"`      // Play sound effects (if supported)
}

// StatsConfig contains statistics configuration options
type StatsConfig struct {
	AutoExport      bool               `json:"auto_export"`      // Auto-export stats periodically
	ExportFormat    stats.ExportFormat `json:"export_format"`    // Default export format
	MaxHistorySize  int                `json:"max_history_size"` // Maximum number of games to keep in history
	ShowDailyStats  bool               `json:"show_daily_stats"` // Show daily statistics breakdown
	ShowStreaks     bool               `json:"show_streaks"`     // Show win/loss streaks
	ShowAdvanced    bool               `json:"show_advanced"`    // Show advanced statistics
	ExportDirectory string             `json:"export_directory"` // Directory for exported files
}

// EducationConfig contains educational feature configuration
type EducationConfig struct {
	ShowExplanations bool `json:"show_explanations"` // Show probability explanations
	ShowMath         bool `json:"show_math"`         // Show mathematical details
	InteractiveMode  bool `json:"interactive_mode"`  // Enable interactive tutorials
	SkipTutorial     bool `json:"skip_tutorial"`     // Skip tutorial on startup
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	exportDir := filepath.Join(homeDir, "Documents", "MontyHall")

	return &Config{
		UI: UIConfig{
			ColorScheme:    "default",
			AnimationSpeed: 2, // Normal speed
			ShowTutorial:   true,
			AutoSave:       true,
			TerminalWidth:  0, // Auto-detect
			TerminalHeight: 0, // Auto-detect
			ShowAnimations: true,
			ReducedMotion:  false,
			HighContrast:   false,
			LargeText:      false,
		},
		Game: GameConfig{
			AutoAdvance:     false,
			ConfirmChoices:  false,
			ShowProbability: true,
			DefaultStrategy: "ask", // Ask user each time
			ShowHints:       true,
			PlaySounds:      false, // Disabled by default for terminal app
		},
		Stats: StatsConfig{
			AutoExport:      false,
			ExportFormat:    stats.ExportJSON,
			MaxHistorySize:  10000,
			ShowDailyStats:  true,
			ShowStreaks:     true,
			ShowAdvanced:    false,
			ExportDirectory: exportDir,
		},
		Education: EducationConfig{
			ShowExplanations: true,
			ShowMath:         false, // Keep it simple by default
			InteractiveMode:  true,
			SkipTutorial:     false,
		},
		Version: "1.0.0",
	}
}

// GetConfigDir returns the configuration directory for the application
func GetConfigDir() (string, error) {
	var configDir string

	switch runtime.GOOS {
	case "windows":
		configDir = os.Getenv("APPDATA")
		if configDir == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("failed to get user home directory: %w", err)
			}
			configDir = filepath.Join(homeDir, "AppData", "Roaming")
		}
		configDir = filepath.Join(configDir, "MontyHall")
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		configDir = filepath.Join(homeDir, "Library", "Application Support", "MontyHall")
	default: // Linux and other Unix-like systems
		configDir = os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("failed to get user home directory: %w", err)
			}
			configDir = filepath.Join(homeDir, ".config")
		}
		configDir = filepath.Join(configDir, "monty-hall")
	}

	return configDir, nil
}

// GetConfigPath returns the full path to the configuration file
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.json"), nil
}

// Validate validates the configuration and returns any errors
func (c *Config) Validate() error {
	// Validate UI config
	validColorSchemes := map[string]bool{
		"default":         true,
		"high-contrast":   true,
		"colorblind-safe": true,
	}
	if !validColorSchemes[c.UI.ColorScheme] {
		return fmt.Errorf("invalid color scheme: %s", c.UI.ColorScheme)
	}

	if c.UI.AnimationSpeed < 0 || c.UI.AnimationSpeed > 3 {
		return fmt.Errorf("animation speed must be between 0 and 3, got %d", c.UI.AnimationSpeed)
	}

	if c.UI.TerminalWidth < 0 || c.UI.TerminalHeight < 0 {
		return fmt.Errorf("terminal dimensions cannot be negative")
	}

	// Validate Game config
	validStrategies := map[string]bool{
		"switch": true,
		"stay":   true,
		"ask":    true,
	}
	if !validStrategies[c.Game.DefaultStrategy] {
		return fmt.Errorf("invalid default strategy: %s", c.Game.DefaultStrategy)
	}

	// Validate Stats config
	if c.Stats.MaxHistorySize < 0 {
		return fmt.Errorf("max history size cannot be negative")
	}

	// Validate export format
	validFormats := []stats.ExportFormat{stats.ExportJSON, stats.ExportCSV, stats.ExportText}
	validFormat := false
	for _, format := range validFormats {
		if c.Stats.ExportFormat == format {
			validFormat = true
			break
		}
	}
	if !validFormat {
		return fmt.Errorf("invalid export format: %v", c.Stats.ExportFormat)
	}

	return nil
}

// ApplyDefaults fills in any missing values with defaults
func (c *Config) ApplyDefaults() {
	defaults := DefaultConfig()

	// Apply UI defaults
	if c.UI.ColorScheme == "" {
		c.UI.ColorScheme = defaults.UI.ColorScheme
	}
	if c.UI.AnimationSpeed == 0 && !c.UI.ReducedMotion {
		c.UI.AnimationSpeed = defaults.UI.AnimationSpeed
	}

	// Apply Game defaults
	if c.Game.DefaultStrategy == "" {
		c.Game.DefaultStrategy = defaults.Game.DefaultStrategy
	}

	// Apply Stats defaults
	if c.Stats.MaxHistorySize == 0 {
		c.Stats.MaxHistorySize = defaults.Stats.MaxHistorySize
	}
	if c.Stats.ExportDirectory == "" {
		c.Stats.ExportDirectory = defaults.Stats.ExportDirectory
	}

	// Apply version if missing
	if c.Version == "" {
		c.Version = defaults.Version
	}
}

// Clone creates a deep copy of the configuration
func (c *Config) Clone() *Config {
	data, _ := json.Marshal(c)
	var clone Config
	json.Unmarshal(data, &clone)
	return &clone
}

// String returns a string representation of the configuration
func (c *Config) String() string {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Sprintf("Config{error: %v}", err)
	}
	return string(data)
}
