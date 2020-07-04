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

// addNodes adds new nodes in the list of the nodes
// the shape belongs to.
func (d *domain) addNodes(nodes ...*quadTreeNode) {
	d.treeNodes = append(d.treeNodes, nodes...)
}

// containsNode returns true if the node is included
// in the shape's domain, and false otherwise.
func (d *domain) containsNode(node *quadTreeNode) bool {
	for _, treeNode := range d.treeNodes {
		if treeNode == node {
			return true
		}
	}

	return false
}

// removeNodes removes the quad tree nodes from
// the list of the nodes the shape belongs to.
//
// This method does nothing if the node
// is not included in the shape's domain.
func (d *domain) removeNodes(nodes ...*quadTreeNode) {
	for _, node := range nodes {
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
}

// clearNodes removes all the nodes
// from the shape's domain.
func (d *domain) clearNodes() {
	d.treeNodes = []*quadTreeNode{}
}
