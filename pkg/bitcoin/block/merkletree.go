package block

import (
	"fmt"
	"math"

	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
)

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

func MakeMerkleTree(total int) *MerkleTree {
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

func (m *MerkleTree) PopulateTree(flagBits []byte, hashes [][]byte) error {
	// iterate until we have a root populated in the tree
	for {
		// break condition, we have populated the merkle root, hurray
		if !m.RootIsEmpty() {
			break
		}

		// left nodes are always given a hash e.g. the rights are the ones
		// that we need to worry about duplicating
		if m.IsLeaf() {
			// dequeue the last flag bit because reasons - the book doesnt explain this
			flagBits = flagBits[:len(flagBits)-2]

			// get the last hash in the array
			m.SetCurrentNode(hashes[len(hashes)-1])

			// dequeue the last hash so we have n-1 hashes left
			hashes = hashes[:len(hashes)-2]

			// move up the tree
			m.Up()
		} else {
			leftHash := m.GetLeftNode()

			//If we don’t have the left child value, there are two possibilities.
			// This node’s value may be in the hashes field, or it might need calculation.
			if utils.IsNull(leftHash) {
				// flag bit for this tells whether or not to calculate the node.
				// it it is 0, next hash is the value for the node. If the bit is
				// set to a 1, need to calculate the left and maybe the right as well
				bit := flagBits[len(flagBits)-1]

				// simulate pop
				flagBits = flagBits[:len(flagBits)-2]

				if bit == 0 {
					m.SetCurrentNode(hashes[len(hashes)-1])

					// simulate pop
					hashes = hashes[:len(hashes)-2]
				} else {
					// know we have a left node, so move over to it
					m.Left()
				}
			} else if m.RightExists() {
				rightHash := m.GetRightNode()

				// Check that we have a value for the right node
				if utils.IsNull(rightHash) {
					// move over the the right node now
					m.Right()
				} else {
					// we have both the left and right nodes now, so we can
					// calculate the merkle parent for the current node
					m.SetCurrentNode(utils.MerkleParent(leftHash, rightHash))

					// move up to the next parent level
					m.Up()
				}
			} else {
				// we have an odd number of nodes, so hash the left and left
				// to make a parent
				m.SetCurrentNode(utils.MerkleParent(leftHash, leftHash))

				// move up to the next parent level
				m.Up()
			}
		}
	}

	if len(hashes) != 0 {
		// have not consumed all the hashes so something went wrong
		return fmt.Errorf("not all the hashes have been consumed")
	}

	// all the flag bits should have been consumed, so check that it is true
	for _, bit := range flagBits {
		if bit != 0x00 {
			return fmt.Errorf("not all the flag bits have been consumed (%v)", flagBits)
		}
	}

	return nil
}
