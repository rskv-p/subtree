package subtree

//-------------------
// Node48 Definition
//-------------------

// node48 represents a node with 48 possible children. It is optimized for memory usage,
// as the child array is 16 bytes per node entry, resulting in smaller memory usage
// compared to node256. The key array is used for mapping keys to children,
// with 0 meaning no entry and thus effectively making the key array 1-indexed.
// The struct is optimized for memory alignment according to govet/fieldalignment recommendations.
type node48 struct {
	child [48]node  // Array of child nodes (up to 48 children)
	meta            // Inherited metadata (prefix and size)
	key   [256]byte // Array of keys, 1-indexed (0 means no entry)
}

//-------------------
// Node48 Methods
//-------------------

// newNode48 creates a new node48 with the specified prefix and returns a pointer to it.
func newNode48(prefix []byte) *node48 {
	nn := &node48{}
	nn.setPrefix(prefix) // Set the prefix for the node
	return nn
}

// addChild adds a child node to the current node. It appends the node at the next available position.
// It will panic if the node already has 48 children (node is full).
func (n *node48) addChild(c byte, nn node) {
	if n.size >= 48 {
		// Panic if the node has reached its maximum capacity of 48 children
		panic("node48 full!")
	}
	n.child[n.size] = nn        // Store the child node
	n.key[c] = byte(n.size + 1) // 1-indexed key (0 means no entry)
	n.size++                    // Increment the size to reflect the added child
}

// findChild looks for a child node by its key (byte). If found, it returns a pointer to the child node.
func (n *node48) findChild(c byte) *node {
	i := n.key[c]
	if i == 0 {
		return nil // Return nil if the child doesn't exist
	}
	return &n.child[i-1] // Adjust for 1-indexing and return the child node
}

// isFull checks if the node has reached its maximum capacity of 48 children.
func (n *node48) isFull() bool { return n.size >= 48 }

// grow converts this node48 into a node256 (a larger node type) when more children are needed.
// It copies over the existing children to the new node256.
func (n *node48) grow() node {
	nn := newNode256(n.prefix) // Create a new node256 with the same prefix
	for c := 0; c < len(n.key); c++ {
		if i := n.key[byte(c)]; i > 0 {
			nn.addChild(byte(c), n.child[i-1]) // Add each child to the new node256
		}
	}
	return nn // Return the newly grown node
}

// deleteChild removes a child node by its key. It adjusts the remaining children accordingly.
func (n *node48) deleteChild(c byte) {
	i := n.key[c]
	if i == 0 {
		return // If no child exists with the key, do nothing
	}
	i-- // Adjust for 1-indexing
	last := byte(n.size - 1)
	if i < last {
		n.child[i] = n.child[last] // Swap the child with the last one
		// Update the key array to reflect the swap
		for ic := 0; ic < len(n.key); ic++ {
			if n.key[byte(ic)] == last+1 {
				n.key[byte(ic)] = i + 1
				break
			}
		}
	}
	n.child[last] = nil // Set the last child to nil
	n.key[c] = 0        // Remove the key
	n.size--            // Decrease the size to reflect the removal
}

// shrink attempts to shrink the node if possible. If the node has 16 or fewer children, it converts to node16.
// Otherwise, it returns nil to indicate shrinking is not possible.
func (n *node48) shrink() node {
	if n.size > 16 {
		return nil // Return nil if shrinking is not possible (more than 16 children)
	}
	nn := newNode16(nil) // Create a new node16 with no prefix
	for c := 0; c < len(n.key); c++ {
		if i := n.key[byte(c)]; i > 0 {
			nn.addChild(byte(c), n.child[i-1]) // Add each child to the new node16
		}
	}
	return nn // Return the newly shrunk node (node16)
}

// iter iterates over all children nodes and applies the function f to each of them.
// If the function returns false, the iteration stops.
func (n *node48) iter(f func(node) bool) {
	for _, c := range n.child {
		if c != nil && !f(c) { // Call the function for each child, stop if it returns false
			return
		}
	}
}

// children returns a slice containing all the child nodes.
func (n *node48) children() []node {
	return n.child[:n.size] // Return only the children that are currently in use (up to 'size')
}
