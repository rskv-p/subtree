package subtree

//-------------------
// Node Interface
//-------------------

// The node interface defines the common methods that both leaf and internal nodes must implement.
// These methods handle various operations like matching, adding/removing children, and managing node data.
type node interface {
	isLeaf() bool                               // Returns true if the node is a leaf, false otherwise
	base() *meta                                // Returns the base metadata of the node
	setPrefix(pre []byte)                       // Sets the prefix for the node
	addChild(c byte, n node)                    // Adds a child node for the given character
	findChild(c byte) *node                     // Finds and returns a child node for the given character
	deleteChild(c byte)                         // Deletes a child node for the given character
	isFull() bool                               // Returns true if the node is full (i.e., can no longer hold more children)
	grow() node                                 // Expands the node (e.g., converting it to a larger node type)
	shrink() node                               // Shrinks the node (e.g., converting it to a smaller node type)
	matchParts(parts [][]byte) ([][]byte, bool) // Matches parts against the node's prefix
	kind() string                               // Returns a string identifying the type of the node
	iter(f func(node) bool)                     // Iterates over the children of the node
	children() []node                           // Returns the children of the node
	numChildren() uint16                        // Returns the number of children the node has
	path() []byte                               // Returns the path (or prefix) associated with the node
}

//-------------------
// Meta Data Structure
//-------------------

// The meta struct holds metadata about a node, specifically the prefix and the number of children it has.
type meta struct {
	prefix []byte // The prefix associated with this node
	size   uint16 // The number of children this node has
}

//-------------------
// Meta Methods
//-------------------

// isLeaf returns false because meta nodes are internal nodes.
func (n *meta) isLeaf() bool { return false }

// base returns the meta node itself as the base of internal nodes.
func (n *meta) base() *meta { return n }

// setPrefix sets the prefix for this node by copying the provided byte slice.
func (n *meta) setPrefix(pre []byte) {
	n.prefix = append([]byte(nil), pre...) // Safely copy the prefix to avoid modifying the original slice
}

// numChildren returns the number of children for this meta node.
func (n *meta) numChildren() uint16 { return n.size }

// path returns the prefix of the node.
func (n *meta) path() []byte { return n.prefix }

//-------------------
// Meta Node Matching
//-------------------

// matchParts compares the given parts with the node's prefix and returns the result.
func (n *meta) matchParts(parts [][]byte) ([][]byte, bool) {
	return matchParts(parts, n.prefix) // Delegate the comparison to matchParts function
}
