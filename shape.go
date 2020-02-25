package cirno

// Shape represents a shape in the space.
type Shape interface {
	Center() Vector
	Move(Vector) Vector
	SetPosition(Vector) Vector
	ContainsPoint(Vector) bool
	GetTag() Tag
	SetTag(Tag)
	ShouldCollide(Shape) bool
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
