package cirno_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zergon321/cirno"
)

func TestGetRectangleVertices(t *testing.T) {
	rect := cirno.NewRectangle(cirno.NewVector(18, 13), 4, 4, 45)
	vertices := rect.Vertices()

	assert.Equal(t, vertices[0], cirno.NewVector(18, 10))
	assert.Equal(t, vertices[1], cirno.NewVector(15, 13))
	assert.Equal(t, vertices[2], cirno.NewVector(18, 16))
	assert.Equal(t, vertices[3], cirno.NewVector(21, 13))
}
