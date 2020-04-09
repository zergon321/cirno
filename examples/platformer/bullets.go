package main

import (
	"github.com/faiface/pixel"
	"github.com/zergon321/cirno"
)

type bloodBullet struct {
	spawner   *beholder
	hitLine   *cirno.Line
	sprite    *pixel.Sprite
	direction cirno.Vector
	speed     float64
	transform pixel.Matrix
}

func (bb *bloodBullet) update(space *cirno.Space, deltaTime float64) error {
	// Move the bullet.
	movement := bb.direction.MultiplyByScalar(bb.speed * deltaTime)
	bb.hitLine.Move(movement)
	_, err := space.Update(bb.hitLine)

	if err != nil {
		return err
	}

	bb.transform = bb.transform.Moved(cirnoToPixel(movement))

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
		err = space.Remove(bb.hitLine)

		if err != nil {
			return err
		}
	}

	return nil
}

func (bb *bloodBullet) draw(target pixel.Target) {
	bb.sprite.Draw(target, bb.transform)
}

type electroBullet struct {
	spawner   *player
	hitCircle *cirno.Circle
	sprite    *pixel.Sprite
	direction cirno.Vector
	speed     float64
	transform pixel.Matrix
}

func (eb *electroBullet) update(space *cirno.Space, deltaTime float64) error {
	// Move the bullet.
	movement := eb.direction.MultiplyByScalar(eb.speed * deltaTime)
	eb.hitCircle.Move(movement)
	_, err := space.Update(eb.hitCircle)

	if err != nil {
		return err
	}

	eb.transform = eb.transform.Moved(cirnoToPixel(movement))

	// Remove the bullet if it's out of bounds.
	if eb.hitCircle.Center().X > width || eb.hitCircle.Center().X < 0 {
		ind := -1

		for i, bullet := range eb.spawner.spawnedBullets {
			if bullet == eb {
				ind = i
				break
			}
		}

		eb.spawner.spawnedBullets = append(eb.spawner.spawnedBullets[:ind],
			eb.spawner.spawnedBullets[ind+1:]...)
		err = space.Remove(eb.hitCircle)

		if err != nil {
			return err
		}
	}

	return nil
}

func (eb *electroBullet) draw(target pixel.Target) {
	eb.sprite.Draw(target, eb.transform)
}
