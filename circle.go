package cirno

// Circle represents a geometric euclidian circle.
type Circle struct {
	center Vector
	radius float64
}

// Center returns the coordinates of the center
// of the circle.
func (c *Circle) Center() Vector {
	return c.center
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
