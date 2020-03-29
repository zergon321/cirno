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
	speed     float64
	jumpSpeed float64
	rect      *cirno.Rectangle
	sprite    *pixel.Sprite
	animation []*pixel.Sprite
	transform pixel.Matrix
	dead      bool
}

func (p *player) update(win *pixelgl.Window, space *cirno.Space, deltaTime float64) error {
	movement := cirno.Zero

	// Read inputs.
	if win.Pressed(pixelgl.KeyLeft) {
		movement.X--
	}

	if win.Pressed(pixelgl.KeyRight) {
		movement.X++
	}

	if win.JustPressed(pixelgl.KeyUp) {
		movement.Y++
	}

	// Adjust movement with framerate.
	movement = movement.MultiplyByScalar(p.speed * deltaTime)
	movement.Y = gravity * deltaTime

	if movement != cirno.Zero {
		// Update player sprite.
		if movement.X > 0 {
			p.sprite = p.animation[1]
		} else if movement.X < 0 {
			p.sprite = p.animation[0]
		}

		shapes, err := space.WouldBeColliding(p.rect, movement, 0)

		if err != nil {
			return err
		}

		// Resolve collision.
		if len(shapes) > 0 {
			temp := movement.Y
			movement.Y = 0
			newShapes, err := space.WouldBeColliding(p.rect, movement, 0)

			if err != nil {
				return err
			}

			// If we can't go up or down.
			if len(newShapes) > 0 {
				movement.Y = temp
				pos, _, err := cirno.Approximate(p.rect, movement, 0, shapes,
					intensity, space.UseTags())

				if err != nil {
					return err
				}

				movement = pos.Subtract(p.rect.Center())
			}
		}

		// Move sprite and hitbox.
		prev := p.rect.Center()
		p.rect.Move(movement)
		space.AdjustShapePosition(p.rect)
		p.transform = p.transform.Moved(cirnoToPixel(p.rect.Center().Subtract(prev)))
		_, err = space.Update(p.rect)

		if err != nil {
			return err
		}
	}

	return nil
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
	testmanRightSprite := pixel.NewSprite(testmanPic, pixel.R(32, 0, 64, 64))
	//electroBulletSprite := pixel.NewSprite(projectileSheet, pixel.R(0, 0, 64, 64))
	//bloodBulletSprite := pixel.NewSprite(projectileSheet, pixel.R(64, 0, 192, 64))
	platformSprite := pixel.NewSprite(platformPic, pixel.R(0, 0, 128, 32))
	beholderLeftSprite := pixel.NewSprite(beholderPic, pixel.R(0, 0, 129, 315))
	beholderRightSprite := pixel.NewSprite(beholderPic, pixel.R(129, 0, 258, 315))

	//bulletBatch := pixel.NewBatch(new(pixel.TrianglesData), projectileSheet)

	// Create platforms.
	lowerPlatform := &platform{
		rect:   cirno.NewRectangle(cirno.NewVector(640, 40), 384, 32, 0),
		sprite: platformSprite,
		transform: pixel.IM.ScaledXY(pixel.ZV, pixel.V(3, 1)).
			Moved(pixel.V(640, 40)),
	}
	middlePlatform := &platform{
		rect:   cirno.NewRectangle(cirno.NewVector(320, 220), 384, 32, 0),
		sprite: platformSprite,
		transform: pixel.IM.ScaledXY(pixel.ZV, pixel.V(3, 1)).
			Moved(pixel.V(320, 220)),
	}
	higherPlatform := &platform{
		rect:   cirno.NewRectangle(cirno.NewVector(960, 400), 384, 32, 0),
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

	// Create beholders.
	lowerBeholder := &beholder{
		rect:      cirno.NewRectangle(cirno.NewVector(320, 314.75), 64.5, 157.5, 0),
		hitCircle: cirno.NewCircle(cirno.NewVector(304, 377.5), 16),
		sprite:    beholderLeftSprite,
		transform: pixel.IM.Scaled(pixel.ZV, 0.5).Moved(pixel.V(320, 314.75)),
		dead:      false,
	}
	higherBeholder := &beholder{
		rect:      cirno.NewRectangle(cirno.NewVector(960, 494.75), 64.5, 157.5, 0),
		hitCircle: cirno.NewCircle(cirno.NewVector(976, 557.5), 16),
		sprite:    beholderRightSprite,
		transform: pixel.IM.Scaled(pixel.ZV, 0.5).Moved(pixel.V(960, 494.75)),
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
		speed:     500,
		jumpSpeed: 2300 * 2,
		rect:      cirno.NewRectangle(cirno.NewVector(640, 121), 64, 128, 0),
		sprite:    testmanLeftSprite,
		animation: []*pixel.Sprite{testmanLeftSprite, testmanRightSprite},
		transform: pixel.IM.Scaled(pixel.V(0, 0), 2).Moved(pixel.V(640, 121)),
		dead:      false,
	}

	hero.rect.SetIdentity(playerID)
	hero.rect.SetMask(platformID | electroBulletID)
	hero.rect.SetData(hero)

	// Create a new collision space.
	space, err := cirno.NewSpace(5, 20, width*4, height*4,
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
		imd.Color = colors.Lightgreen
	}

	for !win.Closed() {
		deltaTime := time.Since(last).Seconds()
		last = time.Now()

		err = hero.update(win, space, deltaTime)
		handleError(err)

		wallSprite.Draw(win, pixel.IM.Moved(pixel.V(width/2, height/2)))

		// Draw platforms.
		lowerPlatform.draw(win)
		middlePlatform.draw(win)
		higherPlatform.draw(win)

		// Draw beholders.
		lowerBeholder.draw(win)
		higherBeholder.draw(win)

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

func cirnoToPixel(vec cirno.Vector) pixel.Vec {
	return pixel.V(vec.X, vec.Y)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
