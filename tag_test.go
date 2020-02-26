package cirno_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tracer8086/cirno"
)

func TestShouldCollide(t *testing.T) {
	c1 := cirno.NewCircle(cirno.NewVector(4, 4), 3)
	c1.SetIdentity(32)
	c1.SetMask(32)
	c2 := cirno.NewCircle(cirno.NewVector(8, 4), 2)
	c2.SetIdentity(32 | 64)
	c2.SetMask(32)

	assert.True(t, cirno.ResolveCollision(c1, c2, true))
	assert.True(t, cirno.ResolveCollision(c2, c2, true))
}

func TestShouldntCollide(t *testing.T) {
	c1 := cirno.NewCircle(cirno.NewVector(4, 4), 3)
	c1.SetIdentity(32)
	c1.SetMask(32)
	c2 := cirno.NewCircle(cirno.NewVector(8, 4), 2)
	c2.SetIdentity(64)
	c2.SetMask(64)

	assert.False(t, cirno.ResolveCollision(c1, c2, true))
}

func TestControversy(t *testing.T) {
	c1 := cirno.NewCircle(cirno.NewVector(4, 4), 3)
	c1.SetIdentity(32)
	c1.SetMask(32)
	c2 := cirno.NewCircle(cirno.NewVector(8, 4), 2)
	c2.SetIdentity(32)
	c2.SetMask(64)

	assert.True(t, cirno.ResolveCollision(c1, c2, true))
	assert.False(t, cirno.ResolveCollision(c2, c1, true))
}
