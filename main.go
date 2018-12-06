package main

import (
    "deus.solita.fi/Solita/projects/drone_code_camp/repositories/git/ddr.git"
    "fmt"
    "gobot.io/x/gobot/platforms/keyboard"
    "gocv.io/x/gocv"
    "image"
    "image/color"
    "strconv"
)

func main() {

    window := gocv.NewWindow("Drone")
    drone := ddr.NewDrone(ddr.DroneReal, "../examples/drone-camera-calibration-400.yaml")
    //drone := ddr.NewDrone(ddr.DroneFake, "../examples/camera-calibration.yaml")
    err := drone.Init()

    if err != nil {
        fmt.Printf("error while initializing drone: %v\n", err)
        return
    }

    state := LANDED
    flightStatus := "Initialized"
    currentKey := -1
    speed := 10;

    fly := func(key int, operation func(d ddr.Drone, val int) error, status string) {
        if flying(state) {
            if currentKey == key {
                speed = min(2*speed, 100)
            } else {
                speed = 10
            }
            currentKey = key;
            operation(drone, speed)
            state = MOVING
            flightStatus = status + " " + strconv.Itoa(speed)
        }
    }

    for {
        frame := <- drone.VideoStream()

        gocv.PutText(&frame, flightStatus,
            image.Pt(50,50),
            gocv.FontHersheyComplex, 0.8, color.RGBA{255, 255, 255, 0}, 2)

        window.IMShow(frame)
        key := window.WaitKey(1)

        //fmt.Printf(" %d ", key)

        switch key {
        case keyboard.Spacebar: // space
            switch state {
            case MOVING:
                drone.Hover()
                drone.CeaseRotation()
                speed = 10
                state = HOVERING;
                flightStatus = "Hovering"
            case HOVERING:
                drone.Land()
                state = LANDED
                flightStatus = "Land"
            case LANDED:
                drone.TakeOff()
                flightStatus = "Take off"
                state = MOVING;
            }
        case keyboard.W:
            fly(key, ddr.Drone.Forward, "Going forward")
        case keyboard.S:
            fly(key, ddr.Drone.Backward, "Going backward")
        case ArrowRight:
            fly(key, ddr.Drone.Right, "Going right")
        case ArrowLeft:
            fly(key, ddr.Drone.Left, "Going left")
        case keyboard.D:
            fly(key, ddr.Drone.Clockwise, "Rotating clockwise")
        case keyboard.A:
            fly(key, ddr.Drone.CounterClockwise, "Rotating counter clockwise")
        case ArrowUp:
            fly(key, ddr.Drone.Up, "Going up")
        case ArrowDown:
            fly(key, ddr.Drone.Down, "Going down")
        }
    }
}

func min(a, b int) int {
    if a <= b {
        return a
    }
    return b
}

const (
    ArrowLeft = iota + 81
    ArrowUp
    ArrowRight
    ArrowDown
)

func flying(state FlyMode) bool {
    return state == HOVERING || state == MOVING
}

type FlyMode int
const (
    LANDED FlyMode = iota
    HOVERING
    MOVING
)