package cirno_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zergon321/cirno"
)

func TestRaycast(t *testing.T) {
	// Create new shapes.
	circle, _ := cirno.NewCircle(cirno.NewVector(7, 21), 3)
	rect, _ := cirno.NewRectangle(cirno.NewVector(7.5, 3.5), 11, 5, 0)
	line, _ := cirno.NewLine(cirno.NewVector(24, 24), cirno.NewVector(33, 18))
	cube, _ := cirno.NewRectangle(cirno.NewVector(30, 5), 6, 6, 0)
	rhombus, _ := cirno.NewRectangle(cirno.NewVector(18, 13), 4, 4, 45)
	littleCircle, _ := cirno.NewCircle(cirno.NewVector(32, 24), 2)

	// Create a new space.
	space, err := cirno.NewSpace(1, 10, 64, 64,
		cirno.NewVector(0, 0), cirno.NewVector(64, 64), false)
	assert.Nil(t, err)

	// Fill the space with the shapes.
	err = space.Add(circle, rect, line, cube, rhombus, littleCircle)
	assert.Nil(t, err)

	// Do raycast.
	shape, _, _ := space.Raycast(rhombus.Center(), cirno.NewVector(1, 1), 0, 0)
	assert.Equal(t, shape, line)
	shape, _, _ = space.Raycast(rhombus.Center(), cirno.NewVector(-1, 1), 0, 0)
	assert.Equal(t, shape, circle)
	shape, _, _ = space.Raycast(rhombus.Center(), cirno.NewVector(-1, -1), 0, 0)
	assert.Equal(t, shape, rect)
	shape, _, _ = space.Raycast(rhombus.Center(), cirno.NewVector(1, -1), 0, 0)
	assert.Equal(t, shape, cube)
	shape, _, _ = space.Raycast(rhombus.Center(), cirno.NewVector(0, -1), 0, 0)
	assert.Nil(t, shape)
}
