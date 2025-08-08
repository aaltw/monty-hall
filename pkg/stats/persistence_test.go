package stats

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/westhuis/monty-hall/pkg/game"
)

func TestNewPersistenceManager(t *testing.T) {
	pm := NewPersistenceManager()

	if pm == nil {
		t.Fatal("PersistenceManager should not be nil")
	}

	if pm.filePath == "" {
		t.Error("File path should not be empty")
	}
}

func TestNewPersistenceManagerCustomPath(t *testing.T) {
	customPath := "/tmp/test_stats.json"
	pm := NewPersistenceManager(customPath)

	if pm.GetFilePath() != customPath {
		t.Errorf("Expected custom path %s, got %s", customPath, pm.GetFilePath())
	}
}

func TestSaveAndLoad(t *testing.T) {
	tempFile := filepath.Join(os.TempDir(), "test_stats.json")
	defer os.Remove(tempFile)

	pm := NewPersistenceManager(tempFile)

	stats := &GameStats{
		TotalGames: 5,
		TotalWins:  3,
		DailyStats: make(map[string]DailyStats),
	}

	err := pm.Save(stats)
	if err != nil {
		t.Errorf("Unexpected error saving: %v", err)
	}

	if !pm.Exists() {
		t.Error("File should exist after save")
	}

	loadedStats, err := pm.Load()
	if err != nil {
		t.Errorf("Unexpected error loading: %v", err)
	}

	if loadedStats.TotalGames != stats.TotalGames {
		t.Errorf("Expected %d total games, got %d", stats.TotalGames, loadedStats.TotalGames)
	}

	if loadedStats.TotalWins != stats.TotalWins {
		t.Errorf("Expected %d total wins, got %d", stats.TotalWins, loadedStats.TotalWins)
	}
}

func TestSaveNilStats(t *testing.T) {
	tempFile := filepath.Join(os.TempDir(), "test_stats.json")
	defer os.Remove(tempFile)

	pm := NewPersistenceManager(tempFile)

	err := pm.Save(nil)
	if err == nil {
		t.Error("Expected error for nil stats")
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	tempFile := filepath.Join(os.TempDir(), "non_existent_stats.json")

	pm := NewPersistenceManager(tempFile)

	stats, err := pm.Load()
	if err != nil {
		t.Errorf("Unexpected error loading non-existent file: %v", err)
	}

	if stats == nil {
		t.Error("Stats should not be nil for non-existent file")
	}

	if stats.TotalGames != 0 {
		t.Error("New stats should have 0 total games")
	}

	if stats.DailyStats == nil {
		t.Error("DailyStats should be initialized")
	}
}

func TestDelete(t *testing.T) {
	tempFile := filepath.Join(os.TempDir(), "test_delete_stats.json")

	pm := NewPersistenceManager(tempFile)

	stats := &GameStats{
		TotalGames: 1,
		DailyStats: make(map[string]DailyStats),
	}

	pm.Save(stats)

	if !pm.Exists() {
		t.Error("File should exist before delete")
	}

	err := pm.Delete()
	if err != nil {
		t.Errorf("Unexpected error deleting: %v", err)
	}

	if pm.Exists() {
		t.Error("File should not exist after delete")
	}
}

func TestGetFileSize(t *testing.T) {
	tempFile := filepath.Join(os.TempDir(), "test_size_stats.json")
	defer os.Remove(tempFile)

	pm := NewPersistenceManager(tempFile)

	size, err := pm.GetFileSize()
	if err != nil {
		t.Errorf("Unexpected error getting size: %v", err)
	}

	if size != 0 {
		t.Errorf("Expected size 0 for non-existent file, got %d", size)
	}

	stats := &GameStats{
		TotalGames: 1,
		DailyStats: make(map[string]DailyStats),
	}

	pm.Save(stats)

	size, err = pm.GetFileSize()
	if err != nil {
		t.Errorf("Unexpected error getting size: %v", err)
	}

	if size <= 0 {
		t.Errorf("Expected positive size, got %d", size)
	}
}

func TestBackupAndRestore(t *testing.T) {
	tempFile := filepath.Join(os.TempDir(), "test_backup_stats.json")
	backupFile := filepath.Join(os.TempDir(), "test_backup_stats_backup.json")
	defer os.Remove(tempFile)
	defer os.Remove(backupFile)

	pm := NewPersistenceManager(tempFile)

	stats := &GameStats{
		TotalGames: 10,
		TotalWins:  6,
		DailyStats: make(map[string]DailyStats),
	}

	pm.Save(stats)

	err := pm.Backup(backupFile)
	if err != nil {
		t.Errorf("Unexpected error creating backup: %v", err)
	}

	stats.TotalGames = 20
	pm.Save(stats)

	err = pm.Restore(backupFile)
	if err != nil {
		t.Errorf("Unexpected error restoring backup: %v", err)
	}

	restoredStats, err := pm.Load()
	if err != nil {
		t.Errorf("Unexpected error loading restored stats: %v", err)
	}

	if restoredStats.TotalGames != 10 {
		t.Errorf("Expected 10 total games after restore, got %d", restoredStats.TotalGames)
	}
}

func TestNewStatsManager(t *testing.T) {
	tempFile := filepath.Join(os.TempDir(), "test_manager_stats.json")
	defer os.Remove(tempFile)

	sm := NewStatsManager(tempFile)

	if sm == nil {
		t.Fatal("StatsManager should not be nil")
	}

	if sm.collector == nil {
		t.Error("Collector should not be nil")
	}

	if sm.persistence == nil {
		t.Error("Persistence should not be nil")
	}
}

func TestStatsManagerRecordGame(t *testing.T) {
	tempFile := filepath.Join(os.TempDir(), "test_manager_record_stats.json")
	defer os.Remove(tempFile)

	sm := NewStatsManager(tempFile)

	result := &game.GameResult{
		Won:            true,
		Strategy:       game.Switch,
		InitialChoice:  1,
		FinalChoice:    2,
		CarPosition:    2,
		HostOpenedDoor: 3,
		GameDuration:   time.Millisecond * 500,
		Timestamp:      time.Now(),
	}

	err := sm.RecordGame(result)
	if err != nil {
		t.Errorf("Unexpected error recording game: %v", err)
	}

	stats := sm.GetStats()
	if stats.TotalGames != 1 {
		t.Errorf("Expected 1 total game, got %d", stats.TotalGames)
	}

	sm2 := NewStatsManager(tempFile)
	stats2 := sm2.GetStats()
	if stats2.TotalGames != 1 {
		t.Errorf("Expected 1 total game after reload, got %d", stats2.TotalGames)
	}
}

func TestStatsManagerReset(t *testing.T) {
	tempFile := filepath.Join(os.TempDir(), "test_manager_reset_stats.json")
	defer os.Remove(tempFile)

	sm := NewStatsManager(tempFile)

	result := &game.GameResult{
		Won:            true,
		Strategy:       game.Switch,
		InitialChoice:  1,
		FinalChoice:    2,
		CarPosition:    2,
		HostOpenedDoor: 3,
		GameDuration:   time.Millisecond * 500,
		Timestamp:      time.Now(),
	}

	sm.RecordGame(result)

	if sm.GetStats().TotalGames != 1 {
		t.Error("Should have 1 game before reset")
	}

	err := sm.Reset()
	if err != nil {
		t.Errorf("Unexpected error resetting: %v", err)
	}

	if sm.GetStats().TotalGames != 0 {
		t.Error("Should have 0 games after reset")
	}

	sm2 := NewStatsManager(tempFile)
	if sm2.GetStats().TotalGames != 0 {
		t.Error("Should have 0 games after reload")
	}
}

func TestStatsManagerBackupRestore(t *testing.T) {
	tempFile := filepath.Join(os.TempDir(), "test_manager_backup_stats.json")
	backupFile := filepath.Join(os.TempDir(), "test_manager_backup_stats_backup.json")
	defer os.Remove(tempFile)
	defer os.Remove(backupFile)

	sm := NewStatsManager(tempFile)

	result := &game.GameResult{
		Won:            true,
		Strategy:       game.Switch,
		InitialChoice:  1,
		FinalChoice:    2,
		CarPosition:    2,
		HostOpenedDoor: 3,
		GameDuration:   time.Millisecond * 500,
		Timestamp:      time.Now(),
	}

	sm.RecordGame(result)

	err := sm.Backup(backupFile)
	if err != nil {
		t.Errorf("Unexpected error creating backup: %v", err)
	}

	sm.Reset()

	err = sm.Restore(backupFile)
	if err != nil {
		t.Errorf("Unexpected error restoring backup: %v", err)
	}

	if sm.GetStats().TotalGames != 1 {
		t.Error("Should have 1 game after restore")
	}
}
