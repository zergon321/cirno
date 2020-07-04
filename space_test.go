package cirno_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zergon321/cirno"
)

func TestSpaceAddShapes(t *testing.T) {
	space, err := cirno.NewSpace(5, 3, 20, 20, cirno.NewVector(0, 0), cirno.NewVector(10, 10), false)
	assert.Nil(t, err)
	rect := cirno.NewRectangle(cirno.NewVector(9, 9), 4, 6, 0.0)
	circle := cirno.NewCircle(cirno.NewVector(8, 4), 3)
	line := cirno.NewLine(cirno.NewVector(2, 8), cirno.NewVector(6, 12))

	err = space.Add(rect)
	assert.Nil(t, err)
	err = space.Add(circle)
	assert.Nil(t, err)
	err = space.Add(line)
	assert.Nil(t, err)
	err = space.Add(circle)
	assert.Nil(t, err)
	err = space.Add(line)
	assert.Nil(t, err)
	err = space.Add(rect)
	assert.Nil(t, err)

	result := space.Shapes()

	assert.Equal(t, len(result), 3)
	assert.True(t, result.Contains(rect))
	assert.True(t, result.Contains(circle))
	assert.True(t, result.Contains(line))
}

func TestSpaceRemoveShapes(t *testing.T) {
	space, err := cirno.NewSpace(5, 3, 20, 20, cirno.NewVector(0, 0), cirno.NewVector(10, 10), false)
	assert.Nil(t, err)
	rect := cirno.NewRectangle(cirno.NewVector(9, 9), 4, 6, 0.0)
	circle := cirno.NewCircle(cirno.NewVector(8, 4), 3)
	line := cirno.NewLine(cirno.NewVector(2, 8), cirno.NewVector(6, 12))
	extra := cirno.NewCircle(cirno.NewVector(2, 3), 5)

	err = space.Add(rect)
	assert.Nil(t, err)
	err = space.Add(circle)
	assert.Nil(t, err)
	err = space.Add(line)
	assert.Nil(t, err)

	// Remove the shape that's not in the space.
	err = space.Remove(extra)
	assert.Nil(t, err)

	result := space.Shapes()

	assert.Equal(t, len(result), 3)
	assert.True(t, result.Contains(rect))
	assert.True(t, result.Contains(circle))
	assert.True(t, result.Contains(line))

	// Remove all the rest.
	err = space.Remove(line)
	assert.Nil(t, err)
	err = space.Remove(rect)
	assert.Nil(t, err)
	err = space.Remove(circle)
	assert.Nil(t, err)

	result = space.Shapes()
	assert.Equal(t, len(result), 0)
}

func TestSpaceContainsShape(t *testing.T) {
	space, err := cirno.NewSpace(5, 3, 20, 20, cirno.NewVector(0, 0), cirno.NewVector(10, 10), false)
	assert.Nil(t, err)
	rect := cirno.NewRectangle(cirno.NewVector(9, 9), 4, 6, 0.0)
	circle := cirno.NewCircle(cirno.NewVector(8, 4), 3)
	line := cirno.NewLine(cirno.NewVector(2, 8), cirno.NewVector(6, 12))
	extra := cirno.NewCircle(cirno.NewVector(2, 3), 5)

	err = space.Add(rect)
	assert.Nil(t, err)
	err = space.Add(circle)
	assert.Nil(t, err)
	err = space.Add(line)
	assert.Nil(t, err)

	assert.True(t, space.Contains(rect))
	assert.True(t, space.Contains(circle))
	assert.True(t, space.Contains(line))
	assert.False(t, space.Contains(extra))
}

func TestSpaceClear(t *testing.T) {
	space, err := cirno.NewSpace(5, 3, 20, 20, cirno.NewVector(0, 0), cirno.NewVector(10, 10), false)
	assert.Nil(t, err)
	rect := cirno.NewRectangle(cirno.NewVector(9, 9), 4, 6, 0.0)
	circle := cirno.NewCircle(cirno.NewVector(8, 4), 3)
	line := cirno.NewLine(cirno.NewVector(2, 8), cirno.NewVector(6, 12))

	err = space.Add(rect)
	assert.Nil(t, err)
	err = space.Add(circle)
	assert.Nil(t, err)
	err = space.Add(line)
	assert.Nil(t, err)
	err = space.Clear()
	assert.Nil(t, err)

	result := space.Shapes()
	assert.Equal(t, len(result), 0)
}

func TestGetCollidingShapes(t *testing.T) {
	space, err := cirno.NewSpace(4, 1, 20, 20, cirno.NewVector(0, 0), cirno.NewVector(10, 10), false)
	assert.Nil(t, err)
	rect := cirno.NewRectangle(cirno.NewVector(9, 9), 4, 6, 0.0)
	circle := cirno.NewCircle(cirno.NewVector(8, 4), 3)
	line := cirno.NewLine(cirno.NewVector(2, 8), cirno.NewVector(6, 12))

	err = space.Add(rect)
	assert.Nil(t, err)
	err = space.Add(circle)
	assert.Nil(t, err)
	err = space.Add(line)
	assert.Nil(t, err)

	result, err := space.CollidingShapes()

	assert.Nil(t, err)
	assert.Equal(t, len(result), 2)
	assert.Equal(t, len(result[rect]), 1)
	assert.True(t, result[rect].Contains(circle))
	assert.Equal(t, len(result[circle]), 1)
	assert.True(t, result[circle].Contains(rect))
}

func TestGetShapesCollidingWithShape(t *testing.T) {
	space, err := cirno.NewSpace(10, 10, 1280*2, 720*2, cirno.NewVector(0, 0), cirno.NewVector(1280, 720), false)
	assert.Nil(t, err)
	rect := cirno.NewRectangle(cirno.NewVector(600, 228), 150, 50, 0.0)
	circle := cirno.NewCircle(cirno.NewVector(8, 4), 3)
	line := cirno.NewLine(cirno.NewVector(480, 240), cirno.NewVector(720, 240))

	err = space.Add(rect)
	assert.Nil(t, err)
	err = space.Add(circle)
	assert.Nil(t, err)
	err = space.Add(line)
	assert.Nil(t, err)

	shapes, err := space.WouldBeCollidedBy(rect, cirno.NewVector(0, 0.2171355), 0)

	assert.Nil(t, err)
	assert.Equal(t, len(shapes), 1)
	assert.True(t, shapes.Contains(line))

	shapes, err = space.CollidingWith(line)

	assert.Nil(t, err)
	assert.Equal(t, len(shapes), 1)
	assert.True(t, shapes.Contains(rect))
}

func TestOneShapeCollisions(t *testing.T) {
	space, err := cirno.NewSpace(4, 1, 20, 20, cirno.NewVector(0, 0), cirno.NewVector(10, 10), false)
	assert.Nil(t, err)
	rect := cirno.NewRectangle(cirno.NewVector(9, 9), 4, 6, 0.0)

	err = space.Add(rect)
	assert.Nil(t, err)

	shapeGroups, err := space.CollidingShapes()

	assert.Nil(t, err)
	assert.Equal(t, len(shapeGroups), 0)

	shapes, err := space.CollidingWith(rect)

	assert.Nil(t, err)
	assert.Equal(t, len(shapes), 0)
}

func TestQuadTree(t *testing.T) {
	space, err := cirno.NewSpace(6, 1, 20, 20,
		cirno.NewVector(-20, -20), cirno.NewVector(10, 10), false)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(space.Cells()))

	rects := []cirno.Shape{
		cirno.NewRectangle(cirno.NewVector(2.5, -2.5), 0.2, 0.2, 0),
		cirno.NewRectangle(cirno.NewVector(7.5, -2.5), 0.2, 0.2, 0),
		cirno.NewRectangle(cirno.NewVector(2.5, -7.5), 0.2, 0.2, 0),
		cirno.NewRectangle(cirno.NewVector(7.5, -7.5), 0.2, 0.2, 0),
	}

	// Test quad tree split.
	err = space.Add(rects...)

	assert.Nil(t, err)
	assert.Equal(t, 7, len(space.Cells()))

	// Test quad tree assemble.
	rects[0].SetPosition(cirno.NewVector(-7.5, 7.5))
	rects[1].SetPosition(cirno.NewVector(-2.5, 7.5))
	rects[2].SetPosition(cirno.NewVector(-7.5, 2.5))
	rects[3].SetPosition(cirno.NewVector(-2.5, 2.5))

	for _, shape := range rects {
		_, err = space.Update(shape)
		assert.Nil(t, err)
	}

	assert.Equal(t, 7, len(space.Cells()))
}
