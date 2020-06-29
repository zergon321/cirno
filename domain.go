package cirno

// domain contains all the
// quad tree nodes the shape
// belongs to.
type domain struct {
	treeNodes []*quadTreeNode
}

// nodes returns all the quad tree nodes the shape
// belongs to.
func (d *domain) nodes() []*quadTreeNode {
	nodes := make([]*quadTreeNode, len(d.treeNodes))
	copy(nodes, d.treeNodes)

	return nodes
}

// addNode adds a new node in the list of nodes
// the shape belongs to.
func (d *domain) addNode(node *quadTreeNode) {
	d.treeNodes = append(d.treeNodes, node)
}

// removeNode removes the quad tree node from the list
// nodes the shape belongs to.
//
// This method does nothing if the node is not included
// in the shape's domain.
func (d *domain) removeNode(node *quadTreeNode) {
	index := -1

	for i, treeNode := range d.treeNodes {
		if treeNode == node {
			index = i

			break
		}
	}

	if index >= 0 {
		d.treeNodes = append(d.treeNodes[:index],
			d.treeNodes[index+1:]...)
	}
}
