package cirno

import (
	"fmt"
	"math"
	"reflect"
)

// ОШИБКА ОБРАЩЕНИЯ К ДАННЫМ: ни в одной из функций данного файла
// не происходит проверка аргумента на nil (пустой указатель).

const (
	// RadToDeg is a factor to transfrom radians to degrees.
	RadToDeg float64 = 180.0 / math.Pi
	// DegToRad is a factor to transform degrees to radians.
	DegToRad float64 = math.Pi / 180.0
	// Epsilon is the constant for approximate comparisons.
	Epsilon float64 = 0.01
	// CollinearityThreshold is the constant to detect if two vectors
	// are effectively collinear.
	CollinearityThreshold float64 = 10.0
)

type none struct{}

// Sign returns the sign of the number
// according to Epsilon.
func Sign(number float64) float64 {
	if number < 0 {
		return -1
	} else if number < Epsilon {
		return 0
	} else {
		return 1
	}
}

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

// Approximate attempts to move the shape in the specified direction
// to detect the closest point until the shape collides other shapes.
func Approximate(shape Shape, moveDiff Vector, turnDiff float64, shapes Shapes, intensity int, useTags bool) (Vector, float64, Shape, error) {
	var foundShape Shape

	if intensity < 0 {
		return Zero, 0, nil, fmt.Errorf("The value of intensity must be non-zero")
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

				if line.CollinearTo(otherLine) {
					// Compare line tags.
					if useTags && !line.ShouldCollide(otherLine) {
						continue
					}

					// Vice versa.
					movement := currentPos.Subtract(prevPos)
					turn := currentAngle - prevAngle

					if linesWouldCollide(prevPos, prevAngle, movement, turn, line, otherLine) {
						collisionFound = true
						foundShape = otherLine

						break
					}
				}
			}

			if ResolveCollision(shape, other, useTags) {
				collisionFound = true
				foundShape = other

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

	if math.IsNaN(prevPos.X) || math.IsNaN(prevPos.Y) {
		return Zero, 0.0, nil, fmt.Errorf("Couldn't approximate the shape")
	}

	return prevPos, prevAngle, foundShape, nil
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
