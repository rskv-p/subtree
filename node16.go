package subtree

//-------------------
// Node16 Definition
//-------------------

// node16 represents a node with 16 possible children. It is designed to handle situations
// where more children are needed than the previous nodes (e.g., node4 or node10).
// The struct is optimized for memory alignment according to govet/fieldalignment recommendations.
type node16 struct {
	child [16]node // Array of child nodes (up to 16 children)
	meta           // Inherited metadata (prefix and size)
	key   [16]byte // Array of keys corresponding to the children
}

//-------------------
// Node16 Methods
//-------------------

// newNode16 creates a new node16 with the specified prefix and returns a pointer to it.
func newNode16(prefix []byte) *node16 {
	nn := &node16{}
	nn.setPrefix(prefix) // Set the prefix for the node
	return nn
}

// addChild adds a child node to the current node. It appends the node at the next available position.
// It will panic if the node already has 16 children (node is full).
func (n *node16) addChild(c byte, nn node) {
	if n.size >= 16 {
		// Panic if the node has reached its maximum capacity of 16 children
		panic("node16 full!")
	}
	n.key[n.size] = c    // Store the key associated with the child node
	n.child[n.size] = nn // Store the child node itself
	n.size++             // Increment the size to reflect the added child
}

// findChild looks for a child node by its key (byte). If found, it returns a pointer to the child node.
func (n *node16) findChild(c byte) *node {
	for i := uint16(0); i < n.size; i++ {
		if n.key[i] == c {
			return &n.child[i] // Return the pointer to the found child node
		}
	}
	return nil // Return nil if no child with the given key is found
}

// isFull checks if the node has reached its maximum capacity of 16 children.
func (n *node16) isFull() bool { return n.size >= 16 }

// grow converts this node16 into a node48 (a larger node type) when more children are needed.
// It copies over the existing children to the new node48.
func (n *node16) grow() node {
	nn := newNode48(n.prefix) // Create a new node48 with the same prefix
	for i := 0; i < 16; i++ {
		nn.addChild(n.key[i], n.child[i]) // Add each child to the new node48
	}
	return nn // Return the newly grown node
}

// deleteChild removes a child node by its key. It swaps the child with the last one and reduces the size.
func (n *node16) deleteChild(c byte) {
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

// shrink attempts to shrink the node if possible. If the node has 10 or fewer children, it converts to node10.
// Otherwise, it returns nil to indicate shrinking is not possible.
func (n *node16) shrink() node {
	if n.size > 10 {
		return nil // Return nil if shrinking is not possible (more than 10 children)
	}
	nn := newNode10(nil) // Create a new node10 with no prefix
	for i := uint16(0); i < n.size; i++ {
		nn.addChild(n.key[i], n.child[i]) // Add each child to the new node10
	}
	return nn // Return the newly shrunk node (node10)
}

// iter iterates over all children nodes and applies the function f to each of them.
// If the function returns false, the iteration stops.
func (n *node16) iter(f func(node) bool) {
	for i := uint16(0); i < n.size; i++ {
		if !f(n.child[i]) { // Call the function for each child, stop if it returns false
			return
		}
	}
}

// children returns a slice containing all the child nodes.
func (n *node16) children() []node {
	return n.child[:n.size] // Return only the children that are currently in use (up to 'size')
}
