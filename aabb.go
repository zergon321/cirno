package cirno

import (
	"fmt"
)

// aabb is an axis-aligned bounding box,
// a rectangle which has no orientation
// and whose edges are aligned with X and Y axes.
//
// For inner usage only.
type aabb struct {
	min Vector
	max Vector
}

// toRectangle returns an oriented rectangle
// based on the given AABB.
func (bb *aabb) toRectangle() *Rectangle {
	return &Rectangle{
		center: bb.center(),
		extents: NewVector((bb.max.X-bb.min.X)/2.0,
			(bb.max.Y-bb.min.Y)/2.0),
		xAxis: Right(),
		yAxis: Up(),
	}
}

// center returns the central point of the AABB.
func (bb *aabb) center() Vector {
	return bb.min.Add(bb.max).MultiplyByScalar(0.5)
}

// vertices returns the vertices of the AABB.
func (bb *aabb) vertices() [4]Vector {
	a := bb.min
	b := NewVector(bb.min.X, bb.max.Y)
	c := bb.max
	d := NewVector(bb.max.X, bb.min.Y)

	return [4]Vector{a, b, c, d}
}

// containsPoint returns true if the given
// point is located inside the AABB.
func (bb *aabb) containsPoint(point Vector) bool {
	return point.X >= bb.min.X &&
		point.Y >= bb.min.Y &&
		point.X <= bb.max.X &&
		point.Y <= bb.max.Y
}

// collidesShape returns true if the AABB
// overlaps the given shape, and false otherwise.
func (bb *aabb) collidesShape(shape Shape) (bool, error) {
	switch other := shape.(type) {
	case *Rectangle:
		return bb.collidesRectangle(other)

	case *Circle:
		return bb.collidesCircle(other)

	case *Line:
		return bb.collidesLine(other)

	default:
		return false, fmt.Errorf("shape type is unknown")
	}
}

// collidesAABB returns true if the given AABB
// overlaps another AABB, and false otherwise.
func (bb *aabb) collidesAABB(other *aabb) (bool, error) {
	if other == nil {
		return false, fmt.Errorf("the other AABB is nil")
	}

	return bb.min.X <= other.max.X &&
		bb.max.X >= other.min.X &&
		bb.min.Y <= other.max.Y &&
		bb.max.Y >= other.min.Y, nil
}

// collidesRectangle returns true if the given AABB
// overlaps the given rectangle, and false otherwise.
func (bb *aabb) collidesRectangle(rect *Rectangle) (bool, error) {
	bbRect := &Rectangle{
		center: bb.center(),
		extents: NewVector((bb.max.X-bb.min.X)/2.0,
			(bb.max.Y-bb.min.Y)/2.0),
		xAxis: Right(),
		yAxis: Up(),
	}
	bbRect.treeNodes = []*quadTreeNode{}

	return CollisionRectangleToRectangle(bbRect, rect)
}

// collidesLine returns true if the given AABB
// overlaps the given line, and false otherwise.
func (bb *aabb) collidesLine(line *Line) (bool, error) {
	if line == nil {
		return false, fmt.Errorf("the line is nil")
	}

	vertices := bb.vertices()
	ab := &Line{
		p:     vertices[0],
		q:     vertices[1],
		angle: 90.0,
	}
	bc := &Line{
		p:     vertices[1],
		q:     vertices[2],
		angle: 0.0,
	}
	cd := &Line{
		p:     vertices[2],
		q:     vertices[3],
		angle: 270.0,
	}
	ad := &Line{
		p:     vertices[0],
		q:     vertices[3],
		angle: 0.0,
	}

	if bb.containsPoint(line.p) ||
		bb.containsPoint(line.q) {
		return true, nil
	}

	abIntersects, err := IntersectionLineToLine(ab, line)

	if err != nil {
		return false, err
	}

	bcIntersects, err := IntersectionLineToLine(bc, line)

	if err != nil {
		return false, err
	}

	cdIntersects, err := IntersectionLineToLine(cd, line)

	if err != nil {
		return false, err
	}

	adIntersects, err := IntersectionLineToLine(ad, line)

	if err != nil {
		return false, err
	}

	return abIntersects ||
		bcIntersects ||
		cdIntersects ||
		adIntersects, nil
}

// collidesCircle returns true if the AABB collides
// the given circle, and false otherwise.
func (bb *aabb) collidesCircle(circle *Circle) (bool, error) {
	if circle == nil {
		return false, fmt.Errorf("the circle is nil")
	}

	closestPoint := circle.center

	// Find the point of the rectangle which is closest to
	// the center of the circle.
	if closestPoint.X < bb.min.X {
		closestPoint.X = bb.min.X
	} else if closestPoint.X > bb.max.X {
		closestPoint.X = bb.max.X
	}

	if closestPoint.Y < bb.min.Y {
		closestPoint.Y = bb.min.Y
	} else if closestPoint.Y > bb.max.Y {
		closestPoint.Y = bb.max.Y
	}

	// If the closest point is inside the circle,
	// the rectangle and the circle do intersect.
	return circle.ContainsPoint(closestPoint), nil
}

// newAABB creates a new AABB out of min and max points.
func newAABB(min, max Vector) (*aabb, error) {
	if min.X >= max.X || min.Y >= max.Y {
		return nil, fmt.Errorf(
			"invalid points specified for AABB")
	}

	return &aabb{
		min: min,
		max: max,
	}, nil
}
