package main

import (
	"deus.solita.fi/Solita/projects/drone_code_camp/repositories/git/ddr.git"
	"fmt"
	"gobot.io/x/gobot/platforms/keyboard"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"os"
)

const (
	ArrowLeft = iota + 81
	ArrowUp
	ArrowRight
	ArrowDown
)

var keymap = map[int]OperationId{
	keyboard.W: Forward,
	keyboard.S: Backward,
	keyboard.D: TurnRight,
	keyboard.A: TurnLeft,

	ArrowRight: Right,
	ArrowLeft:  Left,
	ArrowUp:    Up,
	ArrowDown:  Down,
}

func getDrone() ddr.Drone {
	if len(os.Args) > 1 && os.Args[1] == "test" {
		fmt.Println("Using fake drone.")
		return ddr.NewDrone(ddr.DroneFake, "camera-calibration.yaml")
	} else {
		fmt.Println("Using real drone.")
		return ddr.NewDrone(ddr.DroneReal, "drone-camera-calibration-400.yaml")
	}
}

func main() {
	window := gocv.NewWindow("Drone")

	drone := getDrone()
	err := drone.Init()

	if err != nil {
		fmt.Printf("error while initializing drone: %v\n", err)
		return
	}

	track := ddr.NewTrack()
	defer track.Close()

	state := DroneState{
		false,
		NOOP,
		MinSpeed,
		"Initialized",
	}

	for {
		frame := <-drone.VideoStream()

		gocv.PutText(&frame, state.message,
			image.Pt(50, 50),
			gocv.FontHersheyComplex, 0.8, color.RGBA{255, 255, 255, 0}, 2)

		markers := track.GetMarkers(&frame)
		rings := track.ExtractRings(markers)

		// detect markers in this frame
		displayRings(rings, frame, drone)

		window.IMShow(frame)
		key := window.WaitKey(1)

		switch key {
		case keyboard.Spacebar: // space
			state = toggleMode(state)
			apply(drone, state)
		case -1:
			ring, exists := rings[3] // TODO ring selection
			if exists {
				state = aiFly(state, ring, drone)
				apply(drone, state)
			}
		default:
			operation, validKey := keymap[key]
			if validKey {
				state = fly(state, operation)
				apply(drone, state)
			}
		}

		//fmt.Println(state)
	}
}

func displayRings(rings map[int]*ddr.Ring, frame gocv.Mat, drone ddr.Drone) {

	for _, ring := range rings {
		pose := ring.EstimatePose(drone)
		ring.Draw(&frame, pose, drone)
	}
}

func aiFly(state DroneState, ring *ddr.Ring, drone ddr.Drone) DroneState {
	pose := ring.EstimatePose(drone)
	fmt.Println(pose)
	return operation(state, NOOP)
}
