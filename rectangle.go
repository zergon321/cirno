package cirno

// Rectangle represents an oriented euclidian rectangle.
type Rectangle struct {
	center  Vector
	extents Vector
	xAxis   Vector
	yAxis   Vector
	angle   float64
	Tag
}

// Center returns the coordinates of the center of the rectangle.
func (r *Rectangle) Center() Vector {
	return r.center
}

// Move moves the rectangle in the specified direction; returns its new position.
func (r *Rectangle) Move(direction Vector) Vector {
	r.center = r.center.Add(direction)

	return r.center
}

// SetPosition sets the position of the rectangle to the given coordinates.
func (r *Rectangle) SetPosition(pos Vector) Vector {
	r.center = pos

	return r.center
}

// ContainsPoint detects if the given point is inside the rectangle.
func (r *Rectangle) ContainsPoint(point Vector) bool {
	localPoint := point.Subtract(r.center)
	theta := -r.angle
	localPoint = localPoint.Rotate(theta)
	localRect := NewRectangle(NewVector(0, 0), r.Width(), r.Height(), 0.0)
	min := localRect.Min()
	max := localRect.Max()

	return min.X <= point.X &&
		min.Y <= point.Y &&
		point.X <= max.X &&
		point.Y <= max.Y
}

// Width returns the width of the rectangle.
func (r *Rectangle) Width() float64 {
	return r.extents.X * 2
}

// Height returns the height of the rectangle.
func (r *Rectangle) Height() float64 {
	return r.extents.Y * 2
}

// Angle returns the angle of the rectangle (in degrees).
func (r *Rectangle) Angle() float64 {
	return r.angle
}

// AngleRadians returns the angle of the rectangle (in radians).
func (r *Rectangle) AngleRadians() float64 {
	return r.angle * DegToRad
}

// Rotate rotates the whole rectangle at the specified angle (in degrees).
func (r *Rectangle) Rotate(angle float64) {
	r.xAxis = r.xAxis.Rotate(angle)
	r.yAxis = r.yAxis.Rotate(angle)

	r.angle += angle
}

// RotateRadians rotates the whole rectangle at the specified angle (in radians).
func (r *Rectangle) RotateRadians(angle float64) {
	r.xAxis = r.xAxis.RotateRadians(angle)
	r.yAxis = r.yAxis.RotateRadians(angle)

	r.angle += angle
}

// Max returns the upper right point of the rectangle with no rotation.
func (r *Rectangle) Max() Vector {
	return r.center.Add(r.xAxis.MultiplyByScalar(r.extents.X).
		Add(r.yAxis.MultiplyByScalar(r.extents.Y)))
}

// Min returns the lower left point of the rectangle with no rotation.
func (r *Rectangle) Min() Vector {
	return r.center.Subtract(r.xAxis.MultiplyByScalar(r.extents.X).
		Add(r.yAxis.MultiplyByScalar(r.extents.Y)))
}

// NewRectangle returns a new rectangle with specified parameters.
func NewRectangle(position Vector, width, height, angle float64) *Rectangle {
	rect := new(Rectangle)

	rect.center = position
	rect.extents = NewVector(width/2.0, height/2.0)
	rect.xAxis = NewVector(1, 0)
	rect.yAxis = NewVector(0, 1)
	rect.Rotate(angle)

	return rect
}
