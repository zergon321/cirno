package cirno_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zergon321/cirno"
)

func TestRectanglesCollision(t *testing.T) {
	r1, _ := cirno.NewRectangle(cirno.NewVector(2, 2), 2, 2, 0)
	r2, _ := cirno.NewRectangle(cirno.NewVector(6, 3), 4, 4, 0)
	r3, _ := cirno.NewRectangle(cirno.NewVector(3, 2), 2, 2, 0)
	r4, _ := cirno.NewRectangle(cirno.NewVector(3, 6), 4, 4, 0)
	r5, _ := cirno.NewRectangle(cirno.NewVector(3, 2), 2, 2, 45)
	r6, _ := cirno.NewRectangle(cirno.NewVector(3, 6), 4, 4, 0)
	r7, _ := cirno.NewRectangle(cirno.NewVector(5, 3), 2, 2, 0)
	r8, _ := cirno.NewRectangle(cirno.NewVector(4, 4), 6, 6, 0)
	r9, _ := cirno.NewRectangle(cirno.NewVector(5, 3), 2, 2, 0)
	r10, _ := cirno.NewRectangle(cirno.NewVector(3, 3), 4, 4, 0)
	r11, _ := cirno.NewRectangle(cirno.NewVector(4, 2), 2, 2, 0)
	r12, _ := cirno.NewRectangle(cirno.NewVector(3, 5), 4, 2, 85)

	result0, _ := cirno.CollisionRectangleToRectangle(r1, r2)
	result1, _ := cirno.CollisionRectangleToRectangle(r3, r4)
	result2, _ := cirno.CollisionRectangleToRectangle(r5, r6)
	result3, _ := cirno.CollisionRectangleToRectangle(r7, r8)
	result4, _ := cirno.CollisionRectangleToRectangle(r9, r10)
	result5, _ := cirno.CollisionRectangleToRectangle(r11, r12)

	assert.False(t, result0)
	assert.False(t, result1)
	assert.False(t, result2)
	assert.True(t, result3)
	assert.True(t, result4)
	assert.True(t, result5)
}

func TestRectangleToCircleCollision(t *testing.T) {
	r1, _ := cirno.NewRectangle(cirno.NewVector(3, 2), 4, 2, 45)
	c1, _ := cirno.NewCircle(cirno.NewVector(6, 5), 1.0)
	r2, _ := cirno.NewRectangle(cirno.NewVector(3, 2), 4, 2, 85)
	c2, _ := cirno.NewCircle(cirno.NewVector(3, 4.5), 1.0)
	r3, _ := cirno.NewRectangle(cirno.NewVector(3, 2), 4, 2, 85)
	c3, _ := cirno.NewCircle(cirno.NewVector(3, 2), 0.5)
	r4, _ := cirno.NewRectangle(cirno.NewVector(3, 2), 4, 2, 85)
	c4, _ := cirno.NewCircle(cirno.NewVector(3, 2), 3)

	res0, _ := cirno.CollisionRectangleToCircle(r1, c1)
	res1, _ := cirno.CollisionRectangleToCircle(r2, c2)
	res2, _ := cirno.CollisionRectangleToCircle(r3, c3)
	res3, _ := cirno.CollisionRectangleToCircle(r4, c4)

	assert.False(t, res0)
	assert.True(t, res1)
	assert.True(t, res2)
	assert.True(t, res3)
}

func TestLinesIntersection(t *testing.T) {
	l1, _ := cirno.NewLine(cirno.NewVector(1, 2), cirno.NewVector(5, 5))
	l2, _ := cirno.NewLine(cirno.NewVector(4, 3), cirno.NewVector(6, 1))
	l3, _ := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(4, 4))
	l4, _ := cirno.NewLine(cirno.NewVector(2, 3), cirno.NewVector(4, 1))
	l5, _ := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(3, 3))
	l6, _ := cirno.NewLine(cirno.NewVector(4, 4), cirno.NewVector(6, 6))
	l7, _ := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(3, 3))
	l8, _ := cirno.NewLine(cirno.NewVector(2, 2), cirno.NewVector(4, 4))
	l9, _ := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(4, 4))
	l10, _ := cirno.NewLine(cirno.NewVector(2, 2), cirno.NewVector(3, 3))
	l11, _ := cirno.NewLine(cirno.NewVector(5, 1), cirno.NewVector(5, 4))
	l12, _ := cirno.NewLine(cirno.NewVector(5, 2), cirno.NewVector(5, 3))
	l13, _ := cirno.NewLine(cirno.NewVector(5, 4), cirno.NewVector(8, 8))
	l14, _ := cirno.NewLine(cirno.NewVector(2, 5), cirno.NewVector(7, 1))
	l15, _ := cirno.NewLine(cirno.NewVector(3, 3), cirno.NewVector(4, 3))
	l16, _ := cirno.NewLine(cirno.NewVector(3, 4), cirno.NewVector(4, 4))
	l17, _ := cirno.NewLine(cirno.NewVector(3, 3), cirno.NewVector(4, 4))

	res0, _ := cirno.IntersectionLineToLine(l1, l2)
	res1, _ := cirno.IntersectionLineToLine(l3, l4)
	res2, _ := cirno.IntersectionLineToLine(l5, l6)
	res3, _ := cirno.IntersectionLineToLine(l9, l10)
	res4, _ := cirno.IntersectionLineToLine(l7, l8)
	res5, _ := cirno.IntersectionLineToLine(l11, l12)
	res6, _ := cirno.IntersectionLineToLine(l13, l14)
	res7, _ := cirno.IntersectionLineToLine(l3, l15)
	res8, _ := cirno.IntersectionLineToLine(l3, l16)
	res9, _ := cirno.IntersectionLineToLine(l3, l3)
	res10, _ := cirno.IntersectionLineToLine(l5, l17)

	assert.False(t, res0)
	assert.True(t, res1)
	assert.False(t, res2)
	assert.True(t, res3)
	assert.True(t, res4)
	assert.True(t, res5)
	assert.False(t, res6)
	assert.True(t, res7)
	assert.True(t, res8)
	assert.True(t, res9)
	assert.True(t, res10)
}

func TestLineCircleIntersection(t *testing.T) {
	l1, _ := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(5, 6))
	c1, _ := cirno.NewCircle(cirno.NewVector(3, 4), 2)
	l2, _ := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(5, 2))
	c2, _ := cirno.NewCircle(cirno.NewVector(3, 4), 2)
	l3, _ := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(1, 6))
	c3, _ := cirno.NewCircle(cirno.NewVector(3, 4), 2)
	l4, _ := cirno.NewLine(cirno.NewVector(1, 4), cirno.NewVector(5, 4))
	c4, _ := cirno.NewCircle(cirno.NewVector(6, 4), 3)

	res0, _ := cirno.IntersectionLineToCircle(l1, c1)
	res1, _ := cirno.IntersectionLineToCircle(l2, c2)
	res2, _ := cirno.IntersectionLineToCircle(l3, c3)
	res3, _ := cirno.IntersectionLineToCircle(l4, c4)

	assert.True(t, res0)
	assert.False(t, res1)
	assert.True(t, res2)
	assert.True(t, res3)
}

func TestLineRectangleIntersection(t *testing.T) {
	l1, _ := cirno.NewLine(cirno.NewVector(4, 1), cirno.NewVector(7, 4))
	r1, _ := cirno.NewRectangle(cirno.NewVector(3, 4), 4, 2, 0.0)
	l2, _ := cirno.NewLine(cirno.NewVector(4, 1), cirno.NewVector(7, 4))
	r2, _ := cirno.NewRectangle(cirno.NewVector(6, 2), 2, 2, 0.0)
	l3, _ := cirno.NewLine(cirno.NewVector(4, 1), cirno.NewVector(7, 4))
	r3, _ := cirno.NewRectangle(cirno.NewVector(9, 4), 6, 2, 30.0)
	l4, _ := cirno.NewLine(cirno.NewVector(2, 2), cirno.NewVector(4, 6))
	r4, _ := cirno.NewRectangle(cirno.NewVector(4, 4), 6, 6, 0.0)
	l5, _ := cirno.NewLine(cirno.NewVector(480, 240), cirno.NewVector(720, 240))
	r5, _ := cirno.NewRectangle(cirno.NewVector(600, 228.2171355), 150, 50, 0.0)

	res0, _ := cirno.IntersectionLineToRectangle(l1, r1)
	res1, _ := cirno.IntersectionLineToRectangle(l2, r2)
	res2, _ := cirno.IntersectionLineToRectangle(l3, r3)
	res3, _ := cirno.IntersectionLineToRectangle(l4, r4)
	res4, _ := cirno.IntersectionLineToRectangle(l5, r5)

	assert.False(t, res0)
	assert.True(t, res1)
	assert.True(t, res2)
	assert.True(t, res3)
	assert.True(t, res4)
}

func TestCirclesCollision(t *testing.T) {
	c1, _ := cirno.NewCircle(cirno.NewVector(6, 3), 1)
	c2, _ := cirno.NewCircle(cirno.NewVector(3, 4), 2)
	c3, _ := cirno.NewCircle(cirno.NewVector(6, 3), 2)
	c4, _ := cirno.NewCircle(cirno.NewVector(3, 4), 2)

	res0, _ := cirno.CollisionCircleToCircle(c1, c2)
	res1, _ := cirno.CollisionCircleToCircle(c3, c4)

	assert.False(t, res0)
	assert.True(t, res1)
}
