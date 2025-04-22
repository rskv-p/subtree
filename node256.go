package subtree

//-------------------
// Node256 Definition
//-------------------

// node256 represents a node with up to 256 possible children. It is designed for situations
// where the node needs to support a larger number of children without requiring additional
// memory optimizations. The child array is directly indexed by the byte value of the key.
// The struct is optimized for memory alignment according to govet/fieldalignment recommendations.
type node256 struct {
	child [256]node // Array of child nodes (up to 256 children)
	meta            // Inherited metadata (prefix and size)
}

//-------------------
// Node256 Methods
//-------------------

// newNode256 creates a new node256 with the specified prefix and returns a pointer to it.
func newNode256(prefix []byte) *node256 {
	nn := &node256{}
	nn.setPrefix(prefix) // Set the prefix for the node
	return nn
}

// addChild adds a child node to the current node. The child is indexed by the byte value of its key.
// This method directly stores the child in the array at the position corresponding to the key.
func (n *node256) addChild(c byte, nn node) {
	n.child[c] = nn // Store the child node at the index corresponding to the key
	n.size++        // Increment the size to reflect the added child
}

// findChild looks for a child node by its key (byte). If found, it returns a pointer to the child node.
func (n *node256) findChild(c byte) *node {
	if n.child[c] != nil {
		return &n.child[c] // Return the pointer to the found child node
	}
	return nil // Return nil if no child with the given key is found
}

// isFull always returns false, as node256 can hold up to 256 children and doesn't need to check for fullness.
func (n *node256) isFull() bool { return false }

// grow attempts to grow the node256, but this operation is not allowed for node256.
// It will panic if called.
func (n *node256) grow() node {
	panic("grow can not be called on node256") // Node256 cannot grow any further
}

// deleteChild removes a child node by its key. It sets the child at the given index to nil and reduces the size.
func (n *node256) deleteChild(c byte) {
	if n.child[c] != nil {
		n.child[c] = nil // Remove the child by setting it to nil
		n.size--         // Decrease the size to reflect the removal
	}
}

// shrink attempts to shrink the node if possible. If the node has 48 or fewer children, it converts to node48.
// Otherwise, it returns nil to indicate shrinking is not possible.
func (n *node256) shrink() node {
	if n.size > 48 {
		return nil // Return nil if shrinking is not possible (more than 48 children)
	}
	nn := newNode48(nil) // Create a new node48 with no prefix
	for c, child := range n.child {
		if child != nil {
			nn.addChild(byte(c), child) // Add each non-nil child to the new node48
		}
	}
	return nn // Return the newly shrunk node (node48)
}

// iter iterates over all children nodes and applies the function f to each of them.
// If the function returns false, the iteration stops.
func (n *node256) iter(f func(node) bool) {
	for i := 0; i < 256; i++ {
		if n.child[i] != nil { // Only call the function for non-nil children
			if !f(n.child[i]) { // Stop iteration if the function returns false
				return
			}
		}
	}
}

// children returns a slice containing all the child nodes. This includes all 256 slots, even if some are nil.
func (n *node256) children() []node {
	return n.child[:256] // Return all children (up to 256)
}
