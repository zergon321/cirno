package main

import (
	"fmt"
	"image"
	_ "image/png"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/zergon321/cirno"
	"golang.org/x/image/colornames"
)

const (
	numberOfParticles = 256
	width             = 1280
	height            = 720
	speed             = 200
	intensity         = 10
)

func loadPicture(pic string) (pixel.Picture, error) {
	file, err := os.Open(pic)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	return pixel.PictureDataFromImage(img), nil
}

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

	// Setup physics.
	rand.Seed(time.Now().UnixNano())

	space, err := cirno.NewSpace(5, 20, width*2, height*2,
		cirno.Zero, cirno.NewVector(width, height), false)
	handleError(err)

	particles := make([]*cirno.Circle, 0)

	// Create borders.
	borderWest := cirno.NewLine(cirno.NewVector(0, 0), cirno.NewVector(0, height))
	borderSouth := cirno.NewLine(cirno.NewVector(0, 0), cirno.NewVector(width, 0))
	borderNorth := cirno.NewLine(cirno.NewVector(0, height), cirno.NewVector(width, height))
	borderEast := cirno.NewLine(cirno.NewVector(width, 0), cirno.NewVector(width, height))

	// Add the borders in the space.
	err = space.Add(borderNorth)
	handleError(err)
	err = space.Add(borderWest)
	handleError(err)
	err = space.Add(borderSouth)
	handleError(err)
	err = space.Add(borderEast)
	handleError(err)

	// Create particles.
	for i := 0; i < numberOfParticles/16; i++ {
		for j := 0; j < numberOfParticles/16; j++ {
			iDist := width / (numberOfParticles / 16)
			jDist := height / (numberOfParticles / 16)

			position := cirno.NewVector(float64((i+1)*iDist), float64((j+1)*jDist)).MultiplyByScalar(0.8)
			particle := cirno.NewCircle(position, float64(rand.Intn(3)+5))

			particles = append(particles, particle)
			err := space.Add(particle)
			handleError(err)
		}
	}

	// Setup graphics.
	pic, err := loadPicture("round_particle.png")
	handleError(err)

	batch := pixel.NewBatch(&pixel.TrianglesData{}, pic)
	sprites := make([]*pixel.Sprite, 0)

	for range particles {
		sprite := pixel.NewSprite(pic, pixel.R(0, 0, 45, 45))
		sprites = append(sprites, sprite)
	}

	// Setup metrics.
	last := time.Now()
	fps := 0
	cps := 0
	perSecond := time.Tick(time.Second)

	// Setup wireframe drawer.
	imd := imdraw.New(nil)

	for !win.Closed() {
		deltaTime := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(colornames.White)

		// Movement.
		for _, particle := range particles {
			angle := math.Floor(rand.Float64() * 1000)
			direction := cirno.NewVector(1, 0).Rotate(angle)
			movement := direction.MultiplyByScalar(speed * deltaTime)

			/*shapes, err := space.WouldBeCollidedBy(particle, movement, 0.0)
			handleError(err)

			// If a collision occurred, the particle won't
			// go the full way.
			pos := particle.Center()

			if len(shapes) > 0 {
				pos, err = cirno.Approximate(particle, movement, shapes, intensity)
				handleError(err)
				movement = pos.Subtract(particle.Center())
			}*/

			particle.Move(movement)
			space.AdjustShapePosition(particle)
			_, err = space.Update(particle)
			handleError(err)

			shapes, err := space.CollidingWith(particle)
			handleError(err)

			if len(shapes) > 0 {
				cps += len(shapes)
				particle.Move(movement.MultiplyByScalar(-1))

				pos, _, _, err := cirno.Approximate(particle, movement, 0,
					shapes, intensity, false)
				handleError(err)
				movement = pos.Subtract(particle.Center())
				particle.Move(movement)
				space.AdjustShapePosition(particle)
				_, err = space.Update(particle)
				handleError(err)
			}
		}

		// Rendering.
		batch.Clear()

		for i, particle := range particles {
			scale := particle.Radius() * 2.0 / 45.0
			angle := math.Floor(rand.Float64() * 100)
			pos := particle.Center()
			transform := pixel.IM.
				Scaled(pixel.ZV, scale).
				Rotated(pixel.ZV, angle).
				Moved(pixel.V(pos.X, pos.Y))

			sprites[i].Draw(batch, transform)
		}

		imd.Clear()

		for cell, cellShapes := range space.Cells() {
			imd.Color = colornames.Green
			vertices := cell.Vertices()

			imd.Push(
				pixel.V(vertices[0].X, vertices[0].Y),
				pixel.V(vertices[1].X, vertices[1].Y),
				pixel.V(vertices[2].X, vertices[2].Y),
				pixel.V(vertices[3].X, vertices[3].Y),
			)

			imd.Polygon(2)

			imd.Push(pixel.V(cell.Center().X, cell.Center().Y))

			imd.Circle(float64(len(cellShapes)), 2)
		}

		batch.Draw(win)
		imd.Draw(win)

		win.Update()

		// Show FPS in the window title.
		fps++

		select {
		case <-perSecond:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d, CPS: %d, CPF: %f", cfg.Title, fps, cps, float64(cps)/float64(fps)))
			fps = 0
			cps = 0

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
