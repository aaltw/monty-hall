package stats

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/westhuis/monty-hall/pkg/game"
)

// Custom error types for better error handling
var (
	ErrNilStats     = errors.New("stats cannot be nil")
	ErrFileNotFound = errors.New("stats file not found")
)

const (
	DefaultStatsFileName = "monty_hall_stats.json"
	DefaultStatsDir      = ".monty-hall"
)

type PersistenceManager struct {
	filePath string
}

func NewPersistenceManager(customPath ...string) *PersistenceManager {
	var filePath string

	if len(customPath) > 0 && customPath[0] != "" {
		filePath = customPath[0]
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			filePath = DefaultStatsFileName
		} else {
			statsDir := filepath.Join(homeDir, DefaultStatsDir)
			filePath = filepath.Join(statsDir, DefaultStatsFileName)
		}
	}

	return &PersistenceManager{
		filePath: filePath,
	}
}

func (pm *PersistenceManager) Save(stats *GameStats) error {
	if stats == nil {
		return ErrNilStats
	}

	dir := filepath.Dir(pm.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal stats: %w", err)
	}

	if err := os.WriteFile(pm.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write stats file: %w", err)
	}

	return nil
}

func (pm *PersistenceManager) Load() (*GameStats, error) {
	if !pm.Exists() {
		return &GameStats{
			DailyStats: make(map[string]DailyStats),
		}, nil
	}

	data, err := os.ReadFile(pm.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read stats file: %w", err)
	}

	var stats GameStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stats: %w", err)
	}

	if stats.DailyStats == nil {
		stats.DailyStats = make(map[string]DailyStats)
	}

	return &stats, nil
}

func (pm *PersistenceManager) Exists() bool {
	_, err := os.Stat(pm.filePath)
	return err == nil
}

func (pm *PersistenceManager) Delete() error {
	if !pm.Exists() {
		return nil
	}

	if err := os.Remove(pm.filePath); err != nil {
		return fmt.Errorf("failed to delete stats file: %w", err)
	}

	return nil
}

func (pm *PersistenceManager) GetFilePath() string {
	return pm.filePath
}

func (pm *PersistenceManager) GetFileSize() (int64, error) {
	if !pm.Exists() {
		return 0, nil
	}

	info, err := os.Stat(pm.filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to get file info: %w", err)
	}

	return info.Size(), nil
}

func (pm *PersistenceManager) Backup(backupPath string) error {
	if !pm.Exists() {
		return fmt.Errorf("stats file does not exist")
	}

	data, err := os.ReadFile(pm.filePath)
	if err != nil {
		return fmt.Errorf("failed to read stats file: %w", err)
	}

	backupDir := filepath.Dir(backupPath)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
}

func (pm *PersistenceManager) Restore(backupPath string) error {
	if _, err := os.Stat(backupPath); err != nil {
		return fmt.Errorf("backup file does not exist: %w", err)
	}

	data, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	var stats GameStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return fmt.Errorf("invalid backup file format: %w", err)
	}

	dir := filepath.Dir(pm.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(pm.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to restore stats file: %w", err)
	}

	return nil
}

type StatsManager struct {
	collector   *Collector
	persistence *PersistenceManager
}

func NewStatsManager(customPath ...string) *StatsManager {
	persistence := NewPersistenceManager(customPath...)

	stats, err := persistence.Load()
	if err != nil {
		// Create fresh stats if loading fails
		stats = &GameStats{
			DailyStats: make(map[string]DailyStats),
		}
		// Try to save the fresh stats to ensure the file system is writable
		if saveErr := persistence.Save(stats); saveErr != nil {
			// If we can't save, at least log the issue (in a real app, you'd use proper logging)
			// For now, we'll continue with in-memory stats
		}
	}

	collector := &Collector{stats: stats}

	return &StatsManager{
		collector:   collector,
		persistence: persistence,
	}
}

func (sm *StatsManager) RecordGame(result *game.GameResult) error {
	if err := sm.collector.RecordGame(result); err != nil {
		return err
	}

	return sm.persistence.Save(sm.collector.GetStats())
}

func (sm *StatsManager) GetStats() *GameStats {
	return sm.collector.GetStats()
}

func (sm *StatsManager) GetSummary() StatsSummary {
	return sm.collector.GetSummary()
}

func (sm *StatsManager) GetFilteredGames(filter StatsFilter) []GameRecord {
	return sm.collector.GetFilteredGames(filter)
}

func (sm *StatsManager) GetStatsFilePath() string {
	return sm.persistence.filePath
}

func (sm *StatsManager) Reset() error {
	sm.collector.Reset()
	return sm.persistence.Save(sm.collector.GetStats())
}

func (sm *StatsManager) Backup(backupPath string) error {
	return sm.persistence.Backup(backupPath)
}

func (sm *StatsManager) Restore(backupPath string) error {
	if err := sm.persistence.Restore(backupPath); err != nil {
		return err
	}

	stats, err := sm.persistence.Load()
	if err != nil {
		return err
	}

	sm.collector = &Collector{stats: stats}
	return nil
}

func (sm *StatsManager) GetFilePath() string {
	return sm.persistence.GetFilePath()
}

func (sm *StatsManager) GetFileSize() (int64, error) {
	return sm.persistence.GetFileSize()
}
