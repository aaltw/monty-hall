package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Manager handles configuration loading, saving, and management
type Manager struct {
	config     *Config
	configPath string
	mutex      sync.RWMutex
	watchers   []func(*Config)
}

// NewManager creates a new configuration manager
func NewManager() (*Manager, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	manager := &Manager{
		configPath: configPath,
		watchers:   make([]func(*Config), 0),
	}

	// Try to load existing config, create default if not found
	if err := manager.Load(); err != nil {
		if os.IsNotExist(err) {
			// Create default config and save it
			manager.config = DefaultConfig()
			if saveErr := manager.Save(); saveErr != nil {
				return nil, fmt.Errorf("failed to save default config: %w", saveErr)
			}
		} else {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
	}

	return manager, nil
}

// Load loads the configuration from disk
func (m *Manager) Load() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults for any missing values
	config.ApplyDefaults()

	// Validate the configuration
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	m.config = &config

	// Notify watchers
	for _, watcher := range m.watchers {
		watcher(m.config)
	}

	return nil
}

// Save saves the current configuration to disk
func (m *Manager) Save() error {
	m.mutex.RLock()
	config := m.config.Clone()
	m.mutex.RUnlock()

	// Ensure config directory exists
	configDir := filepath.Dir(m.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to JSON with indentation
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to temporary file first, then rename (atomic operation)
	tempPath := m.configPath + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	if err := os.Rename(tempPath, m.configPath); err != nil {
		os.Remove(tempPath) // Clean up temp file
		return fmt.Errorf("failed to save config file: %w", err)
	}

	return nil
}

// Get returns a copy of the current configuration
func (m *Manager) Get() *Config {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.config.Clone()
}

// Update updates the configuration with the provided config
func (m *Manager) Update(newConfig *Config) error {
	// Validate the new configuration
	if err := newConfig.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	m.mutex.Lock()
	m.config = newConfig.Clone()
	m.mutex.Unlock()

	// Save to disk
	if err := m.Save(); err != nil {
		return fmt.Errorf("failed to save updated config: %w", err)
	}

	// Notify watchers
	for _, watcher := range m.watchers {
		watcher(m.config)
	}

	return nil
}

// UpdateUI updates only the UI configuration
func (m *Manager) UpdateUI(uiConfig UIConfig) error {
	m.mutex.Lock()
	m.config.UI = uiConfig
	config := m.config.Clone()
	m.mutex.Unlock()

	return m.Update(config)
}

// UpdateGame updates only the game configuration
func (m *Manager) UpdateGame(gameConfig GameConfig) error {
	m.mutex.Lock()
	m.config.Game = gameConfig
	config := m.config.Clone()
	m.mutex.Unlock()

	return m.Update(config)
}

// UpdateStats updates only the stats configuration
func (m *Manager) UpdateStats(statsConfig StatsConfig) error {
	m.mutex.Lock()
	m.config.Stats = statsConfig
	config := m.config.Clone()
	m.mutex.Unlock()

	return m.Update(config)
}

// UpdateEducation updates only the education configuration
func (m *Manager) UpdateEducation(educationConfig EducationConfig) error {
	m.mutex.Lock()
	m.config.Education = educationConfig
	config := m.config.Clone()
	m.mutex.Unlock()

	return m.Update(config)
}

// Reset resets the configuration to defaults
func (m *Manager) Reset() error {
	defaultConfig := DefaultConfig()
	return m.Update(defaultConfig)
}

// Backup creates a backup of the current configuration
func (m *Manager) Backup() error {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupPath := m.configPath + ".backup." + timestamp

	m.mutex.RLock()
	config := m.config.Clone()
	m.mutex.RUnlock()

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config for backup: %w", err)
	}

	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
}

// RestoreFromBackup restores configuration from a backup file
func (m *Manager) RestoreFromBackup(backupPath string) error {
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse backup file: %w", err)
	}

	return m.Update(&config)
}

// AddWatcher adds a function to be called when configuration changes
func (m *Manager) AddWatcher(watcher func(*Config)) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.watchers = append(m.watchers, watcher)
}

// RemoveWatcher removes a configuration watcher (not implemented for simplicity)
// In a production app, you might want to implement this with watcher IDs

// GetConfigPath returns the path to the configuration file
func (m *Manager) GetConfigPath() string {
	return m.configPath
}

// Exists checks if the configuration file exists
func (m *Manager) Exists() bool {
	_, err := os.Stat(m.configPath)
	return err == nil
}

// Delete removes the configuration file
func (m *Manager) Delete() error {
	if err := os.Remove(m.configPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete config file: %w", err)
	}
	return nil
}

// GetSize returns the size of the configuration file in bytes
func (m *Manager) GetSize() (int64, error) {
	info, err := os.Stat(m.configPath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// GetModTime returns the last modification time of the configuration file
func (m *Manager) GetModTime() (time.Time, error) {
	info, err := os.Stat(m.configPath)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

// IsDefault checks if the current configuration matches the default configuration
func (m *Manager) IsDefault() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	defaultConfig := DefaultConfig()

	// Compare JSON representations for simplicity
	currentData, _ := json.Marshal(m.config)
	defaultData, _ := json.Marshal(defaultConfig)

	return string(currentData) == string(defaultData)
}

// GetColorSchemes returns available color schemes
func GetColorSchemes() []string {
	return []string{"default", "high-contrast", "colorblind-safe"}
}

// GetAnimationSpeeds returns available animation speeds with descriptions
func GetAnimationSpeeds() map[int]string {
	return map[int]string{
		0: "Disabled",
		1: "Slow",
		2: "Normal",
		3: "Fast",
	}
}

// GetDefaultStrategies returns available default strategies with descriptions
func GetDefaultStrategies() map[string]string {
	return map[string]string{
		"switch": "Always switch doors",
		"stay":   "Always stay with initial choice",
		"ask":    "Ask user each time",
	}
}
