package main

import (
	"deus.solita.fi/Solita/projects/drone_code_camp/repositories/git/ddr.git"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
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

func isInTestMode() bool {
	return len(os.Args) > 1 && os.Args[1] == "test"
}

func getDrone() ddr.Drone {
	if isInTestMode() {
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
		0,
	}

	for {
		frame := <-drone.VideoStream()

		ddr.DrawControls(drone, &frame)

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
			ring, exists := rings[state.nextRingId]
			if exists {
				state = aiFly(state, ring, drone)
				apply(drone, state)
			} /*else {
				state = next(state, TurnRight, 10, "Seeking ring turning right")
				apply(drone, state)
			}*/
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
	// fmt.Println(pose.Position)
	// fmt.Println(pose.Rotation.Mul3x1(mgl32.Vec3{0.0, 0.0, 1.0}))

	x := pose.Rotation.Mul3x1(mgl32.Vec3{0.0, 0.0, 1.0}).X()
	fmt.Println(x)
	position := drone.CameraToDroneMatrix().Mul3x1(pose.Position)
	if x < -0.2 {
		return next(state, TurnRight, 10, "AI turn left")
	} else if x > 0.2 {
		return next(state, TurnLeft, 10, "AI turn right")
	} else {
		drone.CeaseRotation()
		if position.X() > 0.05 {
			return next(state, Right, 10, "AI go right")
		} else if position.X() < -0.05 {
			return next(state, Left, 10, "AI go left")
		} else if position.Y() > 0.05 {
			return next(state, Down, 10, "AI go down")
		} else if position.Y() < -0.05 {
			return next(state, Up, 10, "AI go up")
		} else {
			return next(state, Forward, 10, "GOAL")
		}
	}
}
