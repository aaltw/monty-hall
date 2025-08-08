package game

import (
	"fmt"
	"math/rand"
	"time"
)

type DoorState int

const (
	Closed DoorState = iota
	Opened
	Selected
)

type DoorContent int

const (
	Goat DoorContent = iota
	Car
)

type Door struct {
	ID       int
	State    DoorState
	Content  DoorContent
	Position int
}

func NewDoor(id, position int, content DoorContent) *Door {
	return &Door{
		ID:       id,
		State:    Closed,
		Content:  content,
		Position: position,
	}
}

func (d *Door) Open() {
	d.State = Opened
}

func (d *Door) Select() {
	d.State = Selected
}

func (d *Door) Reset() {
	d.State = Closed
}

func (d *Door) IsOpen() bool {
	return d.State == Opened
}

func (d *Door) IsSelected() bool {
	return d.State == Selected
}

func (d *Door) IsClosed() bool {
	return d.State == Closed
}

func (d *Door) HasCar() bool {
	return d.Content == Car
}

func (d *Door) HasGoat() bool {
	return d.Content == Goat
}

func (d *Door) String() string {
	var state string
	switch d.State {
	case Closed:
		state = "Closed"
	case Opened:
		state = "Opened"
	case Selected:
		state = "Selected"
	}

	var content string
	if d.State == Opened {
		if d.Content == Car {
			content = " (Car)"
		} else {
			content = " (Goat)"
		}
	}

	return fmt.Sprintf("Door %d: %s%s", d.ID, state, content)
}

func CreateDoorsWithRandomCar() []*Door {
	doors := make([]*Door, NumDoors)

	// Use a properly seeded random source for better randomness
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	carPosition := rng.Intn(NumDoors)

	for i := range NumDoors {
		content := Goat
		if i == carPosition {
			content = Car
		}
		doors[i] = NewDoor(i+1, i, content)
	}

	return doors
}
