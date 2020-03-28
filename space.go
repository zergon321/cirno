package cirno

import (
	"fmt"
	"reflect"
)

// Space represents a geometric space
// with shapes within it.
type Space struct {
	min     Vector
	max     Vector
	shapes  Shapes
	tree    *quadTree
	useTags bool
}

// Max returns the max point of the space.
func (space *Space) Max() Vector {
	return space.max
}

// Min returns the min point of the space.
func (space *Space) Min() Vector {
	return space.min
}

// UseTags indicates whether the space relies on
// tags for collision detection.
func (space *Space) UseTags() bool {
	return space.useTags
}

// InBounds detects if the center the given shape
// is within the space bounds.
func (space *Space) InBounds(shape Shape) bool {
	pos := shape.Center()

	return pos.X >= space.min.X ||
		pos.Y >= space.min.Y || pos.X <= space.max.X ||
		pos.Y <= space.max.Y
}

// AdjustShapePosition changes the position of the shape
// if it's out of bounds.
func (space *Space) AdjustShapePosition(shape Shape) {
	pos := shape.Center()

	if pos.X < space.min.X {
		pos = shape.SetPosition(NewVector(space.min.X, pos.Y))
	}

	if pos.Y < space.min.Y {
		pos = shape.SetPosition(NewVector(pos.X, space.min.Y))
	}

	if pos.X > space.max.X {
		pos = shape.SetPosition(NewVector(space.max.X, pos.Y))
	}

	if pos.Y > space.max.Y {
		pos = shape.SetPosition(NewVector(pos.X, space.max.Y))
	}
}

// Add adds a new shape in the space.
func (space *Space) Add(shapes ...Shape) error {
	for _, shape := range shapes {
		if shape != nil {
			if !space.InBounds(shape) {
				return fmt.Errorf("The shape is out of bounds")
			}

			space.shapes.Insert(shape)
			_, err := space.tree.insert(shape)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Remove removes the shape from the space.
func (space *Space) Remove(shapes ...Shape) error {
	for _, shape := range shapes {
		space.shapes.Remove(shape)
		err := space.tree.remove(shape)

		if err != nil {
			return err
		}
	}

	return nil
}

// Contains returns true if the shape is within the space,
// and false otherwise.
func (space *Space) Contains(shape Shape) bool {
	return space.shapes.Contains(shape)
}

// Clear removes all shapes from
// the space.
func (space *Space) Clear() error {
	space.shapes = make(Shapes, 0)

	return space.tree.clear()
}

// Shapes returns the list of all shapes
// within the space.
func (space *Space) Shapes() Shapes {
	return space.shapes.Copy()
}

// Update should be called on the shape whenever
// it's moved within the space.
func (space *Space) Update(shape Shape) (map[Vector]Shapes, error) {
	if !space.Contains(shape) {
		return nil, fmt.Errorf("The space doesn't contain the given shape")
	}

	err := space.tree.remove(shape)

	if err != nil {
		return nil, err
	}

	nodes, err := space.tree.insert(shape)

	if err != nil {
		return nil, err
	}

	shapes := make(map[Vector]Shapes)

	for _, node := range nodes {
		shapes[node.boundary.center] = node.shapes.Copy()
	}

	return shapes, nil
}

// CollidingShapes returns the dictionary where key
// is a shape and value is the set of shapes
// colliding with the key shape.
func (space *Space) CollidingShapes() (map[Shape]Shapes, error) {
	collidingShapes := make(map[Shape]Shapes)
	shapeGroups := space.tree.shapeGroups()

	for _, area := range shapeGroups {
		shapes := area.Items()

		for i, shape := range shapes {
			if _, ok := collidingShapes[shape]; !ok {
				collidingShapes[shape] = make(Shapes, 0)
			}

			for _, otherShape := range shapes[i+1:] {
				if ResolveCollision(shape, otherShape, space.useTags) {
					collidingShapes[shape].Insert(otherShape)

					if _, ok := collidingShapes[otherShape]; !ok {
						collidingShapes[otherShape] = make(Shapes, 0)
					}

					collidingShapes[otherShape].Insert(shape)
				}
			}

			if len(collidingShapes[shape]) == 0 {
				delete(collidingShapes, shape)
			}
		}
	}

	return collidingShapes, nil
}

// CollidingWith returns the set of shapes colliding with the given shape.
func (space *Space) CollidingWith(shape Shape) (Shapes, error) {
	shapes := make(Shapes, 0)
	nodes, err := space.tree.search(shape)

	if err != nil {
		return nil, err
	}

	for _, area := range nodes {
		for item := range area.shapes {
			if item != shape && ResolveCollision(shape, item, space.useTags) {
				shapes.Insert(item)
			}
		}
	}

	return shapes, nil
}

// WouldBeColliding checks if the given would shape would be colliding any other shape
// if it moved in the specified direction.
func (space *Space) WouldBeColliding(shape Shape, moveDiff Vector, turnDiff float64) (Shapes, error) {
	shapes := make(Shapes, 0)
	originalPos := shape.Center()
	originalAngle := shape.Angle()
	shapeType := reflect.TypeOf(shape).Elem()

	// Get all the nodes where the shape is located
	// before movement.
	areas, err := space.tree.search(shape)

	if err != nil {
		return nil, err
	}

	shape.Move(moveDiff)
	shape.Rotate(turnDiff)
	// Make sure the shape is in bounds.
	space.AdjustShapePosition(shape)
	// Update the shape's position in the quad tree.
	nodes, err := space.Update(shape)

	if err != nil {
		return nil, err
	}

	// Add shapes from the previous nodes.
	for _, area := range areas {
		nodes[area.boundary.center] = area.shapes.Copy()
	}

	// Search for collisions in the nodes
	// the shape belongs to.
	for _, area := range nodes {
		for item := range area {
			if item == shape {
				continue
			}

			itemType := reflect.TypeOf(item).Elem()
			id := shapeType.Name() + "_" + itemType.Name()

			// Make sure lines will collide.
			if id == "Line_Line" {
				lineShape := shape.(*Line)
				lineItem := item.(*Line)

				if linesWouldCollide(originalPos, originalAngle, moveDiff, turnDiff, lineShape, lineItem) {
					shapes.Insert(lineItem)

					continue
				}
			}

			if ResolveCollision(shape, item, space.useTags) {
				shapes.Insert(item)
			}
		}
	}

	// Move the shape back.
	shape.SetPosition(originalPos)
	shape.SetAngle(originalAngle)
	_, err = space.Update(shape)

	if err != nil {
		return nil, err
	}

	return shapes, nil
}

// NewSpace creates a new empty space with the given parameters.
func NewSpace(subdivisionFactor, shapesInArea int, width, height float64, min, max Vector, useTags bool) (*Space, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("Space must have positive values for width and height")
	}

	if min.X >= max.X {
		return nil, fmt.Errorf("The value of min X is incorrect")
	}

	if min.Y >= max.Y {
		return nil, fmt.Errorf("The value of min Y is incorrect")
	}

	space := new(Space)
	space.min = min
	space.max = max
	space.useTags = useTags
	space.shapes = make(Shapes, 0)
	tree, err := newQuadTree(NewRectangle(NewVector(0, 0), width, height, 0.0),
		subdivisionFactor, shapesInArea)

	if err != nil {
		return nil, err
	}

	space.tree = tree

	return space, nil
}
