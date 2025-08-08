package game

import (
	"testing"
)

func TestNewGame(t *testing.T) {
	game := NewGame()

	if len(game.Doors) != 3 {
		t.Errorf("Expected 3 doors, got %d", len(game.Doors))
	}

	if game.Phase != InitialChoice {
		t.Errorf("Expected InitialChoice phase, got %v", game.Phase)
	}

	if game.PlayerInitialChoice != -1 {
		t.Errorf("Expected PlayerInitialChoice -1, got %d", game.PlayerInitialChoice)
	}

	if game.PlayerFinalChoice != -1 {
		t.Errorf("Expected PlayerFinalChoice -1, got %d", game.PlayerFinalChoice)
	}

	if game.HostOpenedDoor != -1 {
		t.Errorf("Expected HostOpenedDoor -1, got %d", game.HostOpenedDoor)
	}

	if game.Host == nil {
		t.Error("Host should not be nil")
	}

	carCount := 0
	for _, door := range game.Doors {
		if door.HasCar() {
			carCount++
		}
	}

	if carCount != 1 {
		t.Errorf("Expected exactly 1 car, got %d", carCount)
	}
}

func TestMakeInitialChoice(t *testing.T) {
	game := NewGame()

	err := game.MakeInitialChoice(1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if game.PlayerInitialChoice != 1 {
		t.Errorf("Expected PlayerInitialChoice 1, got %d", game.PlayerInitialChoice)
	}

	if !game.Doors[1].IsSelected() {
		t.Error("Door 1 should be selected")
	}

	if game.Phase != FinalChoice {
		t.Errorf("Expected FinalChoice phase, got %v", game.Phase)
	}

	if game.HostOpenedDoor == -1 {
		t.Error("Host should have opened a door")
	}

	if game.Doors[game.HostOpenedDoor].HasCar() {
		t.Error("Host should not open a door with a car")
	}

	if game.HostOpenedDoor == game.PlayerInitialChoice {
		t.Error("Host should not open the player's chosen door")
	}
}

func TestMakeInitialChoiceInvalidInputs(t *testing.T) {
	game := NewGame()

	err := game.MakeInitialChoice(-1)
	if err == nil {
		t.Error("Expected error for invalid door index")
	}

	err = game.MakeInitialChoice(3)
	if err == nil {
		t.Error("Expected error for invalid door index")
	}

	game.Phase = GameOver
	err = game.MakeInitialChoice(1)
	if err == nil {
		t.Error("Expected error for wrong phase")
	}
}

func TestMakeFinalChoice(t *testing.T) {
	game := NewGame()
	game.MakeInitialChoice(0)

	availableChoices := game.GetAvailableChoices()
	if len(availableChoices) != 2 {
		t.Errorf("Expected 2 available choices, got %d", len(availableChoices))
	}

	var finalChoice int
	for _, choice := range availableChoices {
		if choice != game.PlayerInitialChoice {
			finalChoice = choice
			break
		}
	}

	err := game.MakeFinalChoice(finalChoice)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if game.PlayerFinalChoice != finalChoice {
		t.Errorf("Expected PlayerFinalChoice %d, got %d", finalChoice, game.PlayerFinalChoice)
	}

	if game.Phase != GameOver {
		t.Errorf("Expected GameOver phase, got %v", game.Phase)
	}

	if game.Result == nil {
		t.Error("Game result should not be nil")
	}
}

func TestMakeFinalChoiceInvalidInputs(t *testing.T) {
	game := NewGame()
	game.MakeInitialChoice(0)

	err := game.MakeFinalChoice(game.HostOpenedDoor)
	if err == nil {
		t.Error("Expected error for choosing opened door")
	}

	err = game.MakeFinalChoice(-1)
	if err == nil {
		t.Error("Expected error for invalid door index")
	}

	err = game.MakeFinalChoice(3)
	if err == nil {
		t.Error("Expected error for invalid door index")
	}
}

func TestSwitchChoice(t *testing.T) {
	game := NewGame()
	game.MakeInitialChoice(0)

	initialChoice := game.PlayerInitialChoice

	err := game.SwitchChoice()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if game.PlayerFinalChoice == initialChoice {
		t.Error("Final choice should be different from initial choice when switching")
	}

	if game.Result.Strategy != Switch {
		t.Errorf("Expected Switch strategy, got %v", game.Result.Strategy)
	}
}

func TestStayWithChoice(t *testing.T) {
	game := NewGame()
	game.MakeInitialChoice(1)

	initialChoice := game.PlayerInitialChoice

	err := game.StayWithChoice()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if game.PlayerFinalChoice != initialChoice {
		t.Error("Final choice should be same as initial choice when staying")
	}

	if game.Result.Strategy != Stay {
		t.Errorf("Expected Stay strategy, got %v", game.Result.Strategy)
	}
}

func TestGameResult(t *testing.T) {
	game := NewGame()
	game.MakeInitialChoice(0)
	game.StayWithChoice()

	result := game.Result
	if result == nil {
		t.Fatal("Result should not be nil")
	}

	if result.InitialChoice != 1 {
		t.Errorf("Expected InitialChoice 1 (1-indexed), got %d", result.InitialChoice)
	}

	if result.FinalChoice != 1 {
		t.Errorf("Expected FinalChoice 1 (1-indexed), got %d", result.FinalChoice)
	}

	if result.Strategy != Stay {
		t.Errorf("Expected Stay strategy, got %v", result.Strategy)
	}

	if result.GameDuration <= 0 {
		t.Error("Game duration should be positive")
	}

	if result.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}

	expectedWin := game.Doors[0].HasCar()
	if result.Won != expectedWin {
		t.Errorf("Expected Won %v, got %v", expectedWin, result.Won)
	}
}

func TestGetAvailableChoices(t *testing.T) {
	game := NewGame()

	choices := game.GetAvailableChoices()
	if len(choices) != 3 {
		t.Errorf("Expected 3 available choices initially, got %d", len(choices))
	}

	game.MakeInitialChoice(0)

	choices = game.GetAvailableChoices()
	if len(choices) != 2 {
		t.Errorf("Expected 2 available choices after host opens door, got %d", len(choices))
	}

	for _, choice := range choices {
		if game.Doors[choice].IsOpen() {
			t.Error("Available choices should not include opened doors")
		}
	}
}

func TestGetSelectedDoor(t *testing.T) {
	game := NewGame()

	selected := game.GetSelectedDoor()
	if selected != -1 {
		t.Errorf("Expected no selected door initially, got %d", selected)
	}

	game.MakeInitialChoice(1)

	selected = game.GetSelectedDoor()
	if selected != 1 {
		t.Errorf("Expected selected door 1, got %d", selected)
	}
}

func TestGetOpenDoors(t *testing.T) {
	game := NewGame()

	openDoors := game.GetOpenDoors()
	if len(openDoors) != 0 {
		t.Errorf("Expected no open doors initially, got %d", len(openDoors))
	}

	game.MakeInitialChoice(0)

	openDoors = game.GetOpenDoors()
	if len(openDoors) != 1 {
		t.Errorf("Expected 1 open door after host action, got %d", len(openDoors))
	}

	if openDoors[0] != game.HostOpenedDoor {
		t.Errorf("Expected open door %d, got %d", game.HostOpenedDoor, openDoors[0])
	}
}

func TestIsGameOver(t *testing.T) {
	game := NewGame()

	if game.IsGameOver() {
		t.Error("Game should not be over initially")
	}

	game.MakeInitialChoice(0)

	if game.IsGameOver() {
		t.Error("Game should not be over after initial choice")
	}

	game.StayWithChoice()

	if !game.IsGameOver() {
		t.Error("Game should be over after final choice")
	}
}

func TestGetPhaseDescription(t *testing.T) {
	game := NewGame()

	description := game.GetPhaseDescription()
	if description == "" {
		t.Error("Phase description should not be empty")
	}

	game.MakeInitialChoice(0)
	description = game.GetPhaseDescription()
	if description == "" {
		t.Error("Phase description should not be empty")
	}

	game.StayWithChoice()
	description = game.GetPhaseDescription()
	if description == "" {
		t.Error("Phase description should not be empty")
	}
}

func TestReset(t *testing.T) {
	game := NewGame()
	game.MakeInitialChoice(0)
	game.StayWithChoice()

	if !game.IsGameOver() {
		t.Error("Game should be over before reset")
	}

	game.Reset()

	if game.IsGameOver() {
		t.Error("Game should not be over after reset")
	}

	if game.Phase != InitialChoice {
		t.Errorf("Expected InitialChoice phase after reset, got %v", game.Phase)
	}

	if game.PlayerInitialChoice != -1 {
		t.Errorf("Expected PlayerInitialChoice -1 after reset, got %d", game.PlayerInitialChoice)
	}
}

func TestGetGameState(t *testing.T) {
	game := NewGame()

	state := game.GetGameState()
	if state == nil {
		t.Error("Game state should not be nil")
	}

	if state["phase"] != game.Phase {
		t.Error("Game state should include current phase")
	}

	if state["doors"] == nil {
		t.Error("Game state should include doors")
	}

	if state["phaseDescription"] == "" {
		t.Error("Game state should include phase description")
	}
}

func TestProbabilityDistribution(t *testing.T) {
	const numGames = 1000
	switchWins := 0
	stayWins := 0

	for i := 0; i < numGames; i++ {
		game := NewGame()
		game.MakeInitialChoice(0)

		switchGame := *game
		switchGame.SwitchChoice()
		if switchGame.Result.Won {
			switchWins++
		}

		stayGame := *game
		stayGame.StayWithChoice()
		if stayGame.Result.Won {
			stayWins++
		}
	}

	switchRate := float64(switchWins) / float64(numGames)
	stayRate := float64(stayWins) / float64(numGames)

	if switchRate < 0.6 || switchRate > 0.7 {
		t.Errorf("Switch win rate should be around 2/3 (0.667), got %.3f", switchRate)
	}

	if stayRate < 0.3 || stayRate > 0.4 {
		t.Errorf("Stay win rate should be around 1/3 (0.333), got %.3f", stayRate)
	}

	t.Logf("Switch win rate: %.3f, Stay win rate: %.3f", switchRate, stayRate)
}
