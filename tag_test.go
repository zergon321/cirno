package cirno_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zergon321/cirno"
)

func TestShouldCollide(t *testing.T) {
	c1, _ := cirno.NewCircle(cirno.NewVector(4, 4), 3)
	c1.SetIdentity(32)
	c1.SetMask(32)
	c2, _ := cirno.NewCircle(cirno.NewVector(8, 4), 2)
	c2.SetIdentity(32 | 64)
	c2.SetMask(32)

	res0, _ := cirno.ResolveCollision(c1, c2, true)
	res1, _ := cirno.ResolveCollision(c2, c2, true)

	assert.True(t, res0)
	assert.True(t, res1)
}

func TestShouldntCollide(t *testing.T) {
	c1, _ := cirno.NewCircle(cirno.NewVector(4, 4), 3)
	c1.SetIdentity(32)
	c1.SetMask(32)
	c2, _ := cirno.NewCircle(cirno.NewVector(8, 4), 2)
	c2.SetIdentity(64)
	c2.SetMask(64)

	res, _ := cirno.ResolveCollision(c1, c2, true)

	assert.False(t, res)
}

func TestControversy(t *testing.T) {
	c1, _ := cirno.NewCircle(cirno.NewVector(4, 4), 3)
	c1.SetIdentity(32)
	c1.SetMask(32)
	c2, _ := cirno.NewCircle(cirno.NewVector(8, 4), 2)
	c2.SetIdentity(32)
	c2.SetMask(64)

	res0, _ := cirno.ResolveCollision(c1, c2, true)
	res1, _ := cirno.ResolveCollision(c2, c1, true)

	assert.True(t, res0)
	assert.False(t, res1)
}
