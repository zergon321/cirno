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
func Approximate(shape Shape, moveDiff Vector, turnDiff float64, shapes Shapes, intensity int, useTags bool) (Vector, float64, error) {
	if intensity < 0 {
		return Zero, 0, fmt.Errorf("The value of intensity must be non-zero")
	}

	step := 1.0 / float64(intensity)
	originalPos := shape.Center()
	originalAngle := shape.Angle()
	prevPos := originalPos
	prevAngle := originalAngle
	shapeType := reflect.TypeOf(shape).Elem()

	for i := 0; i < intensity; i++ {
		currentPos := prevPos.Add(moveDiff.MultiplyByScalar(step))
		currentAngle := prevAngle + turnDiff*step
		shape.SetPosition(currentPos)
		shape.SetAngle(currentAngle)
		collisionFound := false

		for other := range shapes {
			otherShapeType := reflect.TypeOf(other).Elem()
			id := shapeType.Name() + "_" + otherShapeType.Name()

			// Make sure lines will collide.
			if id == "Line_Line" {
				line := shape.(*Line)
				otherLine := other.(*Line)
				// Vice versa.
				movement := currentPos.Subtract(prevPos)
				turn := currentAngle - prevAngle

				if linesWouldCollide(prevPos, prevAngle, movement, turn, line, otherLine) {
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

		prevPos = currentPos
		prevAngle = currentAngle
	}

	shape.SetPosition(originalPos)
	shape.SetAngle(originalAngle)

	return prevPos, prevAngle, nil
}

// AdjustAngle adjusts the value of the angle so it
// is bettween 0 and 360.
func AdjustAngle(angle float64) float64 {
	// Adjust the angle so its value is between 0 and 360.
	if angle >= 360 {
		angle = angle - float64(int64(angle/360))*360
	} else if angle < 0 {
		if angle <= -360 {
			angle = angle - float64(int64(angle/360))*360
		}

		angle += 360

		if angle >= 360 {
			angle = angle - float64(int64(angle/360))*360
		}
	}

	return angle
}

// linesWouldCollide returns true if the first line moved in the specified
// direction from its original position would collide the second line on the way.
func linesWouldCollide(originalPos Vector, originalAngle float64, moveDiff Vector, turnDiff float64, line, otherLine *Line) bool {
	tmpPos := line.Center()
	tmpAngle := line.Angle()
	line.SetPosition(originalPos)
	line.SetAngle(originalAngle)
	origP := line.p
	origQ := line.q

	otherTmpPos := otherLine.Center()
	otherTmpAngle := otherLine.Angle()
	otherLine.Move(moveDiff.MultiplyByScalar(-1))
	otherLine.Rotate(-turnDiff)
	movedP := otherLine.p
	movedQ := otherLine.q
	otherLine.SetPosition(otherTmpPos)
	otherLine.SetAngle(otherTmpAngle)

	pp := NewLine(movedP, otherLine.p)
	qq := NewLine(movedQ, otherLine.q)

	if IntersectionLineToLine(pp, line) ||
		IntersectionLineToLine(qq, line) {
		return true
	}

	line.SetPosition(tmpPos)
	line.SetAngle(tmpAngle)
	pp = NewLine(origP, line.p)
	qq = NewLine(origQ, line.q)

	if IntersectionLineToLine(pp, otherLine) ||
		IntersectionLineToLine(qq, otherLine) ||
		IntersectionLineToLine(line, otherLine) {
		return true
	}

	return false
}
