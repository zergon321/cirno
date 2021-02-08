package cirno

import (
	"fmt"
	"math"
)

// ResolveCollision assumes types of the given shapes
// and detects if they collide.
func ResolveCollision(one, another Shape, useTags bool) (bool, error) {
	if one == nil {
		return false, fmt.Errorf("the first shape is nil")
	}

	if another == nil {
		return false, fmt.Errorf("the second shape is nil")
	}

	if useTags && !one.ShouldCollide(another) {
		return false, nil
	}

	id := one.TypeName() + "_" + another.TypeName()

	switch id {
	case "Rectangle_Rectangle":
		return CollisionRectangleToRectangle(one.(*Rectangle), another.(*Rectangle))

	case "Rectangle_Circle":
		return CollisionRectangleToCircle(one.(*Rectangle), another.(*Circle))

	case "Circle_Rectangle":
		return CollisionRectangleToCircle(another.(*Rectangle), one.(*Circle))

	case "Circle_Circle":
		return CollisionCircleToCircle(one.(*Circle), another.(*Circle))

	case "Line_Line":
		return IntersectionLineToLine(one.(*Line), another.(*Line))

	case "Line_Circle":
		return IntersectionLineToCircle(one.(*Line), another.(*Circle))

	case "Circle_Line":
		return IntersectionLineToCircle(another.(*Line), one.(*Circle))

	case "Line_Rectangle":
		return IntersectionLineToRectangle(one.(*Line), another.(*Rectangle))

	case "Rectangle_Line":
		return IntersectionLineToRectangle(another.(*Line), one.(*Rectangle))
	}

	return false, fmt.Errorf(
		"unknown shape type combination: '%s' and '%s'",
		one.TypeName(), another.TypeName())
}

// CollisionRectangleToRectangle detects if there is an intersection
// between two oriented rectangles.
func CollisionRectangleToRectangle(a, b *Rectangle) (bool, error) {
	if a == nil {
		return false, fmt.Errorf("the first rectangle is nil")
	}

	if b == nil {
		return false, fmt.Errorf("the second rectangle is nil")
	}

	// A vector from the center of rectangle A to the center of rectangle B.
	t := b.center.Subtract(a.center)

	// Check if Ax is parallel to the separating axis and hence the separating axis exists.
	sepAx := math.Abs(Dot(t, a.xAxis)) > a.extents.X+
		math.Abs(Dot(b.xAxis.MultiplyByScalar(b.extents.X), a.xAxis))+
		math.Abs(Dot(b.yAxis.MultiplyByScalar(b.extents.Y), a.xAxis))

	// Check if Ay is parallel to the separating axis and hence the separating axis exists.
	sepAy := math.Abs(Dot(t, a.yAxis)) > a.extents.Y+
		math.Abs(Dot(b.xAxis.MultiplyByScalar(b.extents.X), a.yAxis))+
		math.Abs(Dot(b.yAxis.MultiplyByScalar(b.extents.Y), a.yAxis))

	// Check if Bx is parallel to the separating axis and hence the separating axis exists.
	sepBx := math.Abs(Dot(t, b.xAxis)) > b.extents.X+
		math.Abs(Dot(a.xAxis.MultiplyByScalar(a.extents.X), b.xAxis))+
		math.Abs(Dot(a.yAxis.MultiplyByScalar(a.extents.Y), b.xAxis))

	// Check if By is parallel to the separating axis and hence the separating axis exists.
	sepBy := math.Abs(Dot(t, b.yAxis)) > b.extents.Y+
		math.Abs(Dot(a.xAxis.MultiplyByScalar(a.extents.X), b.yAxis))+
		math.Abs(Dot(a.yAxis.MultiplyByScalar(a.extents.Y), b.yAxis))

	// If the separating axis exists.
	if sepAx || sepAy || sepBx || sepBy {
		return false, nil
	}

	// If the separating axis doesn't exist, the rectangles intersect.
	return true, nil
}

// TODO: use AABB.

// CollisionRectangleToCircle detects if there's an intersection between
// an oriented rectangle and a circle.
func CollisionRectangleToCircle(rect *Rectangle, circle *Circle) (bool, error) {
	if rect == nil {
		return false, fmt.Errorf("the rectangle is nil")
	}

	if circle == nil {
		return false, fmt.Errorf("the circle is nil")
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

	// If the closest point is inside the circle,
	// the rectangle and the circle do intersect.
	return localCircle.ContainsPoint(closestPoint), nil
}

// CollisionCircleToCircle detects if there's an intersection
// between two circles.
func CollisionCircleToCircle(a, b *Circle) (bool, error) {
	if a == nil {
		return false, fmt.Errorf("the first circle is nil")
	}

	if b == nil {
		return false, fmt.Errorf("the second circle is nil")
	}

	t := a.center.Subtract(b.center)
	radiiSum := a.radius + b.radius

	return t.SquaredMagnitude() <= radiiSum*radiiSum, nil
}

// IntersectionLineToLine detects if two lines
// intersect.
func IntersectionLineToLine(a, b *Line) (bool, error) {
	if a == nil {
		return false, fmt.Errorf("the first line is nil")
	}

	if b == nil {
		return false, fmt.Errorf("the second line is nil")
	}

	aMin, aMax := a.GetBoundingBox()
	bMin, bMax := b.GetBoundingBox()

	// If one line fully contains another.
	if a.ContainsPoint(b.p) && a.ContainsPoint(b.q) ||
		b.ContainsPoint(a.p) && b.ContainsPoint(a.q) {
		return true, nil
	}

	aBox := &aabb{min: aMin, max: aMax}
	bBox := &aabb{min: bMin, max: bMax}
	overlap, err := aBox.collidesAABB(bBox)

	if err != nil {
		return false, err
	}

	aTouchesB, err := a.touchesOrCrosses(b)
	bTouchesA, err := b.touchesOrCrosses(a)

	return overlap &&
		aTouchesB && bTouchesA, nil
}

// IntersectionLineToCircle detects if a line and a circle do intersect.
func IntersectionLineToCircle(line *Line, circle *Circle) (bool, error) {
	if line == nil {
		return false, fmt.Errorf("the line is nil")
	}

	if circle == nil {
		return false, fmt.Errorf("the circle is nil")
	}

	if circle.ContainsPoint(line.p) || circle.ContainsPoint(line.q) {
		return true, nil
	}

	pq := line.q.Subtract(line.p)
	t := Dot(circle.center.Subtract(line.p), pq) / Dot(pq, pq)

	if t < 0.0 || t > 1.0 {
		return false, nil
	}

	closestPoint := line.p.Add(pq.MultiplyByScalar(t))

	return SquaredDistance(circle.center, closestPoint) <= circle.radius*circle.radius, nil
}

// IntersectionLineToRectangle detects if there's an intersection between
// a line and a rectangle.
func IntersectionLineToRectangle(line *Line, rect *Rectangle) (bool, error) {
	if line == nil {
		return false, fmt.Errorf("the line is nil")
	}

	if rect == nil {
		return false, fmt.Errorf("the rectangle is nil")
	}

	// The method for two rectangles is as well appropriate
	// for line and rectangle because line segment is just
	// a rectangle with no Y extent.
	lineAxisX := line.q.Subtract(line.Center()).Normalize()
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

	if sepAx || sepAy || sepLineX || sepLineY {
		return false, nil
	}

	return true, nil
}
