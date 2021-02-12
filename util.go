package cirno

import (
	"fmt"
	"math"
)

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
	if shape == nil {
		return Zero(), -1, nil,
			fmt.Errorf("the shape being approximated is nil")
	}

	if shape == nil {
		return Zero(), -1, nil,
			fmt.Errorf("the set of shapes is nil")
	}

	var foundShape Shape

	if intensity < 0 {
		return Zero(), 0, nil,
			fmt.Errorf("the value of intensity must be positive")
	}

	step := 1.0 / float64(intensity)
	originalPos := shape.Center()
	originalAngle := shape.Angle()
	prevPos := originalPos
	prevAngle := originalAngle

	for i := 0; i < intensity; i++ {
		currentPos := prevPos.Add(moveDiff.MultiplyByScalar(step))
		currentAngle := prevAngle + turnDiff*step
		shape.SetPosition(currentPos)
		shape.SetAngle(currentAngle)
		collisionFound := false

		for other := range shapes {
			id := shape.TypeName() + "_" + other.TypeName()

			// Make sure lines will collide.
			if id == "Line_Line" {
				line := shape.(*Line)
				otherLine := other.(*Line)
				linesCollinear, err := line.CollinearTo(otherLine)

				if intensity < 0 {
					return Zero(), 0, nil, err
				}

				if linesCollinear {
					// Compare line tags.
					shouldCollide, err := line.ShouldCollide(otherLine)

					if err != nil {
						return Zero(), 0, nil, err
					}

					if useTags && !shouldCollide {
						continue
					}

					// Vice versa.
					movement := currentPos.Subtract(prevPos)
					turn := currentAngle - prevAngle
					linesWouldIntersect, err := linesWouldCollide(
						prevPos, prevAngle, movement, turn, line, otherLine)

					if err != nil {
						return Zero(), -1, nil, err
					}

					if linesWouldIntersect {
						collisionFound = true
						foundShape = otherLine

						break
					}
				}
			}

			overlapped, err := ResolveCollision(shape, other, useTags)

			if err != nil {
				return Zero(), -1, nil, err
			}

			if overlapped {
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
		return Zero(), 0.0, nil, fmt.Errorf("couldn't approximate the shape")
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
func linesWouldCollide(originalPos Vector, originalAngle float64, moveDiff Vector, turnDiff float64, line, otherLine *Line) (bool, error) {
	if line == nil {
		return false, fmt.Errorf("the first line is nil")
	}

	if otherLine == nil {
		return false, fmt.Errorf("the second line is nil")
	}

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

	pp := &Line{
		p: movedP,
		q: otherLine.p,
	}
	qq := &Line{
		p: movedQ,
		q: otherLine.q,
	}

	ppIntersects, err := IntersectionLineToLine(pp, line)

	if err != nil {
		return false, err
	}

	qqIntersects, err := IntersectionLineToLine(qq, line)

	if err != nil {
		return false, err
	}

	if ppIntersects ||
		qqIntersects {
		return true, nil
	}

	line.SetPosition(tmpPos)
	line.SetAngle(tmpAngle)

	pp = &Line{
		p: origP,
		q: line.p,
	}
	qq = &Line{
		p: origQ,
		q: line.q,
	}

	ppIntersects, err = IntersectionLineToLine(pp, otherLine)

	if err != nil {
		return false, err
	}

	qqIntersects, err = IntersectionLineToLine(qq, otherLine)

	if err != nil {
		return false, err
	}

	linesIntersect, err := IntersectionLineToLine(line, otherLine)

	if err != nil {
		return false, err
	}

	if ppIntersects ||
		qqIntersects ||
		linesIntersect {
		return true, nil
	}

	return false, nil
}
