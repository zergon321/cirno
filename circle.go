package cirno

// Circle represents a geometric euclidian circle.
type Circle struct {
	center Vector
	radius float64
	tag
	data
	domain
}

// TypeName returns the name of the shape type.
func (c *Circle) TypeName() string {
	return "Circle"
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

// RotateAround rotates the circle around the specified point
// changing the circle's position.
func (c *Circle) RotateAround(angle float64, base Vector) Vector {
	c.center = c.center.RotateAround(angle, base)

	return c.center
}

// RotateAroundRadians rotates the circle around the base point at the
// specified angle in radians.
func (c *Circle) RotateAroundRadians(angle float64, base Vector) Vector {
	c.center = c.center.RotateAroundRadians(angle, base)

	return c.center
}

// SetPosition sets the circle position to the given
// coordinates.
func (c *Circle) SetPosition(pos Vector) Vector {
	c.center = pos

	return c.center
}

// SetAngle does nothing to the circle because there is
// no sense to rotate it.
//
// This method is added just to match the
// Shape interface.
func (c *Circle) SetAngle(angle float64) float64 {
	return 0
}

// SetAngleRadians does nothing to the circle because there is
// no sense to rotate it.
//
// This method is added just to match the
// Shape interface.
func (c *Circle) SetAngleRadians(angle float64) float64 {
	return 0
}

// ContainsPoint detects if the given point is inside the circle.
func (c *Circle) ContainsPoint(point Vector) bool {
	d := c.center.Subtract(point)

	return d.SquaredMagnitude() <= c.radius*c.radius
}

// ОШИБКА ИНТЕРФЕЙСА: возможно указание отрицательного
//  или нулевого значения для радиуса окружности.

// NewCircle create a new circle with the given parameters.
func NewCircle(position Vector, radius float64) *Circle {
	circle := &Circle{}

	circle.center = position
	circle.radius = radius
	circle.treeNodes = []*quadTreeNode{}

	return circle
}
