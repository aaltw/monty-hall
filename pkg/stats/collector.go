package stats

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/westhuis/monty-hall/pkg/game"
)

const (
	MaxHistorySize = 10000 // Maximum number of games to keep in memory
	TrimSize       = 1000  // Number of games to remove when trimming
)

type Collector struct {
	stats *GameStats
}

func NewCollector() *Collector {
	return &Collector{
		stats: &GameStats{
			DailyStats: make(map[string]DailyStats),
		},
	}
}

func (c *Collector) RecordGame(result *game.GameResult) error {
	if result == nil {
		return fmt.Errorf("game result cannot be nil")
	}

	record := c.createGameRecord(result)

	c.stats.GameHistory = append(c.stats.GameHistory, record)

	// Manage memory by trimming old games if history gets too large
	if len(c.stats.GameHistory) > MaxHistorySize {
		// Keep the most recent games, remove the oldest ones
		c.stats.GameHistory = c.stats.GameHistory[TrimSize:]
	}

	c.updateAggregateStats(record)
	c.updateDailyStats(record)
	c.updateStreakStats(record)
	c.updateTimeStats(record)

	return nil
}

func (c *Collector) createGameRecord(result *game.GameResult) GameRecord {
	id := c.generateGameID()

	return GameRecord{
		ID:             id,
		Timestamp:      result.Timestamp,
		Strategy:       result.Strategy,
		Won:            result.Won,
		InitialChoice:  result.InitialChoice,
		FinalChoice:    result.FinalChoice,
		CarPosition:    result.CarPosition,
		HostOpenedDoor: result.HostOpenedDoor,
		GameDuration:   result.GameDuration,
		DayOfWeek:      result.Timestamp.Weekday().String(),
		HourOfDay:      result.Timestamp.Hour(),
	}
}

func (c *Collector) generateGameID() string {
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		// Fallback to timestamp-based ID if crypto/rand fails
		return fmt.Sprintf("%x", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x", bytes)
}

func (c *Collector) updateAggregateStats(record GameRecord) {
	c.stats.TotalGames++

	if record.Won {
		c.stats.TotalWins++
	} else {
		c.stats.TotalLosses++
	}

	if record.Strategy == game.Switch {
		c.stats.SwitchStats.GamesPlayed++
		if record.Won {
			c.stats.SwitchStats.Wins++
		} else {
			c.stats.SwitchStats.Losses++
		}
		c.stats.SwitchStats.WinRate = float64(c.stats.SwitchStats.Wins) / float64(c.stats.SwitchStats.GamesPlayed)
	} else {
		c.stats.StayStats.GamesPlayed++
		if record.Won {
			c.stats.StayStats.Wins++
		} else {
			c.stats.StayStats.Losses++
		}
		c.stats.StayStats.WinRate = float64(c.stats.StayStats.Wins) / float64(c.stats.StayStats.GamesPlayed)
	}
}

func (c *Collector) updateDailyStats(record GameRecord) {
	dateKey := record.Timestamp.Format("2006-01-02")

	daily, exists := c.stats.DailyStats[dateKey]
	if !exists {
		daily = DailyStats{
			Date: dateKey,
		}
	}

	daily.GamesPlayed++
	daily.TotalTime += record.GameDuration
	daily.AverageTime = daily.TotalTime / time.Duration(daily.GamesPlayed)

	if record.Won {
		daily.Wins++
	} else {
		daily.Losses++
	}

	daily.WinRate = float64(daily.Wins) / float64(daily.GamesPlayed)

	if record.Strategy == game.Switch {
		daily.SwitchGames++
	} else {
		daily.StayGames++
	}

	c.stats.DailyStats[dateKey] = daily
}

func (c *Collector) updateStreakStats(record GameRecord) {
	if record.Won {
		c.stats.StreakStats.CurrentWinStreak++
		c.stats.StreakStats.CurrentLossStreak = 0

		if c.stats.StreakStats.CurrentWinStreak > c.stats.StreakStats.LongestWinStreak {
			c.stats.StreakStats.LongestWinStreak = c.stats.StreakStats.CurrentWinStreak
		}
	} else {
		c.stats.StreakStats.CurrentLossStreak++
		c.stats.StreakStats.CurrentWinStreak = 0

		if c.stats.StreakStats.CurrentLossStreak > c.stats.StreakStats.LongestLossStreak {
			c.stats.StreakStats.LongestLossStreak = c.stats.StreakStats.CurrentLossStreak
		}
	}

	if record.Strategy == game.Switch {
		c.stats.StreakStats.CurrentSwitchStreak++
		c.stats.StreakStats.CurrentStayStreak = 0
	} else {
		c.stats.StreakStats.CurrentStayStreak++
		c.stats.StreakStats.CurrentSwitchStreak = 0
	}
}

func (c *Collector) updateTimeStats(record GameRecord) {
	c.stats.TotalGameTime += record.GameDuration
	c.stats.AverageGameTime = c.stats.TotalGameTime / time.Duration(c.stats.TotalGames)

	if c.stats.FirstGameTime == nil {
		c.stats.FirstGameTime = &record.Timestamp
	}

	c.stats.LastGameTime = &record.Timestamp
}

func (c *Collector) GetStats() *GameStats {
	return c.stats
}

func (c *Collector) GetSummary() StatsSummary {
	stats := c.stats

	summary := StatsSummary{
		TotalGames:      stats.TotalGames,
		AverageGameTime: c.formatDuration(stats.AverageGameTime),
		TotalPlayTime:   c.formatDuration(stats.TotalGameTime),
		BestStreak:      stats.StreakStats.LongestWinStreak,
	}

	if stats.TotalGames > 0 {
		summary.OverallWinRate = float64(stats.TotalWins) / float64(stats.TotalGames)
	}

	if stats.SwitchStats.GamesPlayed > 0 {
		summary.SwitchWinRate = stats.SwitchStats.WinRate
	}

	if stats.StayStats.GamesPlayed > 0 {
		summary.StayWinRate = stats.StayStats.WinRate
	}

	summary.SwitchAdvantage = summary.SwitchWinRate - summary.StayWinRate

	if stats.SwitchStats.GamesPlayed > stats.StayStats.GamesPlayed {
		summary.FavoriteStrategy = "Switch"
	} else if stats.StayStats.GamesPlayed > stats.SwitchStats.GamesPlayed {
		summary.FavoriteStrategy = "Stay"
	} else {
		summary.FavoriteStrategy = "Balanced"
	}

	summary.RecentForm = c.getRecentForm()

	return summary
}

func (c *Collector) formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%.0fms", float64(d.Nanoseconds())/1e6)
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	return fmt.Sprintf("%.1fh", d.Hours())
}

func (c *Collector) getRecentForm() string {
	const recentGames = 5
	history := c.stats.GameHistory

	if len(history) < recentGames {
		return "Insufficient data"
	}

	recent := history[len(history)-recentGames:]
	wins := 0

	for _, game := range recent {
		if game.Won {
			wins++
		}
	}

	switch wins {
	case 5:
		return "Excellent (5/5)"
	case 4:
		return "Very Good (4/5)"
	case 3:
		return "Good (3/5)"
	case 2:
		return "Fair (2/5)"
	case 1:
		return "Poor (1/5)"
	default:
		return "Very Poor (0/5)"
	}
}

func (c *Collector) GetFilteredGames(filter StatsFilter) []GameRecord {
	var filtered []GameRecord

	for _, record := range c.stats.GameHistory {
		if c.matchesFilter(record, filter) {
			filtered = append(filtered, record)
		}
	}

	if filter.Limit > 0 && len(filtered) > filter.Limit {
		start := len(filtered) - filter.Limit
		filtered = filtered[start:]
	}

	return filtered
}

func (c *Collector) matchesFilter(record GameRecord, filter StatsFilter) bool {
	if filter.Strategy != nil && record.Strategy != *filter.Strategy {
		return false
	}

	if filter.TimeRange != nil {
		if record.Timestamp.Before(filter.TimeRange.Start) || record.Timestamp.After(filter.TimeRange.End) {
			return false
		}
	}

	if filter.WonOnly && !record.Won {
		return false
	}

	if filter.LostOnly && record.Won {
		return false
	}

	return true
}

func (c *Collector) Reset() {
	c.stats = &GameStats{
		DailyStats: make(map[string]DailyStats),
	}
}
