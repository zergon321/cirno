package main

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/tracer8086/cirno"
	"golang.org/x/image/colornames"
)

const (
	width  = 1280
	height = 720
)

func run() {
	// Create a new window.
	cfg := pixelgl.WindowConfig{
		Title:  "Cirno demo",
		Bounds: pixel.R(0, 0, width, height),
		//VSync:  true,
		//Undecorated: true,
		//Monitor: pixelgl.PrimaryMonitor(),
	}
	win, err := pixelgl.NewWindow(cfg)
	handleError(err)

	// Create new shapes.
	circle := cirno.NewCircle(cirno.NewVector(7, 21), 3)
	rect := cirno.NewRectangle(cirno.NewVector(7.5, 3.5), 11, 5, 0)
	line := cirno.NewLine(cirno.NewVector(24, 24), cirno.NewVector(33, 18))
	cube := cirno.NewRectangle(cirno.NewVector(30, 5), 6, 6, 0)
	rhombus := cirno.NewRectangle(cirno.NewVector(18, 13), 4, 4, 45)
	littleCircle := cirno.NewCircle(cirno.NewVector(32, 24), 2)

	// Create a new space.
	space, err := cirno.NewSpace(1, 10, 64, 64,
		cirno.NewVector(0, 0), cirno.NewVector(64, 64), false)
	handleError(err)
	// Fill the space with the shapes.
	err = space.Add(circle, rect, line, cube, rhombus, littleCircle)
	handleError(err)

	// Raycast parameters.
	origin := rhombus.Center()
	direction := cirno.NewVector(0, 1)

	// A shape hit by raycast.
	var hitShape cirno.Shape

	// Setup metrics.
	last := time.Now()
	fps := 0
	perSecond := time.Tick(time.Second)

	for !win.Closed() {
		deltaTime := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(colornames.White)

		// Raycast angle turn.
		turn := 0.0

		// Reading inputs.
		if win.Pressed(pixelgl.KeyLeft) {
			turn++
		}

		if win.Pressed(pixelgl.KeyRight) {
			turn--
		}

		if turn != 0 {
			direction.Rotate(turn * deltaTime)
		}

		hitShape = space.Raycast(origin, direction, 0, 0)

		win.Update()

		fps++

		select {
		case <-perSecond:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, fps))
			fps = 0

		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
