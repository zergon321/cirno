package cirno

import "fmt"

// Shape represents a shape in the space.
type Shape interface {
	// Common methods.
	TypeName() string
	Center() Vector
	Angle() float64
	AngleRadians() float64
	Move(Vector) Vector
	Rotate(float64) float64
	RotateRadians(float64) float64
	RotateAround(float64, Vector) Vector
	RotateAroundRadians(float64, Vector) Vector
	SetPosition(Vector) Vector
	SetAngle(float64) float64
	SetAngleRadians(float64) float64
	ContainsPoint(Vector) bool
	NormalTo(Shape) (Vector, error)

	// Tag-related methods.
	GetIdentity() int32
	SetIdentity(int32)
	GetMask() int32
	SetMask(int32)
	ShouldCollide(Shape) (bool, error)

	// Data-related methods.
	Data() interface{}
	SetData(data interface{})

	// Domain-related methods.
	nodes() []*quadTreeNode
	addNodes(...*quadTreeNode)
	containsNode(*quadTreeNode) bool
	removeNodes(...*quadTreeNode)
	clearNodes()
}

// Shapes represents a list of shapes.
type Shapes map[Shape]none

// Contains checks if the list of shapes contains
// the given shape.
func (shapes Shapes) Contains(shape Shape) (bool, error) {
	if shape == nil {
		return false, fmt.Errorf("the shape is nil")
	}

	_, exists := shapes[shape]

	return exists, nil
}

// Insert adds the given shape in the set of shapes.
func (shapes Shapes) Insert(shapesToInsert ...Shape) error {
	for _, shape := range shapesToInsert {
		if shape == nil {
			return fmt.Errorf("the shape is nil")
		}

		shapes[shape] = none{}
	}

	return nil
}

// Merge adds all the shapes from
// the other set to the current set.
func (shapes Shapes) Merge(otherShapes Shapes) error {
	if otherShapes == nil {
		return fmt.Errorf("the shape set is nil")
	}

	for shape := range otherShapes {
		shapes[shape] = none{}
	}

	return nil
}

// Remove removes the specified shape from the set.
func (shapes Shapes) Remove(shape Shape) error {
	if shape == nil {
		return fmt.Errorf("the shape is nil")
	}

	delete(shapes, shape)

	return nil
}

// Items returns the list all the shapes
// from the set.
func (shapes Shapes) Items() []Shape {
	items := make([]Shape, 0)

	for shape := range shapes {
		items = append(items, shape)
	}

	return items
}

// Copy returns a new hash set with same shapes.
func (shapes Shapes) Copy() Shapes {
	setCopy := make(Shapes, len(shapes))

	for shape := range shapes {
		setCopy[shape] = none{}
	}

	return setCopy
}
