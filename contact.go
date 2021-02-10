package cirno

import (
	"fmt"
	"math"
)

// Contact returns the contact points between two given shapes
// (if they exist).
func Contact(one, other Shape) ([]Vector, error) {
	if one == nil {
		return nil, fmt.Errorf("the first shape is nil")
	}

	if other == nil {
		return nil, fmt.Errorf("the second shape is nil")
	}

	id := one.TypeName() + "_" + other.TypeName()

	switch id {
	case "Rectangle_Rectangle":
		return ContactRectangleToRectangle(one.(*Rectangle), other.(*Rectangle))

	case "Rectangle_Circle":
		return ContactRectangleToCircle(one.(*Rectangle), other.(*Circle))

	case "Circle_Rectangle":
		return ContactRectangleToCircle(other.(*Rectangle), one.(*Circle))

	case "Circle_Circle":
		return ContactCircleToCircle(one.(*Circle), other.(*Circle))

	case "Line_Line":
		return ContactLineToLine(one.(*Line), other.(*Line))

	case "Line_Circle":
		return ContactLineToCircle(one.(*Line), other.(*Circle))

	case "Circle_Line":
		return ContactLineToCircle(other.(*Line), one.(*Circle))

	case "Line_Rectangle":
		return ContactLineToRectangle(one.(*Line), other.(*Rectangle))

	case "Rectangle_Line":
		return ContactLineToRectangle(other.(*Line), one.(*Rectangle))
	}

	return nil, fmt.Errorf(
		"unknown shape type combination: '%s' and '%s'",
		one.TypeName(), other.TypeName())
}

// ContactLineToCircle returns the contact points between the
// line and the circle (if they exist).
func ContactLineToCircle(line *Line, circle *Circle) ([]Vector, error) {
	if line == nil {
		return nil, fmt.Errorf("the line is nil")
	}

	if circle == nil {
		return nil, fmt.Errorf("the circle is nil")
	}

	ax := math.Pow(line.p.X-line.q.X, 2)
	ay := math.Pow(line.p.Y-line.q.Y, 2)

	bx := 2 * (line.p.X - line.q.X) * (line.q.X - circle.center.X)
	by := 2 * (line.p.Y - line.q.Y) * (line.q.Y - circle.center.Y)

	cx := math.Pow(line.q.X-circle.center.X, 2)
	cy := math.Pow(line.q.Y-circle.center.Y, 2)

	a := ax + ay
	b := bx + by
	c := cx + cy - circle.radius*circle.radius

	d := b*b - 4*a*c

	// If there is no intersection between the line and
	// the circle.
	if d < 0.0 {
		return []Vector{}, nil
	} else if d < Epsilon {
		// There is probably one point of intersection.
		t := -b / (2 * a)

		// If we really have an intersection point.
		if t >= 0.0 && t <= 1.0 {
			contact := Vector{
				X: t*line.p.X + (1-t)*line.q.X,
				Y: t*line.p.Y + (1-t)*line.q.Y,
			}

			return []Vector{contact}, nil
		}

		return []Vector{}, nil
	} else {
		// There is probably two points of intersection.
		contacts := make([]Vector, 0)
		t1 := (-b + math.Sqrt(d)) / (2 * a)
		t2 := (-b - math.Sqrt(d)) / (2 * a)

		// Check the first contact point.
		if t1 >= 0.0 && t1 <= 1.0 {
			contact := Vector{
				X: t1*line.p.X + (1-t1)*line.q.X,
				Y: t1*line.p.Y + (1-t1)*line.q.Y,
			}

			contacts = append(contacts, contact)
		}

		// Check the second contact point.
		if t2 >= 0.0 && t2 <= 1.0 {
			contact := Vector{
				X: t2*line.p.X + (1-t2)*line.q.X,
				Y: t2*line.p.Y + (1-t2)*line.q.Y,
			}

			contacts = append(contacts, contact)
		}

		return contacts, nil
	}
}

// ContactLineToLine returns the contact point between
// two lines (if it exists).
func ContactLineToLine(one, other *Line) ([]Vector, error) {
	if one == nil {
		return nil, fmt.Errorf("the first line is nil")
	}

	if other == nil {
		return nil, fmt.Errorf("the second line is nil")
	}

	innerContact := func(a, b, c, d Vector) []Vector {
		cmp := c.Subtract(a)
		r := b.Subtract(a)
		s := d.Subtract(c)

		cmpxr := cmp.X*r.Y - cmp.Y*r.X
		cmpxs := cmp.X*s.Y - cmp.Y*s.X
		rxs := r.X*s.Y - r.Y*s.X

		if cmpxr < Epsilon {
			return []Vector{}
		}

		if rxs < Epsilon {
			return []Vector{}
		}

		rxsr := 1.0 / rxs
		t := cmpxs * rxsr
		u := cmpxr * rxsr

		if t >= 0.0 && t <= 1.0 && u >= 0.0 && u <= 1.0 {
			contact := a.Add(r.MultiplyByScalar(t))

			return []Vector{contact}
		}

		return []Vector{}
	}

	contacts := innerContact(one.p, one.q, other.p, other.q)
	contacts = append(contacts,
		innerContact(one.q, one.p, other.p, other.q)...)

	return contacts, nil
}

// ContactLineToRectangle returns the contacts between the line and
// the rectangle (if they exist).
func ContactLineToRectangle(line *Line, rect *Rectangle) ([]Vector, error) {
	if line == nil {
		return nil, fmt.Errorf("the line is nil")
	}

	if rect == nil {
		return nil, fmt.Errorf("the rectangle is nil")
	}

	vertices := rect.Vertices()
	sideAB, err := NewLine(vertices[0], vertices[1])

	if err != nil {
		return nil, err
	}

	sideBC, err := NewLine(vertices[1], vertices[2])

	if err != nil {
		return nil, err
	}

	sideCD, err := NewLine(vertices[2], vertices[3])

	if err != nil {
		return nil, err
	}

	sideAD, err := NewLine(vertices[0], vertices[3])

	if err != nil {
		return nil, err
	}

	sides := []*Line{
		sideAB,
		sideBC,
		sideCD,
		sideAD,
	}
	contacts := make([]Vector, 0)

	for _, side := range sides {
		sideContacts, err := ContactLineToLine(line, side)

		if err != nil {
			return nil, err
		}

		contacts = append(contacts, sideContacts...)
	}

	return contacts, nil
}

// ContactRectangleToCircle returns the contacts between the rectangle and
// the circle (if they exist).
func ContactRectangleToCircle(rect *Rectangle, circle *Circle) ([]Vector, error) {
	if rect == nil {
		return nil, fmt.Errorf("the rectangle is nil")
	}

	if circle == nil {
		return nil, fmt.Errorf("the circle is nil")
	}

	vertices := rect.Vertices()
	sideAB, err := NewLine(vertices[0], vertices[1])

	if err != nil {
		return nil, err
	}

	sideBC, err := NewLine(vertices[1], vertices[2])

	if err != nil {
		return nil, err
	}

	sideCD, err := NewLine(vertices[2], vertices[3])

	if err != nil {
		return nil, err
	}

	sideAD, err := NewLine(vertices[0], vertices[3])

	if err != nil {
		return nil, err
	}

	sides := []*Line{
		sideAB,
		sideBC,
		sideCD,
		sideAD,
	}
	contacts := make([]Vector, 0)

	for _, side := range sides {
		sideContacts, err := ContactLineToCircle(side, circle)

		if err != nil {
			return nil, err
		}

		contacts = append(contacts, sideContacts...)
	}

	return contacts, nil
}

// ContactRectangleToRectangle returns the contacts between two rectangles
// (if they exist).
func ContactRectangleToRectangle(one, other *Rectangle) ([]Vector, error) {
	if one == nil {
		return nil, fmt.Errorf("the first rectangle is nil")
	}

	if other == nil {
		return nil, fmt.Errorf("the second rectangle is nil")
	}

	oneVertices := one.Vertices()
	oneSideAB, err := NewLine(oneVertices[0], oneVertices[1])

	if err != nil {
		return nil, err
	}

	oneSideBC, err := NewLine(oneVertices[1], oneVertices[2])

	if err != nil {
		return nil, err
	}

	oneSideCD, err := NewLine(oneVertices[2], oneVertices[3])

	if err != nil {
		return nil, err
	}

	oneSideAD, err := NewLine(oneVertices[0], oneVertices[3])

	if err != nil {
		return nil, err
	}

	oneSides := []*Line{
		oneSideAB,
		oneSideBC,
		oneSideCD,
		oneSideAD,
	}

	otherVertices := other.Vertices()
	otherSideAB, err := NewLine(otherVertices[0], otherVertices[1])

	if err != nil {
		return nil, err
	}

	otherSideBC, err := NewLine(otherVertices[1], otherVertices[2])

	if err != nil {
		return nil, err
	}

	otherSideCD, err := NewLine(otherVertices[2], otherVertices[3])

	if err != nil {
		return nil, err
	}

	otherSideAD, err := NewLine(otherVertices[0], otherVertices[3])

	if err != nil {
		return nil, err
	}

	otherSides := []*Line{
		otherSideAB,
		otherSideBC,
		otherSideCD,
		otherSideAD,
	}

	contacts := make([]Vector, 0)

	for _, oneSide := range oneSides {
		for _, otherSide := range otherSides {
			sideContacts, err := ContactLineToLine(oneSide, otherSide)

			if err != nil {
				return nil, err
			}

			contacts = append(contacts, sideContacts...)
		}
	}

	return contacts, nil
}

// ContactCircleToCircle returns the contact point between
// two circles (if it exists).
func ContactCircleToCircle(one, other *Circle) ([]Vector, error) {
	if one == nil {
		return nil, fmt.Errorf("the first circle is nil")
	}

	if other == nil {
		return nil, fmt.Errorf("the second circle is nil")
	}

	var (
		r  float64
		R  float64
		cx float64
		cy float64
		Cx float64
		Cy float64
	)

	if one.radius < other.radius {
		r = one.radius
		R = other.radius

		cx = one.center.X
		cy = one.center.Y

		Cx = other.center.X
		Cy = other.center.Y
	} else {
		r = other.radius
		R = one.radius

		cx = other.center.X
		cy = other.center.Y

		Cx = one.center.X
		Cy = one.center.Y
	}

	dx := cx - Cx
	dy := cy - Cy
	d := math.Sqrt(dx*dx + dy*dy)

	// Infinite number of contacts.
	if d < Epsilon && math.Abs(r-R) < Epsilon {
		return []Vector{}, nil
	}

	// No instersection between the circles.
	// No contacts.
	if d < Epsilon {
		return []Vector{}, nil
	}

	// One circle within the other.
	// No contacts.
	if d+r < R || R+r < d {
		return []Vector{}, nil
	}

	p := Vector{
		X: (dx/d)*R + Cx,
		Y: (dy/d)*R + Cy,
	}

	// One intersection point.
	if math.Abs(r+R-d) < Epsilon {
		return []Vector{p}, nil
	}

	// Compute two intersection points.
	c := NewVector(Cx, Cy)
	arg := (r*r - d*d - R*R) / (-2.0 * d * R)
	var angle float64

	if arg >= 1.0 {
		angle = 0.0
	} else if arg <= -1.0 {
		angle = math.Pi
	} else {
		angle = math.Acos(arg)
	}

	altRotate := func(fp, pt Vector, angle float64) Vector {
		x := pt.X - fp.X
		y := pt.Y - fp.Y

		cosA := math.Cos(angle)
		sinA := math.Sin(angle)

		xRot := x*cosA + y*sinA
		yRot := y*cosA - x*sinA

		return NewVector(fp.X+xRot, fp.Y+yRot)
	}
	pointA := altRotate(c, p, angle)
	pointB := altRotate(c, p, -angle)

	return []Vector{pointA, pointB}, nil
}
