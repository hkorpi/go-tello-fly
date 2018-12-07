package main

import (
	"deus.solita.fi/Solita/projects/drone_code_camp/repositories/git/ddr.git"
	"strconv"
)

const MinSpeed = 10

type DroneState struct {
	flying    bool
	operation OperationId
	speed     int
	message   string
}

func operation(state DroneState, next OperationId) DroneState {
	return DroneState{
		state.flying,
		next,
		state.speed,
		state.message,
	}
}

type OperationId int

const (
	Up OperationId = iota
	Down
	Left
	Right
	Forward
	Backward
	TurnLeft
	TurnRight

	Hover
	Land
	Takeoff

	NOOP
)

type Operation struct {
	fn          func(d ddr.Drone, val int) error
	description string
}

var operations = map[OperationId]Operation{
	Up:    Operation{ddr.Drone.Up, "Going up"},
	Down:  Operation{ddr.Drone.Down, "Going down"},
	Left:  Operation{ddr.Drone.Left, "Going left"},
	Right: Operation{ddr.Drone.Right, "Going right"},

	Forward:   Operation{ddr.Drone.Forward, "Going forward"},
	Backward:  Operation{ddr.Drone.Backward, "Going backward"},
	TurnLeft:  Operation{ddr.Drone.CounterClockwise, "Rotating left"},
	TurnRight: Operation{ddr.Drone.Clockwise, "Rotating right"},
}

func fly(state DroneState, next OperationId) DroneState {
	if state.flying {
		speed := MinSpeed
		if state.operation == next {
			speed = min(2*state.speed, 100)
		}
		description := operations[next].description
		return DroneState{
			true,
			next,
			speed,
			description + " " + strconv.Itoa(speed),
		}
	} else {
		return operation(state, NOOP)
	}
}

func toggleMode(state DroneState) DroneState {
	if state.flying && state.operation != Hover {
		return DroneState{
			true,
			Hover,
			MinSpeed,
			"Hovering",
		}
	} else if state.flying && state.operation == Hover {
		return DroneState{
			false,
			Land,
			MinSpeed,
			"Landing",
		}
	} else {
		return DroneState{
			true,
			Takeoff,
			MinSpeed,
			"Takeoff",
		}
	}
}

func apply(drone ddr.Drone, state DroneState) {
	switch state.operation {
	case Hover:
		drone.Hover()
		drone.CeaseRotation()
	case Land:
		drone.Land()
	case Takeoff:
		drone.TakeOff()
	case NOOP:
		// do nothing
	default:
		operation := operations[state.operation].fn
		operation(drone, state.speed)
	}
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
