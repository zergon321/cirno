package cirno_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zergon321/cirno"
)

func TestShouldCollide(t *testing.T) {
	c1, err := cirno.NewCircle(cirno.NewVector(4, 4), 3)
	assert.Nil(t, err)
	c1.SetIdentity(32)
	c1.SetMask(32)

	c2, err := cirno.NewCircle(cirno.NewVector(8, 4), 2)
	assert.Nil(t, err)
	c2.SetIdentity(32 | 64)
	c2.SetMask(32)

	res0, err := cirno.ResolveCollision(c1, c2, true)
	assert.Nil(t, err)
	res1, err := cirno.ResolveCollision(c2, c2, true)
	assert.Nil(t, err)

	assert.True(t, res0)
	assert.True(t, res1)
}

func TestShouldntCollide(t *testing.T) {
	c1, err := cirno.NewCircle(cirno.NewVector(4, 4), 3)
	assert.Nil(t, err)
	c1.SetIdentity(32)
	c1.SetMask(32)

	c2, err := cirno.NewCircle(cirno.NewVector(8, 4), 2)
	assert.Nil(t, err)
	c2.SetIdentity(64)
	c2.SetMask(64)

	res, err := cirno.ResolveCollision(c1, c2, true)
	assert.Nil(t, err)

	assert.False(t, res)
}

func TestControversy(t *testing.T) {
	c1, err := cirno.NewCircle(cirno.NewVector(4, 4), 3)
	assert.Nil(t, err)
	c1.SetIdentity(32)
	c1.SetMask(32)

	c2, err := cirno.NewCircle(cirno.NewVector(8, 4), 2)
	assert.Nil(t, err)
	c2.SetIdentity(32)
	c2.SetMask(64)

	res0, err := cirno.ResolveCollision(c1, c2, true)
	assert.Nil(t, err)
	res1, err := cirno.ResolveCollision(c2, c1, true)
	assert.Nil(t, err)

	assert.True(t, res0)
	assert.False(t, res1)
}
