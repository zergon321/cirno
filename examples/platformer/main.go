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
	width     = 1280
	height    = 720
	intensity = 1000
	gravity   = -500
)

const (
	platformID      = 1
	beholderID      = 1 << 1
	beholderEyeID   = 1 << 2
	playerID        = 1 << 3
	electroBulletID = 1 << 4
	bloodBulletID   = 1 << 5
)

var (
	vsync     bool
	drawWires bool
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
	flag.BoolVar(&vsync, "vsync", true, "Enable vertical synchronization.")
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
	wallPic, err := loadPicture("wooden-wall.png")
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
	bloodBulletSprite := pixel.NewSprite(projectileSheet, pixel.R(64, 0, 256, 64))
	platformSprite := pixel.NewSprite(platformPic, pixel.R(0, 0, 128, 32))
	beholderLeftSprite := pixel.NewSprite(beholderPic, pixel.R(0, 0, 129, 315))
	beholderRightSprite := pixel.NewSprite(beholderPic, pixel.R(129, 0, 258, 315))

	//bulletBatch := pixel.NewBatch(new(pixel.TrianglesData), projectileSheet)

	// Create platforms.
	lowerRect, err := cirno.NewRectangle(cirno.NewVector(640, 40), 384, 32, 0)
	handleError(err)
	lowerPlatform := &platform{
		rect:   lowerRect,
		sprite: platformSprite,
		transform: pixel.IM.ScaledXY(pixel.ZV, pixel.V(3, 1)).
			Moved(pixel.V(640, 40)),
	}

	middleRect, err := cirno.NewRectangle(cirno.NewVector(320, 220), 384, 32, 0)
	handleError(err)
	middlePlatform := &platform{
		rect:   middleRect,
		sprite: platformSprite,
		transform: pixel.IM.ScaledXY(pixel.ZV, pixel.V(3, 1)).
			Moved(pixel.V(320, 220)),
	}

	higherRect, err := cirno.NewRectangle(cirno.NewVector(960, 400), 384, 32, 0)
	handleError(err)
	higherPlatform := &platform{
		rect:   higherRect,
		sprite: platformSprite,
		transform: pixel.IM.ScaledXY(pixel.ZV, pixel.V(3, 1)).
			Moved(pixel.V(960, 400)),
	}

	lowerPlatform.rect.SetIdentity(platformID)
	middlePlatform.rect.SetIdentity(platformID)
	higherPlatform.rect.SetIdentity(platformID)

	lowerPlatform.rect.SetData(lowerPlatform)
	middlePlatform.rect.SetData(middlePlatform)
	higherPlatform.rect.SetData(higherPlatform)

	lowerBeholderRect, err := cirno.NewRectangle(cirno.NewVector(320, 316), 64.5, 157.5, 0)
	handleError(err)
	lowerBeholderCircle, err := cirno.NewCircle(cirno.NewVector(304, 378.75), 16)
	handleError(err)
	// Create beholders.
	lowerBeholder := &beholder{
		rect:           lowerBeholderRect,
		hitCircle:      lowerBeholderCircle,
		sprite:         beholderLeftSprite,
		anim:           []*pixel.Sprite{beholderLeftSprite, beholderRightSprite},
		bulletSprite:   bloodBulletSprite,
		bulletSpeed:    400,
		bulletCooldown: 600 * time.Millisecond,
		bulletTimer:    nil,
		spawnedBullets: make([]*bloodBullet, 0),
		speed:          250,
		direction:      cirno.Left(),
		transform:      pixel.IM.Scaled(pixel.ZV, 0.5).Moved(pixel.V(320, 316)),
		dead:           false,
	}

	higherBeholderRect, err := cirno.NewRectangle(cirno.NewVector(960, 496), 64.5, 157.5, 0)
	handleError(err)
	higherBeholderCircle, err := cirno.NewCircle(cirno.NewVector(976, 558.75), 16)
	handleError(err)
	higherBeholder := &beholder{
		rect:           higherBeholderRect,
		hitCircle:      higherBeholderCircle,
		sprite:         beholderRightSprite,
		anim:           []*pixel.Sprite{beholderLeftSprite, beholderRightSprite},
		bulletSprite:   bloodBulletSprite,
		bulletSpeed:    300,
		bulletCooldown: 400 * time.Millisecond,
		bulletTimer:    nil,
		spawnedBullets: make([]*bloodBullet, 0),
		speed:          300,
		direction:      cirno.Right(),
		transform:      pixel.IM.Scaled(pixel.ZV, 0.5).Moved(pixel.V(960, 496)),
		dead:           false,
	}

	lowerBeholder.rect.SetIdentity(beholderID)
	lowerBeholder.hitCircle.SetIdentity(beholderEyeID)
	higherBeholder.rect.SetIdentity(beholderID)
	higherBeholder.hitCircle.SetIdentity(beholderEyeID)

	lowerBeholder.rect.SetMask(platformID)
	higherBeholder.rect.SetMask(platformID)

	lowerBeholder.rect.SetData(lowerBeholder)
	lowerBeholder.hitCircle.SetData(lowerBeholder)
	higherBeholder.rect.SetData(higherBeholder)
	higherBeholder.hitCircle.SetData(higherBeholder)

	// Create hero.
	playerHitbox, err := cirno.NewRectangle(cirno.NewVector(640, 121), 64, 128, 0)
	handleError(err)
	hero := &player{
		speed:            500,
		jumpAcceleration: 80,
		verticalSpeed:    gravity,
		terminalSpeed:    gravity,
		aim:              cirno.Left(),
		bulletSprite:     electroBulletSprite,
		bulletSpeed:      200,
		spawnedBullets:   make([]*electroBullet, 0),
		rect:             playerHitbox,
		sprite:           testmanLeftSprite,
		animation:        []*pixel.Sprite{testmanLeftSprite, testmanRightSprite},
		transform:        pixel.IM.Scaled(pixel.V(0, 0), 2).Moved(pixel.V(640, 121)),
		dead:             false,
	}

	hero.rect.SetIdentity(playerID)
	hero.rect.SetMask(platformID)
	hero.rect.SetData(hero)

	// Create a new collision space.
	space, err := cirno.NewSpace(5, 20, width*4, height*4,
		cirno.Zero(), cirno.NewVector(width, height), true)
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
		imd.Color = colors.Lightgreen
	}

	for !win.Closed() {
		deltaTime := time.Since(last).Seconds()
		last = time.Now()

		// Update beholders.
		if !lowerBeholder.dead {
			err = lowerBeholder.update(space, deltaTime)
			handleError(err)
		}

		if !higherBeholder.dead {
			err = higherBeholder.update(space, deltaTime)
			handleError(err)
		}

		// Update beholder bullets.
		for _, bullet := range lowerBeholder.spawnedBullets {
			err = bullet.update(space, deltaTime)
			handleError(err)
		}

		for _, bullet := range higherBeholder.spawnedBullets {
			err = bullet.update(space, deltaTime)
			handleError(err)
		}

		// Update hero.
		if !hero.dead {
			err = hero.update(win, space, deltaTime)
			handleError(err)
		}

		// Update hero bullets.
		for _, bullet := range hero.spawnedBullets {
			err = bullet.update(space, deltaTime)
			handleError(err)
		}

		wallSprite.Draw(win, pixel.IM.Moved(pixel.V(width/2, height/2)))

		// Draw platforms.
		lowerPlatform.draw(win)
		middlePlatform.draw(win)
		higherPlatform.draw(win)

		// Draw beholders.
		if !lowerBeholder.dead {
			lowerBeholder.draw(win)
		}

		if !higherBeholder.dead {
			higherBeholder.draw(win)
		}

		// Draw hero.
		if !hero.dead {
			hero.draw(win)
		}

		// Draw beholder bullets.
		for _, bullet := range lowerBeholder.spawnedBullets {
			bullet.draw(win)
		}

		for _, bullet := range higherBeholder.spawnedBullets {
			bullet.draw(win)
		}

		for _, bullet := range hero.spawnedBullets {
			bullet.draw(win)
		}

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

func cirnoToPixel(vec cirno.Vector) pixel.Vec {
	return pixel.V(vec.X, vec.Y)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
