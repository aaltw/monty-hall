package game

import (
	"testing"
)

func TestNewDoor(t *testing.T) {
	door := NewDoor(1, 0, Car)

	if door.ID != 1 {
		t.Errorf("Expected ID 1, got %d", door.ID)
	}

	if door.Position != 0 {
		t.Errorf("Expected Position 0, got %d", door.Position)
	}

	if door.Content != Car {
		t.Errorf("Expected Car content, got %v", door.Content)
	}

	if door.State != Closed {
		t.Errorf("Expected Closed state, got %v", door.State)
	}
}

func TestDoorStates(t *testing.T) {
	door := NewDoor(1, 0, Goat)

	if !door.IsClosed() {
		t.Error("New door should be closed")
	}

	door.Open()
	if !door.IsOpen() {
		t.Error("Door should be open after Open()")
	}

	door.Reset()
	if !door.IsClosed() {
		t.Error("Door should be closed after Reset()")
	}

	door.Select()
	if !door.IsSelected() {
		t.Error("Door should be selected after Select()")
	}
}

func TestDoorContent(t *testing.T) {
	carDoor := NewDoor(1, 0, Car)
	goatDoor := NewDoor(2, 1, Goat)

	if !carDoor.HasCar() {
		t.Error("Car door should have car")
	}

	if carDoor.HasGoat() {
		t.Error("Car door should not have goat")
	}

	if !goatDoor.HasGoat() {
		t.Error("Goat door should have goat")
	}

	if goatDoor.HasCar() {
		t.Error("Goat door should not have car")
	}
}

func TestDoorString(t *testing.T) {
	door := NewDoor(1, 0, Car)

	expected := "Door 1: Closed"
	if door.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, door.String())
	}

	door.Open()
	expected = "Door 1: Opened (Car)"
	if door.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, door.String())
	}

	goatDoor := NewDoor(2, 1, Goat)
	goatDoor.Open()
	expected = "Door 2: Opened (Goat)"
	if goatDoor.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, goatDoor.String())
	}
}

func TestCreateDoorsWithRandomCar(t *testing.T) {
	doors := CreateDoorsWithRandomCar()

	if len(doors) != 3 {
		t.Errorf("Expected 3 doors, got %d", len(doors))
	}

	carCount := 0
	goatCount := 0

	for i, door := range doors {
		if door.ID != i+1 {
			t.Errorf("Door %d has wrong ID: %d", i, door.ID)
		}

		if door.Position != i {
			t.Errorf("Door %d has wrong position: %d", i, door.Position)
		}

		if door.HasCar() {
			carCount++
		} else {
			goatCount++
		}

		if !door.IsClosed() {
			t.Errorf("Door %d should start closed", i)
		}
	}

	if carCount != 1 {
		t.Errorf("Expected exactly 1 car, got %d", carCount)
	}

	if goatCount != 2 {
		t.Errorf("Expected exactly 2 goats, got %d", goatCount)
	}
}
