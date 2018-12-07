package main

import (
	"deus.solita.fi/Solita/projects/drone_code_camp/repositories/git/ddr.git"
	"fmt"
	"gobot.io/x/gobot/platforms/keyboard"
	"gocv.io/x/gocv"
	"image"
	"image/color"
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

func main() {

	window := gocv.NewWindow("Drone")
	//drone := ddr.NewDrone(ddr.DroneReal, "drone-camera-calibration-400.yaml")
	drone := ddr.NewDrone(ddr.DroneFake, "camera-calibration.yaml")

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
		default:
			operation, validKey := keymap[key]
			if validKey {
				state = fly(state, operation)
				apply(drone, state)
			} else {

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

func aiFly(state DroneState) {

}
