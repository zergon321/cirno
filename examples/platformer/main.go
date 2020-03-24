package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/zergon321/cirno"
	colors "golang.org/x/image/colornames"
)

const (
	width           = 1280
	height          = 720
	platformID      = 1
	beholderID      = 1 << 1
	beholderEyeID   = 1 << 2
	playerID        = 1 << 3
	electroBulletID = 1 << 4
)

var (
	vsync     bool
	drawWires bool
)

type platform struct {
	rect      *cirno.Rectangle
	sprite    *pixel.Sprite
	transform pixel.Matrix
}

func (pl *platform) draw(target pixel.Target) {
	pl.sprite.Draw(target, pl.transform)
}

type beholder struct {
	rect      *cirno.Rectangle
	hitCircle *cirno.Circle
	sprite    *pixel.Sprite
	transform pixel.Matrix
	dead      bool
}

func (br *beholder) draw(target pixel.Target) {
	br.sprite.Draw(target, br.transform)
}

type bullet struct {
	hitShape  *cirno.Shape
	sprite    *pixel.Sprite
	transform pixel.Matrix
	direction cirno.Vector
}

func (b *bullet) draw(target pixel.Target) {
	b.sprite.Draw(target, b.transform)
}

type player struct {
	rect      *cirno.Rectangle
	sprite    *pixel.Sprite
	transform pixel.Matrix
	dead      bool
}

func (p *player) draw(target pixel.Target) {
	p.sprite.Draw(target, p.transform)
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

func drawShape(imd *imdraw.IMDraw, shape cirno.Shape) {
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
		imd.Circle(circleShape.Radius(), 2)

	case *cirno.Rectangle:
		rectShape := shape.(*cirno.Rectangle)
		vertices := rectShape.Vertices()

		imd.Push(
			pixel.V(vertices[0].X, vertices[0].Y),
			pixel.V(vertices[1].X, vertices[1].Y),
			pixel.V(vertices[2].X, vertices[2].Y),
			pixel.V(vertices[3].X, vertices[3].Y),
		)
		imd.Polygon(2)
	}
}

func parseFlags() {
	flag.BoolVar(&vsync, "vsync", false, "Enable vertical synchronization.")
	flag.BoolVar(&drawWires, "draw-wires", false, "Enable hit shape drawing.")

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
	//projectileSheet, err := loadPicture("projectiles.png")
	//handleError(err)
	beholderPic, err := loadPicture("beholders.png")
	handleError(err)
	testmanPic, err := loadPicture("testmen.png")
	handleError(err)

	// Create sprites and batches.
	wallSprite := pixel.NewSprite(wallPic, pixel.R(0, 0, width, height))
	testmanLeftSprite := pixel.NewSprite(testmanPic, pixel.R(0, 0, 32, 64))
	//testmanRightSprite := pixel.NewSprite(testmanPic, pixel.R(32, 0, 64, 64))
	//electroBulletSprite := pixel.NewSprite(projectileSheet, pixel.R(0, 0, 64, 64))
	//bloodBulletSprite := pixel.NewSprite(projectileSheet, pixel.R(64, 0, 192, 64))
	platformSprite := pixel.NewSprite(platformPic, pixel.R(0, 0, 128, 32))
	beholderLeftSprite := pixel.NewSprite(beholderPic, pixel.R(0, 0, 129, 315))
	beholderRightSprite := pixel.NewSprite(beholderPic, pixel.R(129, 0, 258, 315))

	platformBatch := pixel.NewBatch(new(pixel.TrianglesData), platformPic)
	//bulletBatch := pixel.NewBatch(new(pixel.TrianglesData), projectileSheet)
	beholderBatch := pixel.NewBatch(new(pixel.TrianglesData), beholderPic)

	// Create platforms.
	lowerPlatform := &platform{
		rect:      cirno.NewRectangle(cirno.NewVector(640, 180), 128, 32, 0),
		sprite:    platformSprite,
		transform: pixel.IM.Moved(pixel.V(640, 180)),
	}
	middlePlatform := &platform{
		rect:      cirno.NewRectangle(cirno.NewVector(320, 360), 128, 32, 0),
		sprite:    platformSprite,
		transform: pixel.IM.Moved(pixel.V(320, 360)),
	}
	higherPlatform := &platform{
		rect:      cirno.NewRectangle(cirno.NewVector(960, 540), 128, 32, 0),
		sprite:    platformSprite,
		transform: pixel.IM.Moved(pixel.V(960, 540)),
	}

	lowerPlatform.rect.SetIdentity(platformID)
	middlePlatform.rect.SetIdentity(platformID)
	higherPlatform.rect.SetIdentity(platformID)

	lowerPlatform.rect.SetData(lowerPlatform)
	middlePlatform.rect.SetData(middlePlatform)
	higherPlatform.rect.SetData(higherPlatform)

	// Create beholders.
	lowerBeholder := &beholder{
		rect:      cirno.NewRectangle(cirno.NewVector(320, 534), 129, 315, 0),
		hitCircle: cirno.NewCircle(cirno.NewVector(352, 628), 32),
		sprite:    beholderLeftSprite,
		transform: pixel.IM.Moved(pixel.V(320, 534)),
		dead:      false,
	}
	higherBeholder := &beholder{
		rect:      cirno.NewRectangle(cirno.NewVector(960, 714), 129, 315, 0),
		hitCircle: cirno.NewCircle(cirno.NewVector(1055, 808), 32),
		sprite:    beholderRightSprite,
		transform: pixel.IM.Moved(pixel.V(320, 534)),
		dead:      false,
	}

	lowerBeholder.rect.SetIdentity(beholderID)
	lowerBeholder.hitCircle.SetIdentity(beholderEyeID)
	higherBeholder.rect.SetIdentity(beholderID)
	higherBeholder.hitCircle.SetIdentity(beholderEyeID)

	lowerBeholder.rect.SetData(lowerBeholder)
	lowerBeholder.hitCircle.SetData(lowerBeholder)
	higherBeholder.rect.SetData(higherBeholder)
	higherBeholder.hitCircle.SetData(higherBeholder)

	// Create hero.
	hero := &player{
		rect:      cirno.NewRectangle(cirno.NewVector(640, 228), 32, 64, 0),
		sprite:    testmanLeftSprite,
		transform: pixel.IM.Moved(pixel.V(640, 228)),
		dead:      false,
	}

	hero.rect.SetIdentity(playerID)
	hero.rect.SetMask(platformID | electroBulletID)
	hero.rect.SetData(hero)

	// Create a new collision space.
	space, err := cirno.NewSpace(5, 20, width*2, height*2,
		cirno.Zero, cirno.NewVector(width, height), true)
	handleError(err)
	// Add hit shapes to the space.
	err = space.Add(lowerPlatform.rect, middlePlatform.rect, higherPlatform.rect,
		lowerBeholder.rect, higherBeholder.rect, lowerBeholder.hitCircle,
		higherBeholder.hitCircle, hero.rect)
	handleError(err)

	// Setup metrics.
	last := time.Now()
	fps := 0
	perSecond := time.Tick(time.Second)

	var imd *imdraw.IMDraw

	if drawWires {
		imd = imdraw.New(nil)
		imd.Color = colors.Green
	}

	for !win.Closed() {
		deltaTime := time.Since(last).Seconds()
		last = time.Now()

		fmt.Println("Delta:", deltaTime)

		wallSprite.Draw(win, pixel.IM.Moved(pixel.V(width/2, height/2)))

		// Draw platforms.
		lowerPlatform.draw(platformBatch)
		middlePlatform.draw(platformBatch)
		higherPlatform.draw(platformBatch)

		// Draw beholders.
		lowerBeholder.draw(beholderBatch)
		higherBeholder.draw(beholderBatch)

		// Draw batches.
		platformBatch.Draw(win)
		beholderBatch.Draw(win)

		// Draw hero.
		hero.draw(win)

		if drawWires {
			imd.Clear()

			for shape := range space.Shapes() {
				drawShape(imd, shape)
			}

			imd.Draw(win)
		}

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
