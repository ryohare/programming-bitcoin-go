package block

import "math"

type MerkleTree struct {
	// Total number of elements in the data structure
	Total int

	// Maximum depth that the tree will support.
	// Set by either the Make function or GetMaxDepth()
	MaxDepth int

	// Nodes

	// Current depth of the tree
	CurrentDepth int

	// Current index into the underlying array
	CurrentIndex int
}

func (m *MerkleTree) SetMaxDepth() {
	m.MaxDepth = m.GetMaxDepth()
}

func (m *MerkleTree) GetMaxDepth() int {
	return int(math.Ceil(math.Log2(float64(m.Total))))
}
