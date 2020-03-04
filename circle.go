package cirno

// Circle represents a geometric euclidian circle.
type Circle struct {
	center Vector
	radius float64
	Tag
}

// Center returns the coordinates of the center
// of the circle.
func (c *Circle) Center() Vector {
	return c.center
}

// Angle doesn't return any valuable
// data because there is no sense to
// rotate the circle.
//
// This method is added just to match the
// Shape interface.
func (c *Circle) Angle() float64 {
	return 0
}

// AngleRadians doesn't return any valuable
// data because there is no sense to
// rotate the circle.
//
// This method is added just to match the
// Shape interface.
func (c *Circle) AngleRadians() float64 {
	return 0
}

// Radius returns the radius og the circle.
func (c *Circle) Radius() float64 {
	return c.radius
}

// Move moves the circle in the specified direction.
func (c *Circle) Move(direction Vector) Vector {
	c.center = c.center.Add(direction)

	return c.center
}

// Rotate does nothing to the circle because there
// is no sense to rotate it.
//
// This method is added just to match the
// Shape interface.
func (c *Circle) Rotate(angle float64) float64 {
	return 0
}

// RotateRadians does nothing to the circle because there
// is no sense to rotate it.
//
// This method is added just to match the
// Shape interface.
func (c *Circle) RotateRadians(angle float64) float64 {
	return 0
}

// SetPosition sets the circle position to the given
// coordinates.
func (c *Circle) SetPosition(pos Vector) Vector {
	c.center = pos

	return c.center
}

// ContainsPoint detects if the given point is inside the circle.
func (c *Circle) ContainsPoint(point Vector) bool {
	d := c.center.Subtract(point)

	return d.SquaredMagnitude() <= c.radius*c.radius
}

// NewCircle create a new circle with the given parameters.
func NewCircle(position Vector, radius float64) *Circle {
	return &Circle{
		center: position,
		radius: radius,
	}
}
