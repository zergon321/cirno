package cirno

import (
	"math"
)

// NormalTo returns the normal from the given circle
// to the other shape.
func (circle *Circle) NormalTo(shape Shape) Vector {
	switch other := shape.(type) {
	case *Circle:
		return circle.NormalToCircle(other)

	case *Line:
		return circle.NormalToLine(other)
	}

	return Zero
}

// NormalTo returns the normal from the given rectangle
// to the other shape.
func (rect *Rectangle) NormalTo(shape Shape) Vector {
	return Zero
}

// NormalTo returns the normal from the given line to
// the other shape.
func (line *Line) NormalTo(shape Shape) Vector {
	switch other := shape.(type) {
	case *Circle:
		return line.NormalToCircle(other)

	case *Line:
		return line.NormalToLine(other)
	}

	return Zero
}

// NormalToCircle returns the normal from the given circle
// to the other circle.
func (circle *Circle) NormalToCircle(other *Circle) Vector {
	return other.center.Subtract(circle.center).Normalize()
}

// NormalToLine returns the normal from the given circle
// to the line.
func (circle *Circle) NormalToLine(line *Line) Vector {
	closestPoint := line.ProjectPoint(circle.center)

	if !line.ContainsPoint(closestPoint) {
		cp := line.P().Subtract(circle.Center())
		cq := line.Q().Subtract(circle.Center())

		if cp.SquaredMagnitude() < cq.SquaredMagnitude() {
			closestPoint = line.P()
		} else {
			closestPoint = line.Q()
		}
	}

	normal := closestPoint.Subtract(circle.center).Normalize()

	if math.IsNaN(normal.X) {
		normal.X = 0.0
	}

	if math.IsNaN(normal.Y) {
		normal.Y = 0.0
	}

	return normal
}

// NormalToCircle returns the normal from the given line
// to the circle.
func (line *Line) NormalToCircle(circle *Circle) Vector {
	return circle.NormalToLine(line).MultiplyByScalar(-1)
}

// TODO: fix line to line normal.

// NormalToLine returns the normal from the given line
// to the other line.
func (line *Line) NormalToLine(other *Line) Vector {
	pp := other.ProjectPoint(line.p)
	qp := other.ProjectPoint(line.q)
	pv := pp.Subtract(line.p)
	qv := qp.Subtract(line.q)

	if pv.SquaredMagnitude() < qv.SquaredMagnitude() {
		return pv.Normalize()
	} else {
		return qv.Normalize()
	}
}
