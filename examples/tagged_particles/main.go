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
		cirno.Zero(), cirno.NewVector(width, height), true)
	handleError(err)

	particles := make([]*cirno.Circle, 0)

	// Create borders.
	borderWest, err := cirno.NewLine(cirno.NewVector(0, 0), cirno.NewVector(0, height))
	handleError(err)
	borderSouth, err := cirno.NewLine(cirno.NewVector(0, 0), cirno.NewVector(width, 0))
	handleError(err)
	borderNorth, err := cirno.NewLine(cirno.NewVector(0, height), cirno.NewVector(width, height))
	handleError(err)
	borderEast, err := cirno.NewLine(cirno.NewVector(width, 0), cirno.NewVector(width, height))
	handleError(err)

	// Borders shouldn't collide each other.
	borderWest.SetIdentity(1)
	borderSouth.SetIdentity(1)
	borderNorth.SetIdentity(1)
	borderEast.SetIdentity(1)

	// But they should collide anything else.
	borderWest.SetMask(^1)
	borderSouth.SetMask(^1)
	borderNorth.SetMask(^1)
	borderEast.SetMask(^1)

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
			particle, err := cirno.NewCircle(position, float64(rand.Intn(3)+5))
			handleError(err)

			if j%2 == 0 {
				particle.SetIdentity(2)
				particle.SetMask(3)
			} else {
				particle.SetIdentity(4)
				particle.SetMask(5)
			}

			particles = append(particles, particle)
			err = space.Add(particle)
			handleError(err)
		}
	}

	// Setup graphics.
	pic, err := loadPicture("particles.png")
	handleError(err)

	batch := pixel.NewBatch(&pixel.TrianglesData{}, pic)
	sprites := make([]*pixel.Sprite, 0)

	for i := range particles {
		var sprite *pixel.Sprite

		if i%2 == 0 {
			sprite = pixel.NewSprite(pic, pixel.R(0, 0, 64, 64))
		} else {
			sprite = pixel.NewSprite(pic, pixel.R(64, 0, 128, 64))
		}

		sprites = append(sprites, sprite)
	}

	// Setup metrics.
	last := time.Now()
	fps := 0
	cps := 0
	perSecond := time.Tick(time.Second)

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
			scale := particle.Radius() * 2.0 / 64.0
			angle := math.Floor(rand.Float64() * 100)
			pos := particle.Center()
			transform := pixel.IM.
				Scaled(pixel.ZV, scale).
				Rotated(pixel.ZV, angle).
				Moved(pixel.V(pos.X, pos.Y))

			sprites[i].Draw(batch, transform)
		}

		batch.Draw(win)

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
