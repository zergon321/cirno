package cirno_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zergon321/cirno"
)

func TestVectorAdd(t *testing.T) {
	v1 := cirno.NewVector(2, -34)
	v2 := cirno.NewVector(-5, 12)
	res := cirno.NewVector(-3, -22)

	assert.Equal(t, v1.Add(v2), res)
}

func TestVectorSubtract(t *testing.T) {
	v1 := cirno.NewVector(2, -34)
	v2 := cirno.NewVector(-5, 12)
	res := cirno.NewVector(7, -46)

	assert.Equal(t, v1.Subtract(v2), res)
}

func TestVectorMagnitude(t *testing.T) {
	v := cirno.NewVector(1, -2)

	assert.Equal(t, v.Magnitude(), math.Sqrt(5))
}

func TestVectorRotate(t *testing.T) {
	v := cirno.NewVector(2, 0)
	res := cirno.NewVector(0, 2)

	assert.True(t, v.Rotate(90).X-res.X < 0.00000000000001)
	assert.True(t, v.Rotate(90).Y-res.Y < 0.00000000000001)
}

func TestVectorProject(t *testing.T) {
	v := cirno.NewVector(2, 6)
	axis := cirno.NewVector(2, 2)
	res := cirno.NewVector(4, 4)

	assert.Equal(t, v.Project(axis), res)
}

func TestVectorNormalize(t *testing.T) {
	v := cirno.NewVector(22, -14)

	assert.Equal(t, math.Ceil(v.Normalize().Magnitude()), 1.0)
	assert.Equal(t, cirno.Angle(v, v.Normalize()), 0.0)
}

func TestVectorDotProduct(t *testing.T) {
	v1 := cirno.NewVector(2, 6)
	v2 := cirno.NewVector(2, 2)
	res := 16.0

	assert.Equal(t, cirno.Dot(v1, v2), res)
}

func TestVectorAngle(t *testing.T) {
	v1 := cirno.NewVector(2, 0)
	v2 := cirno.NewVector(0, 2)
	res := 90.0

	assert.Equal(t, cirno.Angle(v1, v2), res)
}

func TestVectorPerpendicular(t *testing.T) {
	v1 := cirno.NewVector(2, 3)
	vp1 := v1.PerpendicularClockwise()
	vp2 := v1.PerpendicularCounterClockwise()

	fmt.Println(vp1)
	fmt.Println(vp2)

	assert.Equal(t, cirno.Angle(v1, vp1), 90.0)
	assert.Equal(t, cirno.Angle(v1, vp2), 90.0)
}
