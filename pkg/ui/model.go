package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/westhuis/monty-hall/pkg/config"
	"github.com/westhuis/monty-hall/pkg/game"
	"github.com/westhuis/monty-hall/pkg/stats"
)

// NewModel creates a new TUI model
func NewModel() *Model {
	statsManager := stats.NewStatsManager()

	return &Model{
		CurrentView:           MainMenuView,
		Width:                 80,
		Height:                24,
		ConfigManager:         nil, // Will be nil for backward compatibility
		Game:                  nil,
		StatsManager:          statsManager,
		MenuCursor:            0,
		DoorCursor:            0,
		ShowHelp:              false,
		ErrorMessage:          "",
		SuccessMessage:        "",
		GamePhase:             game.Setup,
		ShowResult:            false,
		StatsPage:             0,
		MaxStatsPages:         1,
		AnimationManager:      NewAnimationManager(),
		DoorAnimations:        make(map[int]*DoorOpenAnimation),
		ShowAnimations:        true,
		IsRevealing:           false,
		ShowResetConfirmation: false,
		CurrentInputIndex:     0,
	}
}

// NewModelWithConfig creates a new TUI model with configuration support
func NewModelWithConfig(configManager *config.Manager) *Model {
	statsManager := stats.NewStatsManager()
	cfg := configManager.Get()

	// Apply configuration settings
	width := 80
	height := 24
	if cfg.UI.TerminalWidth > 0 {
		width = cfg.UI.TerminalWidth
	}
	if cfg.UI.TerminalHeight > 0 {
		height = cfg.UI.TerminalHeight
	}

	return &Model{
		CurrentView:           MainMenuView,
		Width:                 width,
		Height:                height,
		ConfigManager:         configManager,
		Game:                  nil,
		StatsManager:          statsManager,
		MenuCursor:            0,
		DoorCursor:            0,
		ShowHelp:              false,
		ErrorMessage:          "",
		SuccessMessage:        "",
		GamePhase:             game.Setup,
		ShowResult:            false,
		StatsPage:             0,
		MaxStatsPages:         1,
		AnimationManager:      NewAnimationManager(),
		DoorAnimations:        make(map[int]*DoorOpenAnimation),
		ShowAnimations:        cfg.UI.ShowAnimations && !cfg.UI.ReducedMotion,
		IsRevealing:           false,
		ShowResetConfirmation: false,
		CurrentInputIndex:     0,
	}
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		// Adjust layout for small screens
		if m.Width < 80 {
			// Compact mode for small terminals
		}
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case GameUpdateMsg:
		m.Game = msg.Game
		if m.Game != nil {
			m.GamePhase = m.Game.Phase
		}
		return m, nil

	case StatsUpdateMsg:
		return m, nil

	case ErrorMsg:
		m.ErrorMessage = msg.Error
		return m, nil

	case SuccessMsg:
		m.SuccessMessage = msg.Message
		return m, nil

	case AnimationTickMsg:
		// Update animations
		return m, m.AnimationManager.Update()

	case RevealDelayMsg:
		// End the revealing state and show results
		m.IsRevealing = false
		m.ShowResult = true

		// Record the game result
		if m.Game.Result != nil {
			if err := m.StatsManager.RecordGame(m.Game.Result); err != nil {
				m.ErrorMessage = fmt.Sprintf("Failed to save statistics: %v", err)
			}
		}

		// Start winning animation if player won
		if m.ShowAnimations && m.Game.Result != nil && m.Game.Result.Won {
			return m, m.startWinningAnimation()
		}

		return m, nil
	}

	return m, nil
}

// handleKeyPress processes keyboard input
func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Clear messages on any key press
	m.ErrorMessage = ""
	m.SuccessMessage = ""

	// Handle reset confirmation input first (highest priority)
	if m.ShowResetConfirmation {
		return m.handleResetConfirmationKeys(msg)
	}

	// Global key bindings
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit

	case KeyQ:
		// Context-aware 'q' behavior
		if m.ShowHelp {
			m.ShowHelp = false
			return m, nil
		}
		if m.CurrentView == MainMenuView {
			// Quit application from main menu
			return m, tea.Quit
		} else {
			// Return to main menu from other screens
			m.CurrentView = MainMenuView
			m.MenuCursor = 0
			return m, nil
		}

	case KeyH:
		m.ShowHelp = !m.ShowHelp
		return m, nil

	case KeyEscape:
		if m.ShowHelp {
			m.ShowHelp = false
			return m, nil
		}
		if m.CurrentView != MainMenuView {
			m.CurrentView = MainMenuView
			m.MenuCursor = 0
			return m, nil
		}
	}

	// View-specific key bindings
	switch m.CurrentView {
	case MainMenuView:
		return m.handleMainMenuKeys(msg)
	case GameView:
		return m.handleGameKeys(msg)
	case StatsView:
		return m.handleStatsKeys(msg)
	}

	return m, nil
}

// handleMainMenuKeys processes main menu navigation
func (m *Model) handleMainMenuKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case KeyUp, "k":
		if m.MenuCursor > 0 {
			m.MenuCursor--
		}

	case KeyDown, "j":
		maxOptions := 4 // Play, Stats, Help, Exit
		if m.MenuCursor < maxOptions-1 {
			m.MenuCursor++
		}

	case KeyEnter, KeySpace:
		return m.executeMenuAction()
	}

	return m, nil
}

// executeMenuAction performs the selected menu action
func (m *Model) executeMenuAction() (tea.Model, tea.Cmd) {
	switch m.MenuCursor {
	case 0: // Play Game
		m.Game = game.NewGame()
		m.CurrentView = GameView
		m.DoorCursor = 0
		m.ShowResult = false
		return m, nil

	case 1: // View Statistics
		m.CurrentView = StatsView
		m.StatsPage = 0
		return m, nil

	case 2: // Help
		m.ShowHelp = true
		return m, nil

	case 3: // Exit
		return m, tea.Quit
	}

	return m, nil
}

// handleGameKeys processes game view input with door selection restrictions
func (m *Model) handleGameKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.Game == nil {
		return m, nil
	}

	switch msg.String() {
	case KeyLeft, "h":
		m.moveCursorLeft()

	case KeyRight, "l":
		m.moveCursorRight()

	case Key1:
		if m.isDoorSelectable(0) {
			m.DoorCursor = 0
		}

	case Key2:
		if m.isDoorSelectable(1) {
			m.DoorCursor = 1
		}

	case Key3:
		if m.isDoorSelectable(2) {
			m.DoorCursor = 2
		}

	case KeyEnter, KeySpace:
		if m.Game.IsGameOver() {
			// Play again
			m.Game = game.NewGame()
			m.DoorCursor = 0
			m.ShowResult = false
			return m, nil
		}
		return m.selectDoor()

	case KeyS:
		if m.Game.Phase == game.FinalChoice {
			return m.switchChoice()
		} else {
			// View statistics (available in all phases except FinalChoice)
			m.CurrentView = StatsView
			return m, nil
		}

	case KeyR:
		if m.Game.IsGameOver() {
			m.Game = game.NewGame()
			m.DoorCursor = 0
			m.ShowResult = false
			return m, nil
		}
	}

	return m, nil
}

// selectDoor handles door selection logic
func (m *Model) selectDoor() (tea.Model, tea.Cmd) {
	if m.Game == nil {
		return m, nil
	}

	// Validate that the door is selectable
	if !m.isDoorSelectable(m.DoorCursor) {
		m.ErrorMessage = "Cannot select this door"
		return m, nil
	}

	switch m.Game.Phase {
	case game.InitialChoice:
		err := m.Game.MakeInitialChoice(m.DoorCursor)
		if err != nil {
			m.ErrorMessage = err.Error()
		}
		return m, nil

	case game.FinalChoice:
		err := m.Game.MakeFinalChoice(m.DoorCursor)
		if err != nil {
			m.ErrorMessage = err.Error()
		} else {
			// Start dramatic reveal delay before showing results
			return m, m.startRevealDelay()
		}
		return m, nil
	}

	return m, nil
}

// switchChoice handles the switch action
func (m *Model) switchChoice() (tea.Model, tea.Cmd) {
	if m.Game == nil || m.Game.Phase != game.FinalChoice {
		return m, nil
	}

	err := m.Game.SwitchChoice()
	if err != nil {
		m.ErrorMessage = err.Error()
	} else {
		// Start dramatic reveal delay before showing results
		return m, m.startRevealDelay()
	}

	return m, nil
}

// handleStatsKeys processes statistics view input
func (m *Model) handleStatsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {

	switch msg.String() {
	case KeyLeft, "h":
		if m.StatsPage > 0 {
			m.StatsPage--
		}

	case KeyRight, "l":
		if m.StatsPage < m.MaxStatsPages-1 {
			m.StatsPage++
		}

	case KeyEnter, KeySpace:
		// Start a new game
		m.Game = game.NewGame()
		m.CurrentView = GameView
		m.DoorCursor = 0
		m.ShowResult = false
		return m, nil

	case KeyR:
		// Reset statistics with confirmation
		return m.confirmResetStats()

	case KeyE:
		// Export statistics
		return m.exportStats()

	case KeyQ:
		// Return to main menu (same as ESC)
		m.CurrentView = MainMenuView
		m.MenuCursor = 0
		return m, nil
	}

	return m, nil
}

// confirmResetStats handles statistics reset confirmation
func (m *Model) confirmResetStats() (tea.Model, tea.Cmd) {
	// Generate 4 random numbers for confirmation
	for i := 0; i < 4; i++ {
		m.ResetConfirmationNumbers[i] = game.SecureIntn(9) + 1 // Numbers 1-9
	}

	// Reset user input
	m.UserInputNumbers = [4]int{0, 0, 0, 0}
	m.CurrentInputIndex = 0

	// Show the confirmation dialog
	m.ShowResetConfirmation = true

	return m, nil
}

// handleResetConfirmationKeys processes input during reset confirmation
func (m *Model) handleResetConfirmationKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case KeyEscape:
		// Cancel reset confirmation
		m.ShowResetConfirmation = false
		m.CurrentInputIndex = 0
		m.UserInputNumbers = [4]int{0, 0, 0, 0}
		return m, nil

	case "backspace":
		// Delete current input and move back
		if m.CurrentInputIndex > 0 {
			m.CurrentInputIndex--
			m.UserInputNumbers[m.CurrentInputIndex] = 0
		}
		return m, nil

	case "1", "2", "3", "4", "5", "6", "7", "8", "9":
		// Input a number
		if m.CurrentInputIndex < 4 {
			num := int(msg.String()[0] - '0') // Convert character to number
			m.UserInputNumbers[m.CurrentInputIndex] = num
			m.CurrentInputIndex++

			// Check if all numbers are entered
			if m.CurrentInputIndex >= 4 {
				return m.validateAndResetStats()
			}
		}
		return m, nil

	case KeyEnter:
		// Validate and reset if all numbers are entered
		if m.CurrentInputIndex >= 4 {
			return m.validateAndResetStats()
		}
		return m, nil
	}

	return m, nil
}

// validateAndResetStats validates the confirmation numbers and resets stats if correct
func (m *Model) validateAndResetStats() (tea.Model, tea.Cmd) {
	// Check if all numbers match
	for i := 0; i < 4; i++ {
		if m.UserInputNumbers[i] != m.ResetConfirmationNumbers[i] {
			// Numbers don't match - show error and reset input
			expectedNums := fmt.Sprintf("%d %d %d %d",
				m.ResetConfirmationNumbers[0], m.ResetConfirmationNumbers[1],
				m.ResetConfirmationNumbers[2], m.ResetConfirmationNumbers[3])
			enhancedErr := CreateInvalidInputError("confirmation numbers", expectedNums)
			m.ErrorMessage = FormatErrorForDisplay(enhancedErr)
			m.CurrentInputIndex = 0
			m.UserInputNumbers = [4]int{0, 0, 0, 0}
			return m, nil
		}
	}

	// Numbers match - reset statistics
	err := m.StatsManager.Reset()
	if err != nil {
		enhancedErr := WrapError(err, "reset statistics")
		m.ErrorMessage = FormatErrorForDisplay(enhancedErr)
	} else {
		m.SuccessMessage = "Statistics reset successfully!"
	}

	// Hide confirmation dialog
	m.ShowResetConfirmation = false
	m.CurrentInputIndex = 0
	m.UserInputNumbers = [4]int{0, 0, 0, 0}

	return m, nil
}

// exportStats handles statistics export
func (m *Model) exportStats() (tea.Model, tea.Cmd) {
	// Use default export options (JSON format)
	options := stats.DefaultExportOptions()

	err := m.StatsManager.ExportStats(options)
	if err != nil {
		enhancedErr := WrapError(err, "export statistics")
		m.ErrorMessage = FormatErrorForDisplay(enhancedErr)
	} else {
		m.SuccessMessage = fmt.Sprintf("Statistics exported to: %s", options.Filename)
	}

	return m, nil
}

// View renders the current view
func (m *Model) View() string {
	if m.ShowHelp {
		return m.renderHelp()
	}

	switch m.CurrentView {
	case MainMenuView:
		return m.renderMainMenu()
	case GameView:
		return m.renderGame()
	case StatsView:
		return m.renderStats()
	default:
		return "Unknown view"
	}
}

// renderHelp renders the help screen
func (m *Model) renderHelp() string {
	helpContent := []string{
		"",
		"🎯 The Monty Hall Problem:",
		"You're on a game show with 3 doors. Behind one is a car, behind the",
		"others are goats. After you pick a door, the host opens a door with a",
		"goat. You can then switch your choice or stay with your original pick.",
		"",
		"🎮 Controls:",
		"• Arrow keys / hjkl - Navigate",
		"• Enter / Space - Select",
		"• q - Quit application",
		"• h - Toggle help",
		"• r - Reset statistics",
		"• s - Switch choice (during final decision)",
		"",
		"🎲 Game Flow:",
		"1. Choose a door (1, 2, or 3)",
		"2. Host reveals a goat behind another door",
		"3. Decide to switch or stay",
		"4. See the result and updated statistics",
		"",
		"🧮 Mathematical Insight:",
		"Switching gives you a 2/3 chance of winning!",
		"Staying gives you only a 1/3 chance of winning.",
		"",
		"Play multiple games to see this probability in action!",
		"",
		"📁 Statistics File:",
		fmt.Sprintf("Stats are saved to: %s", m.StatsManager.GetStatsFilePath()),
	}

	helpBox := NewHelpBox("HELP - Monty Hall Simulator", helpContent, GetLayoutWidth(m.Width))

	footer := RenderFooter([]KeyBinding{
		{"Enter", "Play game"},
		{"r", "Reset stats"},
		{"q", "Main menu"},
	})

	return lipgloss.JoinVertical(lipgloss.Center,
		Spacer(2),
		Center(helpBox.Render(), m.Width, 1),
		footer,
	)
}

// renderMainMenu renders the main menu with clean, functional layout
func (m *Model) renderMainMenu() string {
	// Banner - use ASCII art for larger screens
	banner := CreateASCIIBanner(m.Width)

	// Subtitle
	subtitle := SubtitleStyle.Render("Test your intuition against probability theory")

	// Menu options
	options := []string{
		"Play Game",
		"View Statistics",
		"Help",
		"Quit",
	}

	// Create flat menu items
	var menuItems []string
	for i, option := range options {
		button := NewMenuButton(option, i == m.MenuCursor)
		menuItems = append(menuItems, button.Render())
	}

	// Arrange menu vertically
	menu := lipgloss.JoinVertical(lipgloss.Center, menuItems...)

	// Messages
	var messages []string
	if m.ErrorMessage != "" {
		messages = append(messages, ErrorStyle.Render("❌ "+m.ErrorMessage))
	}
	if m.SuccessMessage != "" {
		messages = append(messages, SuccessStyle.Render("✅ "+m.SuccessMessage))
	}

	// Footer
	footer := RenderFooter([]KeyBinding{
		{"Enter", "Select"},
		{"↑↓", "Navigate"},
		{"q", "Quit"},
	})

	// Combine all elements
	var content []string
	content = append(content, banner)
	content = append(content, Spacer(1))
	content = append(content, subtitle)
	content = append(content, Spacer(2))
	content = append(content, menu)

	if len(messages) > 0 {
		content = append(content, Spacer(1))
		content = append(content, lipgloss.JoinVertical(lipgloss.Center, messages...))
	}

	content = append(content, footer)

	// Join all content vertically and center horizontally in the terminal
	menuContent := lipgloss.JoinVertical(lipgloss.Center, content...)
	// Use SafeCenter for horizontal centering, then center vertically
	horizontallyCentered := SafeCenter(menuContent, m.Width)
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, horizontallyCentered)
}

// renderGame renders the game view with fixed-height content area above doors
func (m *Model) renderGame() string {
	if m.Game == nil {
		return ErrorStyle.Render("Error: No game in progress")
	}

	// Header (always present) - use ASCII art for larger screens
	header := CreateGameBanner(m.Width)

	// Phase indicator (always present)
	phaseIndicator := NewGamePhaseIndicator(m.Game.Phase)

	// Create fixed-height content area above doors (8 lines total)
	var contentLines []string

	// Handle revealing state with dramatic pause
	if m.IsRevealing {
		contentLines = append(contentLines, Center(TitleStyle.Render("The host is opening a door..."), m.Width, 1))
		contentLines = append(contentLines, Center(SubtitleStyle.Render("..."), m.Width, 1))
		contentLines = append(contentLines, "") // Empty line
		contentLines = append(contentLines, "") // Empty line
		contentLines = append(contentLines, "") // Empty line
		contentLines = append(contentLines, "") // Empty line
		contentLines = append(contentLines, "") // Empty line
		contentLines = append(contentLines, "") // Empty line
	} else {
		switch m.Game.Phase {
		case game.InitialChoice:
			contentLines = append(contentLines, Center(TitleStyle.Render("Choose a door (1, 2, or 3):"), m.Width, 1))
			contentLines = append(contentLines, Center(SubtitleStyle.Render(fmt.Sprintf("Currently highlighting: Door %d", m.DoorCursor+1)), m.Width, 1))
			contentLines = append(contentLines, "") // Empty line
			contentLines = append(contentLines, "") // Empty line
			contentLines = append(contentLines, "") // Empty line
			contentLines = append(contentLines, "") // Empty line
			contentLines = append(contentLines, "") // Empty line
			contentLines = append(contentLines, "") // Empty line

		case game.FinalChoice:
			instruction1 := fmt.Sprintf("You initially chose door %d.", m.Game.PlayerInitialChoice+1)
			instruction2 := fmt.Sprintf("The host opened door %d, revealing a goat!", m.Game.HostOpenedDoor+1)
			contentLines = append(contentLines, Center(TitleStyle.Render(instruction1), m.Width, 1))
			contentLines = append(contentLines, Center(SubtitleStyle.Render(instruction2), m.Width, 1))
			contentLines = append(contentLines, "") // Empty line
			contentLines = append(contentLines, Center(lipgloss.NewStyle().Foreground(WarningColor).Bold(true).Render("Final Decision: Do you want to switch or stay?"), m.Width, 1))

			// Add clear instructions with cursor info
			availableDoors := m.getSelectableDoors()
			var doorOptions []string
			for _, doorIdx := range availableDoors {
				if doorIdx == m.Game.PlayerInitialChoice {
					doorOptions = append(doorOptions, fmt.Sprintf("Door %d (STAY)", doorIdx+1))
				} else {
					doorOptions = append(doorOptions, fmt.Sprintf("Door %d (SWITCH)", doorIdx+1))
				}
			}
			cursorInfo := fmt.Sprintf("Use ←→ to choose between: %s", lipgloss.JoinHorizontal(lipgloss.Left, doorOptions...))
			contentLines = append(contentLines, Center(SubtitleStyle.Render(cursorInfo), m.Width, 1))
			contentLines = append(contentLines, Center(lipgloss.NewStyle().Foreground(PrimaryColor).Render("Press 's' to SWITCH to the other door"), m.Width, 1))
			contentLines = append(contentLines, Center(lipgloss.NewStyle().Foreground(SecondaryColor).Render("Press Enter to confirm your choice"), m.Width, 1))
			contentLines = append(contentLines, "") // Empty line

		case game.GameOver:
			if m.Game.Result != nil {
				summary1 := fmt.Sprintf("You initially chose door %d", m.Game.Result.InitialChoice+1)
				summary2 := fmt.Sprintf("The host opened door %d, revealing a goat", m.Game.Result.HostOpenedDoor+1)

				var strategy string
				if m.Game.Result.Strategy == game.Switch {
					strategy = "You decided to SWITCH! 🔄"
				} else {
					strategy = "You decided to STAY! 🛡️"
				}

				contentLines = append(contentLines, Center(TitleStyle.Render("GAME COMPLETE"), m.Width, 1))
				contentLines = append(contentLines, Center(SubtitleStyle.Render(summary1), m.Width, 1))
				contentLines = append(contentLines, Center(SubtitleStyle.Render(summary2), m.Width, 1))
				contentLines = append(contentLines, Center(lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true).Render(strategy), m.Width, 1))
				contentLines = append(contentLines, "") // Empty line
				contentLines = append(contentLines, "") // Empty line
				contentLines = append(contentLines, "") // Empty line
				contentLines = append(contentLines, "") // Empty line
			}
		}
	}

	// Build the complete layout
	var content []string
	content = append(content, header)
	content = append(content, phaseIndicator.Render())
	content = append(content, Spacer(1))

	// Add fixed-height content area (8 lines)
	content = append(content, contentLines...)
	content = append(content, Spacer(1))

	// Add doors (always in the same position)
	var doors string
	if m.IsRevealing {
		doors = RenderDoorsRow(m.Game.Doors, m.Game.PlayerInitialChoice, -1, -1, false)
	} else {
		switch m.Game.Phase {
		case game.InitialChoice:
			doors = RenderDoorsRow(m.Game.Doors, -1, -1, m.DoorCursor, false)
		case game.FinalChoice:
			doors = RenderDoorsRow(m.Game.Doors, m.Game.PlayerInitialChoice, m.Game.HostOpenedDoor, m.DoorCursor, false)
		case game.GameOver:
			doors = RenderDoorsRow(m.Game.Doors, m.Game.PlayerInitialChoice, m.Game.HostOpenedDoor, -1, true)
		}
	}
	content = append(content, SafeCenter(doors, m.Width))

	// Add result message for GameOver phase (only after reveal delay is complete)
	if m.Game.Phase == game.GameOver && m.Game.Result != nil && m.ShowResult && !m.IsRevealing {
		content = append(content, Spacer(1))
		if m.Game.Result.Won {
			winMessage := "🎉 CONGRATULATIONS! You won the car! 🎉"
			enhancedWinMessage := CreateWinningMessage(winMessage)
			content = append(content, Center(enhancedWinMessage, m.Width, 1))
		} else {
			loseMessage := "😔 Sorry, you got a goat. Better luck next time!"
			content = append(content, Center(MutedStyle.Render(loseMessage), m.Width, 1))
		}
	}

	// Add footer based on phase
	var footer string
	switch m.Game.Phase {
	case game.InitialChoice:
		footer = RenderFooter([]KeyBinding{
			{"Enter", "Select door"},
			{"s", "Statistics"},
			{"←→", "Navigate"},
			{"q", "Main menu"},
		})
	case game.FinalChoice:
		footer = RenderFooter([]KeyBinding{
			{"Enter", "Confirm choice"},
			{"s", "Switch doors"},
			{"←→", "Choose door"},
			{"q", "Main menu"},
		})
	case game.GameOver:
		footer = RenderFooter([]KeyBinding{
			{"Enter", "Play again"},
			{"s", "Statistics"},
			{"q", "Main menu"},
		})
	}
	if footer != "" {
		content = append(content, footer)
	}

	// Error message
	if m.ErrorMessage != "" {
		content = append(content, Spacer(1))
		content = append(content, Center(ErrorStyle.Render("❌ "+m.ErrorMessage), m.Width, 1))
	}

	// Join all content - consistent top alignment for all phases
	gameContent := lipgloss.JoinVertical(lipgloss.Center, content...)
	return gameContent
}

// renderStats renders the statistics view
func (m *Model) renderStats() string {
	stats := m.StatsManager.GetStats()

	var content []string

	// Header - use ASCII art for larger screens
	header := CreateStatsBanner(m.Width)
	content = append(content, header)
	content = append(content, Spacer(1))

	if stats.TotalGames == 0 {
		// No games played yet
		noGamesMsg := "No games played yet. Start playing to see statistics!"
		content = append(content, Center(SubtitleStyle.Render(noGamesMsg), m.Width, 1))

		footer := RenderFooter([]KeyBinding{
			{"Enter", "Play game"},
			{"q", "Main menu"},
		})
		content = append(content, footer)

		// Join all content vertically and center consistently
		menuContent := lipgloss.JoinVertical(lipgloss.Center, content...)
		horizontallyCentered := SafeCenter(menuContent, m.Width)
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, horizontallyCentered)
	}

	// Stats cards row
	totalCard := NewStatsCard(
		"Total Games",
		fmt.Sprintf("%d", stats.TotalGames),
		fmt.Sprintf("%.1f%% win rate", float64(stats.TotalWins)/float64(stats.TotalGames)*100),
		PrimaryColor,
	)

	winsCard := NewStatsCard(
		"Total Wins",
		fmt.Sprintf("%d", stats.TotalWins),
		fmt.Sprintf("%d losses", stats.TotalLosses),
		SecondaryColor,
	)

	streakCard := NewStatsCard(
		"Best Streak",
		fmt.Sprintf("%d", stats.StreakStats.LongestWinStreak),
		fmt.Sprintf("Current: %d", stats.StreakStats.CurrentWinStreak),
		AccentColor,
	)

	cardsRow := lipgloss.JoinHorizontal(lipgloss.Top,
		totalCard.Render(),
		" ",
		winsCard.Render(),
		" ",
		streakCard.Render(),
	)
	content = append(content, Center(cardsRow, m.Width, 1))
	content = append(content, Spacer(1))

	// Strategy comparison
	strategyTitle := StatsHeaderStyle.Render("STRATEGY PERFORMANCE")
	content = append(content, Center(strategyTitle, m.Width, 1))
	content = append(content, Spacer(1))

	// Progress bars for strategies
	if stats.StayStats.GamesPlayed > 0 {
		stayBar := NewProgressBar(
			stats.StayStats.Wins,
			stats.StayStats.GamesPlayed,
			40,
			fmt.Sprintf("Stay Strategy (%.1f%%)", stats.StayStats.WinRate*100),
		)
		content = append(content, Center(stayBar.Render(), m.Width, 1))
	}

	if stats.SwitchStats.GamesPlayed > 0 {
		switchBar := NewProgressBar(
			stats.SwitchStats.Wins,
			stats.SwitchStats.GamesPlayed,
			40,
			fmt.Sprintf("Switch Strategy (%.1f%%)", stats.SwitchStats.WinRate*100),
		)
		content = append(content, Center(switchBar.Render(), m.Width, 1))
	}

	content = append(content, Spacer(1))

	// Theoretical vs Actual
	theoryTitle := StatsHeaderStyle.Render("THEORETICAL vs ACTUAL")
	content = append(content, Center(theoryTitle, m.Width, 1))

	theoryLines := []string{
		"Stay should win:   33.3% (1/3 probability)",
		"Switch should win: 66.7% (2/3 probability)",
	}

	for _, line := range theoryLines {
		content = append(content, Center(MutedStyle.Render(line), m.Width, 1))
	}

	// Insights
	if stats.TotalGames >= 10 {
		content = append(content, Spacer(1))
		insightTitle := StatsHeaderStyle.Render("📈 INSIGHTS")
		content = append(content, Center(insightTitle, m.Width, 1))

		var insight string
		if stats.SwitchStats.WinRate > 0.6 {
			insight = "✅ Switching is proving more successful!"
		} else if stats.StayStats.WinRate > 0.4 {
			insight = "🎲 Results are still converging to theory."
		} else {
			insight = "📊 Play more games to see clearer patterns."
		}

		content = append(content, Center(SuccessStyle.Render(insight), m.Width, 1))
	}

	// Footer
	footer := RenderFooter([]KeyBinding{
		{"e", "Export stats"},
		{"r", "Reset stats"},
		{"ESC/q", "Return"},
	})
	content = append(content, footer)

	statsContent := lipgloss.JoinVertical(lipgloss.Center, content...)

	// Show reset confirmation popover if active
	if m.ShowResetConfirmation {
		popover := NewResetConfirmationPopover(
			m.ResetConfirmationNumbers,
			m.UserInputNumbers,
			m.CurrentInputIndex,
			60, // Width of the popover
		)

		// Overlay the popover on top of the stats content
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, popover.Render())
	}

	return statsContent
}

// Helper methods for door navigation and selection

// isDoorSelectable returns true if the door can be selected in the current game phase
func (m *Model) isDoorSelectable(doorIndex int) bool {
	if m.Game == nil {
		return false
	}

	// No doors are selectable during reveal countdown
	if m.IsRevealing {
		return false
	}

	switch m.Game.Phase {
	case game.InitialChoice:
		// All doors are selectable during initial choice
		return doorIndex >= 0 && doorIndex < game.NumDoors

	case game.HostReveal:
		// No doors are selectable during host reveal phase (countdown)
		return false

	case game.FinalChoice:
		// Only original choice and the other unopened door are selectable
		// Host-opened door should not be selectable
		return doorIndex != m.Game.HostOpenedDoor

	case game.GameOver:
		// No doors are selectable when game is over
		return false

	default:
		return false
	}
}

// getSelectableDoors returns a slice of door indices that can be selected
func (m *Model) getSelectableDoors() []int {
	if m.Game == nil {
		return []int{}
	}

	var selectable []int
	for i := 0; i < game.NumDoors; i++ {
		if m.isDoorSelectable(i) {
			selectable = append(selectable, i)
		}
	}
	return selectable
}

// moveCursorLeft moves cursor to the previous selectable door
func (m *Model) moveCursorLeft() {
	selectableDoors := m.getSelectableDoors()
	if len(selectableDoors) == 0 {
		return
	}

	// Find current cursor position in selectable doors
	currentIndex := -1
	for i, door := range selectableDoors {
		if door == m.DoorCursor {
			currentIndex = i
			break
		}
	}

	// Move to previous selectable door (wrap around)
	if currentIndex > 0 {
		m.DoorCursor = selectableDoors[currentIndex-1]
	} else if currentIndex == 0 {
		// Wrap to last selectable door
		m.DoorCursor = selectableDoors[len(selectableDoors)-1]
	} else {
		// Current cursor is not on a selectable door, move to first selectable
		m.DoorCursor = selectableDoors[0]
	}
}

// moveCursorRight moves cursor to the next selectable door
func (m *Model) moveCursorRight() {
	selectableDoors := m.getSelectableDoors()
	if len(selectableDoors) == 0 {
		return
	}

	// Find current cursor position in selectable doors
	currentIndex := -1
	for i, door := range selectableDoors {
		if door == m.DoorCursor {
			currentIndex = i
			break
		}
	}

	// Move to next selectable door (wrap around)
	if currentIndex >= 0 && currentIndex < len(selectableDoors)-1 {
		m.DoorCursor = selectableDoors[currentIndex+1]
	} else if currentIndex == len(selectableDoors)-1 {
		// Wrap to first selectable door
		m.DoorCursor = selectableDoors[0]
	} else {
		// Current cursor is not on a selectable door, move to first selectable
		m.DoorCursor = selectableDoors[0]
	}
}

// Animation helper methods

// startDoorOpenAnimation starts a door opening animation for the specified door
func (m *Model) startDoorOpenAnimation(doorIndex int) tea.Cmd {
	if !m.ShowAnimations || m.AnimationManager == nil {
		return nil
	}

	// Create and start door opening animation
	doorAnim := NewDoorOpenAnimation(doorIndex)
	m.DoorAnimations[doorIndex] = doorAnim
	m.AnimationManager.AddAnimation(doorAnim.Animation)
	m.AnimationManager.StartAnimation(doorAnim.ID)

	// Return the animation update command
	return m.AnimationManager.Update()
}

// startWinningAnimation starts a winning celebration animation
func (m *Model) startWinningAnimation() tea.Cmd {
	if !m.ShowAnimations || m.AnimationManager == nil {
		return nil
	}

	// Create pulse animation for winning door
	if m.Game != nil && m.Game.Result != nil {
		pulseAnim := NewPulseAnimation(
			"winning_pulse",
			WinningDoorStyle,
			CarColor,
		)
		m.AnimationManager.AddAnimation(pulseAnim.Animation)
		m.AnimationManager.StartAnimation(pulseAnim.ID)

		// Start the animation loop
		return m.AnimationManager.Update()
	}

	return nil
}

// getDoorAnimationFrame returns the current animation frame for a door
func (m *Model) getDoorAnimationFrame(doorIndex int) (string, lipgloss.Color) {
	if anim, exists := m.DoorAnimations[doorIndex]; exists && anim.IsRunning() {
		return anim.GetCurrentFrame()
	}
	return "🚪", DoorColor
}

// isAnimationRunning checks if any animations are currently running
func (m *Model) isAnimationRunning() bool {
	return m.AnimationManager != nil && m.AnimationManager.HasRunningAnimations()
}

// startRevealDelay starts the dramatic reveal delay
func (m *Model) startRevealDelay() tea.Cmd {
	m.IsRevealing = true
	m.RevealStartTime = time.Now()

	// Return a command that will send RevealDelayMsg after 2 seconds
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		return RevealDelayMsg{}
	})
}
