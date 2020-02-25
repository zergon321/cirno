package cirno_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tracer8086/cirno"
)

func TestLinesAreCollinear(t *testing.T) {
	l1 := cirno.NewLine(cirno.NewVector(1, 3), cirno.NewVector(3, 5))
	l2 := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(5, 5))
	l3 := cirno.NewLine(cirno.NewVector(3, 1), cirno.NewVector(5, 1))

	assert.True(t, l1.CollinearTo(l2))
	assert.True(t, l2.CollinearTo(l1))
	assert.False(t, l1.CollinearTo(l3))
	assert.False(t, l2.CollinearTo(l3))
	assert.False(t, l3.CollinearTo(l1))
	assert.False(t, l3.CollinearTo(l2))
}

func TestProjectPointOntoLine(t *testing.T) {
	line := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(8, 1))
	p1 := cirno.NewVector(3, 4)
	r1 := cirno.NewVector(3, 1)
	p2 := cirno.NewVector(9, 5)
	r2 := cirno.NewVector(9, 1)
	p3 := cirno.NewVector(-1, -3)
	r3 := cirno.NewVector(-1, 1)

	assert.Equal(t, line.ProjectPoint(p1), r1)
	assert.Equal(t, line.ProjectPoint(p2), r2)
	assert.Equal(t, line.ProjectPoint(p3), r3)

	line = cirno.NewLine(cirno.NewVector(3, 2), cirno.NewVector(6, 4))
	t.Log("Projection of P1:", line.ProjectPoint(p1))
	t.Log("Projection of P2:", line.ProjectPoint(p2))
	t.Log("Projection of P3:", line.ProjectPoint(p3))
}

func TestLineSegmentsOnSameLine(t *testing.T) {
	l1 := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(3, 3))
	l2 := cirno.NewLine(cirno.NewVector(4, 4), cirno.NewVector(6, 6))
	l3 := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(4, 4))
	l4 := cirno.NewLine(cirno.NewVector(2, 3), cirno.NewVector(4, 1))

	assert.True(t, l1.SameLineWith(l2))
	assert.True(t, l2.SameLineWith(l1))
	assert.False(t, l3.SameLineWith(l4))
}

func TestLinesAreParallel(t *testing.T) {
	l1 := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(3, 3))
	l2 := cirno.NewLine(cirno.NewVector(4, 4), cirno.NewVector(6, 6))
	l3 := cirno.NewLine(cirno.NewVector(3, 1), cirno.NewVector(6, 4))

	assert.True(t, l1.ParallelTo(l3))
	assert.True(t, l3.ParallelTo(l1))
	assert.True(t, l2.ParallelTo(l3))
	assert.True(t, l3.ParallelTo(l2))
	assert.False(t, l1.ParallelTo(l2))
	assert.False(t, l2.ParallelTo(l1))
}
