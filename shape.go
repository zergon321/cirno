package cirno

// Shape represents a shape in the space.
type Shape interface {
	Center() Vector
	Angle() float64
	AngleRadians() float64
	Move(Vector) Vector
	Rotate(float64) float64
	RotateRadians(float64) float64
	SetPosition(Vector) Vector
	SetAngle(float64) float64
	SetAngleRadians(float64) float64
	ContainsPoint(Vector) bool

	GetIdentity() int32
	SetIdentity(int32)

	GetMask() int32
	SetMask(int32)

	ShouldCollide(Shape) bool

	Data() interface{}
	SetData(data interface{})

	NormalTo(Shape) Vector

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
func (shapes Shapes) Contains(shape Shape) bool {
	_, exists := shapes[shape]

	return exists
}

// Insert adds the given shape in the set of shapes.
func (shapes Shapes) Insert(shape Shape) {
	shapes[shape] = none{}
}

// Remove removes the specified shape from the set.
func (shapes Shapes) Remove(shape Shape) {
	delete(shapes, shape)
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
