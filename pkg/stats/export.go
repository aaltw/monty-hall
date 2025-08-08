package stats

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/westhuis/monty-hall/pkg/game"
)

// ExportFormat represents the available export formats
type ExportFormat int

const (
	ExportJSON ExportFormat = iota
	ExportCSV
	ExportText
)

// String returns the string representation of the export format
func (ef ExportFormat) String() string {
	switch ef {
	case ExportJSON:
		return "JSON"
	case ExportCSV:
		return "CSV"
	case ExportText:
		return "Text"
	default:
		return "Unknown"
	}
}

// GetFileExtension returns the file extension for the export format
func (ef ExportFormat) GetFileExtension() string {
	switch ef {
	case ExportJSON:
		return ".json"
	case ExportCSV:
		return ".csv"
	case ExportText:
		return ".txt"
	default:
		return ".txt"
	}
}

// ExportOptions contains options for exporting statistics
type ExportOptions struct {
	Format            ExportFormat
	Filename          string
	IncludeHistory    bool
	IncludeDailyStats bool
	TimeRange         *TimeRange
}

// DefaultExportOptions returns default export options
func DefaultExportOptions() ExportOptions {
	return ExportOptions{
		Format:            ExportJSON,
		Filename:          "",
		IncludeHistory:    true,
		IncludeDailyStats: true,
		TimeRange:         nil,
	}
}

// ExportStats exports statistics to a file in the specified format
func (sm *StatsManager) ExportStats(options ExportOptions) error {
	stats := sm.GetStats()

	// Generate filename if not provided
	if options.Filename == "" {
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		options.Filename = fmt.Sprintf("monty-hall-stats_%s%s", timestamp, options.Format.GetFileExtension())
	}

	// Ensure filename has correct extension
	if !strings.HasSuffix(options.Filename, options.Format.GetFileExtension()) {
		options.Filename += options.Format.GetFileExtension()
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(options.Filename)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Export based on format
	switch options.Format {
	case ExportJSON:
		return sm.exportJSON(stats, options)
	case ExportCSV:
		return sm.exportCSV(stats, options)
	case ExportText:
		return sm.exportText(stats, options)
	default:
		return fmt.Errorf("unsupported export format: %v", options.Format)
	}
}

// exportJSON exports statistics as JSON
func (sm *StatsManager) exportJSON(stats *GameStats, options ExportOptions) error {
	// Create export data structure
	exportData := map[string]interface{}{
		"export_info": map[string]interface{}{
			"timestamp":   time.Now().Format(time.RFC3339),
			"format":      "JSON",
			"version":     "1.0",
			"total_games": stats.TotalGames,
		},
		"summary": sm.GetSummary(),
		"aggregate_stats": map[string]interface{}{
			"total_games":       stats.TotalGames,
			"total_wins":        stats.TotalWins,
			"total_losses":      stats.TotalLosses,
			"switch_stats":      stats.SwitchStats,
			"stay_stats":        stats.StayStats,
			"average_game_time": stats.AverageGameTime.String(),
			"total_game_time":   stats.TotalGameTime.String(),
			"first_game_time":   stats.FirstGameTime,
			"last_game_time":    stats.LastGameTime,
			"streak_stats":      stats.StreakStats,
		},
	}

	// Include game history if requested
	if options.IncludeHistory {
		history := stats.GameHistory
		if options.TimeRange != nil {
			history = sm.filterGamesByTimeRange(stats.GameHistory, *options.TimeRange)
		}
		exportData["game_history"] = history
	}

	// Include daily stats if requested
	if options.IncludeDailyStats {
		exportData["daily_stats"] = stats.DailyStats
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(options.Filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}

// exportCSV exports game history as CSV
func (sm *StatsManager) exportCSV(stats *GameStats, options ExportOptions) error {
	file, err := os.Create(options.Filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Game ID",
		"Timestamp",
		"Strategy",
		"Won",
		"Initial Choice",
		"Final Choice",
		"Car Position",
		"Host Opened Door",
		"Game Duration (ms)",
		"Day of Week",
		"Hour of Day",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Filter games by time range if specified
	games := stats.GameHistory
	if options.TimeRange != nil {
		games = sm.filterGamesByTimeRange(stats.GameHistory, *options.TimeRange)
	}

	// Write game records
	for _, gameRecord := range games {
		strategyStr := "STAY"
		if gameRecord.Strategy == game.Switch {
			strategyStr = "SWITCH"
		}

		record := []string{
			gameRecord.ID,
			gameRecord.Timestamp.Format(time.RFC3339),
			strategyStr,
			fmt.Sprintf("%t", gameRecord.Won),
			fmt.Sprintf("%d", gameRecord.InitialChoice+1),  // Convert to 1-based
			fmt.Sprintf("%d", gameRecord.FinalChoice+1),    // Convert to 1-based
			fmt.Sprintf("%d", gameRecord.CarPosition+1),    // Convert to 1-based
			fmt.Sprintf("%d", gameRecord.HostOpenedDoor+1), // Convert to 1-based
			fmt.Sprintf("%d", gameRecord.GameDuration.Milliseconds()),
			gameRecord.DayOfWeek,
			fmt.Sprintf("%d", gameRecord.HourOfDay),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}

// exportText exports statistics as human-readable text
func (sm *StatsManager) exportText(stats *GameStats, options ExportOptions) error {
	var content strings.Builder

	// Header
	content.WriteString("MONTY HALL GAME STATISTICS REPORT\n")
	content.WriteString("==================================\n\n")
	content.WriteString(fmt.Sprintf("Generated: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	content.WriteString(fmt.Sprintf("Total Games: %d\n\n", stats.TotalGames))

	// Overall Statistics
	content.WriteString("OVERALL STATISTICS\n")
	content.WriteString("------------------\n")
	if stats.TotalGames > 0 {
		overallWinRate := float64(stats.TotalWins) / float64(stats.TotalGames) * 100
		content.WriteString(fmt.Sprintf("Total Games Played: %d\n", stats.TotalGames))
		content.WriteString(fmt.Sprintf("Total Wins: %d\n", stats.TotalWins))
		content.WriteString(fmt.Sprintf("Total Losses: %d\n", stats.TotalLosses))
		content.WriteString(fmt.Sprintf("Overall Win Rate: %.1f%%\n", overallWinRate))
		content.WriteString(fmt.Sprintf("Average Game Time: %s\n", stats.AverageGameTime))
		content.WriteString(fmt.Sprintf("Total Play Time: %s\n", stats.TotalGameTime))
		if stats.FirstGameTime != nil {
			content.WriteString(fmt.Sprintf("First Game: %s\n", stats.FirstGameTime.Format("2006-01-02 15:04:05")))
		}
		if stats.LastGameTime != nil {
			content.WriteString(fmt.Sprintf("Last Game: %s\n", stats.LastGameTime.Format("2006-01-02 15:04:05")))
		}
	} else {
		content.WriteString("No games played yet.\n")
	}
	content.WriteString("\n")

	// Strategy Statistics
	content.WriteString("STRATEGY COMPARISON\n")
	content.WriteString("-------------------\n")
	content.WriteString(fmt.Sprintf("STAY Strategy:\n"))
	content.WriteString(fmt.Sprintf("  Games Played: %d\n", stats.StayStats.GamesPlayed))
	content.WriteString(fmt.Sprintf("  Wins: %d\n", stats.StayStats.Wins))
	content.WriteString(fmt.Sprintf("  Losses: %d\n", stats.StayStats.Losses))
	if stats.StayStats.GamesPlayed > 0 {
		content.WriteString(fmt.Sprintf("  Win Rate: %.1f%%\n", stats.StayStats.WinRate*100))
	}
	content.WriteString("\n")

	content.WriteString(fmt.Sprintf("SWITCH Strategy:\n"))
	content.WriteString(fmt.Sprintf("  Games Played: %d\n", stats.SwitchStats.GamesPlayed))
	content.WriteString(fmt.Sprintf("  Wins: %d\n", stats.SwitchStats.Wins))
	content.WriteString(fmt.Sprintf("  Losses: %d\n", stats.SwitchStats.Losses))
	if stats.SwitchStats.GamesPlayed > 0 {
		content.WriteString(fmt.Sprintf("  Win Rate: %.1f%%\n", stats.SwitchStats.WinRate*100))
	}
	content.WriteString("\n")

	// Theoretical vs Actual
	content.WriteString("THEORETICAL vs ACTUAL\n")
	content.WriteString("---------------------\n")
	content.WriteString("Theoretical Probabilities:\n")
	content.WriteString("  STAY Strategy: 33.3% (1/3)\n")
	content.WriteString("  SWITCH Strategy: 66.7% (2/3)\n\n")

	if stats.StayStats.GamesPlayed > 0 || stats.SwitchStats.GamesPlayed > 0 {
		content.WriteString("Actual Results:\n")
		if stats.StayStats.GamesPlayed > 0 {
			content.WriteString(fmt.Sprintf("  STAY Strategy: %.1f%% (%d/%d games)\n",
				stats.StayStats.WinRate*100, stats.StayStats.Wins, stats.StayStats.GamesPlayed))
		}
		if stats.SwitchStats.GamesPlayed > 0 {
			content.WriteString(fmt.Sprintf("  SWITCH Strategy: %.1f%% (%d/%d games)\n",
				stats.SwitchStats.WinRate*100, stats.SwitchStats.Wins, stats.SwitchStats.GamesPlayed))
		}
	}
	content.WriteString("\n")

	// Streak Statistics
	content.WriteString("STREAK STATISTICS\n")
	content.WriteString("-----------------\n")
	content.WriteString(fmt.Sprintf("Current Win Streak: %d\n", stats.StreakStats.CurrentWinStreak))
	content.WriteString(fmt.Sprintf("Current Loss Streak: %d\n", stats.StreakStats.CurrentLossStreak))
	content.WriteString(fmt.Sprintf("Longest Win Streak: %d\n", stats.StreakStats.LongestWinStreak))
	content.WriteString(fmt.Sprintf("Longest Loss Streak: %d\n", stats.StreakStats.LongestLossStreak))
	content.WriteString(fmt.Sprintf("Current Switch Streak: %d\n", stats.StreakStats.CurrentSwitchStreak))
	content.WriteString(fmt.Sprintf("Current Stay Streak: %d\n", stats.StreakStats.CurrentStayStreak))
	content.WriteString("\n")

	// Recent Games (if history is included)
	if options.IncludeHistory && len(stats.GameHistory) > 0 {
		content.WriteString("RECENT GAMES (Last 10)\n")
		content.WriteString("-----------------------\n")

		games := stats.GameHistory
		if options.TimeRange != nil {
			games = sm.filterGamesByTimeRange(stats.GameHistory, *options.TimeRange)
		}

		// Show last 10 games
		start := len(games) - 10
		if start < 0 {
			start = 0
		}

		for i := start; i < len(games); i++ {
			gameRecord := games[i]
			result := "LOSS"
			if gameRecord.Won {
				result = "WIN"
			}
			strategyStr := "STAY"
			if gameRecord.Strategy == game.Switch {
				strategyStr = "SWITCH"
			}
			content.WriteString(fmt.Sprintf("%s | %s | %s | Door %d→%d | %s\n",
				gameRecord.Timestamp.Format("2006-01-02 15:04"),
				strategyStr,
				result,
				gameRecord.InitialChoice+1,
				gameRecord.FinalChoice+1,
				gameRecord.GameDuration.Round(time.Millisecond)))
		}
		content.WriteString("\n")
	}

	// Footer
	content.WriteString("Generated by Monty Hall Terminal Application\n")
	content.WriteString("For more information, visit: https://github.com/westhuis/monty-hall\n")

	// Write to file
	if err := os.WriteFile(options.Filename, []byte(content.String()), 0644); err != nil {
		return fmt.Errorf("failed to write text file: %w", err)
	}

	return nil
}

// filterGamesByTimeRange filters games by the specified time range
func (sm *StatsManager) filterGamesByTimeRange(games []GameRecord, timeRange TimeRange) []GameRecord {
	var filtered []GameRecord
	for _, game := range games {
		if game.Timestamp.After(timeRange.Start) && game.Timestamp.Before(timeRange.End) {
			filtered = append(filtered, game)
		}
	}
	return filtered
}

// GetExportFormats returns all available export formats
func GetExportFormats() []ExportFormat {
	return []ExportFormat{ExportJSON, ExportCSV, ExportText}
}

// ParseExportFormat parses a string into an ExportFormat
func ParseExportFormat(format string) (ExportFormat, error) {
	switch strings.ToLower(format) {
	case "json":
		return ExportJSON, nil
	case "csv":
		return ExportCSV, nil
	case "text", "txt":
		return ExportText, nil
	default:
		return ExportJSON, fmt.Errorf("unknown export format: %s", format)
	}
}
