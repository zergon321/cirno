package cirno_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zergon321/cirno"
)

func TestContactLineToCircle(t *testing.T) {
	circle, _ := cirno.NewCircle(cirno.NewVector(4, 6), 2)
	line, _ := cirno.NewLine(cirno.NewVector(2, 2), cirno.NewVector(2, 10))
	contacts, _ := cirno.Contact(circle, line)

	assert.Equal(t, len(contacts), 1)
}

func TestContactCircleToCircle(t *testing.T) {
	circle1, _ := cirno.NewCircle(cirno.NewVector(4, 6), 2)
	circle2, _ := cirno.NewCircle(cirno.NewVector(9, 6), 3)
	contacts, _ := cirno.Contact(circle1, circle2)

	assert.Equal(t, len(contacts), 1)
}
