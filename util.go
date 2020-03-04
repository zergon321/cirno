package cirno

import (
	"fmt"
	"math"
	"reflect"
)

const (
	// RadToDeg is a factor to transfrom radians to degrees.
	RadToDeg float64 = 180.0 / math.Pi
	// DegToRad is a factor to transform degrees to radians.
	DegToRad float64 = math.Pi / 180.0
	// Epsilon is the constant for approximate comparisons.
	Epsilon float64 = 0.000001
)

type none struct{}

// Distance returns the value of distance
// between two points represented as vectors.
func Distance(a, b Vector) float64 {
	return b.Subtract(a).Magnitude()
}

// SquaredDistance returns the value of distance in the power of 2
// between two points represented as vectors.
func SquaredDistance(a, b Vector) float64 {
	return b.Subtract(a).SquaredMagnitude()
}

// boundingBoxesIntersect checks if 2 bounding boxes formed by the
// given points do intersect.
func boundingBoxesIntersect(a0, a1, b0, b1 Vector) bool {
	return a0.X <= b1.X && a1.X >= b0.X && a0.Y <= b1.Y &&
		a1.Y >= b0.Y
}

// Approximate attempts to move the shape in the specified direction
// to detect the closest point until the shape collides other shapes.
func Approximate(shape Shape, diff Vector, shapes Shapes, intensity int, useTags bool) (Vector, error) {
	if intensity < 0 {
		return Zero, fmt.Errorf("The value of intensity must be non-zero")
	}

	step := 1.0 / float64(intensity)
	originalPos := shape.Center()
	prev := originalPos
	shapeType := reflect.TypeOf(shape).Elem()

	for i := 0; i < intensity; i++ {
		current := prev.Add(diff.MultiplyByScalar(step))
		shape.SetPosition(current)
		collisionFound := false

		for other := range shapes {
			otherShapeType := reflect.TypeOf(other).Elem()
			id := shapeType.Name() + "_" + otherShapeType.Name()

			// Make sure collinear lines will collide.
			if id == "Line_Line" {
				line := shape.(*Line)
				otherLine := other.(*Line)
				movement := prev.Subtract(current)

				if line.CollinearTo(otherLine) &&
					collinearLinesWouldCollide(prev, movement, line, otherLine) {
					collisionFound = true
					break
				}
			}

			if ResolveCollision(shape, other, useTags) {
				collisionFound = true
				break
			}
		}

		if collisionFound {
			break
		}

		prev = current
	}

	shape.SetPosition(originalPos)

	return prev, nil
}

// collinearLinesWouldCollide returns true if the first line moved in the specified
// direction from its original position would collide the second line on the way.
func collinearLinesWouldCollide(originalPos, diff Vector, line, otherLine *Line) bool {
	tmpPos := line.Center()
	line.SetPosition(originalPos)
	origP := line.p
	origQ := line.q

	otherTmpPos := otherLine.Center()
	otherLine.Move(diff.MultiplyByScalar(-1))
	movedP := otherLine.p
	movedQ := otherLine.q
	otherLine.SetPosition(otherTmpPos)

	pp := NewLine(movedP, otherLine.p)
	qq := NewLine(movedQ, otherLine.q)

	if IntersectionLineToLine(pp, line) ||
		IntersectionLineToLine(qq, line) {
		return true
	}

	line.SetPosition(tmpPos)
	pp = NewLine(origP, line.p)
	qq = NewLine(origQ, line.q)

	if IntersectionLineToLine(pp, otherLine) ||
		IntersectionLineToLine(qq, otherLine) {
		return true
	}

	return false
}
