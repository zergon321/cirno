package cirno_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zergon321/cirno"
)

func TestLinesAreCollinear(t *testing.T) {
	l1, err := cirno.NewLine(cirno.NewVector(1, 3), cirno.NewVector(3, 5))
	assert.Nil(t, err)
	l2, err := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(5, 5))
	assert.Nil(t, err)
	l3, err := cirno.NewLine(cirno.NewVector(3, 1), cirno.NewVector(5, 1))
	assert.Nil(t, err)

	res0, err := l1.CollinearTo(l2)
	assert.Nil(t, err)
	res1, err := l2.CollinearTo(l1)
	assert.Nil(t, err)
	res2, err := l1.CollinearTo(l3)
	assert.Nil(t, err)
	res3, err := l2.CollinearTo(l3)
	assert.Nil(t, err)
	res4, err := l3.CollinearTo(l1)
	assert.Nil(t, err)
	res5, err := l3.CollinearTo(l2)
	assert.Nil(t, err)

	assert.True(t, res0)
	assert.True(t, res1)
	assert.False(t, res2)
	assert.False(t, res3)
	assert.False(t, res4)
	assert.False(t, res5)
}

func TestProjectPointOntoLine(t *testing.T) {
	line, err := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(8, 1))
	assert.Nil(t, err)

	p1 := cirno.NewVector(3, 4)
	r1 := cirno.NewVector(3, 1)
	p2 := cirno.NewVector(9, 5)
	r2 := cirno.NewVector(9, 1)
	p3 := cirno.NewVector(-1, -3)
	r3 := cirno.NewVector(-1, 1)

	assert.Equal(t, line.ProjectPoint(p1), r1)
	assert.Equal(t, line.ProjectPoint(p2), r2)
	assert.Equal(t, line.ProjectPoint(p3), r3)

	line, err = cirno.NewLine(cirno.NewVector(3, 2), cirno.NewVector(6, 4))
	assert.Nil(t, err)

	t.Log("Projection of P1:", line.ProjectPoint(p1))
	t.Log("Projection of P2:", line.ProjectPoint(p2))
	t.Log("Projection of P3:", line.ProjectPoint(p3))
}

func TestLineSegmentsOnSameLine(t *testing.T) {
	l1, err := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(3, 3))
	assert.Nil(t, err)
	l2, err := cirno.NewLine(cirno.NewVector(4, 4), cirno.NewVector(6, 6))
	assert.Nil(t, err)
	l3, err := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(4, 4))
	assert.Nil(t, err)
	l4, err := cirno.NewLine(cirno.NewVector(2, 3), cirno.NewVector(4, 1))
	assert.Nil(t, err)

	res0, err := l1.SameLineWith(l2)
	assert.Nil(t, err)
	res1, err := l2.SameLineWith(l1)
	assert.Nil(t, err)
	res2, err := l3.SameLineWith(l4)
	assert.Nil(t, err)

	assert.True(t, res0)
	assert.True(t, res1)
	assert.False(t, res2)
}

func TestLinesAreParallel(t *testing.T) {
	l1, err := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(3, 3))
	assert.Nil(t, err)
	l2, err := cirno.NewLine(cirno.NewVector(4, 4), cirno.NewVector(6, 6))
	assert.Nil(t, err)
	l3, err := cirno.NewLine(cirno.NewVector(3, 1), cirno.NewVector(6, 4))
	assert.Nil(t, err)

	res0, err := l1.ParallelTo(l3)
	assert.Nil(t, err)
	res1, err := l3.ParallelTo(l1)
	assert.Nil(t, err)
	res2, err := l2.ParallelTo(l3)
	assert.Nil(t, err)
	res3, err := l3.ParallelTo(l2)
	assert.Nil(t, err)
	res4, err := l1.ParallelTo(l2)
	assert.Nil(t, err)
	res5, err := l2.ParallelTo(l1)
	assert.Nil(t, err)

	assert.True(t, res0)
	assert.True(t, res1)
	assert.True(t, res2)
	assert.True(t, res3)
	assert.False(t, res4)
	assert.False(t, res5)
}
