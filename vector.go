package cirno

import (
	"fmt"
	"math"
)

// Vector represents a 2-dimensional vector.
type Vector struct {
	X float64
	Y float64
}

// Up is the vector {0; 1}.
func Up() Vector {
	return Vector{X: 0, Y: 1}
}

// Down is the vector {0; -1}.
func Down() Vector {
	return Vector{X: 0, Y: -1}
}

// Left is the vector {-1; 0}.
func Left() Vector {
	return Vector{X: -1, Y: 0}
}

// Right is the vector {1; 0}.
func Right() Vector {
	return Vector{X: 1, Y: 0}
}

// Zero is the vector {0; 0}.
func Zero() Vector {
	return Vector{X: 0, Y: 0}
}

// MultiplyByScalar returns the vector multiplied by the specified scalar.
func (v Vector) MultiplyByScalar(scalar float64) Vector {
	return NewVector(v.X*scalar, v.Y*scalar)
}

// Add returns the sum of two vectors.
func (v Vector) Add(other Vector) Vector {
	return NewVector(v.X+other.X, v.Y+other.Y)
}

// Subtract returns the difference of two vectors.
func (v Vector) Subtract(other Vector) Vector {
	return NewVector(v.X-other.X, v.Y-other.Y)
}

// MultiplyBy multiplies components of the first vector
// by the components of the other vector.
func (v Vector) MultiplyBy(other Vector) Vector {
	return NewVector(v.X*other.X, v.Y*other.Y)
}

// Magnitude returns the length of the vector.
func (v Vector) Magnitude() float64 {
	return math.Hypot(v.X, v.Y)
}

// SquaredMagnitude returns the vectors's magnitude
// in the power of 2.
func (v Vector) SquaredMagnitude() float64 {
	return v.X*v.X + v.Y*v.Y
}

// Rotate returns the vector rotated at an angle of n degrees.
func (v Vector) Rotate(angle float64) Vector {
	return v.RotateRadians(angle * DegToRad)
}

// RotateRadians returns the vector totated at an angle of n radians.
func (v Vector) RotateRadians(angle float64) Vector {
	cosA := math.Cos(angle)
	sinA := math.Sin(angle)

	return NewVector(v.X*cosA-v.Y*sinA, v.X*sinA+v.Y*cosA)
}

// RotateAround returns the vector rotated at an angle of n degrees
// around the base vector.
func (v Vector) RotateAround(angle float64, base Vector) Vector {
	return v.RotateAroundRadians(angle*DegToRad, base)
}

// RotateAroundRadians returns the vector totated at an angle of n radians
// around the base vector.
func (v Vector) RotateAroundRadians(angle float64, base Vector) Vector {
	x := v.X - base.X
	y := v.Y - base.Y

	cosA := math.Cos(angle)
	sinA := math.Sin(angle)

	xRot := x*cosA - y*sinA
	yRot := y*cosA + x*sinA

	return NewVector(base.X+xRot, base.Y+yRot)
}

// CollinearTo indicates if the given vector is collinear to the
// other vector.
func (v Vector) CollinearTo(other Vector) (bool, error) {
	angle, err := Angle(v, other)

	if err != nil {
		return false, err
	}

	return AdjustAngle(angle) < CollinearityThreshold, nil
}

// Project returns vector v projected onto the axis.
func (v Vector) Project(axis Vector) Vector {
	factor := Dot(v, axis) / Dot(axis, axis)

	return axis.MultiplyByScalar(factor)
}

// Normalize returns a normalized vector with magnitude of 1.
func (v Vector) Normalize() (Vector, error) {
	if math.Abs(v.X) < Epsilon && math.Abs(v.Y) < Epsilon {
		return Zero(),
			fmt.Errorf("tried to normalize zero vector")
	}

	length := v.Magnitude()

	return NewVector(v.X/length, v.Y/length), nil
}

// ApproximatelyEqual returns true if the vector is
// approximately equal to the other one, and false
// otherwise.
func (v Vector) ApproximatelyEqual(other Vector) bool {
	return math.Abs(v.X-other.X) < Epsilon &&
		math.Abs(v.Y-other.Y) < Epsilon
}

// Dot returns the dot product of two vectors.
func Dot(a, b Vector) float64 {
	return a.X*b.X + a.Y*b.Y
}

// Cross returns the cross product of two vectors.
func Cross(a, b Vector) float64 {
	return a.X*b.Y - b.X*a.Y
}

// Angle returns the angle between two vectors (in degrees).
func Angle(a, b Vector) (float64, error) {
	angle, err := AngleRadians(a, b)

	if err != nil {
		return math.NaN(), err
	}

	return angle * RadToDeg, nil
}

// AngleRadians returns the angle between two vectors (in radians).
func AngleRadians(a, b Vector) (float64, error) {
	magA := a.Magnitude()
	magB := b.Magnitude()

	if magA < Epsilon || magB < Epsilon {
		return math.NaN(),
			fmt.Errorf("one of the vectors is zero")
	}

	cosine := Dot(a, b) / (magA * magB)

	if cosine > 1.0 {
		cosine = 1.0
	} else if cosine < -1.0 {
		cosine = -1.0
	}

	return math.Acos(cosine), nil
}

// PerpendicularClockwise returns a new vector
// perpendicular to the given one and turned in
// a clockwise direction relatively to it.
func (v Vector) PerpendicularClockwise() Vector {
	return NewVector(v.Y, -v.X)
}

// PerpendicularCounterClockwise returns a new vector
// perpendicular to the given one and turned in
// a clockwise direction relatively to it.
func (v Vector) PerpendicularCounterClockwise() Vector {
	return NewVector(-v.Y, v.X)
}

// NewVector returns a new vector with the given coordinates.
func NewVector(x, y float64) Vector {
	return Vector{X: x, Y: y}
}
