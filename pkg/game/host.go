package game

import (
	"errors"
	"fmt"
)

type Host struct {
	Name string
}

func NewHost() *Host {
	return &Host{
		Name: "Monty",
	}
}

func (h *Host) ChooseDoorToOpen(doors []*Door, playerChoice int) (int, error) {
	if len(doors) != NumDoors {
		return -1, fmt.Errorf("invalid number of doors: expected %d, got %d", NumDoors, len(doors))
	}

	if playerChoice < 0 || playerChoice >= len(doors) {
		return -1, errors.New("invalid player choice")
	}

	var validChoices []int
	for i, door := range doors {
		if i != playerChoice && door.HasGoat() {
			validChoices = append(validChoices, i)
		}
	}

	if len(validChoices) == 0 {
		return -1, errors.New("no valid doors to open")
	}

	if len(validChoices) == 1 {
		return validChoices[0], nil
	}

	randomIndex := SecureIntn(len(validChoices))
	return validChoices[randomIndex], nil
}

func (h *Host) GetSwitchRecommendation(doors []*Door, playerChoice int) (int, error) {
	if len(doors) != NumDoors {
		return -1, fmt.Errorf("invalid number of doors: expected %d, got %d", NumDoors, len(doors))
	}

	if playerChoice < 0 || playerChoice >= len(doors) {
		return -1, errors.New("invalid player choice")
	}

	for i, door := range doors {
		if i != playerChoice && !door.IsOpen() {
			return i, nil
		}
	}

	return -1, errors.New("no door available to switch to")
}

func (h *Host) ExplainProbability() string {
	return `The probability explanation:

When you first chose a door, you had a 1/3 chance of picking the car.
That means there was a 2/3 chance the car was behind one of the other two doors.

When I opened one of those doors to reveal a goat, I didn't change the 
probability that your original door has the car (still 1/3).

But now all of that 2/3 probability is concentrated on the remaining 
unopened door. So switching gives you a 2/3 chance of winning!`
}

func (h *Host) GetHint(doors []*Door, playerChoice int) string {
	if playerChoice < 0 || playerChoice >= len(doors) {
		return "Choose a door first!"
	}

	switchDoor, err := h.GetSwitchRecommendation(doors, playerChoice)
	if err != nil {
		return "I can't give you a hint right now."
	}

	return fmt.Sprintf("Statistically, switching to door %d gives you better odds!", switchDoor+1)
}
