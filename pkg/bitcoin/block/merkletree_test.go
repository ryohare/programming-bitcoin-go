package block

import (
	"encoding/hex"
	"testing"

	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
)

func TestMake(t *testing.T) {
	if MakeMerkleTree(16) == nil {
		t.Fatal("failed to make merkle tree")
	}
}

func TestMerkleTreeRaw(t *testing.T) {
	hashes := []string{
		"9745f7173ef14ee4155722d1cbf13304339fd00d900b759c6f9d58579b5765fb",
		"5573c8ede34936c29cdfdfe743f7f5fdfbd4f54ba0705259e62f39917065cb9b",
		"82a02ecbb6623b4274dfcab82b336dc017a27136e08521091e443e62582e8f05",
		"507ccae5ed9b340363a0e6d765af148be9cb1c8766ccc922f83e4ae681658308",
		"a7a4aec28e7162e1e9ef33dfa30f0bc0526e6cf4b11a576f6c5de58593898330",
		"bb6267664bd833fd9fc82582853ab144fece26b7a8a5bf328f8a059445b59add",
		"ea6d7ac1ee77fbacee58fc717b990c4fcccf1b19af43103c090f601677fd8836",
		"457743861de496c429912558a106b810b0507975a49773228aa788df40730d41",
		"7688029288efc9e9a0011c960a6ed9e5466581abf3e3a6c26ee317461add619a",
		"b1ae7f15836cb2286cdd4e2c37bf9bb7da0a2846d06867a429f654b2e7f383c9",
		"9b74f89fa3f93e71ff2c241f32945d877281a6a50a6bf94adac002980aafe5ab",
		"b3a92b5b255019bdaf754875633c2de9fec2ab03e6b8ce669d07cb5b18804638",
		"b5c0b915312b9bdaedd2b86aa2d0f8feffc73a2d37668fd9010179261e25e263",
		"c9d52c5cb1e557b92c84c52e7c4bfbce859408bedffc8a5560fd6e35e10b8800",
		"c555bc5fc3bc096df0a0c9532f07640bfb76bfe4fc1ace214b8b228a1297a4c2",
		"f9dbfafc3af3400954975da24eb325e326960a25b87fffe23eef3e7ed2fb610e",
	}

	// create the tree
	tree := MakeMerkleTree(len(hashes))

	hashesBytes := make([][]byte, 0, len(hashes))

	for _, h := range hashes {
		b, _ := hex.DecodeString(h)
		hashesBytes = append(hashesBytes, b)
	}

	tree.Nodes[4] = hashesBytes

	b, err := utils.MerkleParentLevel(tree.Nodes[4])
	if err != nil {
		t.Fatal("failed to parse parent level")
	}
	tree.Nodes[3] = b
	b, err = utils.MerkleParentLevel(tree.Nodes[3])
	if err != nil {
		t.Fatal("failed to parse parent level")
	}
	tree.Nodes[2] = b
	b, err = utils.MerkleParentLevel(tree.Nodes[2])
	if err != nil {
		t.Fatal("failed to parse parent level")
	}
	tree.Nodes[1] = b
	b, err = utils.MerkleParentLevel(tree.Nodes[1])
	if err != nil {
		t.Fatal("failed to parse parent level")
	}
	tree.Nodes[0] = b
}

func TestMerkleTreeNormal(t *testing.T) {
	hashes := []string{
		"9745f7173ef14ee4155722d1cbf13304339fd00d900b759c6f9d58579b5765fb",
		"5573c8ede34936c29cdfdfe743f7f5fdfbd4f54ba0705259e62f39917065cb9b",
		"82a02ecbb6623b4274dfcab82b336dc017a27136e08521091e443e62582e8f05",
		"507ccae5ed9b340363a0e6d765af148be9cb1c8766ccc922f83e4ae681658308",
		"a7a4aec28e7162e1e9ef33dfa30f0bc0526e6cf4b11a576f6c5de58593898330",
		"bb6267664bd833fd9fc82582853ab144fece26b7a8a5bf328f8a059445b59add",
		"ea6d7ac1ee77fbacee58fc717b990c4fcccf1b19af43103c090f601677fd8836",
		"457743861de496c429912558a106b810b0507975a49773228aa788df40730d41",
		"7688029288efc9e9a0011c960a6ed9e5466581abf3e3a6c26ee317461add619a",
		"b1ae7f15836cb2286cdd4e2c37bf9bb7da0a2846d06867a429f654b2e7f383c9",
		"9b74f89fa3f93e71ff2c241f32945d877281a6a50a6bf94adac002980aafe5ab",
		"b3a92b5b255019bdaf754875633c2de9fec2ab03e6b8ce669d07cb5b18804638",
		"b5c0b915312b9bdaedd2b86aa2d0f8feffc73a2d37668fd9010179261e25e263",
		"c9d52c5cb1e557b92c84c52e7c4bfbce859408bedffc8a5560fd6e35e10b8800",
		"c555bc5fc3bc096df0a0c9532f07640bfb76bfe4fc1ace214b8b228a1297a4c2",
		"f9dbfafc3af3400954975da24eb325e326960a25b87fffe23eef3e7ed2fb610e",
	}

	// load the first parent level of the tree
	tree := MakeMerkleTree(len(hashes))

	hashesBytes := make([][]byte, 0, len(hashes))

	for _, h := range hashes {
		b, _ := hex.DecodeString(h)
		hashesBytes = append(hashesBytes, b)
	}

	tree.Nodes[4] = hashesBytes

	// now use the functions and iterate over the tree and fill it out
	for {
		// if we populate the root, time to get the hell out of this loop
		if !tree.RootIsEmpty() {
			break
		}

		if tree.IsLeaf() {
			tree.Up()
		} else {
			leftHash := tree.GetLeftNode()
			rightHash := tree.GetRightNode()

			if utils.IsNull(leftHash) {
				tree.Left()
			} else if tree.RightExists() {
				if utils.IsNull(rightHash) {
					tree.Right()
				} else {
					// make the merkle parent for the specified node
					tree.SetCurrentNode(utils.MerkleParent(leftHash, rightHash))

					// move up a level
					tree.Up()
				}
			} else {
				tree.SetCurrentNode(utils.MerkleParent(leftHash, rightHash))

				// move up a level
				tree.Up()
			}
		}
	}
}
