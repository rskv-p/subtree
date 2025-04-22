package subtree

import (
	"fmt"
	"io"
	"strings"
)

//-------------------
// Dumping a tree structure
//-------------------

// Dump outputs a text representation of the entire tree to the given writer.
// It starts by calling the private 'dump' method with the root node.
func (t *SubjectTree[T]) Dump(w io.Writer) {
	t.dump(w, t.root, 0)
	fmt.Fprintln(w) // Add a newline after dumping the tree
}

//-------------------
// Recursive node dumping
//-------------------

// dump is a recursive function that traverses and prints the nodes of the tree.
// It prints a detailed representation of the current node, whether it's a leaf or another node type.
func (t *SubjectTree[T]) dump(w io.Writer, n node, depth int) {
	if n == nil {
		// If the node is nil, print "EMPTY"
		fmt.Fprintf(w, "EMPTY\n")
		return
	}

	// If the node is a leaf, print its details and stop recursion for this branch.
	if n.isLeaf() {
		leaf := n.(*leaf[T]) // Type assertion to a leaf type
		fmt.Fprintf(w, "%s LEAF: Suffix: %q Value: %+v\n", dumpPre(depth), leaf.suffix, leaf.value)
		n = nil // No further traversal for leaf nodes
	} else {
		// If it's not a leaf, it's a node, so print the prefix of the base node.
		bn := n.base() // Get the base node information
		fmt.Fprintf(w, "%s %s Prefix: %q\n", dumpPre(depth), n.kind(), bn.prefix)
		depth++ // Increase depth for child nodes

		// Iterate through child nodes and recursively call dump for each.
		n.iter(func(n node) bool {
			t.dump(w, n, depth)
			return true
		})
	}
}

//-------------------
// Node type definitions
//-------------------

// The following methods define the "kind" of each node type,
// which is useful for distinguishing between different node classes (e.g., NODE4, NODE16, etc.).
func (n *leaf[T]) kind() string { return "LEAF" }
func (n *node4) kind() string   { return "NODE4" }
func (n *node10) kind() string  { return "NODE10" }
func (n *node16) kind() string  { return "NODE16" }
func (n *node48) kind() string  { return "NODE48" }
func (n *node256) kind() string { return "NODE256" }

//-------------------
// Formatting the tree output
//-------------------

// dumpPre calculates the indentation for the current depth of the tree.
// It helps visually represent the tree structure by adding appropriate spaces and symbols.
func dumpPre(depth int) string {
	if depth == 0 {
		return "-- " // Root node, no indentation
	} else {
		var b strings.Builder
		for i := 0; i < depth; i++ {
			b.WriteString("  ") // Add spaces for each level of depth
		}
		b.WriteString("|__ ") // Visual separator for child nodes
		return b.String()
	}
}
