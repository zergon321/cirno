package cirno_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zergon321/cirno"
)

func TestGetRectangleVertices(t *testing.T) {
	rect, err := cirno.NewRectangle(cirno.NewVector(18, 13), 4, 4, 45)
	assert.Nil(t, err)
	vertices := rect.Vertices()

	assert.Equal(t, cirno.NewVector(18, 10), vertices[0])
	assert.Equal(t, cirno.NewVector(15, 13), vertices[1])
	assert.Equal(t, cirno.NewVector(18, 16), vertices[2])
	assert.Equal(t, cirno.NewVector(21, 13), vertices[3])
}
