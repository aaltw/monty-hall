package stats

import (
	"testing"
	"time"

	"github.com/westhuis/monty-hall/pkg/game"
)

func createTestGameResult(strategy game.PlayerStrategy, won bool) *game.GameResult {
	return &game.GameResult{
		Won:            won,
		Strategy:       strategy,
		InitialChoice:  1,
		FinalChoice:    2,
		CarPosition:    2,
		HostOpenedDoor: 3,
		GameDuration:   time.Millisecond * 500,
		Timestamp:      time.Now(),
	}
}

func TestNewCollector(t *testing.T) {
	collector := NewCollector()

	if collector == nil {
		t.Fatal("Collector should not be nil")
	}

	if collector.stats == nil {
		t.Fatal("Stats should not be nil")
	}

	if collector.stats.DailyStats == nil {
		t.Fatal("DailyStats should not be nil")
	}
}

func TestRecordGame(t *testing.T) {
	collector := NewCollector()

	result := createTestGameResult(game.Switch, true)

	err := collector.RecordGame(result)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	stats := collector.GetStats()
	if stats.TotalGames != 1 {
		t.Errorf("Expected 1 total game, got %d", stats.TotalGames)
	}

	if stats.TotalWins != 1 {
		t.Errorf("Expected 1 win, got %d", stats.TotalWins)
	}

	if stats.SwitchStats.GamesPlayed != 1 {
		t.Errorf("Expected 1 switch game, got %d", stats.SwitchStats.GamesPlayed)
	}

	if stats.SwitchStats.Wins != 1 {
		t.Errorf("Expected 1 switch win, got %d", stats.SwitchStats.Wins)
	}
}

func TestRecordGameNilResult(t *testing.T) {
	collector := NewCollector()

	err := collector.RecordGame(nil)
	if err == nil {
		t.Error("Expected error for nil result")
	}
}

func TestUpdateAggregateStats(t *testing.T) {
	collector := NewCollector()

	collector.RecordGame(createTestGameResult(game.Switch, true))
	collector.RecordGame(createTestGameResult(game.Switch, false))
	collector.RecordGame(createTestGameResult(game.Stay, true))
	collector.RecordGame(createTestGameResult(game.Stay, false))

	stats := collector.GetStats()

	if stats.TotalGames != 4 {
		t.Errorf("Expected 4 total games, got %d", stats.TotalGames)
	}

	if stats.TotalWins != 2 {
		t.Errorf("Expected 2 wins, got %d", stats.TotalWins)
	}

	if stats.TotalLosses != 2 {
		t.Errorf("Expected 2 losses, got %d", stats.TotalLosses)
	}

	if stats.SwitchStats.GamesPlayed != 2 {
		t.Errorf("Expected 2 switch games, got %d", stats.SwitchStats.GamesPlayed)
	}

	if stats.SwitchStats.WinRate != 0.5 {
		t.Errorf("Expected 0.5 switch win rate, got %f", stats.SwitchStats.WinRate)
	}

	if stats.StayStats.GamesPlayed != 2 {
		t.Errorf("Expected 2 stay games, got %d", stats.StayStats.GamesPlayed)
	}

	if stats.StayStats.WinRate != 0.5 {
		t.Errorf("Expected 0.5 stay win rate, got %f", stats.StayStats.WinRate)
	}
}

func TestUpdateDailyStats(t *testing.T) {
	collector := NewCollector()

	now := time.Now()
	result := createTestGameResult(game.Switch, true)
	result.Timestamp = now

	collector.RecordGame(result)

	dateKey := now.Format("2006-01-02")
	daily, exists := collector.stats.DailyStats[dateKey]

	if !exists {
		t.Error("Daily stats should exist for today")
	}

	if daily.GamesPlayed != 1 {
		t.Errorf("Expected 1 game played today, got %d", daily.GamesPlayed)
	}

	if daily.Wins != 1 {
		t.Errorf("Expected 1 win today, got %d", daily.Wins)
	}

	if daily.SwitchGames != 1 {
		t.Errorf("Expected 1 switch game today, got %d", daily.SwitchGames)
	}
}

func TestUpdateStreakStats(t *testing.T) {
	collector := NewCollector()

	collector.RecordGame(createTestGameResult(game.Switch, true))
	collector.RecordGame(createTestGameResult(game.Switch, true))
	collector.RecordGame(createTestGameResult(game.Stay, false))

	stats := collector.GetStats()

	if stats.StreakStats.CurrentWinStreak != 0 {
		t.Errorf("Expected 0 current win streak, got %d", stats.StreakStats.CurrentWinStreak)
	}

	if stats.StreakStats.CurrentLossStreak != 1 {
		t.Errorf("Expected 1 current loss streak, got %d", stats.StreakStats.CurrentLossStreak)
	}

	if stats.StreakStats.LongestWinStreak != 2 {
		t.Errorf("Expected 2 longest win streak, got %d", stats.StreakStats.LongestWinStreak)
	}

	if stats.StreakStats.CurrentStayStreak != 1 {
		t.Errorf("Expected 1 current stay streak, got %d", stats.StreakStats.CurrentStayStreak)
	}

	if stats.StreakStats.CurrentSwitchStreak != 0 {
		t.Errorf("Expected 0 current switch streak, got %d", stats.StreakStats.CurrentSwitchStreak)
	}
}

func TestUpdateTimeStats(t *testing.T) {
	collector := NewCollector()

	result1 := createTestGameResult(game.Switch, true)
	result1.GameDuration = time.Millisecond * 100

	result2 := createTestGameResult(game.Stay, false)
	result2.GameDuration = time.Millisecond * 200

	collector.RecordGame(result1)
	collector.RecordGame(result2)

	stats := collector.GetStats()

	expectedTotal := time.Millisecond * 300
	expectedAverage := time.Millisecond * 150

	if stats.TotalGameTime != expectedTotal {
		t.Errorf("Expected total time %v, got %v", expectedTotal, stats.TotalGameTime)
	}

	if stats.AverageGameTime != expectedAverage {
		t.Errorf("Expected average time %v, got %v", expectedAverage, stats.AverageGameTime)
	}

	if stats.FirstGameTime == nil {
		t.Error("First game time should not be nil")
	}

	if stats.LastGameTime == nil {
		t.Error("Last game time should not be nil")
	}
}

func TestGetSummary(t *testing.T) {
	collector := NewCollector()

	for i := 0; i < 10; i++ {
		strategy := game.Switch
		if i%2 == 0 {
			strategy = game.Stay
		}
		won := i < 7
		collector.RecordGame(createTestGameResult(strategy, won))
	}

	summary := collector.GetSummary()

	if summary.TotalGames != 10 {
		t.Errorf("Expected 10 total games, got %d", summary.TotalGames)
	}

	if summary.OverallWinRate != 0.7 {
		t.Errorf("Expected 0.7 overall win rate, got %f", summary.OverallWinRate)
	}

	if summary.FavoriteStrategy != "Balanced" {
		t.Errorf("Expected Balanced strategy, got %s", summary.FavoriteStrategy)
	}
}

func TestGetFilteredGames(t *testing.T) {
	collector := NewCollector()

	switchStrategy := game.Switch
	stayStrategy := game.Stay

	collector.RecordGame(createTestGameResult(game.Switch, true))
	collector.RecordGame(createTestGameResult(game.Stay, false))
	collector.RecordGame(createTestGameResult(game.Switch, false))

	filter := StatsFilter{
		Strategy: &switchStrategy,
	}

	filtered := collector.GetFilteredGames(filter)
	if len(filtered) != 2 {
		t.Errorf("Expected 2 switch games, got %d", len(filtered))
	}

	filter = StatsFilter{
		Strategy: &stayStrategy,
	}

	filtered = collector.GetFilteredGames(filter)
	if len(filtered) != 1 {
		t.Errorf("Expected 1 stay game, got %d", len(filtered))
	}

	filter = StatsFilter{
		WonOnly: true,
	}

	filtered = collector.GetFilteredGames(filter)
	if len(filtered) != 1 {
		t.Errorf("Expected 1 won game, got %d", len(filtered))
	}

	filter = StatsFilter{
		Limit: 2,
	}

	filtered = collector.GetFilteredGames(filter)
	if len(filtered) != 2 {
		t.Errorf("Expected 2 games with limit, got %d", len(filtered))
	}
}

func TestReset(t *testing.T) {
	collector := NewCollector()

	collector.RecordGame(createTestGameResult(game.Switch, true))

	if collector.stats.TotalGames != 1 {
		t.Error("Should have 1 game before reset")
	}

	collector.Reset()

	if collector.stats.TotalGames != 0 {
		t.Error("Should have 0 games after reset")
	}

	if collector.stats.DailyStats == nil {
		t.Error("DailyStats should not be nil after reset")
	}
}

func TestFormatDuration(t *testing.T) {
	collector := NewCollector()

	tests := []struct {
		duration time.Duration
		expected string
	}{
		{time.Millisecond * 500, "500ms"},
		{time.Second * 2, "2.0s"},
		{time.Minute * 3, "3.0m"},
		{time.Hour * 2, "2.0h"},
	}

	for _, test := range tests {
		result := collector.formatDuration(test.duration)
		if result != test.expected {
			t.Errorf("Expected %s, got %s for duration %v", test.expected, result, test.duration)
		}
	}
}

func TestGetRecentForm(t *testing.T) {
	collector := NewCollector()

	form := collector.getRecentForm()
	if form != "Insufficient data" {
		t.Errorf("Expected 'Insufficient data', got %s", form)
	}

	for i := 0; i < 5; i++ {
		collector.RecordGame(createTestGameResult(game.Switch, true))
	}

	form = collector.getRecentForm()
	if form != "Excellent (5/5)" {
		t.Errorf("Expected 'Excellent (5/5)', got %s", form)
	}

	collector.RecordGame(createTestGameResult(game.Switch, false))

	form = collector.getRecentForm()
	if form != "Very Good (4/5)" {
		t.Errorf("Expected 'Very Good (4/5)', got %s", form)
	}
}
