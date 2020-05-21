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
	width            = 1280
	height           = 720
	moveSpeed        = 400
	rayVisibleLength = 30
	intensity        = 500
)

func cirnoToPixel(vector cirno.Vector) pixel.Vec {
	return pixel.V(vector.X, vector.Y)
}

func run() {
	// Create a new window.
	cfg := pixelgl.WindowConfig{
		Title:  "Cirno demo",
		Bounds: pixel.R(0, 0, width, height),
		VSync:  true,
		//Undecorated: true,
		//Monitor: pixelgl.PrimaryMonitor(),
	}
	win, err := pixelgl.NewWindow(cfg)
	handleError(err)

	// Setup physics.
	circleBig := cirno.NewCircle(cirno.NewVector(350, 250), 50)
	circleLittle := cirno.NewCircle(cirno.NewVector(1000, 600), 20)
	space, err := cirno.NewSpace(1, 10, width*2, height*2,
		cirno.Zero, cirno.NewVector(width, height), false)
	handleError(err)
	err = space.Add(circleBig, circleLittle)
	handleError(err)

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

		if movement != cirno.Zero {
			movement = movement.MultiplyByScalar(moveSpeed * deltaTime)

			shapes, err := space.WouldBeCollidedBy(circleLittle, movement, 0.0)
			handleError(err)

			// If a collision occurres, the shape
			// will slide.
			pos := circleLittle.Center()

			if len(shapes) > 0 {
				pos, _, err = cirno.Approximate(circleLittle, movement, 0.0,
					shapes, intensity, false)
				handleError(err)

				// If there's no opportunity to approximate,
				// do sliding.
				if circleLittle.Center().Subtract(pos).Magnitude() < cirno.Epsilon {
					normal := circleBig.NormalToCircle(circleLittle)
					movement = movement.Subtract(normal.
						MultiplyByScalar(cirno.Dot(movement, normal)))
				} else {
					movement = pos.Subtract(circleLittle.Center())
				}
			}

			circleLittle.Move(movement)
			space.AdjustShapePosition(circleLittle)
			_, err = space.Update(circleLittle)
			handleError(err)
		}

		// Rendering.
		imd.Clear()

		imd.Color = colors.Red
		imd.Push(cirnoToPixel(circleBig.Center()))
		imd.Circle(circleBig.Radius(), 0)

		imd.Color = colors.Blue
		imd.Push(cirnoToPixel(circleLittle.Center()))
		imd.Circle(circleLittle.Radius(), 0)

		if movement != cirno.Zero {
			imd.Color = colors.Brown
			imd.Push(cirnoToPixel(circleLittle.Center()))
			imd.Push(cirnoToPixel(circleLittle.Center().
				Add(movement.MultiplyByScalar(8))))
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
	pixelgl.Run(run)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
