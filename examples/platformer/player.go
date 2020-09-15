package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/zergon321/cirno"
)

type player struct {
	speed            float64
	jumpAcceleration float64
	verticalSpeed    float64
	terminalSpeed    float64
	aim              cirno.Vector
	bulletSpeed      float64
	bulletSprite     *pixel.Sprite
	spawnedBullets   []*electroBullet
	rect             *cirno.Rectangle
	sprite           *pixel.Sprite
	animation        []*pixel.Sprite
	transform        pixel.Matrix
	dead             bool
}

func (p *player) update(win *pixelgl.Window, space *cirno.Space, deltaTime float64) error {
	movement := cirno.Zero

	// Read movement inputs to change aim
	// and animation.
	if win.Pressed(pixelgl.KeyLeft) {
		movement.X--
		p.aim = cirno.Left
		p.sprite = p.animation[1]
	}

	if win.Pressed(pixelgl.KeyRight) {
		movement.X++
		p.aim = cirno.Right
		p.sprite = p.animation[0]
	}

	// If player shoots a bullet.
	if win.JustPressed(pixelgl.KeyZ) {
		bulletPos := p.rect.Center().
			Add(p.aim.MultiplyByScalar(p.rect.Width() / 4)).
			Add(p.aim.MultiplyByScalar(p.bulletSprite.Frame().W() / 4))
		bullet := &electroBullet{
			spawner:   p,
			hitCircle: cirno.NewCircle(bulletPos, p.bulletSprite.Frame().W()/4),
			sprite:    p.bulletSprite,
			direction: p.aim,
			speed:     p.bulletSpeed,
			transform: pixel.IM.Moved(cirnoToPixel(bulletPos)),
		}
		bullet.transform = bullet.transform.Scaled(cirnoToPixel(bulletPos), 0.5)

		bullet.hitCircle.SetData(bullet)
		bullet.hitCircle.SetIdentity(electroBulletID)
		bullet.hitCircle.SetMask(beholderEyeID)

		err := space.Add(bullet.hitCircle)

		if err != nil {
			return err
		}

		p.spawnedBullets = append(p.spawnedBullets, bullet)
	}

	// Find out if player is grounded.
	leftRayOrigin := cirno.NewVector(p.rect.Center().X-p.rect.Width()/2,
		p.rect.Center().Y)
	rightRayOrigin := cirno.NewVector(p.rect.Center().X+p.rect.Width()/2,
		p.rect.Center().Y)
	leftShape, _ := space.Raycast(leftRayOrigin, cirno.Down,
		p.rect.Height()/2+4, p.rect.GetMask())
	rightShape, _ := space.Raycast(rightRayOrigin, cirno.Down,
		p.rect.Height()/2+4, p.rect.GetMask())
	grounded := leftShape != nil || rightShape != nil

	// Compute vertical speed.
	if grounded {
		if win.JustPressed(pixelgl.KeyUp) {
			p.verticalSpeed = p.jumpAcceleration
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

		shapes, err := space.WouldBeCollidedBy(p.rect, movement, 0)

		if err != nil {
			return err
		}

		// Resolve collision.
		if len(shapes) > 0 {
			pos, _, _, err := cirno.Approximate(p.rect, movement, 0, shapes,
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
	}

	// Check collision with bullets.
	shapes, err := space.CollidingWith(p.rect)

	if err != nil {
		return err
	}

	bulletShapes := shapes.FilterByIdentity(bloodBulletID)

	// If a bullet or more hit the player.
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
			err = space.Remove(bullet.hitLine)

			if err != nil {
				return err
			}
		}

		// Remove player's hitbox.
		err = space.Remove(p.rect)

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *player) draw(target pixel.Target) {
	p.sprite.Draw(target, p.transform)
}
