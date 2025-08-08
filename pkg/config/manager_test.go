package config

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/westhuis/monty-hall/pkg/stats"
)

func TestNewManager(t *testing.T) {
	// Test that NewManager creates a manager successfully
	// Note: This will use the actual config directory, which is fine for testing
	manager, err := NewManager()
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	if manager == nil {
		t.Fatal("Manager should not be nil")
	}

	// Should have created default config
	config := manager.Get()
	if config == nil {
		t.Fatal("Config should not be nil")
	}

	// Validate the config
	if err := config.Validate(); err != nil {
		t.Errorf("Default config should be valid: %v", err)
	}
}
func TestManagerLoadSave(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create manager with custom path
	manager := &Manager{
		configPath: configPath,
		watchers:   make([]func(*Config), 0),
	}

	// Set a config and save it
	testConfig := DefaultConfig()
	testConfig.UI.ColorScheme = "high-contrast"
	testConfig.Game.ShowHints = false

	manager.config = testConfig
	if err := manager.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Create new manager and load
	manager2 := &Manager{
		configPath: configPath,
		watchers:   make([]func(*Config), 0),
	}

	if err := manager2.Load(); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify loaded config matches saved config
	loadedConfig := manager2.Get()
	if loadedConfig.UI.ColorScheme != testConfig.UI.ColorScheme {
		t.Errorf("Color scheme mismatch: expected '%s', got '%s'",
			testConfig.UI.ColorScheme, loadedConfig.UI.ColorScheme)
	}

	if loadedConfig.Game.ShowHints != testConfig.Game.ShowHints {
		t.Errorf("Show hints mismatch: expected %t, got %t",
			testConfig.Game.ShowHints, loadedConfig.Game.ShowHints)
	}
}

func TestManagerUpdate(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	manager := &Manager{
		configPath: configPath,
		config:     DefaultConfig(),
		watchers:   make([]func(*Config), 0),
	}

	// Test watcher functionality
	watcherCalled := false
	var watchedConfig *Config
	manager.AddWatcher(func(config *Config) {
		watcherCalled = true
		watchedConfig = config
	})

	// Update config
	newConfig := DefaultConfig()
	newConfig.UI.ColorScheme = "colorblind-safe"

	if err := manager.Update(newConfig); err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	// Verify update
	currentConfig := manager.Get()
	if currentConfig.UI.ColorScheme != "colorblind-safe" {
		t.Errorf("Config not updated: expected 'colorblind-safe', got '%s'",
			currentConfig.UI.ColorScheme)
	}

	// Verify watcher was called
	if !watcherCalled {
		t.Error("Watcher should have been called")
	}

	if watchedConfig == nil || watchedConfig.UI.ColorScheme != "colorblind-safe" {
		t.Error("Watcher received incorrect config")
	}
}

func TestManagerUpdateSections(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	manager := &Manager{
		configPath: configPath,
		config:     DefaultConfig(),
		watchers:   make([]func(*Config), 0),
	}

	// Test UpdateUI
	newUI := UIConfig{
		ColorScheme:    "high-contrast",
		AnimationSpeed: 1,
		ShowTutorial:   false,
		AutoSave:       true,
	}

	if err := manager.UpdateUI(newUI); err != nil {
		t.Fatalf("Failed to update UI config: %v", err)
	}

	config := manager.Get()
	if config.UI.ColorScheme != "high-contrast" {
		t.Errorf("UI config not updated: expected 'high-contrast', got '%s'",
			config.UI.ColorScheme)
	}

	// Test UpdateGame
	newGame := GameConfig{
		AutoAdvance:     true,
		ConfirmChoices:  true,
		ShowProbability: false,
		DefaultStrategy: "switch",
		ShowHints:       false,
		PlaySounds:      true,
	}

	if err := manager.UpdateGame(newGame); err != nil {
		t.Fatalf("Failed to update game config: %v", err)
	}

	config = manager.Get()
	if config.Game.DefaultStrategy != "switch" {
		t.Errorf("Game config not updated: expected 'switch', got '%s'",
			config.Game.DefaultStrategy)
	}

	// Test UpdateStats
	newStats := StatsConfig{
		AutoExport:      true,
		ExportFormat:    stats.ExportCSV,
		MaxHistorySize:  5000,
		ShowDailyStats:  false,
		ShowStreaks:     false,
		ShowAdvanced:    true,
		ExportDirectory: "/tmp/exports",
	}

	if err := manager.UpdateStats(newStats); err != nil {
		t.Fatalf("Failed to update stats config: %v", err)
	}

	config = manager.Get()
	if config.Stats.ExportFormat != stats.ExportCSV {
		t.Errorf("Stats config not updated: expected CSV format, got %v",
			config.Stats.ExportFormat)
	}

	// Test UpdateEducation
	newEducation := EducationConfig{
		ShowExplanations: false,
		ShowMath:         true,
		InteractiveMode:  false,
		SkipTutorial:     true,
	}

	if err := manager.UpdateEducation(newEducation); err != nil {
		t.Fatalf("Failed to update education config: %v", err)
	}

	config = manager.Get()
	if !config.Education.SkipTutorial {
		t.Error("Education config not updated: expected SkipTutorial to be true")
	}
}

func TestManagerReset(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	manager := &Manager{
		configPath: configPath,
		config:     DefaultConfig(),
		watchers:   make([]func(*Config), 0),
	}

	// Modify config
	manager.config.UI.ColorScheme = "high-contrast"
	manager.config.Game.DefaultStrategy = "switch"

	// Reset to defaults
	if err := manager.Reset(); err != nil {
		t.Fatalf("Failed to reset config: %v", err)
	}

	// Verify reset
	config := manager.Get()
	defaultConfig := DefaultConfig()

	if config.UI.ColorScheme != defaultConfig.UI.ColorScheme {
		t.Errorf("Config not reset: expected '%s', got '%s'",
			defaultConfig.UI.ColorScheme, config.UI.ColorScheme)
	}

	if config.Game.DefaultStrategy != defaultConfig.Game.DefaultStrategy {
		t.Errorf("Config not reset: expected '%s', got '%s'",
			defaultConfig.Game.DefaultStrategy, config.Game.DefaultStrategy)
	}
}

func TestManagerBackupRestore(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	manager := &Manager{
		configPath: configPath,
		config:     DefaultConfig(),
		watchers:   make([]func(*Config), 0),
	}

	// Modify config
	manager.config.UI.ColorScheme = "high-contrast"
	if err := manager.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Create backup
	if err := manager.Backup(); err != nil {
		t.Fatalf("Failed to create backup: %v", err)
	}

	// Modify config again
	manager.config.UI.ColorScheme = "colorblind-safe"
	if err := manager.Save(); err != nil {
		t.Fatalf("Failed to save modified config: %v", err)
	}

	// Find backup file
	files, err := filepath.Glob(configPath + ".backup.*")
	if err != nil {
		t.Fatalf("Failed to find backup files: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("No backup files found")
	}

	backupPath := files[0]

	// Restore from backup
	if err := manager.RestoreFromBackup(backupPath); err != nil {
		t.Fatalf("Failed to restore from backup: %v", err)
	}

	// Verify restoration
	config := manager.Get()
	if config.UI.ColorScheme != "high-contrast" {
		t.Errorf("Config not restored: expected 'high-contrast', got '%s'",
			config.UI.ColorScheme)
	}
}

func TestManagerUtilityMethods(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	manager := &Manager{
		configPath: configPath,
		config:     DefaultConfig(),
		watchers:   make([]func(*Config), 0),
	}

	// Test Exists (should be false initially)
	if manager.Exists() {
		t.Error("Config file should not exist initially")
	}

	// Save config
	if err := manager.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Test Exists (should be true now)
	if !manager.Exists() {
		t.Error("Config file should exist after saving")
	}

	// Test GetConfigPath
	if manager.GetConfigPath() != configPath {
		t.Errorf("Config path mismatch: expected '%s', got '%s'",
			configPath, manager.GetConfigPath())
	}

	// Test GetSize
	size, err := manager.GetSize()
	if err != nil {
		t.Fatalf("Failed to get config size: %v", err)
	}

	if size <= 0 {
		t.Error("Config file size should be greater than 0")
	}

	// Test GetModTime
	modTime, err := manager.GetModTime()
	if err != nil {
		t.Fatalf("Failed to get mod time: %v", err)
	}

	if modTime.IsZero() {
		t.Error("Mod time should not be zero")
	}

	// Test IsDefault (should be true for default config)
	if !manager.IsDefault() {
		t.Error("Default config should be identified as default")
	}

	// Modify config and test IsDefault again
	manager.config.UI.ColorScheme = "high-contrast"
	if manager.IsDefault() {
		t.Error("Modified config should not be identified as default")
	}

	// Test Delete
	if err := manager.Delete(); err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}

	if manager.Exists() {
		t.Error("Config file should not exist after deletion")
	}
}

func TestManagerInvalidConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	manager := &Manager{
		configPath: configPath,
		config:     DefaultConfig(),
		watchers:   make([]func(*Config), 0),
	}

	// Try to update with invalid config
	invalidConfig := DefaultConfig()
	invalidConfig.UI.ColorScheme = "invalid-scheme"

	err := manager.Update(invalidConfig)
	if err == nil {
		t.Error("Should have failed to update with invalid config")
	}

	// Verify original config is unchanged
	config := manager.Get()
	if config.UI.ColorScheme != "default" {
		t.Error("Original config should be unchanged after failed update")
	}
}

func TestManagerConcurrency(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	manager := &Manager{
		configPath: configPath,
		config:     DefaultConfig(),
		watchers:   make([]func(*Config), 0),
	}

	// Test concurrent reads
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			config := manager.Get()
			if config == nil {
				t.Error("Config should not be nil")
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		select {
		case <-done:
			// Success
		case <-time.After(time.Second):
			t.Fatal("Timeout waiting for concurrent reads")
		}
	}
}

// Benchmark tests
func BenchmarkManagerGet(b *testing.B) {
	manager := &Manager{
		config:   DefaultConfig(),
		watchers: make([]func(*Config), 0),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Get()
	}
}

func BenchmarkManagerUpdate(b *testing.B) {
	tempDir := b.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	manager := &Manager{
		configPath: configPath,
		config:     DefaultConfig(),
		watchers:   make([]func(*Config), 0),
	}

	config := DefaultConfig()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		config.UI.AnimationSpeed = i % 4
		manager.Update(config)
	}
}
