package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/tracer8086/cirno"
)

const (
	width  = 1280
	height = 720
)

var (
	vsync bool
)

type platform struct {
	rect      *cirno.Rectangle
	sprite    *pixel.Sprite
	transform pixel.Matrix
}

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

func parseFlags() {
	flag.BoolVar(&vsync, "vsync", false, "Enable vertical synchronization.")

	flag.Parse()
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

	// Load pictures.
	wallPic, err := loadPicture("wall.png")
	handleError(err)
	platformPic, err := loadPicture("platform.png")
	handleError(err)
	projectileSheet, err := loadPicture("projectiles.png")
	handleError(err)

	// Create sprites and batches.
	wallSprite := pixel.NewSprite(wallPic, pixel.R(0, 0, width, height))
	platformBatch := pixel.NewBatch(new(pixel.TrianglesData), platformPic)
	projectileBatch := pixel.NewBatch(new(pixel.TrianglesData), projectileSheet)

	// Setup metrics.
	last := time.Now()
	fps := 0
	perSecond := time.Tick(time.Second)

	for !win.Closed() {
		deltaTime := time.Since(last).Seconds()
		last = time.Now()

		fmt.Println("Delta:", deltaTime)

		wallSprite.Draw(win, pixel.IM.Moved(pixel.V(width/2, height/2)))

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
