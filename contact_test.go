package cirno_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zergon321/cirno"
)

func TestContactLineToCircle(t *testing.T) {
	circle, err := cirno.NewCircle(cirno.NewVector(4, 6), 2)
	assert.Nil(t, err)
	line, err := cirno.NewLine(cirno.NewVector(2, 2), cirno.NewVector(2, 10))
	assert.Nil(t, err)
	contacts, err := cirno.Contact(circle, line)
	assert.Nil(t, err)

	assert.Equal(t, len(contacts), 1)
}

func TestContactCircleToCircle(t *testing.T) {
	circle1, err := cirno.NewCircle(cirno.NewVector(4, 6), 2)
	assert.Nil(t, err)
	circle2, err := cirno.NewCircle(cirno.NewVector(9, 6), 3)
	assert.Nil(t, err)
	contacts, err := cirno.Contact(circle1, circle2)
	assert.Nil(t, err)

	assert.Equal(t, len(contacts), 1)
}
