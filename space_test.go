package cirno_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zergon321/cirno"
)

func TestSpaceAddShapes(t *testing.T) {
	space, err := cirno.NewSpace(5, 3, 20, 20, cirno.NewVector(0, 0), cirno.NewVector(10, 10), false)
	assert.Nil(t, err)
	rect, err := cirno.NewRectangle(cirno.NewVector(9, 9), 4, 6, 0.0)
	assert.Nil(t, err)
	circle, err := cirno.NewCircle(cirno.NewVector(8, 4), 3)
	assert.Nil(t, err)
	line, err := cirno.NewLine(cirno.NewVector(2, 8), cirno.NewVector(6, 12))
	assert.Nil(t, err)

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

	res0, err := result.Contains(rect)
	assert.Nil(t, err)
	res1, err := result.Contains(circle)
	assert.Nil(t, err)
	res2, err := result.Contains(line)
	assert.Nil(t, err)

	assert.True(t, res0)
	assert.True(t, res1)
	assert.True(t, res2)
	assert.Equal(t, len(result), 3)
}

func TestSpaceRemoveShapes(t *testing.T) {
	space, err := cirno.NewSpace(5, 3, 20, 20, cirno.NewVector(0, 0), cirno.NewVector(10, 10), false)
	assert.Nil(t, err)
	rect, err := cirno.NewRectangle(cirno.NewVector(9, 9), 4, 6, 0.0)
	assert.Nil(t, err)
	circle, err := cirno.NewCircle(cirno.NewVector(8, 4), 3)
	assert.Nil(t, err)
	line, err := cirno.NewLine(cirno.NewVector(2, 8), cirno.NewVector(6, 12))
	assert.Nil(t, err)
	extra, err := cirno.NewCircle(cirno.NewVector(2, 3), 5)
	assert.Nil(t, err)

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

	res0, err := result.Contains(rect)
	assert.Nil(t, err)
	res1, err := result.Contains(circle)
	assert.Nil(t, err)
	res2, err := result.Contains(line)
	assert.Nil(t, err)

	assert.True(t, res0)
	assert.True(t, res1)
	assert.True(t, res2)

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
	rect, err := cirno.NewRectangle(cirno.NewVector(9, 9), 4, 6, 0.0)
	assert.Nil(t, err)
	circle, err := cirno.NewCircle(cirno.NewVector(8, 4), 3)
	assert.Nil(t, err)
	line, err := cirno.NewLine(cirno.NewVector(2, 8), cirno.NewVector(6, 12))
	assert.Nil(t, err)
	extra, err := cirno.NewCircle(cirno.NewVector(2, 3), 5)
	assert.Nil(t, err)

	err = space.Add(rect)
	assert.Nil(t, err)
	err = space.Add(circle)
	assert.Nil(t, err)
	err = space.Add(line)
	assert.Nil(t, err)

	res0, err := space.Contains(rect)
	assert.Nil(t, err)
	res1, err := space.Contains(circle)
	assert.Nil(t, err)
	res2, err := space.Contains(line)
	assert.Nil(t, err)
	res3, err := space.Contains(extra)
	assert.Nil(t, err)

	assert.True(t, res0)
	assert.True(t, res1)
	assert.True(t, res2)
	assert.False(t, res3)
}

func TestSpaceClear(t *testing.T) {
	space, err := cirno.NewSpace(5, 3, 20, 20, cirno.NewVector(0, 0), cirno.NewVector(10, 10), false)
	assert.Nil(t, err)
	rect, err := cirno.NewRectangle(cirno.NewVector(9, 9), 4, 6, 0.0)
	assert.Nil(t, err)
	circle, err := cirno.NewCircle(cirno.NewVector(8, 4), 3)
	assert.Nil(t, err)
	line, err := cirno.NewLine(cirno.NewVector(2, 8), cirno.NewVector(6, 12))
	assert.Nil(t, err)

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
	rect, err := cirno.NewRectangle(cirno.NewVector(9, 9), 4, 6, 0.0)
	assert.Nil(t, err)
	circle, err := cirno.NewCircle(cirno.NewVector(8, 4), 3)
	assert.Nil(t, err)
	line, err := cirno.NewLine(cirno.NewVector(2, 8), cirno.NewVector(6, 12))
	assert.Nil(t, err)

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

	res, err := result[rect].Contains(circle)
	assert.Nil(t, err)
	assert.True(t, res)
	assert.Equal(t, len(result[circle]), 1)

	res, err = result[circle].Contains(rect)
	assert.Nil(t, err)
	assert.True(t, res)
}

func TestGetShapesCollidingWithShape(t *testing.T) {
	space, err := cirno.NewSpace(10, 10, 1280*2, 720*2, cirno.NewVector(0, 0), cirno.NewVector(1280, 720), false)
	assert.Nil(t, err)
	rect, err := cirno.NewRectangle(cirno.NewVector(600, 228), 150, 50, 0.0)
	assert.Nil(t, err)
	circle, err := cirno.NewCircle(cirno.NewVector(8, 4), 3)
	assert.Nil(t, err)
	line, err := cirno.NewLine(cirno.NewVector(480, 240), cirno.NewVector(720, 240))
	assert.Nil(t, err)

	err = space.Add(rect)
	assert.Nil(t, err)
	err = space.Add(circle)
	assert.Nil(t, err)
	err = space.Add(line)
	assert.Nil(t, err)

	shapes, err := space.WouldBeCollidedBy(rect, cirno.NewVector(0, 0.2171355), 0)
	assert.Nil(t, err)
	assert.Equal(t, len(shapes), 1)

	res, err := shapes.Contains(line)
	assert.True(t, res)

	shapes, err = space.CollidingWith(line)
	assert.Nil(t, err)
	assert.Equal(t, len(shapes), 1)

	res, err = shapes.Contains(rect)
	assert.Nil(t, err)
	assert.True(t, res)
}

func TestOneShapeCollisions(t *testing.T) {
	space, err := cirno.NewSpace(4, 1, 20, 20, cirno.NewVector(0, 0), cirno.NewVector(10, 10), false)
	assert.Nil(t, err)
	rect, err := cirno.NewRectangle(cirno.NewVector(9, 9), 4, 6, 0.0)
	assert.Nil(t, err)

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

	r0, err := cirno.NewRectangle(cirno.NewVector(2.5, -2.5), 0.2, 0.2, 0)
	assert.Nil(t, err)
	r1, err := cirno.NewRectangle(cirno.NewVector(7.5, -2.5), 0.2, 0.2, 0)
	assert.Nil(t, err)
	r2, err := cirno.NewRectangle(cirno.NewVector(2.5, -7.5), 0.2, 0.2, 0)
	assert.Nil(t, err)
	r3, err := cirno.NewRectangle(cirno.NewVector(7.5, -7.5), 0.2, 0.2, 0)
	assert.Nil(t, err)

	rects := []cirno.Shape{r0, r1, r2, r3}

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

	err = space.Rebuild()
	assert.Nil(t, err)

	assert.Equal(t, 7, len(space.Cells()))
}
