package cirno_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tracer8086/cirno"
)

func TestAdjustAngle(t *testing.T) {
	var (
		a1  float64 = 360
		a2  float64 = 380
		a3  float64 = 240
		a4  float64 = 720
		a5  float64 = 800
		a6  float64 = -360
		a7  float64 = -380
		a8  float64 = -240
		a9  float64 = -720
		a10 float64 = -800
	)

	assert.Equal(t, cirno.AdjustAngle(a1), 0.0)
	assert.Equal(t, cirno.AdjustAngle(a2), 20.0)
	assert.Equal(t, cirno.AdjustAngle(a3), 240.0)
	assert.Equal(t, cirno.AdjustAngle(a4), 0.0)
	assert.Equal(t, cirno.AdjustAngle(a5), 80.0)
	assert.Equal(t, cirno.AdjustAngle(a6), 0.0)
	assert.Equal(t, cirno.AdjustAngle(a7), 340.0)
	assert.Equal(t, cirno.AdjustAngle(a8), 120.0)
	assert.Equal(t, cirno.AdjustAngle(a9), 0.0)
	assert.Equal(t, cirno.AdjustAngle(a10), 280.0)
}
