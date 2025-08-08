package ui

import (
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/westhuis/monty-hall/pkg/game"
	"github.com/westhuis/monty-hall/pkg/stats"
)

// ViewState represents the current view in the application
type ViewState int

const (
	MainMenuView ViewState = iota
	GameView
	StatsView
	HelpView
	ExitView
)

// Model represents the main application state
type Model struct {
	// Current view state
	CurrentView ViewState

	// Terminal dimensions
	Width  int
	Height int

	// Game state
	Game         *game.Game
	StatsManager *stats.StatsManager

	// UI state
	MenuCursor     int
	DoorCursor     int
	ShowHelp       bool
	ErrorMessage   string
	SuccessMessage string

	// Game flow state
	GamePhase  game.GamePhase
	ShowResult bool

	// Statistics view state
	StatsPage     int
	MaxStatsPages int

	// Animation system
	AnimationManager *AnimationManager
	DoorAnimations   map[int]*DoorOpenAnimation
	ShowAnimations   bool

	// Dramatic reveal system
	IsRevealing     bool
	RevealStartTime time.Time

	// Reset confirmation system
	ShowResetConfirmation    bool
	ResetConfirmationNumbers [4]int
	UserInputNumbers         [4]int
	CurrentInputIndex        int
}

// Msg represents messages that can be sent to update the model
type Msg interface{}

// KeyMsg represents keyboard input
type KeyMsg struct {
	Key string
}

// WindowSizeMsg represents terminal resize events
type WindowSizeMsg struct {
	Width  int
	Height int
}

// GameUpdateMsg represents game state changes
type GameUpdateMsg struct {
	Game *game.Game
}

// StatsUpdateMsg represents statistics updates
type StatsUpdateMsg struct {
	Summary stats.StatsSummary
}

// ErrorMsg represents error messages
type ErrorMsg struct {
	Error string
}

// SuccessMsg represents success messages
type SuccessMsg struct {
	Message string
}

// MenuOption represents a menu item
type MenuOption struct {
	Label       string
	Description string
	Action      func() tea.Cmd
}

// DoorDisplay represents how a door should be displayed
type DoorDisplay struct {
	ID       int
	State    game.DoorState
	Content  game.DoorContent
	Selected bool
	Cursor   bool
	Style    string
}

// StatsDisplay represents formatted statistics for display
type StatsDisplay struct {
	Title       string
	TotalGames  string
	WinRate     string
	SwitchStats string
	StayStats   string
	Streaks     string
	RecentForm  string
	PlayTime    string
}

// Colors and styling constants
const (
	ColorPrimary   = "#FF6B6B"
	ColorSecondary = "#4ECDC4"
	ColorSuccess   = "#45B7D1"
	ColorWarning   = "#FFA07A"
	ColorError     = "#FF6B6B"
	ColorText      = "#FFFFFF"
	ColorMuted     = "#888888"
	ColorBorder    = "#444444"
)

// Key bindings
const (
	KeyUp     = "up"
	KeyDown   = "down"
	KeyLeft   = "left"
	KeyRight  = "right"
	KeyEnter  = "enter"
	KeySpace  = "space"
	KeyEscape = "esc"
	KeyTab    = "tab"
	KeyQ      = "q"
	KeyH      = "h"
	KeyR      = "r"
	KeyS      = "s"
	Key1      = "1"
	Key2      = "2"
	Key3      = "3"
)

// RevealDelayMsg is sent after the reveal delay timer
type RevealDelayMsg struct{}
