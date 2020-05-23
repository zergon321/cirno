package cirno_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zergon321/cirno"
)

func TestRectanglesCollision(t *testing.T) {
	r1 := cirno.NewRectangle(cirno.NewVector(2, 2), 2, 2, 0)
	r2 := cirno.NewRectangle(cirno.NewVector(6, 3), 4, 4, 0)
	r3 := cirno.NewRectangle(cirno.NewVector(3, 2), 2, 2, 0)
	r4 := cirno.NewRectangle(cirno.NewVector(3, 6), 4, 4, 0)
	r5 := cirno.NewRectangle(cirno.NewVector(3, 2), 2, 2, 45)
	r6 := cirno.NewRectangle(cirno.NewVector(3, 6), 4, 4, 0)
	r7 := cirno.NewRectangle(cirno.NewVector(5, 3), 2, 2, 0)
	r8 := cirno.NewRectangle(cirno.NewVector(4, 4), 6, 6, 0)
	r9 := cirno.NewRectangle(cirno.NewVector(5, 3), 2, 2, 0)
	r10 := cirno.NewRectangle(cirno.NewVector(3, 3), 4, 4, 0)
	r11 := cirno.NewRectangle(cirno.NewVector(4, 2), 2, 2, 0)
	r12 := cirno.NewRectangle(cirno.NewVector(3, 5), 4, 2, 85)

	assert.False(t, cirno.CollisionRectangleToRectangle(r1, r2))
	assert.False(t, cirno.CollisionRectangleToRectangle(r3, r4))
	assert.False(t, cirno.CollisionRectangleToRectangle(r5, r6))
	assert.True(t, cirno.CollisionRectangleToRectangle(r7, r8))
	assert.True(t, cirno.CollisionRectangleToRectangle(r9, r10))
	assert.True(t, cirno.CollisionRectangleToRectangle(r11, r12))
}

func TestRectangleToCircleCollision(t *testing.T) {
	r1 := cirno.NewRectangle(cirno.NewVector(3, 2), 4, 2, 45)
	c1 := cirno.NewCircle(cirno.NewVector(6, 5), 1.0)
	r2 := cirno.NewRectangle(cirno.NewVector(3, 2), 4, 2, 85)
	c2 := cirno.NewCircle(cirno.NewVector(3, 4.5), 1.0)
	r3 := cirno.NewRectangle(cirno.NewVector(3, 2), 4, 2, 85)
	c3 := cirno.NewCircle(cirno.NewVector(3, 2), 0.5)
	r4 := cirno.NewRectangle(cirno.NewVector(3, 2), 4, 2, 85)
	c4 := cirno.NewCircle(cirno.NewVector(3, 2), 3)

	assert.False(t, cirno.CollisionRectangleToCircle(r1, c1))
	assert.True(t, cirno.CollisionRectangleToCircle(r2, c2))
	assert.True(t, cirno.CollisionRectangleToCircle(r3, c3))
	assert.True(t, cirno.CollisionRectangleToCircle(r4, c4))
}

func TestLinesIntersection(t *testing.T) {
	l1 := cirno.NewLine(cirno.NewVector(1, 2), cirno.NewVector(5, 5))
	l2 := cirno.NewLine(cirno.NewVector(4, 3), cirno.NewVector(6, 1))
	l3 := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(4, 4))
	l4 := cirno.NewLine(cirno.NewVector(2, 3), cirno.NewVector(4, 1))
	l5 := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(3, 3))
	l6 := cirno.NewLine(cirno.NewVector(4, 4), cirno.NewVector(6, 6))
	l7 := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(3, 3))
	l8 := cirno.NewLine(cirno.NewVector(2, 2), cirno.NewVector(4, 4))
	l9 := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(4, 4))
	l10 := cirno.NewLine(cirno.NewVector(2, 2), cirno.NewVector(3, 3))
	l11 := cirno.NewLine(cirno.NewVector(5, 1), cirno.NewVector(5, 4))
	l12 := cirno.NewLine(cirno.NewVector(5, 2), cirno.NewVector(5, 3))
	l13 := cirno.NewLine(cirno.NewVector(5, 4), cirno.NewVector(8, 8))
	l14 := cirno.NewLine(cirno.NewVector(2, 5), cirno.NewVector(7, 1))

	assert.False(t, cirno.IntersectionLineToLine(l1, l2))
	assert.True(t, cirno.IntersectionLineToLine(l3, l4))
	assert.False(t, cirno.IntersectionLineToLine(l5, l6))
	assert.True(t, cirno.IntersectionLineToLine(l9, l10))
	assert.True(t, cirno.IntersectionLineToLine(l7, l8))
	assert.True(t, cirno.IntersectionLineToLine(l11, l12))
	assert.False(t, cirno.IntersectionLineToLine(l13, l14))
}

func TestLineCircleIntersection(t *testing.T) {
	l1 := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(5, 6))
	c1 := cirno.NewCircle(cirno.NewVector(3, 4), 2)
	l2 := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(5, 2))
	c2 := cirno.NewCircle(cirno.NewVector(3, 4), 2)
	l3 := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(1, 6))
	c3 := cirno.NewCircle(cirno.NewVector(3, 4), 2)
	l4 := cirno.NewLine(cirno.NewVector(1, 4), cirno.NewVector(5, 4))
	c4 := cirno.NewCircle(cirno.NewVector(6, 4), 3)

	assert.True(t, cirno.IntersectionLineToCircle(l1, c1))
	assert.False(t, cirno.IntersectionLineToCircle(l2, c2))
	assert.True(t, cirno.IntersectionLineToCircle(l3, c3))
	assert.True(t, cirno.IntersectionLineToCircle(l4, c4))
}

func TestLineRectangleIntersection(t *testing.T) {
	l1 := cirno.NewLine(cirno.NewVector(4, 1), cirno.NewVector(7, 4))
	r1 := cirno.NewRectangle(cirno.NewVector(3, 4), 4, 2, 0.0)
	l2 := cirno.NewLine(cirno.NewVector(4, 1), cirno.NewVector(7, 4))
	r2 := cirno.NewRectangle(cirno.NewVector(6, 2), 2, 2, 0.0)
	l3 := cirno.NewLine(cirno.NewVector(4, 1), cirno.NewVector(7, 4))
	r3 := cirno.NewRectangle(cirno.NewVector(9, 4), 6, 2, 30.0)
	l4 := cirno.NewLine(cirno.NewVector(2, 2), cirno.NewVector(4, 6))
	r4 := cirno.NewRectangle(cirno.NewVector(4, 4), 6, 6, 0.0)
	l5 := cirno.NewLine(cirno.NewVector(480, 240), cirno.NewVector(720, 240))
	r5 := cirno.NewRectangle(cirno.NewVector(600, 228.2171355), 150, 50, 0.0)

	assert.False(t, cirno.IntersectionLineToRectangle(l1, r1))
	assert.True(t, cirno.IntersectionLineToRectangle(l2, r2))
	assert.True(t, cirno.IntersectionLineToRectangle(l3, r3))
	assert.True(t, cirno.IntersectionLineToRectangle(l4, r4))
	assert.True(t, cirno.IntersectionLineToRectangle(l5, r5))
}

func TestCirclesCollision(t *testing.T) {
	c1 := cirno.NewCircle(cirno.NewVector(6, 3), 1)
	c2 := cirno.NewCircle(cirno.NewVector(3, 4), 2)
	c3 := cirno.NewCircle(cirno.NewVector(6, 3), 2)
	c4 := cirno.NewCircle(cirno.NewVector(3, 4), 2)

	assert.False(t, cirno.CollisionCircleToCircle(c1, c2))
	assert.True(t, cirno.CollisionCircleToCircle(c3, c4))
}
