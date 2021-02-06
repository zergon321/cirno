package cirno

// TODO: use AABB instead of rectangle for boundary.

// represents a single node in quad tree
// for a certain subarea.
type quadTreeNode struct {
	tree      *quadTree
	parent    *quadTreeNode
	northEast *quadTreeNode
	northWest *quadTreeNode
	southWest *quadTreeNode
	southEast *quadTreeNode
	boundary  *aabb
	shapes    Shapes
	level     int
}

// add adds all the shapes covered by node area
// in the set of node shapes.
func (node *quadTreeNode) add(shapes Shapes) {
	for shape := range shapes {
		if node.boundary.collidesShape(shape) {
			node.shapes.Insert(shape)
			shape.addNodes(node)
		}
	}
}

// remove removes all the shapes from the set that
// have the node in their domains.
func (node *quadTreeNode) remove(shapes Shapes) {
	for shape := range shapes {
		node.shapes.Remove(shape)
		shape.removeNodes(node)
	}
}

// clear removes all the shapes from the node
// removes the node from the shapes' domains.
func (node *quadTreeNode) clear() {
	for shape := range node.shapes {
		shape.removeNodes(node)
	}

	node.shapes = Shapes{}
}

// split subdivides the node area into four subareas
// and creates new nodes for the subareas; reassigns
// shapes to the subnodes.
func (node *quadTreeNode) split() error {
	nextLevel := node.level + 1

	// Compute the center of the original boundary
	// the other points to form new boundaries.
	center := node.boundary.center()
	northPoint := NewVector(center.X, node.boundary.max.Y)
	southPoint := NewVector(center.X, node.boundary.min.Y)
	westPoint := NewVector(node.boundary.min.X, center.Y)
	eastPoint := NewVector(node.boundary.max.X, center.Y)

	// Create new nodes.
	node.northEast = &quadTreeNode{
		tree:     node.tree,
		parent:   node,
		boundary: newAABB(center, node.boundary.max),
		level:    nextLevel,
		shapes:   Shapes{},
	}

	node.northWest = &quadTreeNode{
		tree:     node.tree,
		parent:   node,
		boundary: newAABB(westPoint, northPoint),
		level:    nextLevel,
		shapes:   Shapes{},
	}

	node.southEast = &quadTreeNode{
		tree:     node.tree,
		parent:   node,
		boundary: newAABB(southPoint, eastPoint),
		level:    nextLevel,
		shapes:   Shapes{},
	}

	node.southWest = &quadTreeNode{
		tree:     node.tree,
		parent:   node,
		boundary: newAABB(node.boundary.min, center),
		level:    nextLevel,
		shapes:   Shapes{},
	}

	// Redistribute shapes between subnodes.
	node.northEast.add(node.shapes)
	node.northWest.add(node.shapes)
	node.southEast.add(node.shapes)
	node.southWest.add(node.shapes)
	node.clear()

	// Remove the current node from tree leaves.
	err := node.tree.removeLeaf(node)

	if err != nil {
		return err
	}

	// Add its children to tree leaves.
	err = node.tree.addLeaf(node.northEast)

	if err != nil {
		return err
	}

	err = node.tree.addLeaf(node.northWest)

	if err != nil {
		return err
	}

	err = node.tree.addLeaf(node.southEast)

	if err != nil {
		return err
	}

	err = node.tree.addLeaf(node.southWest)

	return err
}

// assemble adds all the children shapes to the parent
// and removes children.
func (node *quadTreeNode) assemble() error {
	// Add all the shapes in the parent node.
	node.add(node.northWest.shapes)
	node.add(node.northEast.shapes)
	node.add(node.southWest.shapes)
	node.add(node.southEast.shapes)

	// Clear all the child nodes.
	node.northWest.clear()
	node.northEast.clear()
	node.southWest.clear()
	node.southEast.clear()

	// Remove all the children from leaves.
	err := node.tree.removeLeaf(node.northWest)

	if err != nil {
		return err
	}

	err = node.tree.removeLeaf(node.northEast)

	if err != nil {
		return err
	}

	err = node.tree.removeLeaf(node.southWest)

	if err != nil {
		return err
	}

	err = node.tree.removeLeaf(node.southEast)

	if err != nil {
		return err
	}

	// Get rid of child nodes.
	node.northWest.parent = nil
	node.northEast.parent = nil
	node.southWest.parent = nil
	node.southEast.parent = nil

	node.northWest = nil
	node.northEast = nil
	node.southWest = nil
	node.southEast = nil

	// Add the parent node to the tree leaves.
	err = node.tree.addLeaf(node)

	if err != nil {
		return err
	}

	return nil
}
