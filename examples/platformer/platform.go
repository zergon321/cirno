package main

import (
	"github.com/faiface/pixel"
	"github.com/zergon321/cirno"
)

type platform struct {
	rect      *cirno.Rectangle
	sprite    *pixel.Sprite
	transform pixel.Matrix
}

func (pl *platform) draw(target pixel.Target) {
	pl.sprite.Draw(target, pl.transform)
}
