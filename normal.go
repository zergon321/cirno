package cirno

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

	return closestPoint.Subtract(circle.center).Normalize()
}

// NormalToCircle returns the normal from the given line
// to the circle.
func (line *Line) NormalToCircle(circle *Circle) Vector {
	return circle.NormalToLine(line).MultiplyByScalar(-1)
}

// NormalToLine returns the normal from the given line
// to the other line.
func (line *Line) NormalToLine(other *Line) Vector {
	point := line.ProjectPoint(other.Center())

	return other.Center().Subtract(point).Normalize()
}
