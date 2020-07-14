package main

import (
	"flag"
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

var (
	controlledShape string
	vsync           bool
)

func cirnoToPixel(vec cirno.Vector) pixel.Vec {
	return pixel.V(vec.X, vec.Y)
}

func parseFlags() {
	flag.StringVar(&controlledShape, "shape", "rectangle",
		"The shape controlled during execution of the demo.")
	flag.BoolVar(&vsync, "vsync", true, "Enable vertical synchronization.")

	flag.Parse()
}

func run() {
	// Create a new window.
	cfg := pixelgl.WindowConfig{
		Title:  "Cirno demo",
		Bounds: pixel.R(0, 0, width, height),
	}
	win, err := pixelgl.NewWindow(cfg)
	handleError(err)

	base := cirno.NewVector(width/2, height/2)
	var object cirno.Shape

	switch controlledShape {
	case "line":
		object = cirno.NewLine(
			base.Add(cirno.NewVector(40, 0)),
			base.Add(cirno.NewVector(90, 50)),
		)

	case "rectangle":
		object = cirno.NewRectangle(
			base.Add(cirno.NewVector(60, 0)),
			40, 20, 0.0)

	case "circle":
		object = cirno.NewCircle(
			base.Add(cirno.NewVector(60, 0)), 20)

	default:
		handleError(fmt.Errorf("unknown argument: %s",
			controlledShape))
	}

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
		object.RotateAround(angle, base)
		object.Rotate(angle)

		// Rendering.
		imd.Clear()

		imd.Color = colors.Red
		imd.Push(cirnoToPixel(base))
		imd.Circle(20, 3)

		imd.Color = colors.Blue

		switch object.(type) {
		case *cirno.Line:
			lineShape := object.(*cirno.Line)

			imd.Push(
				pixel.V(lineShape.P().X, lineShape.P().Y),
				pixel.V(lineShape.Q().X, lineShape.Q().Y),
			)
			imd.Line(3)

		case *cirno.Circle:
			circleShape := object.(*cirno.Circle)

			imd.Push(pixel.V(circleShape.Center().X,
				circleShape.Center().Y))
			imd.Circle(circleShape.Radius(), 3)

		case *cirno.Rectangle:
			rectShape := object.(*cirno.Rectangle)
			vertices := rectShape.Vertices()

			imd.Push(
				pixel.V(vertices[0].X, vertices[0].Y),
				pixel.V(vertices[1].X, vertices[1].Y),
				pixel.V(vertices[2].X, vertices[2].Y),
				pixel.V(vertices[3].X, vertices[3].Y),
			)
			imd.Polygon(3)
		}

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
	parseFlags()
	pixelgl.Run(run)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
