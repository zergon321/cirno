package cirno

import (
	"fmt"
	"math"
)

// Commonly used normalized vectors.
var (
	Up    = Vector{X: 0, Y: 1}
	Down  = Vector{X: 0, Y: -1}
	Left  = Vector{X: -1, Y: 0}
	Right = Vector{X: 1, Y: 0}
	Zero  = Vector{X: 0, Y: 0}
)

// Vector represents a 2-dimensional vector.
type Vector struct {
	X float64
	Y float64
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
func (v Vector) CollinearTo(other Vector) bool {
	return AdjustAngle(Angle(v, other)) < CollinearityThreshold
}

// Project returns vector v projected onto the axis.
func (v Vector) Project(axis Vector) Vector {
	factor := Dot(v, axis) / Dot(axis, axis)

	return axis.MultiplyByScalar(factor)
}

// Normalize returns a normalized vector with magnitude of 1.
func (v Vector) Normalize() Vector {
	length := v.Magnitude()

	return NewVector(v.X/length, v.Y/length)
}

// ApproximatelyEqual returns true if the vector is
// approximately equal to the other one, and false
// otherwise.
func (v Vector) ApproximatelyEqual(other Vector) bool {
	fmt.Println("dx:", math.Abs(v.X-other.X))
	fmt.Println("dy:", math.Abs(v.Y-other.Y))

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
func Angle(a, b Vector) float64 {
	return AngleRadians(a, b) * RadToDeg
}

// AngleRadians returns the angle between two vectors (in radians).
func AngleRadians(a, b Vector) float64 {
	cosine := Dot(a, b) / (a.Magnitude() * b.Magnitude())

	if cosine > 1.0 {
		cosine = 1.0
	} else if cosine < -1.0 {
		cosine = -1.0
	}

	return math.Acos(cosine)
}

// NewVector returns a new vector with the given coordinates.
func NewVector(x, y float64) Vector {
	return Vector{X: x, Y: y}
}
