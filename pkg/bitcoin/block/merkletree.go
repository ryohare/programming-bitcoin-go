package block

import "math"

type MerkleTree struct {
	// Total number of elements in the data structure
	Total int

	// Maximum depth that the tree will support.
	// Set by either the Make function or GetMaxDepth()
	MaxDepth int

	// Nodes - This is complex - Array of an Array of byte arrays
	// {
	// 	 { 0x01 },
	// 	 { 0x02, 0x03 },
	// 	 { 0x04, 0x05, 0x06, 0x07 },
	// }
	Nodes [][][]byte

	// Current depth of the tree
	CurrentDepth int

	// Current index into the underlying array
	CurrentIndex int
}

func (m *MerkleTree) SetMaxDepth() {
	m.MaxDepth = m.GetMaxDepth()
}

func (m *MerkleTree) GetMaxDepth() int {

	// Since we half at every level, we use the log2 of the number of leaves
	// This gets rounded up because you cant have 1/2 a level after all.
	return int(math.Ceil(math.Log2(float64(m.Total))))
}

func Make(total int) *MerkleTree {
	mt := &MerkleTree{}

	mt.Total = total
	mt.MaxDepth = mt.GetMaxDepth()

	// Allocate all the memory for the nodes up front
	nodes := make([][][]byte, mt.MaxDepth+1)

	for i := 0; i < mt.MaxDepth+1; i++ {

		// number of items at the corresponding level
		totalfF := float64(total)
		maxDepthF := float64(mt.MaxDepth)
		exp := math.Pow(2, maxDepthF-float64(i))
		numItems := math.Ceil(totalfF / exp)

		// make an array of the hashes byte arrays
		nodes[i] = make([][]byte, int(numItems))

		// now create the 32 byte hash for the arrays of bytes
		for j := 0; j < int(numItems); j++ {
			nodes[i][j] = make([]byte, 32)
		}
	}
	mt.Nodes = nodes
	mt.CurrentDepth = 0
	mt.CurrentIndex = 0

	return mt
}

func (m *MerkleTree) Up() {
	m.CurrentDepth -= 1
	m.CurrentIndex /= 2
}

func (m *MerkleTree) Left() {
	m.CurrentDepth += 1
	m.CurrentIndex *= 2
}

func (m *MerkleTree) Right() {
	m.CurrentDepth += 1
	m.CurrentIndex = m.CurrentIndex*2 + 1
}

func (m *MerkleTree) Root() []byte {
	return m.Nodes[0][0]
}

func (m *MerkleTree) SetCurrentNode(val []byte) {
	m.Nodes[m.CurrentDepth][m.CurrentIndex] = val
}

func (m *MerkleTree) GeCurrentNode() []byte {
	return m.Nodes[m.CurrentDepth][m.CurrentIndex]
}

func (m *MerkleTree) GetLeftNode() []byte {
	return m.Nodes[m.CurrentDepth+1][m.CurrentIndex*2]
}

func (m *MerkleTree) GetRightNode() []byte {
	return m.Nodes[m.CurrentDepth+1][m.CurrentIndex*2+1]
}

func (m *MerkleTree) IsLeaf() bool {
	return m.CurrentDepth == m.MaxDepth
}

func (m *MerkleTree) RightExists() bool {
	return len(m.Nodes[m.CurrentDepth+1]) > m.CurrentIndex*2+1
}

func (m *MerkleTree) RootIsEmpty() bool {
	b := m.Root()

	for _, v := range b {
		if v != 0x00 {
			return false
		}
	}

	return true
}
