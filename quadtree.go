package cirno

import (
	"fmt"

	"github.com/golang-collections/collections/queue"
)

// quadTree is an implementation of
// quad tree data structure to subdivide
// space to detect collisions.
type quadTree struct {
	root         *quadTreeNode
	maxLevel     int
	nodeCapacity int
	leaves       map[*quadTreeNode]none
}

// addLeaf adds the quad tree node in the list of quad tree leaves.
func (tree *quadTree) addLeaf(node *quadTreeNode) error {
	if tree.containsLeaf(node) {
		return fmt.Errorf("The leaf {%f, %f} already exists",
			node.boundary.center.X, node.boundary.center.Y)
	}

	if node.northWest != nil {
		return fmt.Errorf("The node {%f, %f} cannot be a leaf",
			node.boundary.center.X, node.boundary.center.Y)
	}

	tree.leaves[node] = none{}

	return nil
}

// removeLeaf removes the specified quad tree node from the list
// of quad tree leaves.
func (tree *quadTree) removeLeaf(node *quadTreeNode) error {
	if !tree.containsLeaf(node) {
		return fmt.Errorf("The leaf {%f, %f} doesn't exist",
			node.boundary.center.X, node.boundary.center.Y)
	}

	delete(tree.leaves, node)

	return nil
}

// containsLeaf returns true if the quad tree contains
// the given node among the leaves, and false otherwise.
func (tree *quadTree) containsLeaf(node *quadTreeNode) bool {
	_, exists := tree.leaves[node]

	return exists
}

// insert inserts the given shape into appropriate nodes
// and returns them.
func (tree *quadTree) insert(shape Shape) ([]*quadTreeNode, error) {
	if shape == nil {
		return nil, fmt.Errorf("The shape cannot be nil")
	}

	if !tree.root.boundary.ContainsPoint(shape.Center()) {
		return nil, fmt.Errorf("The shape is out of bounds")
	}

	nodes := []*quadTreeNode{}
	nodeQueue := queue.New()
	nodeQueue.Enqueue(tree.root)

	for nodeQueue.Len() > 0 {
		node := nodeQueue.Dequeue().(*quadTreeNode)

		// If the shape is not covered by the node area,
		// skip it to the next node.
		if !ResolveCollision(node.boundary, shape, false) {
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
			nodes = append(nodes, node)
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

	// Add nodes to the shape's domain.
	shape.addNodes(nodes...)

	return nodes, nil
}

// search returns all the nodes containing the given shape.
func (tree *quadTree) search(shape Shape) ([]*quadTreeNode, error) {
	if shape == nil {
		return nil, fmt.Errorf("The shape cannot be nil")
	}

	if !tree.root.boundary.ContainsPoint(shape.Center()) {
		return nil, fmt.Errorf("The shape is out of bounds")
	}

	nodes := []*quadTreeNode{}
	nodeQueue := queue.New()
	nodeQueue.Enqueue(tree.root)

	for nodeQueue.Len() > 0 {
		node := nodeQueue.Dequeue().(*quadTreeNode)

		// If the shape is not covered by the node area,
		// skip it to the next node.
		if !ResolveCollision(node.boundary, shape, false) {
			continue
		}

		if node.northWest == nil && node.shapes.Contains(shape) {
			nodes = append(nodes, node)
		} else {
			nodeQueue.Enqueue(node.northEast)
			nodeQueue.Enqueue(node.northWest)
			nodeQueue.Enqueue(node.southEast)
			nodeQueue.Enqueue(node.southWest)
		}
	}

	return nodes, nil
}

// remove removes the specified shape from the quad tree.
func (tree *quadTree) remove(shape Shape) error {
	if shape == nil {
		return fmt.Errorf("The shape cannot be nil")
	}

	if !tree.root.boundary.ContainsPoint(shape.Center()) {
		return fmt.Errorf("The shape is out of bounds")
	}

	for _, node := range shape.nodes() {
		node.shapes.Remove(shape)
	}

	// Remove all the nodes
	// from the shape's domain.
	shape.clearNodes()

	return nil
}

// redistribute removes all the unrequired leafs
// and subtrees containing them.
func (tree *quadTree) redistribute() error {
	nodeQueue := queue.New()

	for leaf := range tree.leaves {
		nodeQueue.Enqueue(leaf)
	}

	for nodeQueue.Len() > 0 {
		node := nodeQueue.Dequeue().(*quadTreeNode)

		if node.parent == nil {
			continue
		}

		parent := node.parent

		// Check if all the child nodes are leafs.
		if tree.containsLeaf(parent.northWest) &&
			tree.containsLeaf(parent.northEast) &&
			tree.containsLeaf(parent.southWest) &&
			tree.containsLeaf(parent.southEast) {

			shapes := Shapes{}

			shapes.Merge(parent.northWest.shapes)
			shapes.Merge(parent.northEast.shapes)
			shapes.Merge(parent.southWest.shapes)
			shapes.Merge(parent.southEast.shapes)

			if len(shapes) <= tree.nodeCapacity {
				err := parent.assemble()

				if err != nil {
					return err
				}

				nodeQueue.Enqueue(parent)
			}
		}
	}

	return nil
}

// clear removes all the shapes from the quad tree.
func (tree *quadTree) clear() error {
	// Remove all the nodes from
	// shapes' domains.
	for node := range tree.leaves {
		for shape := range node.shapes {
			shape.clearNodes()
		}
	}

	tree.root = &quadTreeNode{
		tree:     tree,
		boundary: tree.root.boundary,
		level:    0,
		shapes:   Shapes{},
	}
	tree.leaves = map[*quadTreeNode]none{}

	if err := tree.addLeaf(tree.root); err != nil {
		return err
	}

	return nil
}

// shapeGroups returns the dictionary of shapes grouped
// by their nodes in the quad tree.
func (tree *quadTree) shapeGroups() map[*quadTreeNode]Shapes {
	shapes := map[*quadTreeNode]Shapes{}

	for node := range tree.leaves {
		shapes[node] = node.shapes.Copy()
	}

	return shapes
}

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
	boundary  *Rectangle
	shapes    Shapes
	level     int
}

// add adds all the shapes covered by node area
// in the set of node shapes.
func (node *quadTreeNode) add(shapes Shapes) {
	for shape := range shapes {
		if ResolveCollision(node.boundary, shape, false) {
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

	// Compute centers for new areas.
	neRectCenter := node.boundary.center.Add(NewVector(
		node.boundary.extents.X/2.0,
		node.boundary.extents.Y/2.0,
	))
	nwRectCenter := node.boundary.center.Add(NewVector(
		-node.boundary.extents.X/2.0,
		node.boundary.extents.Y/2.0,
	))
	seRectCenter := node.boundary.center.Add(NewVector(
		node.boundary.extents.X/2.0,
		-node.boundary.extents.Y/2.0,
	))
	swRectCenter := node.boundary.center.Add(NewVector(
		-node.boundary.extents.X/2.0,
		-node.boundary.extents.Y/2.0,
	))

	// Create new nodes.
	node.northEast = &quadTreeNode{
		tree:   node.tree,
		parent: node,
		boundary: NewRectangle(neRectCenter,
			node.boundary.extents.X,
			node.boundary.extents.Y, 0.0),
		level:  nextLevel,
		shapes: Shapes{},
	}

	node.northWest = &quadTreeNode{
		tree:   node.tree,
		parent: node,
		boundary: NewRectangle(nwRectCenter,
			node.boundary.extents.X,
			node.boundary.extents.Y, 0.0),
		level:  nextLevel,
		shapes: Shapes{},
	}

	node.southEast = &quadTreeNode{
		tree:   node.tree,
		parent: node,
		boundary: NewRectangle(seRectCenter,
			node.boundary.extents.X,
			node.boundary.extents.Y, 0.0),
		level:  nextLevel,
		shapes: Shapes{},
	}

	node.southWest = &quadTreeNode{
		tree:   node.tree,
		parent: node,
		boundary: NewRectangle(swRectCenter,
			node.boundary.extents.X,
			node.boundary.extents.Y, 0.0),
		level:  nextLevel,
		shapes: Shapes{},
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

// newQuadTree creates a new empty quad tree.
func newQuadTree(boundary *Rectangle, maxLevel, nodeCapacity int) (*quadTree, error) {
	if maxLevel < 1 {
		return nil, fmt.Errorf("Max depth must be greater or equal to 1")
	}

	tree := new(quadTree)

	tree.maxLevel = maxLevel
	tree.nodeCapacity = nodeCapacity
	tree.leaves = map[*quadTreeNode]none{}
	tree.root = &quadTreeNode{
		tree:     tree,
		parent:   nil,
		boundary: boundary,
		level:    0,
		shapes:   Shapes{},
	}

	if err := tree.addLeaf(tree.root); err != nil {
		return nil, err
	}

	return tree, nil
}
