package game

import (
	"errors"
	"fmt"
	"time"
)

const (
	NumDoors = 3 // Standard Monty Hall problem uses 3 doors
)

type GamePhase int

const (
	Setup GamePhase = iota
	InitialChoice
	HostReveal
	FinalChoice
	GameOver
)

type PlayerStrategy int

const (
	Stay PlayerStrategy = iota
	Switch
)

// GameResult represents the outcome of a completed Monty Hall game
type GameResult struct {
	Won            bool           // Whether the player won the car
	Strategy       PlayerStrategy // Whether the player stayed or switched
	InitialChoice  int            // The door initially chosen by the player (0-2)
	FinalChoice    int            // The door finally chosen by the player (0-2)
	CarPosition    int            // The door where the car was located (0-2)
	HostOpenedDoor int            // The door opened by the host (0-2)
	GameDuration   time.Duration  // How long the game took to complete
	Timestamp      time.Time      // When the game was completed
}

type Game struct {
	Doors               []*Door
	Phase               GamePhase
	PlayerInitialChoice int
	PlayerFinalChoice   int
	HostOpenedDoor      int
	CarPosition         int
	GameStartTime       time.Time
	Result              *GameResult
	Host                *Host
}

func NewGame() *Game {

	game := &Game{
		Doors:               CreateDoorsWithRandomCar(),
		Phase:               Setup,
		PlayerInitialChoice: -1,
		PlayerFinalChoice:   -1,
		HostOpenedDoor:      -1,
		GameStartTime:       time.Now(),
		Host:                NewHost(),
	}

	for i, door := range game.Doors {
		if door.HasCar() {
			game.CarPosition = i
			break
		}
	}

	game.Phase = InitialChoice
	return game
}

func (g *Game) MakeInitialChoice(doorIndex int) error {
	if g.Phase != InitialChoice {
		return errors.New("not in initial choice phase")
	}

	if doorIndex < 0 || doorIndex >= len(g.Doors) {
		return fmt.Errorf("door index %d out of range [0-%d]", doorIndex, len(g.Doors)-1)
	}

	g.PlayerInitialChoice = doorIndex
	g.Doors[doorIndex].Select()
	g.Phase = HostReveal

	hostDoor, err := g.Host.ChooseDoorToOpen(g.Doors, doorIndex)
	if err != nil {
		return fmt.Errorf("host error: %w", err)
	}

	g.HostOpenedDoor = hostDoor
	g.Doors[hostDoor].Open()
	g.Phase = FinalChoice

	return nil
}

func (g *Game) MakeFinalChoice(doorIndex int) error {
	if g.Phase != FinalChoice {
		return errors.New("not in final choice phase")
	}

	if doorIndex < 0 || doorIndex >= len(g.Doors) {
		return fmt.Errorf("door index %d out of range [0-%d]", doorIndex, len(g.Doors)-1)
	}

	if g.Doors[doorIndex].IsOpen() {
		return errors.New("cannot choose an opened door")
	}

	g.Doors[g.PlayerInitialChoice].Reset()
	g.PlayerFinalChoice = doorIndex
	g.Doors[doorIndex].Select()
	g.Phase = GameOver

	g.calculateResult()
	return nil
}

func (g *Game) SwitchChoice() error {
	if g.Phase != FinalChoice {
		return errors.New("not in final choice phase")
	}

	for i, door := range g.Doors {
		if !door.IsOpen() && i != g.PlayerInitialChoice {
			return g.MakeFinalChoice(i)
		}
	}

	return errors.New("no valid door to switch to")
}

func (g *Game) StayWithChoice() error {
	return g.MakeFinalChoice(g.PlayerInitialChoice)
}

func (g *Game) calculateResult() {
	strategy := Stay
	if g.PlayerFinalChoice != g.PlayerInitialChoice {
		strategy = Switch
	}

	won := g.Doors[g.PlayerFinalChoice].HasCar()
	duration := time.Since(g.GameStartTime)

	g.Result = &GameResult{
		Won:            won,
		Strategy:       strategy,
		InitialChoice:  g.PlayerInitialChoice + 1, // 1-indexed for display
		FinalChoice:    g.PlayerFinalChoice + 1,   // 1-indexed for display
		CarPosition:    g.CarPosition + 1,         // 1-indexed for display
		HostOpenedDoor: g.HostOpenedDoor + 1,      // 1-indexed for display
		GameDuration:   duration,
		Timestamp:      time.Now(),
	}
}

func (g *Game) GetAvailableChoices() []int {
	var choices []int
	for i, door := range g.Doors {
		if !door.IsOpen() {
			choices = append(choices, i)
		}
	}
	return choices
}

func (g *Game) GetSelectedDoor() int {
	for i, door := range g.Doors {
		if door.IsSelected() {
			return i
		}
	}
	return -1
}

func (g *Game) GetOpenDoors() []int {
	var openDoors []int
	for i, door := range g.Doors {
		if door.IsOpen() {
			openDoors = append(openDoors, i)
		}
	}
	return openDoors
}

func (g *Game) IsGameOver() bool {
	return g.Phase == GameOver
}

func (g *Game) GetPhaseDescription() string {
	switch g.Phase {
	case Setup:
		return "Setting up the game..."
	case InitialChoice:
		return "Choose your initial door"
	case HostReveal:
		return "Host is revealing a door..."
	case FinalChoice:
		return "Make your final choice: stay or switch?"
	case GameOver:
		return "Game over!"
	default:
		return "Unknown phase"
	}
}

func (g *Game) Reset() {
	*g = *NewGame()
}

func (g *Game) GetGameState() map[string]interface{} {
	return map[string]interface{}{
		"phase":               g.Phase,
		"doors":               g.Doors,
		"playerInitialChoice": g.PlayerInitialChoice,
		"playerFinalChoice":   g.PlayerFinalChoice,
		"hostOpenedDoor":      g.HostOpenedDoor,
		"carPosition":         g.CarPosition,
		"result":              g.Result,
		"availableChoices":    g.GetAvailableChoices(),
		"phaseDescription":    g.GetPhaseDescription(),
	}
}
