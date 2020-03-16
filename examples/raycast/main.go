package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/tracer8086/cirno"
	"golang.org/x/image/colornames"
)

const (
	width            = 1280
	height           = 720
	turnSpeed        = 250
	rayVisibleLength = 30
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
	circle := cirno.NewCircle(cirno.NewVector(140, 472.5), 67.5)
	rect := cirno.NewRectangle(cirno.NewVector(150, 78.75), 220, 112.5, 0)
	line := cirno.NewLine(cirno.NewVector(480, 540), cirno.NewVector(660, 405))
	cube := cirno.NewRectangle(cirno.NewVector(600, 112.5), 135, 135, 0)
	rhombus := cirno.NewRectangle(cirno.NewVector(360, 292.5), 90, 90, 45)
	littleCircle := cirno.NewCircle(cirno.NewVector(640, 540), 45)

	// Create a new space.
	space, err := cirno.NewSpace(1, 10, width*2, height*2,
		cirno.Zero, cirno.NewVector(width, height), false)
	handleError(err)
	// Fill the space with the shapes.
	err = space.Add(circle, rect, line, cube, rhombus, littleCircle)
	handleError(err)

	// Set default color for the shapes.
	for shape := range space.Shapes() {
		shape.SetData(colornames.Blue)
	}

	// Raycast parameters.
	origin := rhombus.Center()
	direction := cirno.NewVector(0, 1)

	// Setup metrics.
	last := time.Now()
	fps := 0
	perSecond := time.Tick(time.Second)

	// IMDraw to render shapes.
	imd := imdraw.New(nil)

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
			direction = direction.Rotate(turn * turnSpeed * deltaTime)
		}

		// If the shape is hit by the ray, its color is set to red.
		hitShape := space.Raycast(origin, direction, 0, 0)

		if hitShape == littleCircle {
			fmt.Println("LOL THE FUCK")
			hitShape = space.Raycast(origin, direction, 0, 0)
		}

		if hitShape != nil {
			hitShape.SetData(colornames.Red)
		}

		ray := cirno.NewLine(origin, origin.Add(direction.MultiplyByScalar(rayVisibleLength)))

		// Rendering.
		imd.Clear()

		for shape := range space.Shapes() {
			imd.Color = shape.Data().(color.RGBA)

			switch shape.(type) {
			case *cirno.Line:
				lineShape := shape.(*cirno.Line)

				imd.Push(
					pixel.V(lineShape.P().X, lineShape.P().Y),
					pixel.V(lineShape.Q().X, lineShape.Q().Y),
				)
				imd.Line(2)

			case *cirno.Circle:
				circleShape := shape.(*cirno.Circle)

				imd.Push(pixel.V(circleShape.Center().X,
					circleShape.Center().Y))
				imd.Circle(circleShape.Radius(), 0)

			case *cirno.Rectangle:
				rectShape := shape.(*cirno.Rectangle)
				vertices := rectShape.Vertices()

				imd.Push(
					pixel.V(vertices[0].X, vertices[0].Y),
					pixel.V(vertices[1].X, vertices[1].Y),
					pixel.V(vertices[2].X, vertices[2].Y),
					pixel.V(vertices[3].X, vertices[3].Y),
				)
				imd.Polygon(0)
			}
		}

		// Draw the ray.
		imd.Color = colornames.Green
		imd.Push(
			pixel.V(ray.P().X, ray.P().Y),
			pixel.V(ray.Q().X, ray.Q().Y),
		)
		imd.Line(2)

		// Restore hit shape color.
		if hitShape != nil {
			hitShape.SetData(colornames.Blue)
		}

		imd.Draw(win)

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
