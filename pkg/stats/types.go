package stats

import (
	"time"

	"github.com/westhuis/monty-hall/pkg/game"
)

type GameStats struct {
	TotalGames      int                   `json:"total_games"`
	TotalWins       int                   `json:"total_wins"`
	TotalLosses     int                   `json:"total_losses"`
	SwitchStats     StrategyStats         `json:"switch_stats"`
	StayStats       StrategyStats         `json:"stay_stats"`
	AverageGameTime time.Duration         `json:"average_game_time"`
	TotalGameTime   time.Duration         `json:"total_game_time"`
	FirstGameTime   *time.Time            `json:"first_game_time,omitempty"`
	LastGameTime    *time.Time            `json:"last_game_time,omitempty"`
	GameHistory     []GameRecord          `json:"game_history"`
	DailyStats      map[string]DailyStats `json:"daily_stats"`
	StreakStats     StreakStats           `json:"streak_stats"`
}

type StrategyStats struct {
	GamesPlayed int     `json:"games_played"`
	Wins        int     `json:"wins"`
	Losses      int     `json:"losses"`
	WinRate     float64 `json:"win_rate"`
}

type GameRecord struct {
	ID             string              `json:"id"`
	Timestamp      time.Time           `json:"timestamp"`
	Strategy       game.PlayerStrategy `json:"strategy"`
	Won            bool                `json:"won"`
	InitialChoice  int                 `json:"initial_choice"`
	FinalChoice    int                 `json:"final_choice"`
	CarPosition    int                 `json:"car_position"`
	HostOpenedDoor int                 `json:"host_opened_door"`
	GameDuration   time.Duration       `json:"game_duration"`
	DayOfWeek      string              `json:"day_of_week"`
	HourOfDay      int                 `json:"hour_of_day"`
}

type DailyStats struct {
	Date        string        `json:"date"`
	GamesPlayed int           `json:"games_played"`
	Wins        int           `json:"wins"`
	Losses      int           `json:"losses"`
	WinRate     float64       `json:"win_rate"`
	SwitchGames int           `json:"switch_games"`
	StayGames   int           `json:"stay_games"`
	TotalTime   time.Duration `json:"total_time"`
	AverageTime time.Duration `json:"average_time"`
}

type StreakStats struct {
	CurrentWinStreak    int `json:"current_win_streak"`
	CurrentLossStreak   int `json:"current_loss_streak"`
	LongestWinStreak    int `json:"longest_win_streak"`
	LongestLossStreak   int `json:"longest_loss_streak"`
	CurrentSwitchStreak int `json:"current_switch_streak"`
	CurrentStayStreak   int `json:"current_stay_streak"`
}

type StatsSummary struct {
	TotalGames       int     `json:"total_games"`
	OverallWinRate   float64 `json:"overall_win_rate"`
	SwitchWinRate    float64 `json:"switch_win_rate"`
	StayWinRate      float64 `json:"stay_win_rate"`
	SwitchAdvantage  float64 `json:"switch_advantage"`
	AverageGameTime  string  `json:"average_game_time"`
	TotalPlayTime    string  `json:"total_play_time"`
	FavoriteStrategy string  `json:"favorite_strategy"`
	BestStreak       int     `json:"best_streak"`
	RecentForm       string  `json:"recent_form"`
}

type TimeRange struct {
	Start time.Time
	End   time.Time
}

type StatsFilter struct {
	Strategy  *game.PlayerStrategy
	TimeRange *TimeRange
	WonOnly   bool
	LostOnly  bool
	Limit     int
}
