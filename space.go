package cirno

import (
	"fmt"
	"reflect"

	"github.com/golang-collections/collections/queue"
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

// Cells returns all the cells the space is subdivided to.
func (space *Space) Cells() map[*Rectangle]Shapes {
	cells := map[*Rectangle]Shapes{}

	for leaf := range space.tree.leaves {
		cell := leaf.boundary.toRectangle()

		cells[cell] = leaf.shapes.Copy()
	}

	return cells
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
func (space *Space) InBounds(shape Shape) (bool, error) {
	if shape == nil {
		return false, fmt.Errorf("the shape is nil")
	}

	pos := shape.Center()

	return pos.X >= space.min.X ||
		pos.Y >= space.min.Y || pos.X <= space.max.X ||
		pos.Y <= space.max.Y, nil
}

// AdjustShapePosition changes the position of the shape
// if it's out of bounds.
func (space *Space) AdjustShapePosition(shape Shape) error {
	if shape == nil {
		return fmt.Errorf("the shape is nil")
	}

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

	return nil
}

// Add adds a new shape in the space.
func (space *Space) Add(shapes ...Shape) error {
	for _, shape := range shapes {
		if shape == nil {
			return fmt.Errorf("the shape is nil")
		}

		inBounds, err := space.InBounds(shape)

		if err != nil {
			return err
		}

		if !inBounds {
			return fmt.Errorf("the shape is out of bounds")
		}

		space.shapes.Insert(shape)
		_, err = space.tree.insert(shape)

		if err != nil {
			return err
		}
	}

	return nil
}

// Remove removes the shape from the space.
func (space *Space) Remove(shapes ...Shape) error {
	for _, shape := range shapes {
		if shape == nil {
			return fmt.Errorf("the shape is nil")
		}

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
func (space *Space) Contains(shape Shape) (bool, error) {
	if shape == nil {
		return false, fmt.Errorf("the shape is nil")
	}

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

// Update should be called on the shape
// whenever it's moved within the space.
func (space *Space) Update(shape Shape) (map[Vector]Shapes, error) {
	if shape == nil {
		return nil, fmt.Errorf("the shape is nil")
	}

	contains, err := space.Contains(shape)

	if err != nil {
		return nil, err
	}

	if !contains {
		return nil, fmt.Errorf("the space doesn't contain the given shape")
	}

	// Remove the shape from all the nodes that don't contain it
	// anymore and remove all these nodes from the shape's domain.
	nodesToRemove := []*quadTreeNode{}

	for _, node := range shape.nodes() {
		overlapped, err := node.boundary.collidesShape(shape)

		if err != nil {
			return nil, err
		}

		if !overlapped {
			nodesToRemove = append(nodesToRemove, node)
		}
	}

	for _, node := range nodesToRemove {
		node.shapes.Remove(shape)
		shape.removeNodes(node)
	}

	// Add the shape in all the nodes
	// that must be in its domain.
	nodeQueue := queue.New()
	nodeQueue.Enqueue(space.tree.root)

	for nodeQueue.Len() > 0 {
		node := nodeQueue.Dequeue().(*quadTreeNode)

		// If the node is already
		// in the domain, skip it.
		if shape.containsNode(node) {
			continue
		}

		// If the shape is not covered by the node area,
		// skip it to the next node.
		overlapped, err := node.boundary.collidesShape(shape)

		if err != nil {
			return nil, err
		}

		if !overlapped {
			continue
		}

		// If the node is not a leaf,
		// skip it.
		if node.northWest != nil {
			nodeQueue.Enqueue(node.northEast)
			nodeQueue.Enqueue(node.northWest)
			nodeQueue.Enqueue(node.southEast)
			nodeQueue.Enqueue(node.southWest)

			continue
		}

		// If the node limit is not exceeded,
		// add the shape in the list of shapes
		// covered by the node area.
		if len(node.shapes) < node.tree.nodeCapacity ||
			node.level >= node.tree.maxLevel {
			node.shapes.Insert(shape)
			shape.addNodes(node)
		} else {
			// Split the node into four subareas
			// and add the subnodes in the queue.
			err := node.split()

			if err != nil {
				return nil, err
			}

			nodeQueue.Enqueue(node.northEast)
			nodeQueue.Enqueue(node.northWest)
			nodeQueue.Enqueue(node.southEast)
			nodeQueue.Enqueue(node.southWest)
		}
	}

	// Return all the nodes where
	// the shape is now located in.
	cells := map[Vector]Shapes{}

	for _, node := range shape.nodes() {
		cells[node.boundary.center()] = node.shapes.Copy()
	}

	return cells, nil
}

// Rebuild rebuilds the space's index
// of fhapes in purpose to optimize it.
func (space *Space) Rebuild() error {
	return space.tree.redistribute()
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
				overlapped, err := ResolveCollision(shape, otherShape, space.useTags)

				if err != nil {
					return nil, err
				}

				if overlapped {
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
	if shape == nil {
		return nil, fmt.Errorf("the shape is nil")
	}

	shapes := make(Shapes, 0)
	nodes := shape.nodes()

	for _, area := range nodes {
		for item := range area.shapes {
			overlapped, err := ResolveCollision(item, shape, space.useTags)

			if err != nil {
				return nil, err
			}

			if item != shape && overlapped {
				shapes.Insert(item)
			}
		}
	}

	return shapes, nil
}

// CollidedBy returns the set of shapes collided by the given shape.
func (space *Space) CollidedBy(shape Shape) (Shapes, error) {
	if shape == nil {
		return nil, fmt.Errorf("the shape is nil")
	}

	shapes := make(Shapes, 0)
	nodes := shape.nodes()

	for _, area := range nodes {
		for item := range area.shapes {
			overlapped, err := ResolveCollision(shape, item, space.useTags)

			if err != nil {
				return nil, err
			}

			if item != shape && overlapped {
				shapes.Insert(item)
			}
		}
	}

	return shapes, nil
}

// WouldBeCollidedBy returns all the shapes that would be collided by
// the given shape if it moved in the specified direction.
func (space *Space) WouldBeCollidedBy(shape Shape, moveDiff Vector, turnDiff float64) (Shapes, error) {
	if shape == nil {
		return nil, fmt.Errorf("the shape is nil")
	}

	shapes := make(Shapes, 0)
	originalPos := shape.Center()
	originalAngle := shape.Angle()
	shapeType := reflect.TypeOf(shape).Elem()

	// Get all the nodes where the shape is located
	// before movement.
	areas := shape.nodes()

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
		nodes[area.boundary.center()] = area.shapes.Copy()
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
				shouldCollide, err := lineShape.ShouldCollide(lineItem)

				if err != nil {
					return nil, err
				}

				if space.useTags && !shouldCollide {
					continue
				}

				linesWouldIntersect, err := linesWouldCollide(originalPos,
					originalAngle, moveDiff, turnDiff, lineShape, lineItem)

				if err != nil {
					return nil, err
				}

				if linesWouldIntersect {
					shapes.Insert(lineItem)

					continue
				}
			}

			overlapped, err := ResolveCollision(shape,
				item, space.useTags)

			if err != nil {
				return nil, err
			}

			if overlapped {
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

// WouldBeCollidingWith returns all the shapes that would be colliding the given one
// if it moved in the specified direction.
func (space *Space) WouldBeCollidingWith(shape Shape, moveDiff Vector, turnDiff float64) (Shapes, error) {
	if shape == nil {
		return nil, fmt.Errorf("the shape is nil")
	}

	shapes := make(Shapes, 0)
	originalPos := shape.Center()
	originalAngle := shape.Angle()
	shapeType := reflect.TypeOf(shape).Elem()

	// Get all the nodes where the shape is located
	// before movement.
	areas := shape.nodes()

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
		nodes[area.boundary.center()] = area.shapes.Copy()
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
				linesCollinear, err := lineShape.CollinearTo(lineItem)

				if err != nil {
					return nil, err
				}

				if linesCollinear {
					// Check line tags.
					shouldCollide, err := lineItem.ShouldCollide(lineShape)

					if err != nil {
						return nil, err
					}

					if space.useTags && !shouldCollide {
						continue
					}

					linesWouldIntersect, err := linesWouldCollide(originalPos,
						originalAngle, moveDiff, turnDiff, lineShape, lineItem)

					if err != nil {
						return nil, err
					}

					if linesWouldIntersect {
						shapes.Insert(lineItem)

						continue
					}
				}
			}

			overlapped, err := ResolveCollision(item,
				shape, space.useTags)

			if err != nil {
				return nil, err
			}

			if overlapped {
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
func NewSpace(
	subdivisionFactor, shapesInArea int, width,
	height float64, min, max Vector, useTags bool,
) (
	*Space, error,
) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf(
			"Space must have positive values for width and height")
	}

	if min.X >= max.X {
		return nil, fmt.Errorf(
			"The value of min X is incorrect")
	}

	if min.Y >= max.Y {
		return nil, fmt.Errorf(
			"The value of min Y is incorrect")
	}

	space := new(Space)
	space.min = min
	space.max = max
	space.useTags = useTags
	space.shapes = make(Shapes, 0)
	boundary, err := newAABB(NewVector(
		-width/2.0, -height/2.0),
		NewVector(width/2.0, height/2.0))

	if err != nil {
		return nil, err
	}

	tree, err := newQuadTree(boundary,
		subdivisionFactor, shapesInArea)

	if err != nil {
		return nil, err
	}

	space.tree = tree

	return space, nil
}
