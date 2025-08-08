package game

import (
	"testing"
)

func TestNewHost(t *testing.T) {
	host := NewHost()

	if host.Name != "Monty" {
		t.Errorf("Expected host name 'Monty', got '%s'", host.Name)
	}
}

func TestHostChooseDoorToOpen(t *testing.T) {
	host := NewHost()

	doors := []*Door{
		NewDoor(1, 0, Goat),
		NewDoor(2, 1, Car),
		NewDoor(3, 2, Goat),
	}

	playerChoice := 1

	doorToOpen, err := host.ChooseDoorToOpen(doors, playerChoice)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if doorToOpen == playerChoice {
		t.Error("Host should not open the player's chosen door")
	}

	if doors[doorToOpen].HasCar() {
		t.Error("Host should not open a door with a car")
	}

	if doorToOpen != 0 && doorToOpen != 2 {
		t.Errorf("Host should open door 0 or 2, got %d", doorToOpen)
	}
}

func TestHostChooseDoorToOpenInvalidInputs(t *testing.T) {
	host := NewHost()

	doors := []*Door{
		NewDoor(1, 0, Goat),
		NewDoor(2, 1, Car),
	}

	_, err := host.ChooseDoorToOpen(doors, 0)
	if err == nil {
		t.Error("Expected error for invalid number of doors")
	}

	doors = []*Door{
		NewDoor(1, 0, Goat),
		NewDoor(2, 1, Car),
		NewDoor(3, 2, Goat),
	}

	_, err = host.ChooseDoorToOpen(doors, -1)
	if err == nil {
		t.Error("Expected error for invalid player choice")
	}

	_, err = host.ChooseDoorToOpen(doors, 3)
	if err == nil {
		t.Error("Expected error for invalid player choice")
	}
}

func TestHostChooseDoorToOpenNoValidChoices(t *testing.T) {
	host := NewHost()

	doors := []*Door{
		NewDoor(1, 0, Car),
		NewDoor(2, 1, Car),
		NewDoor(3, 2, Car),
	}

	_, err := host.ChooseDoorToOpen(doors, 0)
	if err == nil {
		t.Error("Expected error when no valid doors to open")
	}
}

func TestHostGetSwitchRecommendation(t *testing.T) {
	host := NewHost()

	doors := []*Door{
		NewDoor(1, 0, Goat),
		NewDoor(2, 1, Car),
		NewDoor(3, 2, Goat),
	}

	doors[0].Open()
	playerChoice := 2

	switchDoor, err := host.GetSwitchRecommendation(doors, playerChoice)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if switchDoor == playerChoice {
		t.Error("Switch recommendation should not be the player's current choice")
	}

	if doors[switchDoor].IsOpen() {
		t.Error("Switch recommendation should not be an open door")
	}

	if switchDoor != 1 {
		t.Errorf("Expected switch recommendation to be door 1, got %d", switchDoor)
	}
}

func TestHostExplainProbability(t *testing.T) {
	host := NewHost()

	explanation := host.ExplainProbability()
	if explanation == "" {
		t.Error("Explanation should not be empty")
	}

	if len(explanation) < 100 {
		t.Error("Explanation should be substantial")
	}
}

func TestHostGetHint(t *testing.T) {
	host := NewHost()

	doors := []*Door{
		NewDoor(1, 0, Goat),
		NewDoor(2, 1, Car),
		NewDoor(3, 2, Goat),
	}

	doors[0].Open()
	playerChoice := 2

	hint := host.GetHint(doors, playerChoice)
	if hint == "" {
		t.Error("Hint should not be empty")
	}

	hint = host.GetHint(doors, -1)
	expected := "Choose a door first!"
	if hint != expected {
		t.Errorf("Expected '%s', got '%s'", expected, hint)
	}
}
