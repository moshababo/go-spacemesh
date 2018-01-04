package merkle

import (
	"bytes"
	"fmt"
	"errors"
	"github.com/spacemeshos/go-spacemesh/merkle/pb"
)

func (mt *merkleTreeImp) ValidateStructure(root NodeContainer) ([]byte, error) {

	if root == nil {
		return nil, errors.New("expected non-empty root")
	}

	err := root.loadChildren(mt.treeData)
	if err != nil {
		return nil, err
	}

	switch root.getNodeType() {

	case pb.NodeType_branch:

		entries := root.getBranchNode().getAllChildNodePointers()
		children := root.getAllChildren()

		if len(entries) != len(children) {
			return nil, errors.New(fmt.Sprintf("mismatch. entries: %d, children: %d", len(entries), len(children)))
		}

		for _, c := range children {
			_, err := mt.ValidateStructure(c)
			if err != nil {
				return nil, err
			}
		}

		return root.getNodeHash(), nil

	case pb.NodeType_extension:
		children := root.getAllChildren()
		if len(children) != 1 {
			return nil, errors.New("expected 1 child for extension node")
		}

		childHash, err := mt.ValidateStructure(children[0])
		if err != nil {
			return nil, err
		}

		if !bytes.Equal(childHash, root.getExtNode().getValue()) {
			return nil, errors.New("hash mismatch")
		}

		return childHash, nil

	case pb.NodeType_leaf:

		return root.getNodeHash(), nil
	}

	return nil, errors.New("unexpected node type")
}
