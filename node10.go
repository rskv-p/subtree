package subtree

//-------------------
// Node10 Definition
//-------------------

// node10 represents a node with 10 possible children, specifically designed for cases where
// part of the subject is numeric (i.e., keys are limited to the range 0-9).
// The struct is optimized for memory alignment according to govet/fieldalignment recommendations.
type node10 struct {
	child [10]node // Array of child nodes (up to 10 children)
	meta           // Inherited metadata (prefix and size)
	key   [10]byte // Array of keys corresponding to the children (numeric range 0-9)
}

//-------------------
// Node10 Methods
//-------------------

// newNode10 creates a new node10 with the specified prefix and returns a pointer to it.
func newNode10(prefix []byte) *node10 {
	nn := &node10{}
	nn.setPrefix(prefix) // Set the prefix for the node
	return nn
}

// addChild adds a child node to the current node. It appends the node at the next available position.
// It will panic if the node already has 10 children (node is full).
func (n *node10) addChild(c byte, nn node) {
	if n.size >= 10 {
		// Panic if the node has reached its maximum capacity of 10 children
		panic("node10 full!")
	}
	n.key[n.size] = c    // Store the key associated with the child node
	n.child[n.size] = nn // Store the child node itself
	n.size++             // Increment the size to reflect the added child
}

// findChild looks for a child node by its key (byte). If found, it returns a pointer to the child node.
func (n *node10) findChild(c byte) *node {
	for i := uint16(0); i < n.size; i++ {
		if n.key[i] == c {
			return &n.child[i] // Return the pointer to the found child node
		}
	}
	return nil // Return nil if no child with the given key is found
}

// isFull checks if the node has reached its maximum capacity of 10 children.
func (n *node10) isFull() bool { return n.size >= 10 }

// grow converts this node10 into a node16 (a larger node type) when more children are needed.
// It copies over the existing children to the new node16.
func (n *node10) grow() node {
	nn := newNode16(n.prefix) // Create a new node16 with the same prefix
	for i := 0; i < 10; i++ {
		nn.addChild(n.key[i], n.child[i]) // Add each child to the new node16
	}
	return nn // Return the newly grown node
}

// deleteChild removes a child node by its key. It swaps the child with the last one and reduces the size.
func (n *node10) deleteChild(c byte) {
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

// shrink attempts to shrink the node if possible. If the node has 4 or fewer children, it converts to node4.
// Otherwise, it returns nil to indicate shrinking is not possible.
func (n *node10) shrink() node {
	if n.size > 4 {
		return nil // Return nil if shrinking is not possible (more than 4 children)
	}
	nn := newNode4(nil) // Create a new node4 with no prefix
	for i := uint16(0); i < n.size; i++ {
		nn.addChild(n.key[i], n.child[i]) // Add each child to the new node4
	}
	return nn // Return the newly shrunk node (node4)
}

// iter iterates over all children nodes and applies the function f to each of them.
// If the function returns false, the iteration stops.
func (n *node10) iter(f func(node) bool) {
	for i := uint16(0); i < n.size; i++ {
		if !f(n.child[i]) { // Call the function for each child, stop if it returns false
			return
		}
	}
}

// children returns a slice containing all the child nodes.
func (n *node10) children() []node {
	return n.child[:n.size] // Return only the children that are currently in use (up to 'size')
}
