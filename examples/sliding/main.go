package main

import (
	"flag"
	"fmt"
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/zergon321/cirno"
	colors "golang.org/x/image/colornames"
)

const (
	width     = 1280
	height    = 720
	moveSpeed = 400
	intensity = 500
)

var (
	controlledShape string
	vsync           bool
)

func parseFlags() {
	flag.StringVar(&controlledShape, "shape", "line",
		"The shape controlled during execution of the demo.")
	flag.BoolVar(&vsync, "vsync", true, "Enable vertical synchronization.")

	flag.Parse()
}

func cirnoToPixel(vector cirno.Vector) pixel.Vec {
	return pixel.V(vector.X, vector.Y)
}

func run() {
	// Create a new window.
	cfg := pixelgl.WindowConfig{
		Title:  "Cirno demo",
		Bounds: pixel.R(0, 0, width, height),
		VSync:  vsync,
		//Undecorated: true,
		//Monitor: pixelgl.PrimaryMonitor(),
	}
	win, err := pixelgl.NewWindow(cfg)
	handleError(err)

	// Setup physics.
	circleBig := cirno.NewCircle(cirno.NewVector(350, 250), 50)
	circleBig.SetData(colors.Red)
	circleLittle := cirno.NewCircle(cirno.NewVector(1000, 600), 20)
	circleLittle.SetData(colors.Blue)
	circleTemp := cirno.NewCircle(cirno.NewVector(470, 250), 50)
	circleTemp.SetData(colors.Coral)
	line := cirno.NewLine(cirno.NewVector(800, 200), cirno.NewVector(1200, 400))
	line.SetData(colors.Green)
	lineCtrl := cirno.NewLine(cirno.NewVector(750, 500), cirno.NewVector(900, 600))
	lineCtrl.SetData(colors.Chocolate)

	space, err := cirno.NewSpace(1, 10, width*2, height*2,
		cirno.Zero, cirno.NewVector(width, height), false)
	handleError(err)
	err = space.Add(circleBig, circleLittle, line, circleTemp, lineCtrl)
	handleError(err)

	// Choose the shape to control.
	var ctrlShape cirno.Shape

	switch controlledShape {
	case "circle":
		ctrlShape = circleLittle

	case "line":
		ctrlShape = lineCtrl
	}

	// Setup metrics.
	last := time.Now()
	fps := 0
	perSecond := time.Tick(time.Second)

	// IMDraw to render shapes.
	imd := imdraw.New(nil)

	for !win.Closed() {
		deltaTime := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(colors.White)

		// Movement.
		movement := cirno.Zero

		if win.Pressed(pixelgl.KeyUp) {
			movement = movement.Add(cirno.Up)
		}

		if win.Pressed(pixelgl.KeyDown) {
			movement = movement.Add(cirno.Down)
		}

		if win.Pressed(pixelgl.KeyLeft) {
			movement = movement.Add(cirno.Left)
		}

		if win.Pressed(pixelgl.KeyRight) {
			movement = movement.Add(cirno.Right)
		}

		var foundShape cirno.Shape
		var normal cirno.Vector

		if movement != cirno.Zero {
			movement = movement.MultiplyByScalar(moveSpeed * deltaTime)

			shapes, err := space.WouldBeCollidedBy(ctrlShape, movement, 0.0)
			handleError(err)

			// If a collision occurres, the shape
			// will slide.
			pos := ctrlShape.Center()

			if len(shapes) > 0 {
				normal = cirno.Zero
				pos, _, foundShape, err = cirno.Approximate(ctrlShape, movement, 0.0,
					shapes, intensity, false)
				handleError(err)

				// If there's no opportunity to approximate,
				// do sliding.
				if ctrlShape.Center().Subtract(pos).Magnitude() < cirno.Epsilon {
					normal = foundShape.NormalTo(ctrlShape)
					movement = movement.Subtract(normal.
						MultiplyByScalar(cirno.Dot(movement, normal)))

					// Make sure the shape won't collide other shapes
					// while sliding.
					shapes.Remove(foundShape)
					pos, _, _, err = cirno.Approximate(ctrlShape, movement, 0.0,
						shapes, intensity, false)
				}

				movement = pos.Subtract(ctrlShape.Center())
			}

			ctrlShape.Move(movement)
			space.AdjustShapePosition(ctrlShape)
			_, err = space.Update(ctrlShape)
			handleError(err)
		}

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

		if movement != cirno.Zero {
			imd.Color = colors.Magenta
			imd.Push(cirnoToPixel(ctrlShape.Center()))
			imd.Push(cirnoToPixel(ctrlShape.Center().
				Add(movement.MultiplyByScalar(8))))
			imd.Line(2)
		}

		if normal != cirno.Zero {
			imd.Color = colors.Purple
			imd.Push(cirnoToPixel(foundShape.Center()))
			imd.Push(cirnoToPixel(foundShape.Center().
				Add(normal.MultiplyByScalar(16))))
			imd.Line(2)
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
	parseFlags()
	pixelgl.Run(run)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
