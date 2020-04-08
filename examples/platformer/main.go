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

type platform struct {
	rect      *cirno.Rectangle
	sprite    *pixel.Sprite
	transform pixel.Matrix
}

func (pl *platform) draw(target pixel.Target) {
	pl.sprite.Draw(target, pl.transform)
}

type beholder struct {
	rect           *cirno.Rectangle
	hitCircle      *cirno.Circle
	sprite         *pixel.Sprite
	anim           []*pixel.Sprite
	direction      cirno.Vector
	speed          float64
	bulletSprite   *pixel.Sprite
	bulletSpeed    float64
	bulletCooldown time.Duration
	spawnedBullets []*bloodBullet
	bulletTimer    <-chan time.Time
	transform      pixel.Matrix
	dead           bool
}

func (br *beholder) update(space *cirno.Space, deltaTime float64) error {
	// Cast rays to detect if the beholder is on the edge of the platform.
	leftRayOrigin := cirno.NewVector(br.rect.Center().X-br.rect.Width()/2,
		br.rect.Center().Y)
	rightRayOrigin := cirno.NewVector(br.rect.Center().X+br.rect.Width()/2,
		br.rect.Center().Y)
	leftShape := space.Raycast(leftRayOrigin, cirno.Down, br.rect.Height()/2+4,
		br.rect.GetMask())
	rightShape := space.Raycast(rightRayOrigin, cirno.Down, br.rect.Height()/2+4,
		br.rect.GetMask())

	// Change movement direction
	// and hit circle position.
	if leftShape == nil {
		br.direction = cirno.Right
		br.sprite = br.anim[1]
	} else if rightShape == nil {
		br.direction = cirno.Left
		br.sprite = br.anim[0]
	}

	movement := br.direction.MultiplyByScalar(br.speed * deltaTime)

	// Detect player and stop if player detected.
	player := space.Raycast(br.rect.Center(), br.direction, 384, playerID)

	if player != nil {
		movement = cirno.Zero
		// TODO: shoot at the player if he is detected.

		if br.bulletTimer == nil {
			// Spawn the first bullet.
			bulletPos := br.rect.Center().
				Add(cirno.Right.MultiplyByScalar(br.rect.Width() / 2)).
				Add(cirno.Right.MultiplyByScalar(br.bulletSprite.Frame().W() / 2))
			bullet := &bloodBullet{
				spawner: br,
				hitLine: cirno.NewLine(
					bulletPos.Add(cirno.Left.MultiplyByScalar(br.bulletSprite.Frame().W()/2)),
					bulletPos.Add(cirno.Right.MultiplyByScalar(br.bulletSprite.Frame().W()/2)),
				),
				sprite:    br.bulletSprite,
				direction: br.direction,
				speed:     br.bulletSpeed,
				transform: pixel.IM.Moved(cirnoToPixel(bulletPos)),
			}

			bullet.hitLine.SetIdentity(bloodBulletID)
			err := space.Add(bullet.hitLine)

			if err != nil {
				return err
			}

			br.spawnedBullets = append(br.spawnedBullets, bullet)
			// Start the cooldown timer to shoot bullets
			// with timelapse.
			br.bulletTimer = time.Tick(br.bulletCooldown)
		} else {
			select {
			// Shoot a new bullet after the cooldown.
			case <-br.bulletTimer:
				bulletPos := br.rect.Center().
					Add(cirno.Right.MultiplyByScalar(br.rect.Width() / 2)).
					Add(cirno.Right.MultiplyByScalar(br.bulletSprite.Frame().W() / 2))
				bullet := &bloodBullet{
					spawner: br,
					hitLine: cirno.NewLine(
						bulletPos.Add(cirno.Left.MultiplyByScalar(br.bulletSprite.Frame().W()/2)),
						bulletPos.Add(cirno.Right.MultiplyByScalar(br.bulletSprite.Frame().W()/2)),
					),
					sprite:    br.bulletSprite,
					direction: br.direction,
					speed:     br.bulletSpeed,
					transform: pixel.IM.Moved(cirnoToPixel(bulletPos)),
				}

				bullet.hitLine.SetIdentity(bloodBulletID)
				err := space.Add(bullet.hitLine)

				if err != nil {
					return err
				}

			default:
			}
		}
	} else {
		// Stop shooting if player
		// is off sight.
		br.bulletTimer = nil
	}

	if movement != cirno.Zero {
		// Move rect.
		prev := br.rect.Center()
		br.rect.Move(movement)
		space.AdjustShapePosition(br.rect)
		_, err := space.Update(br.rect)

		if err != nil {
			return err
		}

		// Move sprite.
		br.transform = br.transform.Moved(cirnoToPixel(br.rect.Center().Subtract(prev)))
	}

	// Move hit circle.
	hitCirclePos := cirno.NewVector(br.rect.Center().X,
		br.rect.Center().Y+br.rect.Height()/2-br.hitCircle.Radius())

	if br.direction == cirno.Right {
		hitCirclePos.X += br.hitCircle.Radius()
	} else if br.direction == cirno.Left {
		hitCirclePos.X -= br.hitCircle.Radius()
	}

	br.hitCircle.SetPosition(hitCirclePos)
	_, err := space.Update(br.hitCircle)

	return err
}

func (br *beholder) draw(target pixel.Target) {
	br.sprite.Draw(target, br.transform)
}

type bloodBullet struct {
	spawner   *beholder
	hitLine   *cirno.Line
	sprite    *pixel.Sprite
	direction cirno.Vector
	speed     float64
	transform pixel.Matrix
}

func (bb *bloodBullet) update(space *cirno.Space, deltaTime float64) error {
	movement := bb.direction.MultiplyByScalar(bb.speed * deltaTime)
	bb.hitLine.Move(movement)
	_, err := space.Update(bb.hitLine)

	if err != nil {
		return err
	}

	// If the bullet is off screen, remove it.
	if bb.hitLine.Center().X > width || bb.hitLine.Center().X < 0 {
		ind := -1

		for i, bullet := range bb.spawner.spawnedBullets {
			if bullet == bb {
				ind = i
				break
			}
		}

		bb.spawner.spawnedBullets = append(bb.spawner.spawnedBullets[:ind],
			bb.spawner.spawnedBullets[ind+1:]...)
	}

	return nil
}

func (bb *bloodBullet) draw(target pixel.Target) {
	bb.sprite.Draw(target, bb.transform)
}

type player struct {
	speed            float64
	jumpAcceleration float64
	verticalSpeed    float64
	terminalSpeed    float64
	rect             *cirno.Rectangle
	sprite           *pixel.Sprite
	animation        []*pixel.Sprite
	transform        pixel.Matrix
	dead             bool
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

	// Find out if player is grounded.
	leftRayOrigin := cirno.NewVector(p.rect.Center().X-p.rect.Width()/2,
		p.rect.Center().Y)
	rightRayOrigin := cirno.NewVector(p.rect.Center().X+p.rect.Width()/2,
		p.rect.Center().Y)
	leftShape := space.Raycast(leftRayOrigin, cirno.Down,
		p.rect.Height()/2+4, p.rect.GetMask())
	rightShape := space.Raycast(rightRayOrigin, cirno.Down,
		p.rect.Height()/2+4, p.rect.GetMask())
	grounded := leftShape != nil || rightShape != nil

	// Compute vertical speed.
	if grounded {
		if win.JustPressed(pixelgl.KeyUp) {
			p.verticalSpeed = p.jumpAcceleration
			fmt.Println(p.verticalSpeed)
		} else {
			p.verticalSpeed = 0
		}
	} else {
		p.verticalSpeed += gravity * deltaTime

		if p.verticalSpeed < p.terminalSpeed*deltaTime {
			p.verticalSpeed = p.terminalSpeed * deltaTime
		}
	}

	// Adjust movement with framerate.
	movement.X *= p.speed * deltaTime
	movement.Y = p.verticalSpeed

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
			pos, _, err := cirno.Approximate(p.rect, movement, 0, shapes,
				intensity, space.UseTags())

			if err != nil {
				return err
			}

			movement = pos.Subtract(p.rect.Center())
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

		// TODO: check collision with bullets.
		shapes, err = space.CollidingWith(p.rect)

		if err != nil {
			return err
		}

		bulletShapes := shapes.FilterByIdentity(bloodBulletID)

		if len(bulletShapes) > 0 {
			p.dead = true

			// Remove all the bullets that hit the player.
			for bulletShape := range bulletShapes {
				bullet := bulletShape.Data().(*bloodBullet)
				ind := -1

				for i := range bullet.spawner.spawnedBullets {
					if bullet.spawner.spawnedBullets[i] == bullet {
						ind = i
						break
					}
				}

				bullet.spawner.spawnedBullets = append(bullet.spawner.spawnedBullets[:ind],
					bullet.spawner.spawnedBullets[ind+1:]...)
			}
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
	//electroBulletSprite := pixel.NewSprite(projectileSheet, pixel.R(0, 0, 64, 64))
	bloodBulletSprite := pixel.NewSprite(projectileSheet, pixel.R(64, 0, 192, 64))
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
		rect:           cirno.NewRectangle(cirno.NewVector(320, 316), 64.5, 157.5, 0),
		hitCircle:      cirno.NewCircle(cirno.NewVector(304, 378.75), 16),
		sprite:         beholderLeftSprite,
		anim:           []*pixel.Sprite{beholderLeftSprite, beholderRightSprite},
		bulletSprite:   bloodBulletSprite,
		bulletSpeed:    400,
		bulletCooldown: 600 * time.Millisecond,
		bulletTimer:    nil,
		spawnedBullets: make([]*bloodBullet, 0),
		speed:          250,
		direction:      cirno.Left,
		transform:      pixel.IM.Scaled(pixel.ZV, 0.5).Moved(pixel.V(320, 316)),
		dead:           false,
	}
	higherBeholder := &beholder{
		rect:           cirno.NewRectangle(cirno.NewVector(960, 496), 64.5, 157.5, 0),
		hitCircle:      cirno.NewCircle(cirno.NewVector(976, 558.75), 16),
		sprite:         beholderRightSprite,
		anim:           []*pixel.Sprite{beholderLeftSprite, beholderRightSprite},
		bulletSprite:   bloodBulletSprite,
		bulletSpeed:    300,
		bulletCooldown: 400 * time.Millisecond,
		bulletTimer:    nil,
		spawnedBullets: make([]*bloodBullet, 0),
		speed:          300,
		direction:      cirno.Right,
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
	hero := &player{
		speed:            500,
		jumpAcceleration: 80,
		verticalSpeed:    gravity,
		terminalSpeed:    gravity,
		rect:             cirno.NewRectangle(cirno.NewVector(640, 121), 64, 128, 0),
		sprite:           testmanLeftSprite,
		animation:        []*pixel.Sprite{testmanLeftSprite, testmanRightSprite},
		transform:        pixel.IM.Scaled(pixel.V(0, 0), 2).Moved(pixel.V(640, 121)),
		dead:             false,
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

		// Update beholders.
		err = lowerBeholder.update(space, deltaTime)
		handleError(err)
		err = higherBeholder.update(space, deltaTime)
		handleError(err)

		// Update bullets.
		for _, bullet := range lowerBeholder.spawnedBullets {
			err = bullet.update(space, deltaTime)
			handleError(err)
		}

		for _, bullet := range higherBeholder.spawnedBullets {
			err = bullet.update(space, deltaTime)
			handleError(err)
		}

		// Update hero.
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

		// Draw bullets.
		for _, bullet := range lowerBeholder.spawnedBullets {
			bullet.draw(win)
		}

		for _, bullet := range higherBeholder.spawnedBullets {
			bullet.draw(win)
		}

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
