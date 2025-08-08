package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/westhuis/monty-hall/pkg/game"
)

// TestCompleteGameFlow tests the entire game flow from start to finish
func TestCompleteGameFlow(t *testing.T) {
	model := NewModel()
	// Reset statistics for clean test
	model.StatsManager.Reset()

	// Test 1: Start from main menu
	if model.CurrentView != MainMenuView {
		t.Errorf("Expected MainMenuView at start, got %v", model.CurrentView)
	}

	// Test 2: Navigate to Play Game and start a game
	model.MenuCursor = 0 // Play Game option
	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(keyMsg)
	model = updatedModel.(*Model)

	if model.CurrentView != GameView {
		t.Errorf("Expected GameView after selecting Play Game, got %v", model.CurrentView)
	}

	if model.Game == nil {
		t.Fatal("Game should be created when starting new game")
	}

	if model.Game.Phase != game.InitialChoice {
		t.Errorf("Expected InitialChoice phase at game start, got %v", model.Game.Phase)
	}

	// Test 3: Make initial door choice
	model.DoorCursor = 1 // Choose door 2
	keyMsg = tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ = model.Update(keyMsg)
	model = updatedModel.(*Model)

	if model.Game.Phase != game.FinalChoice {
		t.Errorf("Expected FinalChoice phase after initial choice, got %v", model.Game.Phase)
	}

	if model.Game.PlayerInitialChoice != 1 {
		t.Errorf("Expected player initial choice to be 1, got %d", model.Game.PlayerInitialChoice)
	}

	if model.Game.HostOpenedDoor == -1 {
		t.Error("Host should have opened a door after initial choice")
	}

	// Test 4: Make final choice (switch)
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	updatedModel, _ = model.Update(keyMsg)
	model = updatedModel.(*Model)

	if model.Game.Phase != game.GameOver {
		t.Errorf("Expected GameOver phase after final choice, got %v", model.Game.Phase)
	}

	if model.Game.Result == nil {
		t.Error("Game result should not be nil after game completion")
	}

	if model.Game.Result.Strategy != game.Switch {
		t.Errorf("Expected Switch strategy, got %v", model.Game.Result.Strategy)
	}

	// Test 5: Navigate to statistics
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	updatedModel, _ = model.Update(keyMsg)
	model = updatedModel.(*Model)

	if model.CurrentView != StatsView {
		t.Errorf("Expected StatsView after pressing 's', got %v", model.CurrentView)
	}

	// Test 6: Return to main menu
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	updatedModel, _ = model.Update(keyMsg)
	model = updatedModel.(*Model)

	if model.CurrentView != MainMenuView {
		t.Errorf("Expected MainMenuView after pressing 'q' from stats, got %v", model.CurrentView)
	}

	// Test 7: Verify statistics were recorded
	stats := model.StatsManager.GetStats()
	if stats.TotalGames != 1 {
		t.Errorf("Expected 1 total game recorded, got %d", stats.TotalGames)
	}

	if stats.SwitchStats.GamesPlayed != 1 {
		t.Errorf("Expected 1 switch game recorded, got %d", stats.SwitchStats.GamesPlayed)
	}
}

// TestGameFlowWithStay tests the game flow when player chooses to stay
func TestGameFlowWithStay(t *testing.T) {
	model := NewModel()

	// Start game
	model.MenuCursor = 0
	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(keyMsg)
	model = updatedModel.(*Model)

	// Make initial choice
	model.DoorCursor = 0 // Choose door 1
	keyMsg = tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ = model.Update(keyMsg)
	model = updatedModel.(*Model)

	// Stay with original choice
	keyMsg = tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ = model.Update(keyMsg)
	model = updatedModel.(*Model)

	if model.Game.Phase != game.GameOver {
		t.Errorf("Expected GameOver phase after staying, got %v", model.Game.Phase)
	}

	if model.Game.Result.Strategy != game.Stay {
		t.Errorf("Expected Stay strategy, got %v", model.Game.Result.Strategy)
	}

	if model.Game.Result.FinalChoice != model.Game.Result.InitialChoice {
		t.Error("Final choice should equal initial choice when staying")
	}
}

// TestMultipleGames tests playing multiple games in sequence
func TestMultipleGames(t *testing.T) {
	model := NewModel()
	// Reset statistics for clean test
	model.StatsManager.Reset()

	// Play 3 games
	for i := 0; i < 3; i++ {
		// Start new game
		model.MenuCursor = 0
		keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
		updatedModel, _ := model.Update(keyMsg)
		model = updatedModel.(*Model)

		// Make initial choice
		keyMsg = tea.KeyMsg{Type: tea.KeyEnter}
		updatedModel, _ = model.Update(keyMsg)
		model = updatedModel.(*Model)

		// Make final choice (alternate between switch and stay)
		if i%2 == 0 {
			// Switch
			keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
		} else {
			// Stay
			keyMsg = tea.KeyMsg{Type: tea.KeyEnter}
		}
		updatedModel, _ = model.Update(keyMsg)
		model = updatedModel.(*Model)

		if model.Game.Phase != game.GameOver {
			t.Errorf("Game %d should be over, got phase %v", i+1, model.Game.Phase)
		}
	}

	// Check statistics
	stats := model.StatsManager.GetStats()
	if stats.TotalGames != 3 {
		t.Errorf("Expected 3 total games, got %d", stats.TotalGames)
	}

	if stats.SwitchStats.GamesPlayed != 2 {
		t.Errorf("Expected 2 switch games, got %d", stats.SwitchStats.GamesPlayed)
	}

	if stats.StayStats.GamesPlayed != 1 {
		t.Errorf("Expected 1 stay game, got %d", stats.StayStats.GamesPlayed)
	}
}

// TestUIRendering tests that all views render without errors
func TestUIRendering(t *testing.T) {
	model := NewModel()
	model.Width = 100
	model.Height = 30

	// Test main menu rendering
	view := model.View()
	if view == "" {
		t.Error("Main menu view should not be empty")
	}

	// Test help rendering
	model.ShowHelp = true
	helpView := model.View()
	if helpView == "" {
		t.Error("Help view should not be empty")
	}
	model.ShowHelp = false

	// Test game view rendering
	model.CurrentView = GameView
	model.Game = game.NewGame()
	gameView := model.View()
	if gameView == "" {
		t.Error("Game view should not be empty")
	}

	// Test game view in different phases
	model.Game.MakeInitialChoice(0)
	gameViewFinal := model.View()
	if gameViewFinal == "" {
		t.Error("Game view in final choice phase should not be empty")
	}

	// Complete the game
	model.Game.SwitchChoice()
	gameViewComplete := model.View()
	if gameViewComplete == "" {
		t.Error("Game view in complete phase should not be empty")
	}

	// Test stats view rendering
	model.CurrentView = StatsView
	statsView := model.View()
	if statsView == "" {
		t.Error("Stats view should not be empty")
	}
}

// TestErrorHandling tests error scenarios
func TestErrorHandling(t *testing.T) {
	model := NewModel()

	// Test error message display
	errorMsg := ErrorMsg{Error: "Test error message"}
	updatedModel, _ := model.Update(errorMsg)
	model = updatedModel.(*Model)

	if model.ErrorMessage != "Test error message" {
		t.Errorf("Expected error message 'Test error message', got '%s'", model.ErrorMessage)
	}

	// Test error message clearing
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	updatedModel, _ = model.Update(keyMsg)
	model = updatedModel.(*Model)

	if model.ErrorMessage != "" {
		t.Error("Error message should be cleared after key press")
	}

	// Test success message display
	successMsg := SuccessMsg{Message: "Test success message"}
	updatedModel, _ = model.Update(successMsg)
	model = updatedModel.(*Model)

	if model.SuccessMessage != "Test success message" {
		t.Errorf("Expected success message 'Test success message', got '%s'", model.SuccessMessage)
	}
}

// TestKeyboardNavigation tests all keyboard shortcuts
func TestKeyboardNavigation(t *testing.T) {
	model := NewModel()

	// Test menu navigation
	keyMsg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ := model.Update(keyMsg)
	model = updatedModel.(*Model)

	if model.MenuCursor != 1 {
		t.Errorf("Expected menu cursor at 1, got %d", model.MenuCursor)
	}

	keyMsg = tea.KeyMsg{Type: tea.KeyUp}
	updatedModel, _ = model.Update(keyMsg)
	model = updatedModel.(*Model)

	if model.MenuCursor != 0 {
		t.Errorf("Expected menu cursor at 0, got %d", model.MenuCursor)
	}

	// Test help toggle
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}
	updatedModel, _ = model.Update(keyMsg)
	model = updatedModel.(*Model)

	if !model.ShowHelp {
		t.Error("Help should be shown after pressing 'h'")
	}

	// Test help close
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	updatedModel, _ = model.Update(keyMsg)
	model = updatedModel.(*Model)

	if model.ShowHelp {
		t.Error("Help should be hidden after pressing 'q'")
	}

	// Test game navigation
	model.CurrentView = GameView
	model.Game = game.NewGame()

	keyMsg = tea.KeyMsg{Type: tea.KeyRight}
	updatedModel, _ = model.Update(keyMsg)
	model = updatedModel.(*Model)

	if model.DoorCursor != 1 {
		t.Errorf("Expected door cursor at 1, got %d", model.DoorCursor)
	}

	// Test number key navigation
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}}
	updatedModel, _ = model.Update(keyMsg)
	model = updatedModel.(*Model)

	if model.DoorCursor != 2 {
		t.Errorf("Expected door cursor at 2 after pressing '3', got %d", model.DoorCursor)
	}
}

// TestDoorNavigationBugFix tests that the door navigation bug is fixed
func TestDoorNavigationBugFix(t *testing.T) {
	model := NewModel()
	model.StatsManager.Reset()

	// Start a new game
	model.MenuCursor = 0
	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(keyMsg)
	model = updatedModel.(*Model)

	// Make initial choice (door 0)
	model.DoorCursor = 0
	keyMsg = tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ = model.Update(keyMsg)
	model = updatedModel.(*Model)

	// Verify we're in final choice phase
	if model.Game.Phase != game.FinalChoice {
		t.Fatalf("Expected FinalChoice phase, got %v", model.Game.Phase)
	}

	hostOpenedDoor := model.Game.HostOpenedDoor
	if hostOpenedDoor == -1 {
		t.Fatal("Host should have opened a door")
	}

	// Test that arrow keys don't allow navigation to host-opened door
	for i := 0; i < 10; i++ { // Try multiple times
		keyMsg = tea.KeyMsg{Type: tea.KeyRight}
		updatedModel, _ = model.Update(keyMsg)
		model = updatedModel.(*Model)

		if model.DoorCursor == hostOpenedDoor {
			t.Errorf("Arrow key navigation should not allow cursor on host-opened door %d", hostOpenedDoor)
		}

		keyMsg = tea.KeyMsg{Type: tea.KeyLeft}
		updatedModel, _ = model.Update(keyMsg)
		model = updatedModel.(*Model)

		if model.DoorCursor == hostOpenedDoor {
			t.Errorf("Arrow key navigation should not allow cursor on host-opened door %d", hostOpenedDoor)
		}
	}

	// Test that only 2 doors are navigable (original choice and the other unopened door)
	visitedDoors := make(map[int]bool)
	for i := 0; i < 20; i++ { // Navigate extensively
		keyMsg = tea.KeyMsg{Type: tea.KeyRight}
		updatedModel, _ = model.Update(keyMsg)
		model = updatedModel.(*Model)
		visitedDoors[model.DoorCursor] = true
	}

	if len(visitedDoors) != 2 {
		t.Errorf("Should only be able to navigate between 2 doors, but visited %d doors: %v", len(visitedDoors), visitedDoors)
	}

	// Verify host-opened door is never visited
	if visitedDoors[hostOpenedDoor] {
		t.Errorf("Host-opened door %d should never be visitable", hostOpenedDoor)
	}
}
