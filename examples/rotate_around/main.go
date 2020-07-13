package main

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/zergon321/cirno"
	colors "golang.org/x/image/colornames"
)

const (
	width      = 1280
	height     = 720
	angleSpeed = 50
)

func cirnoToPixel(vec cirno.Vector) pixel.Vec {
	return pixel.V(vec.X, vec.Y)
}

// TODO: make object a real shape; let the user
// choose shapes.

func run() {
	// Create a new window.
	cfg := pixelgl.WindowConfig{
		Title:  "Cirno demo",
		Bounds: pixel.R(0, 0, width, height),
	}
	win, err := pixelgl.NewWindow(cfg)
	handleError(err)

	base := cirno.NewVector(width/2, height/2)
	object := base.Add(cirno.NewVector(80, 80))

	// IMDraw instance to draw shapes.
	imd := imdraw.New(nil)

	// Setup metrics.
	last := time.Now()
	fps := 0
	perSecond := time.Tick(time.Second)

	for !win.Closed() {
		deltaTime := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(colors.White)

		// Move the shape around the other shape.
		angle := angleSpeed * deltaTime
		object = object.RotateAround(angle, base)

		// Rendering.
		imd.Clear()

		imd.Color = colors.Red
		imd.Push(cirnoToPixel(base))
		imd.Circle(20, 3)

		imd.Color = colors.Blue
		imd.Push(cirnoToPixel(object))
		imd.Circle(20, 3)

		imd.Draw(win)

		win.Update()

		// Compute the FPS.
		fps++

		select {
		case <-perSecond:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d",
				cfg.Title, fps))
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
