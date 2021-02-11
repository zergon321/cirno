package main

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/zergon321/cirno"
)

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
	leftShape, _, err := space.Raycast(leftRayOrigin, cirno.Down(), br.rect.Height()/2+4,
		br.rect.GetMask())

	if err != nil {
		return err
	}

	rightShape, _, err := space.Raycast(rightRayOrigin, cirno.Down(), br.rect.Height()/2+4,
		br.rect.GetMask())

	if err != nil {
		return err
	}

	// Change movement direction
	// and hit circle position.
	if leftShape == nil {
		br.direction = cirno.Right()
		br.sprite = br.anim[1]
	} else if rightShape == nil {
		br.direction = cirno.Left()
		br.sprite = br.anim[0]
	}

	movement := br.direction.MultiplyByScalar(br.speed * deltaTime)

	// Detect player and stop if player detected.
	player, _, err := space.Raycast(br.rect.Center(), br.direction, 384, playerID)

	if err != nil {
		return err
	}

	if player != nil {
		movement = cirno.Zero()

		// Shoot at the player if he is detected.
		if br.bulletTimer == nil {
			// Spawn the first bullet.
			bulletPos := br.rect.Center().
				Add(br.direction.MultiplyByScalar(br.rect.Width() / 4)).
				Add(br.direction.MultiplyByScalar(br.bulletSprite.Frame().W() / 4))
			bulletHitLine, err := cirno.NewLine(
				bulletPos.Add(cirno.Left().MultiplyByScalar(br.bulletSprite.Frame().W()/4)),
				bulletPos.Add(cirno.Right().MultiplyByScalar(br.bulletSprite.Frame().W()/4)),
			)

			if err != nil {
				return err
			}

			bullet := &bloodBullet{
				spawner:   br,
				hitLine:   bulletHitLine,
				sprite:    br.bulletSprite,
				direction: br.direction,
				speed:     br.bulletSpeed,
				transform: pixel.IM.Moved(cirnoToPixel(bulletPos)),
			}
			bullet.transform = bullet.transform.Scaled(cirnoToPixel(bulletPos), 0.5)

			bullet.hitLine.SetData(bullet)
			bullet.hitLine.SetIdentity(bloodBulletID)
			bullet.hitLine.SetMask(playerID)

			err = space.Add(bullet.hitLine)

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
					Add(br.direction.MultiplyByScalar(br.rect.Width() / 4)).
					Add(br.direction.MultiplyByScalar(br.bulletSprite.Frame().W() / 4))
				bulletHitLine, err := cirno.NewLine(
					bulletPos.Add(cirno.Left().MultiplyByScalar(br.bulletSprite.Frame().W()/4)),
					bulletPos.Add(cirno.Right().MultiplyByScalar(br.bulletSprite.Frame().W()/4)),
				)

				if err != nil {
					return err
				}

				bullet := &bloodBullet{
					spawner:   br,
					hitLine:   bulletHitLine,
					sprite:    br.bulletSprite,
					direction: br.direction,
					speed:     br.bulletSpeed,
					transform: pixel.IM.Moved(cirnoToPixel(bulletPos)),
				}
				bullet.transform = bullet.transform.Scaled(cirnoToPixel(bullet.hitLine.Center()), 0.5)

				bullet.hitLine.SetData(bullet)
				bullet.hitLine.SetIdentity(bloodBulletID)
				bullet.hitLine.SetMask(playerID)

				err = space.Add(bullet.hitLine)

				if err != nil {
					return err
				}

				br.spawnedBullets = append(br.spawnedBullets, bullet)

			default:
			}
		}
	} else {
		// Stop shooting if player
		// is off sight.
		br.bulletTimer = nil
	}

	if movement != cirno.Zero() {
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

	if br.direction == cirno.Right() {
		hitCirclePos.X += br.hitCircle.Radius()
	} else if br.direction == cirno.Left() {
		hitCirclePos.X -= br.hitCircle.Radius()
	}

	br.hitCircle.SetPosition(hitCirclePos)
	_, err = space.Update(br.hitCircle)

	if err != nil {
		return err
	}

	// Check collision with player's bullets.
	shapes, err := space.CollidingWith(br.hitCircle)

	if err != nil {
		return err
	}

	bulletShapes := shapes.FilterByIdentity(electroBulletID)

	if err != nil {
		return err
	}

	// Remove all the bullets that hit the beholder.
	if len(bulletShapes) > 0 {
		br.dead = true

		for shape := range bulletShapes {
			bullet := shape.Data().(*electroBullet)
			ind := -1

			for i, eb := range bullet.spawner.spawnedBullets {
				if eb == bullet {
					ind = i
					break
				}
			}

			bullet.spawner.spawnedBullets = append(bullet.spawner.spawnedBullets[:ind],
				bullet.spawner.spawnedBullets[ind+1:]...)
			err = space.Remove(bullet.hitCircle)

			if err != nil {
				return err
			}
		}

		// Remove beholder's hit shapes.
		err = space.Remove(br.rect, br.hitCircle)

		if err != nil {
			return err
		}
	}

	return nil
}

func (br *beholder) draw(target pixel.Target) {
	br.sprite.Draw(target, br.transform)
}
