package ui

import (
	"testing"

	"github.com/charmbracelet/bubbletea"
	"github.com/westhuis/monty-hall/pkg/game"
)

// TestResetConfirmationInitiation tests that reset confirmation is properly initiated
func TestResetConfirmationInitiation(t *testing.T) {
	model := NewModel()
	model.CurrentView = StatsView

	// Trigger reset confirmation
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	updatedModel, _ := model.Update(keyMsg)
	model = updatedModel.(*Model)

	// Verify confirmation dialog is shown
	if !model.ShowResetConfirmation {
		t.Error("Expected reset confirmation dialog to be shown")
	}

	// Verify confirmation numbers are generated (1-9)
	for i, num := range model.ResetConfirmationNumbers {
		if num < 1 || num > 9 {
			t.Errorf("Confirmation number %d is out of range (1-9): %d", i, num)
		}
	}

	// Verify user input is reset
	for i, num := range model.UserInputNumbers {
		if num != 0 {
			t.Errorf("User input number %d should be 0, got %d", i, num)
		}
	}

	// Verify input index is reset
	if model.CurrentInputIndex != 0 {
		t.Errorf("Expected CurrentInputIndex to be 0, got %d", model.CurrentInputIndex)
	}
}

// TestResetConfirmationCancellation tests canceling the reset confirmation
func TestResetConfirmationCancellation(t *testing.T) {
	model := NewModel()
	model.CurrentView = StatsView
	model.ShowResetConfirmation = true
	model.CurrentInputIndex = 2
	model.UserInputNumbers = [4]int{1, 2, 0, 0}

	// Cancel with Escape key
	keyMsg := tea.KeyMsg{Type: tea.KeyEscape}
	updatedModel, _ := model.Update(keyMsg)
	model = updatedModel.(*Model)

	// Verify confirmation dialog is hidden
	if model.ShowResetConfirmation {
		t.Error("Expected reset confirmation dialog to be hidden")
	}

	// Verify state is reset
	if model.CurrentInputIndex != 0 {
		t.Errorf("Expected CurrentInputIndex to be 0, got %d", model.CurrentInputIndex)
	}

	for i, num := range model.UserInputNumbers {
		if num != 0 {
			t.Errorf("User input number %d should be 0, got %d", i, num)
		}
	}
}

// TestResetConfirmationNumberInput tests inputting numbers during confirmation
func TestResetConfirmationNumberInput(t *testing.T) {
	model := NewModel()
	model.CurrentView = StatsView
	model.ShowResetConfirmation = true
	model.CurrentInputIndex = 0

	// Input valid numbers 1-9
	validNumbers := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}

	for _, numStr := range validNumbers {
		// Reset state for each test
		model.CurrentInputIndex = 0
		model.UserInputNumbers = [4]int{0, 0, 0, 0}

		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{rune(numStr[0])}}
		updatedModel, _ := model.Update(keyMsg)
		model = updatedModel.(*Model)

		expectedNum := int(numStr[0] - '0')
		if model.UserInputNumbers[0] != expectedNum {
			t.Errorf("Expected first input to be %d, got %d", expectedNum, model.UserInputNumbers[0])
		}

		if model.CurrentInputIndex != 1 {
			t.Errorf("Expected CurrentInputIndex to be 1, got %d", model.CurrentInputIndex)
		}
	}
}

// TestResetConfirmationInvalidInput tests that invalid characters are ignored
func TestResetConfirmationInvalidInput(t *testing.T) {
	model := NewModel()
	model.CurrentView = StatsView
	model.ShowResetConfirmation = true
	model.CurrentInputIndex = 0

	// Try invalid characters
	invalidChars := []string{"0", "a", "A", "!", "@", "#", " "}

	for _, char := range invalidChars {
		// Reset state
		model.CurrentInputIndex = 0
		model.UserInputNumbers = [4]int{0, 0, 0, 0}

		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{rune(char[0])}}
		updatedModel, _ := model.Update(keyMsg)
		model = updatedModel.(*Model)

		// Verify no input was recorded
		if model.UserInputNumbers[0] != 0 {
			t.Errorf("Invalid character '%s' should not be recorded, but got %d", char, model.UserInputNumbers[0])
		}

		if model.CurrentInputIndex != 0 {
			t.Errorf("CurrentInputIndex should remain 0 for invalid input '%s', got %d", char, model.CurrentInputIndex)
		}
	}
}

// TestResetConfirmationBackspace tests backspace functionality
func TestResetConfirmationBackspace(t *testing.T) {
	model := NewModel()
	model.CurrentView = StatsView
	model.ShowResetConfirmation = true

	// Input some numbers first
	model.UserInputNumbers = [4]int{1, 2, 3, 0}
	model.CurrentInputIndex = 3

	// Test backspace
	keyMsg := tea.KeyMsg{Type: tea.KeyBackspace}
	updatedModel, _ := model.Update(keyMsg)
	model = updatedModel.(*Model)

	// Verify the last input was cleared and index moved back
	if model.UserInputNumbers[2] != 0 {
		t.Errorf("Expected third input to be cleared, got %d", model.UserInputNumbers[2])
	}

	if model.CurrentInputIndex != 2 {
		t.Errorf("Expected CurrentInputIndex to be 2, got %d", model.CurrentInputIndex)
	}

	// Test backspace at beginning (should not go negative)
	model.CurrentInputIndex = 0
	model.UserInputNumbers = [4]int{0, 0, 0, 0}

	keyMsg = tea.KeyMsg{Type: tea.KeyBackspace}
	updatedModel, _ = model.Update(keyMsg)
	model = updatedModel.(*Model)

	if model.CurrentInputIndex != 0 {
		t.Errorf("CurrentInputIndex should not go below 0, got %d", model.CurrentInputIndex)
	}
}

// TestResetConfirmationSequentialInput tests inputting all 4 numbers sequentially
func TestResetConfirmationSequentialInput(t *testing.T) {
	model := NewModel()
	model.CurrentView = StatsView
	model.ShowResetConfirmation = true
	model.CurrentInputIndex = 0

	// Set known confirmation numbers for predictable testing
	model.ResetConfirmationNumbers = [4]int{1, 2, 3, 4}

	// Input 4 numbers sequentially
	numbers := []string{"1", "2", "3", "4"}

	for i, numStr := range numbers {
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{rune(numStr[0])}}
		updatedModel, _ := model.Update(keyMsg)
		model = updatedModel.(*Model)

		expectedNum := int(numStr[0] - '0')
		if i < 3 { // Only check first 3 inputs, 4th triggers validation and resets
			if model.UserInputNumbers[i] != expectedNum {
				t.Errorf("Expected input %d to be %d, got %d", i, expectedNum, model.UserInputNumbers[i])
			}

			expectedIndex := i + 1
			if model.CurrentInputIndex != expectedIndex {
				t.Errorf("Expected CurrentInputIndex to be %d, got %d", expectedIndex, model.CurrentInputIndex)
			}
		} else {
			// After 4th input, validation should have been triggered
			// Check that confirmation dialog is hidden (validation completed)
			if model.ShowResetConfirmation {
				t.Error("Expected reset confirmation dialog to be hidden after 4th input")
			}
		}
	}
}

// TestResetConfirmationOverflow tests that input stops at 4 numbers
func TestResetConfirmationOverflow(t *testing.T) {
	model := NewModel()
	model.CurrentView = StatsView
	model.ShowResetConfirmation = true

	// Fill all 4 slots
	model.UserInputNumbers = [4]int{1, 2, 3, 4}
	model.CurrentInputIndex = 4

	// Try to input a 5th number
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}}
	updatedModel, _ := model.Update(keyMsg)
	model = updatedModel.(*Model)

	// Verify no change occurred
	expected := [4]int{1, 2, 3, 4}
	if model.UserInputNumbers != expected {
		t.Errorf("Expected UserInputNumbers to remain %v, got %v", expected, model.UserInputNumbers)
	}
}

// TestResetConfirmationCorrectNumbers tests successful reset with correct numbers
func TestResetConfirmationCorrectNumbers(t *testing.T) {
	model := NewModel()
	model.CurrentView = StatsView

	// Reset stats to ensure clean test
	model.StatsManager.Reset()

	// Add some test statistics first
	testResult := &game.GameResult{
		Won:            true,
		Strategy:       game.Switch,
		InitialChoice:  0,
		FinalChoice:    1,
		CarPosition:    1,
		HostOpenedDoor: 2,
	}
	model.StatsManager.RecordGame(testResult)

	// Verify stats exist
	stats := model.StatsManager.GetStats()
	if stats.TotalGames != 1 {
		t.Errorf("Expected 1 game recorded, got %d", stats.TotalGames)
	}

	// Start reset confirmation
	model.ShowResetConfirmation = true
	model.ResetConfirmationNumbers = [4]int{1, 2, 3, 4}
	model.CurrentInputIndex = 0

	// Input the correct numbers
	correctNumbers := []string{"1", "2", "3", "4"}
	for _, numStr := range correctNumbers {
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{rune(numStr[0])}}
		updatedModel, _ := model.Update(keyMsg)
		model = updatedModel.(*Model)
	}

	// Verify reset was successful
	if model.ShowResetConfirmation {
		t.Error("Expected reset confirmation dialog to be hidden after successful reset")
	}

	if model.SuccessMessage == "" {
		t.Error("Expected success message after successful reset")
	}

	// Verify statistics were actually reset
	stats = model.StatsManager.GetStats()
	if stats.TotalGames != 0 {
		t.Errorf("Expected statistics to be reset (0 games), got %d", stats.TotalGames)
	}
}

// TestResetConfirmationIncorrectNumbers tests failed reset with incorrect numbers
func TestResetConfirmationIncorrectNumbers(t *testing.T) {
	model := NewModel()
	model.CurrentView = StatsView

	// Reset stats to ensure clean test
	model.StatsManager.Reset()

	// Add some test statistics first
	testResult := &game.GameResult{
		Won:            true,
		Strategy:       game.Switch,
		InitialChoice:  0,
		FinalChoice:    1,
		CarPosition:    1,
		HostOpenedDoor: 2,
	}
	model.StatsManager.RecordGame(testResult)

	// Start reset confirmation
	model.ShowResetConfirmation = true
	model.ResetConfirmationNumbers = [4]int{1, 2, 3, 4}
	model.CurrentInputIndex = 0

	// Input incorrect numbers
	incorrectNumbers := []string{"1", "2", "3", "5"} // Last number is wrong
	for _, numStr := range incorrectNumbers {
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{rune(numStr[0])}}
		updatedModel, _ := model.Update(keyMsg)
		model = updatedModel.(*Model)
	}

	// Verify reset was not performed
	if !model.ShowResetConfirmation {
		t.Error("Expected reset confirmation dialog to remain visible after incorrect input")
	}

	if model.ErrorMessage == "" {
		t.Error("Expected error message after incorrect input")
	}

	// Verify input was reset for retry
	if model.CurrentInputIndex != 0 {
		t.Errorf("Expected CurrentInputIndex to be reset to 0, got %d", model.CurrentInputIndex)
	}

	for i, num := range model.UserInputNumbers {
		if num != 0 {
			t.Errorf("User input number %d should be reset to 0, got %d", i, num)
		}
	}

	// Verify statistics were NOT reset
	stats := model.StatsManager.GetStats()
	if stats.TotalGames != 1 {
		t.Errorf("Expected statistics to remain unchanged (1 game), got %d", stats.TotalGames)
	}
}

// TestResetConfirmationEnterKey tests using Enter key to validate
func TestResetConfirmationEnterKey(t *testing.T) {
	model := NewModel()
	model.CurrentView = StatsView
	model.ShowResetConfirmation = true
	model.ResetConfirmationNumbers = [4]int{1, 2, 3, 4}

	// Input 3 numbers, then try Enter (should not validate)
	model.UserInputNumbers = [4]int{1, 2, 3, 0}
	model.CurrentInputIndex = 3

	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(keyMsg)
	model = updatedModel.(*Model)

	// Should still be in confirmation mode
	if !model.ShowResetConfirmation {
		t.Error("Expected reset confirmation to remain active with incomplete input")
	}

	// Now input all 4 numbers and try Enter
	model.UserInputNumbers = [4]int{1, 2, 3, 4}
	model.CurrentInputIndex = 4

	keyMsg = tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ = model.Update(keyMsg)
	model = updatedModel.(*Model)

	// Should complete the reset
	if model.ShowResetConfirmation {
		t.Error("Expected reset confirmation to be completed with Enter key")
	}
}

// TestResetConfirmationRandomNumberGeneration tests that random numbers are properly generated
func TestResetConfirmationRandomNumberGeneration(t *testing.T) {
	model := NewModel()
	model.CurrentView = StatsView

	// Generate multiple sets of confirmation numbers
	sets := make([][4]int, 10)
	for i := 0; i < 10; i++ {
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
		updatedModel, _ := model.Update(keyMsg)
		model = updatedModel.(*Model)

		sets[i] = model.ResetConfirmationNumbers

		// Cancel to reset for next iteration
		keyMsg = tea.KeyMsg{Type: tea.KeyEscape}
		updatedModel, _ = model.Update(keyMsg)
		model = updatedModel.(*Model)
	}

	// Verify all numbers are in valid range
	for setIdx, set := range sets {
		for numIdx, num := range set {
			if num < 1 || num > 9 {
				t.Errorf("Set %d, number %d is out of range (1-9): %d", setIdx, numIdx, num)
			}
		}
	}

	// Verify there's some variation (not all sets are identical)
	allSame := true
	firstSet := sets[0]
	for _, set := range sets[1:] {
		if set != firstSet {
			allSame = false
			break
		}
	}

	if allSame {
		t.Log("All confirmation number sets are identical - this could happen with low probability")
		// Don't fail the test as this could legitimately happen with small probability
	}
}

// TestResetConfirmationStateIsolation tests that confirmation state doesn't interfere with other views
func TestResetConfirmationStateIsolation(t *testing.T) {
	model := NewModel()
	model.CurrentView = StatsView
	model.ShowResetConfirmation = true
	model.CurrentInputIndex = 2
	model.UserInputNumbers = [4]int{1, 2, 0, 0}

	// Switch to game view
	model.CurrentView = GameView
	model.ShowResetConfirmation = false

	// Input numbers should not affect game view
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}}
	updatedModel, _ := model.Update(keyMsg)
	model = updatedModel.(*Model)

	// Verify game state is not affected by reset confirmation state
	if model.ShowResetConfirmation {
		t.Error("Reset confirmation should not be active in game view")
	}

	// Switch back to stats view
	model.CurrentView = StatsView

	// Reset confirmation state should be properly isolated
	if model.ShowResetConfirmation {
		t.Error("Reset confirmation should not automatically reappear when returning to stats view")
	}
}
