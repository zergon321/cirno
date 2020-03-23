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
	"github.com/zergon321/cirno"
)

const (
	width         = 1280
	height        = 720
	platformID    = 1
	beholderID    = 1 << 1
	beholderEyeID = 1 << 3
)

var (
	vsync bool
)

type platform struct {
	rect      *cirno.Rectangle
	sprite    *pixel.Sprite
	transform pixel.Matrix
	batch     *pixel.Batch
}

type beholder struct {
	rect      *cirno.Rectangle
	hitCircle *cirno.Circle
	sprite    *pixel.Sprite
	transform pixel.Matrix
	batch     *pixel.Batch
	dead      bool
}

type bullet struct {
	hitShape  *cirno.Shape
	sprite    *pixel.Sprite
	transform pixel.Matrix
	batch     *pixel.Batch
	direction cirno.Vector
}

type player struct {
	rect      *cirno.Rectangle
	sprite    *pixel.Sprite
	transform pixel.Matrix
	dead      bool
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
	beholderPic, err := loadPicture("beholders.png")
	handleError(err)
	testmanPic, err := loadPicture("testmen.png")
	handleError(err)

	// Create sprites and batches.
	wallSprite := pixel.NewSprite(wallPic, pixel.R(0, 0, width, height))
	testmanLeftSprite := pixel.NewSprite(testmanPic, pixel.R(0, 0, 32, 64))
	testmanRightSprite := pixel.NewSprite(testmanPic, pixel.R(32, 0, 64, 64))
	electroBulletSprite := pixel.NewSprite(projectileSheet, pixel.R(0, 0, 64, 64))
	bloodBulletSprite := pixel.NewSprite(projectileSheet, pixel.R(64, 0, 192, 64))
	platformSprite := pixel.NewSprite(platformPic, pixel.R(0, 0, 128, 32))
	beholderLeftSprite := pixel.NewSprite(beholderPic, pixel.R(0, 0, 129, 315))
	beholderRightSprite := pixel.NewSprite(beholderPic, pixel.R(129, 0, 258, 315))

	platformBatch := pixel.NewBatch(new(pixel.TrianglesData), platformPic)
	bulletBatch := pixel.NewBatch(new(pixel.TrianglesData), projectileSheet)
	beholderBatch := pixel.NewBatch(new(pixel.TrianglesData), beholderPic)

	// Create a new collision space.
	space := cirno.NewSpace(5, 20, width*2, height*2,
		cirno.Zero, cirno.NewVector(width, height), true)

	// Create platforms.
	lowerPlatform := platform{
		rect:      cirno.NewRectangle(cirno.NewVector(640, 180), 128, 32, 0),
		sprite:    platformSprite,
		transform: pixel.IM.Moved(pixel.V(640, 180)),
		batch:     platformBatch,
	}
	middlePlatform := platform{
		rect:      cirno.NewRectangle(cirno.NewVector(320, 360), 128, 32, 0),
		sprite:    platformSprite,
		transform: pixel.IM.Moved(pixel.V(320, 360)),
		batch:     platformBatch,
	}
	higherPlatform := platform{
		rect:      cirno.NewRectangle(cirno.NewVector(960, 540), 128, 32, 0),
		sprite:    platformSprite,
		transform: pixel.IM.Moved(pixel.V(960, 540)),
		batch:     platformBatch,
	}

	lowerPlatform.rect.SetIdentity(platformID)
	middlePlatform.rect.SetIdentity(platformID)
	higherPlatform.rect.SetIdentity(platformID)

	// Create beholders.
	lowerBeholder := beholder{
		rect:      cirno.NewRectangle(cirno.NewVector(320, 534), 129, 315, 0),
		hitCircle: cirno.NewCircle(cirno.NewVector(352, 628), 32),
		sprite:    beholderLeftSprite,
		transform: pixel.IM.Moved(pixel.V(320, 534)),
		batch:     beholderBatch,
		dead:      false,
	}
	higherBeholder := beholder{
		rect:      cirno.NewRectangle(cirno.NewVector(960, 714), 129, 315, 0),
		hitCircle: cirno.NewCircle(cirno.NewVector(1055, 808), 32),
		sprite:    beholderRightSprite,
		transform: pixel.IM.Moved(pixel.V(320, 534)),
		batch:     beholderBatch,
		dead:      false,
	}

	lowerBeholder.rect.SetIdentity(beholderID)
	lowerBeholder.hitCircle.SetIdentity(beholderEyeID)
	higherBeholder.rect.SetIdentity(beholderID)
	higherBeholder.hitCircle.SetIdentity(beholderEyeID)

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
