package subtree

//-------------------
// Node4 Definition
//-------------------

// node4 represents a node with 4 possible children. It uses an array of nodes and keys for efficient storage.
// The order of the struct fields is optimized for memory alignment as per govet/fieldalignment recommendations.
type node4 struct {
	child [4]node // Array of child nodes (up to 4 children)
	meta          // Inherited metadata (prefix and size)
	key   [4]byte // Array of keys corresponding to the children
}

//-------------------
// Node4 Methods
//-------------------

// newNode4 creates a new node4 with the specified prefix and returns a pointer to it.
func newNode4(prefix []byte) *node4 {
	nn := &node4{}
	nn.setPrefix(prefix) // Set the prefix for the node
	return nn
}

// addChild adds a child node to the current node. It appends the node at the next available position.
// It will panic if there are already 4 children (node is full).
func (n *node4) addChild(c byte, nn node) {
	if n.size >= 4 {
		// Panic if the node has reached its maximum capacity of 4 children
		panic("node4 full!")
	}
	n.key[n.size] = c    // Store the key associated with the child node
	n.child[n.size] = nn // Store the child node itself
	n.size++             // Increment the size to reflect the added child
}

// findChild looks for a child node by its key. If found, it returns a pointer to the child node.
func (n *node4) findChild(c byte) *node {
	for i := uint16(0); i < n.size; i++ {
		if n.key[i] == c {
			return &n.child[i] // Return the pointer to the found child node
		}
	}
	return nil // Return nil if no child with the given key is found
}

// isFull checks if the node has reached its maximum capacity of 4 children.
func (n *node4) isFull() bool { return n.size >= 4 }

// grow converts this node4 into a node10 (a larger node type) when more children are needed.
// It copies over the existing children to the new node10.
func (n *node4) grow() node {
	nn := newNode10(n.prefix) // Create a new node10 with the same prefix
	for i := 0; i < 4; i++ {
		nn.addChild(n.key[i], n.child[i]) // Add each child to the new node10
	}
	return nn // Return the newly grown node
}

// deleteChild removes a child node by its key. It swaps the child with the last one and reduces the size.
func (n *node4) deleteChild(c byte) {
	for i, last := uint16(0), n.size-1; i < n.size; i++ {
		if n.key[i] == c {
			// If the child to be deleted is not the last one, swap with the last child
			if i < last {
				n.key[i] = n.key[last]
				n.child[i] = n.child[last]
				n.key[last] = 0
				n.child[last] = nil
			} else {
				n.key[i] = 0
				n.child[i] = nil
			}
			n.size-- // Decrease the size to reflect the removal
			return
		}
	}
}

// shrink attempts to shrink the node if possible. If the node has only one child, it returns the child node itself.
// Otherwise, it returns nil.
func (n *node4) shrink() node {
	if n.size == 1 {
		return n.child[0] // Return the single child if the node is reduced to one child
	}
	return nil // Return nil if shrinking is not possible
}

// iter iterates over all children nodes and applies the function f to each of them.
// If the function returns false, the iteration stops.
func (n *node4) iter(f func(node) bool) {
	for i := uint16(0); i < n.size; i++ {
		if !f(n.child[i]) { // Call the function for each child, stop if it returns false
			return
		}
	}
}

// children returns a slice containing all the child nodes.
func (n *node4) children() []node {
	return n.child[:n.size] // Return only the children that are currently in use (up to 'size')
}
