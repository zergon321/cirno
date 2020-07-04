package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/zergon321/cirno"
	colors "golang.org/x/image/colornames"
)

const (
	width          = 1280
	height         = 720
	numberOfShapes = 1000
	speed          = 200
)

func run() {
	// Create a new window.
	cfg := pixelgl.WindowConfig{
		Title:  "Cirno benchmark",
		Bounds: pixel.R(0, 0, width, height),
	}
	win, err := pixelgl.NewWindow(cfg)
	handleError(err)

	// Create a new collision space.
	space, err := cirno.NewSpace(50, 20, width*2, height*2,
		cirno.Zero, cirno.NewVector(width, height), false)
	handleError(err)
	// To prevent flickering and preserve the rendering
	// order shapes.
	shapeList := []*cirno.Rectangle{}

	// Fill the space with random rectangles.
	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < numberOfShapes; i++ {
		rectangle := cirno.NewRectangle(
			cirno.NewVector(float64(rand.Int31n(width)), float64(rand.Int31n(height))),
			float64(rand.Int31n(20)+20), float64(rand.Int31n(10)+10), float64(rand.Int31n(180)))
		rectangle.SetData(color.RGBA{
			R: uint8(rand.Int31n(256)),
			G: uint8(rand.Int31n(256)),
			B: uint8(rand.Int31n(256)),
			A: uint8(rand.Int31n(256)),
		})

		err = space.Add(rectangle)
		handleError(err)
		shapeList = append(shapeList, rectangle)
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

		// Move and rotate all the shapes randomly.
		for _, shape := range shapeList {
			angle := math.Floor(rand.Float64() * 1000)
			direction := cirno.NewVector(1, 0).Rotate(angle)
			movement := direction.MultiplyByScalar(speed * deltaTime)
			rotation := math.Floor(rand.Float64() * 1000)

			shape.Move(movement)
			shape.Rotate(rotation)
			space.AdjustShapePosition(shape)

			_, err = space.Update(shape)
			handleError(err)
		}

		// Find colliding shapes.
		collidingShapes, err := space.CollidingShapes()
		handleError(err)

		imd.Clear()

		// Render all the shapes.
		for _, shape := range shapeList {
			imd.Color = shape.Data().(color.RGBA)
			vertices := shape.Vertices()

			imd.Push(
				pixel.V(vertices[0].X, vertices[0].Y),
				pixel.V(vertices[1].X, vertices[1].Y),
				pixel.V(vertices[2].X, vertices[2].Y),
				pixel.V(vertices[3].X, vertices[3].Y),
			)

			imd.Polygon(0)
		}

		// Render all the space cells.
		for cell := range space.Cells() {
			imd.Color = colors.Green
			vertices := cell.Vertices()

			imd.Push(
				pixel.V(vertices[0].X, vertices[0].Y),
				pixel.V(vertices[1].X, vertices[1].Y),
				pixel.V(vertices[2].X, vertices[2].Y),
				pixel.V(vertices[3].X, vertices[3].Y),
			)

			imd.Polygon(2)
		}

		imd.Draw(win)

		win.Update()

		// Compute the FPS and the number of shapes colliding.
		fps++

		select {
		case <-perSecond:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d | Colliding shapes: %d",
				cfg.Title, fps, len(collidingShapes)))
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
