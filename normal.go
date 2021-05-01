package cirno

import (
	"fmt"
	"math"
)

// NormalTo returns the normal from the given circle
// to the other shape.
func (circle *Circle) NormalTo(shape Shape) (Vector, error) {
	if shape == nil {
		return Zero(), fmt.Errorf("the shape is nil")
	}

	switch other := shape.(type) {
	case *Circle:
		return circle.NormalToCircle(other)

	case *Line:
		return circle.NormalToLine(other)

	case *Rectangle:
		return circle.NormalToRectangle(other)
	}

	return Zero(), fmt.Errorf("unknown shape type")
}

// NormalTo returns the normal from the given rectangle
// to the other shape.
func (rect *Rectangle) NormalTo(shape Shape) (Vector, error) {
	if shape == nil {
		return Zero(), fmt.Errorf("the shape is nil")
	}

	switch other := shape.(type) {
	case *Circle:
		return rect.NormalToCircle(other)

	case *Line:
		return rect.NormalToLine(other)

	case *Rectangle:
		return rect.NormalToRectangle(other)
	}

	return Zero(), fmt.Errorf("unknown shape type")
}

// NormalTo returns the normal from the given line to
// the other shape.
func (line *Line) NormalTo(shape Shape) (Vector, error) {
	if shape == nil {
		return Zero(), fmt.Errorf("the shape is nil")
	}

	switch other := shape.(type) {
	case *Circle:
		return line.NormalToCircle(other)

	case *Line:
		return line.NormalToLine(other)

	case *Rectangle:
		return line.NormalToRectangle(other)
	}

	return Zero(), fmt.Errorf("unknown shape type")
}

// NormalToCircle returns the normal from the given circle
// to the other circle.
func (circle *Circle) NormalToCircle(other *Circle) (Vector, error) {
	if other == nil {
		return Zero(), fmt.Errorf("the other circle is nil")
	}

	return other.center.Subtract(circle.center).Normalize()
}

// NormalToRectangle returns the normal from the given circle
// to the rectangle.
func (circle *Circle) NormalToRectangle(rect *Rectangle) (Vector, error) {
	if rect == nil {
		return Zero(), fmt.Errorf("the rectangle is nil")
	}

	// Transform the circle center coordinates from the world space
	// to the rectangle's local space.
	t := circle.center.Subtract(rect.center)
	theta := -rect.angle
	t = t.Rotate(theta)
	localCircle := &Circle{
		center: t,
		radius: circle.radius,
	}
	localRect := &Rectangle{
		center:  NewVector(0, 0),
		extents: NewVector(rect.Width()/2, rect.Height()/2),
		xAxis:   NewVector(1, 0),
		yAxis:   NewVector(0, 1),
	}
	closestPoint := localCircle.center

	// Find the point of the rectangle which is closest to
	// the center of the circle.
	if closestPoint.X < localRect.Min().X {
		closestPoint.X = localRect.Min().X
	} else if closestPoint.X > localRect.Max().X {
		closestPoint.X = localRect.Max().X
	}

	if closestPoint.Y < localRect.Min().Y {
		closestPoint.Y = localRect.Min().Y
	} else if closestPoint.Y > localRect.Max().Y {
		closestPoint.Y = localRect.Max().Y
	}

	closestPoint = closestPoint.Rotate(-theta).Add(rect.center)
	normal, err := closestPoint.Subtract(circle.center).Normalize()

	return normal, err
}

// NormalToLine returns the normal from the given circle
// to the line.
func (circle *Circle) NormalToLine(line *Line) (Vector, error) {
	if line == nil {
		return Zero(), fmt.Errorf("the line is nil")
	}

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

	normal, err := closestPoint.Subtract(circle.center).Normalize()

	if err != nil {
		return Zero(), err
	}

	if math.IsNaN(normal.X) {
		normal.X = 0.0
	}

	if math.IsNaN(normal.Y) {
		normal.Y = 0.0
	}

	return normal, nil
}

// NormalToCircle returns the normal from the given line
// to the circle.
func (line *Line) NormalToCircle(circle *Circle) (Vector, error) {
	if circle == nil {
		return Zero(), fmt.Errorf("the circle is nil")
	}

	normalToLine, err := circle.NormalToLine(line)

	if err != nil {
		return Zero(), err
	}

	return normalToLine.MultiplyByScalar(-1), nil
}

// NormalToLine returns the normal from the given line
// to the other line.
func (line *Line) NormalToLine(other *Line) (Vector, error) {
	if line == nil {
		return Zero(), fmt.Errorf("the line is nil")
	}

	normal := Zero()
	pRightOfLine, err := line.isPointRightOfLine(other.p)

	if err != nil {
		return Zero(), err
	}

	qRightOfLine, err := line.isPointRightOfLine(other.q)

	if err != nil {
		return Zero(), err
	}

	if pRightOfLine == qRightOfLine {
		pointProj := line.ProjectPoint(other.p)
		normal, err = other.p.Subtract(pointProj).Normalize()

		if err != nil {
			return Zero(), err
		}
	} else {
		pointProj := other.ProjectPoint(line.p)
		normal, err = pointProj.Subtract(line.p).Normalize()

		if err != nil {
			return Zero(), err
		}
	}

	return normal, nil
}

// NormalToRectangle returns the normal from the given line
// to the rectangle.
func (line *Line) NormalToRectangle(rect *Rectangle) (Vector, error) {
	if rect == nil {
		return Zero(), fmt.Errorf("the rectangle is nil")
	}

	normalToLine, err := rect.NormalToLine(line)

	if err != nil {
		return Zero(), err
	}

	return normalToLine.MultiplyByScalar(-1), nil
}

// NormalToCircle returns the normal from the given rectangle
// to the circle.
func (rect *Rectangle) NormalToCircle(circle *Circle) (Vector, error) {
	if circle == nil {
		return Zero(), fmt.Errorf("the circle is nil")
	}

	normalToRect, err := circle.NormalToRectangle(rect)

	if err != nil {
		return Zero(), err
	}

	return normalToRect.MultiplyByScalar(-1), nil
}

// NormalToLine returns the normal between the given rectangle
// and the line.
func (rect *Rectangle) NormalToLine(line *Line) (Vector, error) {
	if line == nil {
		return Zero(), fmt.Errorf("the line is nil")
	}

	lineAxisX, err := line.q.Subtract(line.Center()).Normalize()

	if err != nil {
		return Zero(), err
	}

	lineAxisY := lineAxisX.Rotate(90)
	lineExtent := line.Length() / 2
	t := line.Center().Subtract(rect.center)

	sepAx := math.Abs(Dot(t, rect.xAxis)) > rect.extents.X+
		math.Abs(Dot(lineAxisX.MultiplyByScalar(lineExtent), rect.xAxis))
	sepAy := math.Abs(Dot(t, rect.yAxis)) > rect.extents.Y+
		math.Abs(Dot(lineAxisX.MultiplyByScalar(lineExtent), rect.yAxis))
	sepLineX := math.Abs(Dot(t, lineAxisX)) > lineExtent+
		math.Abs(Dot(rect.xAxis.MultiplyByScalar(rect.extents.X), lineAxisX))+
		math.Abs(Dot(rect.yAxis.MultiplyByScalar(rect.extents.Y), lineAxisX))
	sepLineY := math.Abs(Dot(t, lineAxisY)) >
		math.Abs(Dot(rect.xAxis.MultiplyByScalar(rect.extents.X), lineAxisY))+
			math.Abs(Dot(rect.yAxis.MultiplyByScalar(rect.extents.Y), lineAxisY))

	var normal Vector

	if sepAx {
		normal = rect.xAxis
		sepLine, err := NewLine(rect.center,
			rect.center.Add(rect.yAxis))

		if err != nil {
			return Zero(), err
		}

		if sepLine.Orientation(line.Center()) < 0 {
			normal = normal.MultiplyByScalar(-1)
		}
	} else if sepAy {
		normal = rect.yAxis
		sepLine, err := NewLine(rect.center,
			rect.center.Add(rect.xAxis))

		if err != nil {
			return Zero(), err
		}

		if sepLine.Orientation(line.Center()) > 0 {
			normal = normal.MultiplyByScalar(-1)
		}
	} else if sepLineX {
		normal = lineAxisX
		sepLine, err := NewLine(line.Center(),
			line.Center().Add(lineAxisY))

		if err != nil {
			return Zero(), err
		}

		if sepLine.Orientation(rect.center) > 0 {
			normal = normal.MultiplyByScalar(-1)
		}
	} else if sepLineY {
		normal = lineAxisY
		sepLine, err := NewLine(line.Center(),
			line.Center().Add(lineAxisX))

		if err != nil {
			return Zero(), err
		}

		if sepLine.Orientation(rect.center) < 0 {
			normal = normal.MultiplyByScalar(-1)
		}
	}

	return normal, nil
}

// NormalToRectangle returns the normal from the given rectangle to
// the other rectangle.
func (rect *Rectangle) NormalToRectangle(other *Rectangle) (Vector, error) {
	if other == nil {
		return Zero(), fmt.Errorf("the rectangle is nil")
	}

	// A vector from the center of rectangle A to the center of rectangle B.
	t := other.center.Subtract(rect.center)

	// Check if Ax is parallel to the separating axis and hence the separating axis exists.
	sepAx := math.Abs(Dot(t, rect.xAxis)) > rect.extents.X+
		math.Abs(Dot(other.xAxis.MultiplyByScalar(other.extents.X), rect.xAxis))+
		math.Abs(Dot(other.yAxis.MultiplyByScalar(other.extents.Y), rect.xAxis))

	// Check if Ay is parallel to the separating axis and hence the separating axis exists.
	sepAy := math.Abs(Dot(t, rect.yAxis)) > rect.extents.Y+
		math.Abs(Dot(other.xAxis.MultiplyByScalar(other.extents.X), rect.yAxis))+
		math.Abs(Dot(other.yAxis.MultiplyByScalar(other.extents.Y), rect.yAxis))

	// Check if Bx is parallel to the separating axis and hence the separating axis exists.
	sepBx := math.Abs(Dot(t, other.xAxis)) > other.extents.X+
		math.Abs(Dot(rect.xAxis.MultiplyByScalar(rect.extents.X), other.xAxis))+
		math.Abs(Dot(rect.yAxis.MultiplyByScalar(rect.extents.Y), other.xAxis))

	// Check if By is parallel to the separating axis and hence the separating axis exists.
	sepBy := math.Abs(Dot(t, other.yAxis)) > other.extents.Y+
		math.Abs(Dot(rect.xAxis.MultiplyByScalar(rect.extents.X), other.yAxis))+
		math.Abs(Dot(rect.yAxis.MultiplyByScalar(rect.extents.Y), other.yAxis))

	var normal Vector

	if sepAx {
		normal = rect.xAxis
		sepLine, err := NewLine(rect.center,
			rect.center.Add(rect.yAxis))

		if err != nil {
			return Zero(), err
		}

		if sepLine.Orientation(other.center) < 0 {
			normal = normal.MultiplyByScalar(-1)
		}
	} else if sepAy {
		normal = rect.yAxis
		sepLine, err := NewLine(rect.center,
			rect.center.Add(rect.xAxis))

		if err != nil {
			return Zero(), err
		}

		if sepLine.Orientation(other.center) > 0 {
			normal = normal.MultiplyByScalar(-1)
		}
	} else if sepBx {
		normal = other.xAxis
		sepLine, err := NewLine(other.center,
			other.center.Add(other.yAxis))

		if err != nil {
			return Zero(), err
		}

		if sepLine.Orientation(rect.center) > 0 {
			normal = normal.MultiplyByScalar(-1)
		}
	} else if sepBy {
		normal = other.yAxis
		sepLine, err := NewLine(other.center,
			other.center.Add(other.xAxis))

		if err != nil {
			return Zero(), err
		}

		if sepLine.Orientation(rect.center) < 0 {
			normal = normal.MultiplyByScalar(-1)
		}
	}

	return normal, nil
}
