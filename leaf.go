package subtree

import (
	"bytes"
)

//-------------------
// Leaf Node Definition
//-------------------

// The leaf struct represents a leaf node in the tree.
// It holds the value and suffix for the leaf. The order of fields is optimized for memory alignment.
type leaf[T any] struct {
	value  T      // The value associated with this leaf
	suffix []byte // Suffix portion that we will store, assuming the prefix has been checked already
}

//-------------------
// Leaf Node Methods
//-------------------

// newLeaf creates a new leaf node with the given suffix and value.
// It returns a pointer to the newly created leaf.
func newLeaf[T any](suffix []byte, value T) *leaf[T] {
	return &leaf[T]{value, copyBytes(suffix)} // Use copyBytes to ensure suffix is safely copied
}

// isLeaf returns true as this node is a leaf.
func (n *leaf[T]) isLeaf() bool { return true }

// base returns nil because leaves do not have a base node.
func (n *leaf[T]) base() *meta { return nil }

// match checks if the given subject matches the leaf's suffix.
func (n *leaf[T]) match(subject []byte) bool {
	return bytes.Equal(subject, n.suffix) // Compare subject with the leaf's suffix
}

// setSuffix sets the suffix for this leaf node.
func (n *leaf[T]) setSuffix(suffix []byte) {
	n.suffix = copyBytes(suffix) // Copy the provided suffix to ensure safety
}

// isFull returns true because leaf nodes are considered "full" once they have a value.
func (n *leaf[T]) isFull() bool { return true }

// matchParts checks if the parts of the subject match the leaf's suffix.
// It delegates to the matchParts function for comparison.
func (n *leaf[T]) matchParts(parts [][]byte) ([][]byte, bool) {
	return matchParts(parts, n.suffix)
}

// iter is a no-op for leaf nodes as they don't have children.
func (n *leaf[T]) iter(f func(node) bool) {}

// children returns nil because leaf nodes don't have any children.
func (n *leaf[T]) children() []node { return nil }

// numChildren returns 0 because leaf nodes don't have any children.
func (n *leaf[T]) numChildren() uint16 { return 0 }

// path returns the suffix for this leaf as its path.
func (n *leaf[T]) path() []byte { return n.suffix }

//-------------------
// Methods that should panic when called on a leaf node
//-------------------

// These methods are not applicable to leaf nodes. If they are called, a panic will occur.
func (n *leaf[T]) setPrefix(pre []byte)    { panic("setPrefix called on leaf") }
func (n *leaf[T]) addChild(_ byte, _ node) { panic("addChild called on leaf") }
func (n *leaf[T]) findChild(_ byte) *node  { panic("findChild called on leaf") }
func (n *leaf[T]) grow() node              { panic("grow called on leaf") }
func (n *leaf[T]) deleteChild(_ byte)      { panic("deleteChild called on leaf") }
func (n *leaf[T]) shrink() node            { panic("shrink called on leaf") }
