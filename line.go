package cirno

import (
	"fmt"
	"math"
	"sort"
)

// Line represents a geometric euclidian line segment
// from p to q.
type Line struct {
	p     Vector
	q     Vector
	angle float64
	tag
	data
	domain
}

// TypeName returns the name of the shape type.
func (l *Line) TypeName() string {
	return "Line"
}

// Center returns the coordinates of the middle point
// between p and q.
func (l *Line) Center() Vector {
	return NewVector((l.q.X+l.p.X)/2.0, (l.q.Y+l.p.Y)/2.0)
}

// Angle returns the rotation angle of the line (un degrees).
func (l *Line) Angle() float64 {
	return l.angle
}

// AngleRadians returns the rotation angle of the line (un radians).
func (l *Line) AngleRadians() float64 {
	return l.angle * DegToRad
}

// P returns the starting point of the line.
func (l *Line) P() Vector {
	return l.p
}

// Q returns the ending point of the line.
func (l *Line) Q() Vector {
	return l.q
}

// Move moves the line in the specified direction
// and returns its new position.
func (l *Line) Move(direction Vector) Vector {
	l.q = l.q.Add(direction)
	l.p = l.p.Add(direction)

	return l.Center()
}

// SetPosition sets the position of the line
// to the given coordinates.
func (l *Line) SetPosition(pos Vector) Vector {
	direction := pos.Subtract(l.Center())

	return l.Move(direction)
}

// SetAngle sets the rotation angle of the line
// to the specified value (in degrees).
func (l *Line) SetAngle(angle float64) float64 {
	return l.Rotate(angle - l.angle)
}

// SetAngleRadians sets the rotation angle of the line
// to the specified value (in radians).
func (l *Line) SetAngleRadians(angle float64) float64 {
	return l.RotateRadians(angle - l.angle)
}

// Rotate rotates the line at the
// specified angle (in degrees).
//
// Returns the rotation angle of
// the line (in degrees).
func (l *Line) Rotate(angle float64) float64 {
	center := l.Center()
	pv := l.p.Subtract(center)
	qv := l.q.Subtract(center)

	pv = pv.Rotate(angle)
	qv = qv.Rotate(angle)

	l.p = center.Add(pv)
	l.q = center.Add(qv)

	l.angle += angle
	l.angle = AdjustAngle(l.angle)

	return l.angle
}

// RotateRadians rotates the line at the
// specified angle (in radians).
//
// Returns the rotation angle of
// the line (in radians).
func (l *Line) RotateRadians(angle float64) float64 {
	return l.Rotate(angle*RadToDeg) * DegToRad
}

// RotateAround rotates the line around the base point.
func (l *Line) RotateAround(angle float64, base Vector) Vector {
	center := l.Center()
	cp := l.p.Subtract(center)
	cq := l.q.Subtract(center)

	center = center.RotateAround(angle, base)
	l.p = center.Add(cp)
	l.q = center.Add(cq)

	return center
}

// RotateAroundRadians rotates the line around the base point at the
// angle in radians.
func (l *Line) RotateAroundRadians(angle float64, base Vector) Vector {
	center := l.Center()
	cp := l.p.Subtract(center)
	cq := l.q.Subtract(center)

	center = center.RotateAroundRadians(angle, base)
	l.p = center.Add(cp)
	l.q = center.Add(cq)

	return center
}

// Length returns the length of the line.
func (l *Line) Length() float64 {
	return l.q.Subtract(l.p).Magnitude()
}

// SquaredLength returns the length of the line
// in the power of 2.
func (l *Line) SquaredLength() float64 {
	return l.q.Subtract(l.p).SquaredMagnitude()
}

// ContainsPoint detects if the point lies on the line.
func (l *Line) ContainsPoint(point Vector) bool {
	lTmp := &Line{
		p: Zero(),
		q: l.q.Subtract(l.p),
	}
	pTmp := point.Subtract(l.p)
	r := Cross(lTmp.q, pTmp)
	min, max := l.GetBoundingBox()

	return math.Abs(r) < Epsilon &&
		point.X >= min.X &&
		point.Y >= min.Y &&
		point.X <= max.X &&
		point.Y <= max.Y
}

// Orientation returns 0 if the point is collinear to the line,
// 1 if orientation is clockwise,
// -1 if orientation is counter-clockwise.
func (l *Line) Orientation(point Vector) int {
	val := (l.q.Y-l.p.Y)*(point.X-l.q.X) -
		(l.q.X-l.p.X)*(point.Y-l.q.Y)

	if val == 0 {
		return 0
	}

	if val > 0 {
		return 1
	}

	return -1
}

// GetBoundingBox returns the bounding box for the line.
func (l *Line) GetBoundingBox() (Vector, Vector) {
	min := NewVector(math.Min(l.p.X, l.q.X), math.Min(l.p.Y, l.q.Y))
	max := NewVector(math.Max(l.p.X, l.q.X), math.Max(l.p.Y, l.q.Y))

	return min, max
}

// CollinearTo returns true if the lines are collinear,
// and false otherwise.
func (l *Line) CollinearTo(other *Line) (bool, error) {
	if other == nil {
		return false, fmt.Errorf("the line is nil")
	}

	lVec := l.q.Subtract(l.p)
	otherVec := other.q.Subtract(other.p)

	return lVec.CollinearTo(otherVec)
}

// SameLineWith returns true if two line segments
// lie on the same line.
func (l *Line) SameLineWith(other *Line) (bool, error) {
	if other == nil {
		return false, fmt.Errorf("the line is nil")
	}

	projP := l.ProjectPoint(other.p)
	projQ := l.ProjectPoint(other.q)

	return projP.ApproximatelyEqual(other.p) &&
		projQ.ApproximatelyEqual(other.q), nil
}

// ParallelTo checks if two line segments are collinear but
// don't lie on the same line.
func (l *Line) ParallelTo(other *Line) (bool, error) {
	if other == nil {
		return false, fmt.Errorf("the line is nil")
	}

	collinear, err := l.CollinearTo(other)

	if err != nil {
		return false, err
	}

	sameLine, err := l.SameLineWith(other)

	if err != nil {
		return false, err
	}

	return collinear && !sameLine, nil
}

// ProjectPoint returns the projection of the point
// onto the line.
func (l *Line) ProjectPoint(point Vector) Vector {
	t := ((point.X-l.p.X)*(l.q.X-l.p.X) + (point.Y-l.p.Y)*(l.q.Y-l.p.Y)) /
		((l.q.X-l.p.X)*(l.q.X-l.p.X) + (l.q.Y-l.p.Y)*(l.q.Y-l.p.Y))

	return NewVector(l.p.X+t*(l.q.X-l.p.X), l.p.Y+t*(l.q.Y-l.p.Y))
}

// isPointRightOfLine returns true if the given point
// is located to the right of the line, and false otherwise.
func (l *Line) isPointRightOfLine(p Vector) (bool, error) {
	lTmp, err := NewLine(Zero(), l.q.Subtract(l.p))

	if err != nil {
		return false, err
	}

	pTmp := p.Subtract(l.p)

	return Cross(lTmp.q, pTmp) < 0, nil
}

// touchesOrCrosses returns true if the line touches
// or crosses the otehr line, and false otherwise.
func (l *Line) touchesOrCrosses(other *Line) (bool, error) {
	if other == nil {
		return false, fmt.Errorf("the line is nil")
	}

	pRight, err := l.isPointRightOfLine(other.p)

	if err != nil {
		return false, err
	}

	qRight, err := l.isPointRightOfLine(other.q)

	if err != nil {
		return false, err
	}

	return l.ContainsPoint(other.p) ||
		l.ContainsPoint(other.q) ||
		(pRight != qRight), nil
}

// LinesDistance returns the shortest distance
// between two lines.
func LinesDistance(a, b *Line) (float64, error) {
	if a == nil {
		return math.NaN(),
			fmt.Errorf("the first line is nil")
	}

	if b == nil {
		return math.NaN(),
			fmt.Errorf("the second line is nil")
	}

	apProj := b.ProjectPoint(a.p)
	aqProj := b.ProjectPoint(a.q)
	bpProj := a.ProjectPoint(b.p)
	bqProj := a.ProjectPoint(b.q)
	distances := make([]float64, 0)

	if !b.ContainsPoint(apProj) {
		distances = append(distances,
			math.Min(Distance(a.p, b.p), Distance(a.p, b.q)))
	}

	if !b.ContainsPoint(aqProj) {
		distances = append(distances,
			math.Min(Distance(a.q, b.p), Distance(a.q, b.q)))
	}

	if !a.ContainsPoint(bpProj) {
		distances = append(distances,
			math.Min(Distance(b.p, a.p), Distance(b.p, a.q)))
	}

	if !a.ContainsPoint(bqProj) {
		distances = append(distances,
			math.Min(Distance(b.q, a.p), Distance(b.q, a.q)))
	}

	sort.Slice(distances, func(i, j int) bool {
		return distances[i] < distances[j]
	})

	return distances[0], nil
}

// NewLine returns a new line segment with the given parameters.
func NewLine(p Vector, q Vector) (*Line, error) {
	if Distance(p, q) < Epsilon {
		return nil, fmt.Errorf(
			"the length of the line must be positive")
	}

	line := &Line{
		p: p,
		q: q,
	}

	pq := line.q.Subtract(line.p)
	angle, err := Angle(pq, Right())

	if err != nil {
		return nil, err
	}

	if angle < 0 {
		line.angle = 360 + angle
	} else {
		line.angle = angle
	}

	line.treeNodes = []*quadTreeNode{}

	return line, nil
}
