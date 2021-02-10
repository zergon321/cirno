package cirno_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zergon321/cirno"
)

func TestLinesAreCollinear(t *testing.T) {
	l1, _ := cirno.NewLine(cirno.NewVector(1, 3), cirno.NewVector(3, 5))
	l2, _ := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(5, 5))
	l3, _ := cirno.NewLine(cirno.NewVector(3, 1), cirno.NewVector(5, 1))

	res0, _ := l1.CollinearTo(l2)
	res1, _ := l2.CollinearTo(l1)
	res2, _ := l1.CollinearTo(l3)
	res3, _ := l2.CollinearTo(l3)
	res4, _ := l3.CollinearTo(l1)
	res5, _ := l3.CollinearTo(l2)

	assert.True(t, res0)
	assert.True(t, res1)
	assert.False(t, res2)
	assert.False(t, res3)
	assert.False(t, res4)
	assert.False(t, res5)
}

func TestProjectPointOntoLine(t *testing.T) {
	line, _ := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(8, 1))
	p1 := cirno.NewVector(3, 4)
	r1 := cirno.NewVector(3, 1)
	p2 := cirno.NewVector(9, 5)
	r2 := cirno.NewVector(9, 1)
	p3 := cirno.NewVector(-1, -3)
	r3 := cirno.NewVector(-1, 1)

	assert.Equal(t, line.ProjectPoint(p1), r1)
	assert.Equal(t, line.ProjectPoint(p2), r2)
	assert.Equal(t, line.ProjectPoint(p3), r3)

	line, _ = cirno.NewLine(cirno.NewVector(3, 2), cirno.NewVector(6, 4))
	t.Log("Projection of P1:", line.ProjectPoint(p1))
	t.Log("Projection of P2:", line.ProjectPoint(p2))
	t.Log("Projection of P3:", line.ProjectPoint(p3))
}

func TestLineSegmentsOnSameLine(t *testing.T) {
	l1, _ := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(3, 3))
	l2, _ := cirno.NewLine(cirno.NewVector(4, 4), cirno.NewVector(6, 6))
	l3, _ := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(4, 4))
	l4, _ := cirno.NewLine(cirno.NewVector(2, 3), cirno.NewVector(4, 1))

	res0, _ := l1.SameLineWith(l2)
	res1, _ := l2.SameLineWith(l1)
	res2, _ := l3.SameLineWith(l4)

	assert.True(t, res0)
	assert.True(t, res1)
	assert.False(t, res2)
}

func TestLinesAreParallel(t *testing.T) {
	l1, _ := cirno.NewLine(cirno.NewVector(1, 1), cirno.NewVector(3, 3))
	l2, _ := cirno.NewLine(cirno.NewVector(4, 4), cirno.NewVector(6, 6))
	l3, _ := cirno.NewLine(cirno.NewVector(3, 1), cirno.NewVector(6, 4))

	res0, _ := l1.ParallelTo(l3)
	res1, _ := l3.ParallelTo(l1)
	res2, _ := l2.ParallelTo(l3)
	res3, _ := l3.ParallelTo(l2)
	res4, _ := l1.ParallelTo(l2)
	res5, _ := l2.ParallelTo(l1)

	assert.True(t, res0)
	assert.True(t, res1)
	assert.True(t, res2)
	assert.True(t, res3)
	assert.False(t, res4)
	assert.False(t, res5)
}
