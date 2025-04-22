package subtree

import (
	"flag"
	"fmt"
	"testing"
)

//-------------------
//  Global flags and settings
//-------------------

// Flag to enable results tests (for benchmarking)
var runResults = flag.Bool("results", false, "Enable Results Tests")

//-------------------
//  Test for Node Prefix Mismatch
//-------------------

// Test case to check how the tree handles node prefix mismatches and ensures proper updates during splits
func TestSubjectTreeNodePrefixMismatch(t *testing.T) {
	st := NewSubjectTree[int]()
	st.Insert(b("foo.bar.A"), 11)
	st.Insert(b("foo.bar.B"), 22)
	st.Insert(b("foo.bar.C"), 33)
	// Grab current root. A split should occur after the following insert.
	or := st.root
	// This insert will force a split of the node
	st.Insert(b("foo.foo.A"), 44)
	require_True(t, or != st.root)
	// Now, check if we can correctly retrieve the inserted values
	v, found := st.Find(b("foo.bar.A"))
	require_True(t, found)
	require_Equal(t, *v, 11)
	v, found = st.Find(b("foo.bar.B"))
	require_True(t, found)
	require_Equal(t, *v, 22)
	v, found = st.Find(b("foo.bar.C"))
	require_True(t, found)
	require_Equal(t, *v, 33)
	v, found = st.Find(b("foo.foo.A"))
	require_True(t, found)
	require_Equal(t, *v, 44)
}

//-------------------
//  Test for Nodes and Path Handling
//-------------------

// Test case to ensure the tree correctly handles node prefixes and paths during insertions and deletions
func TestSubjectTreeNodesAndPaths(t *testing.T) {
	st := NewSubjectTree[int]()
	check := func(subj string) {
		t.Helper()
		v, found := st.Find(b(subj))
		require_True(t, found)
		require_Equal(t, *v, 22)
	}
	st.Insert(b("foo.bar.A"), 22)
	st.Insert(b("foo.bar.B"), 22)
	st.Insert(b("foo.bar.C"), 22)
	st.Insert(b("foo.bar"), 22)
	check("foo.bar.A")
	check("foo.bar.B")
	check("foo.bar.C")
	check("foo.bar")
	// This will perform several actions: shrinking and pruning.
	// We want to ensure the prefix is correct after the new top node4 is created.
	st.Delete(b("foo.bar"))
	check("foo.bar.A")
	check("foo.bar.B")
	check("foo.bar.C")
}

//-------------------
//  Test for Tree Construction with Complex Insert Patterns
//-------------------

// Test case to verify that the tree correctly handles complex insert patterns and constructs the tree properly
func TestSubjectTreeConstruction(t *testing.T) {
	st := NewSubjectTree[int]()
	st.Insert(b("foo.bar.A"), 1)
	st.Insert(b("foo.bar.B"), 2)
	st.Insert(b("foo.bar.C"), 3)
	st.Insert(b("foo.baz.A"), 11)
	st.Insert(b("foo.baz.B"), 22)
	st.Insert(b("foo.baz.C"), 33)
	st.Insert(b("foo.bar"), 42)

	checkNode := func(an *node, kind string, pors string, numChildren uint16) {
		t.Helper()
		require_True(t, an != nil)
		n := *an
		require_True(t, n != nil)
		require_Equal(t, n.kind(), kind)
		require_Equal(t, pors, string(n.path()))
		require_Equal(t, numChildren, n.numChildren())
	}

	// Check root node and its children
	checkNode(&st.root, "NODE4", "foo.ba", 2)
	nn := st.root.findChild('r')
	checkNode(nn, "NODE4", "r", 2)
	checkNode((*nn).findChild(noPivot), "LEAF", "", 0)
	rnn := (*nn).findChild('.')
	checkNode(rnn, "NODE4", ".", 3)
	checkNode((*rnn).findChild('A'), "LEAF", "A", 0)
	checkNode((*rnn).findChild('B'), "LEAF", "B", 0)
	checkNode((*rnn).findChild('C'), "LEAF", "C", 0)
	znn := st.root.findChild('z')
	checkNode(znn, "NODE4", "z.", 3)
	checkNode((*znn).findChild('A'), "LEAF", "A", 0)
	checkNode((*znn).findChild('B'), "LEAF", "B", 0)
	checkNode((*znn).findChild('C'), "LEAF", "C", 0)
	// After deletion, ensure the tree is reconstructed correctly
	v, found := st.Delete(b("foo.bar"))
	require_True(t, found)
	require_Equal(t, *v, 42)

	checkNode(&st.root, "NODE4", "foo.ba", 2)
	nn = st.root.findChild('r')
	checkNode(nn, "NODE4", "r.", 3)
	checkNode((*nn).findChild('A'), "LEAF", "A", 0)
	checkNode((*nn).findChild('B'), "LEAF", "B", 0)
	checkNode((*nn).findChild('C'), "LEAF", "C", 0)
	znn = st.root.findChild('z')
	checkNode(znn, "NODE4", "z.", 3)
	checkNode((*znn).findChild('A'), "LEAF", "A", 0)
	checkNode((*znn).findChild('B'), "LEAF", "B", 0)
	checkNode((*znn).findChild('C'), "LEAF", "C", 0)
}

//-------------------
//  Match Helper Function
//-------------------

// Helper function to match a filter against the tree and check the number of matches
func match(t *testing.T, st *SubjectTree[int], filter string, expected int) {
	t.Helper()
	var matches []int
	st.Match(b(filter), func(_ []byte, v *int) {
		matches = append(matches, *v)
	})
	require_Equal(t, expected, len(matches))
}

//-------------------
//  Test for Tree with No Prefix
//-------------------

// Test case to check tree behavior when there is no prefix and many insertions
func TestSubjectTreeNoPrefix(t *testing.T) {
	st := NewSubjectTree[int]()
	for i := 0; i < 26; i++ {
		subj := b(fmt.Sprintf("%c", 'A'+i))
		old, updated := st.Insert(subj, 22)
		require_True(t, old == nil)
		require_False(t, updated)
	}
	n, ok := st.root.(*node48)
	require_True(t, ok)
	require_Equal(t, n.numChildren(), 26)
	v, found := st.Delete(b("B"))
	require_True(t, found)
	require_Equal(t, *v, 22)
	require_Equal(t, n.numChildren(), 25)
	v, found = st.Delete(b("Z"))
	require_True(t, found)
	require_Equal(t, *v, 22)
	require_Equal(t, n.numChildren(), 24)
}

//-------------------
//  Bug Test for Partial Terminal Wildcard Match
//-------------------

// Test case to validate bug with partial terminal wildcard matches
func TestSubjectTreePartialTerminalWildcardBugMatch(t *testing.T) {
	st := NewSubjectTree[int]()
	st.Insert(b("STATE.GLOBAL.CELL1.7PDSGAALXNN000010.PROPERTY-A"), 5)
	st.Insert(b("STATE.GLOBAL.CELL1.7PDSGAALXNN000010.PROPERTY-B"), 1)
	st.Insert(b("STATE.GLOBAL.CELL1.7PDSGAALXNN000010.PROPERTY-C"), 2)
	match(t, st, "STATE.GLOBAL.CELL1.7PDSGAALXNN000010.*", 3)
}

//-------------------
//  Test for Iteration Ordered
//-------------------

// Test case to check the ordered iteration of elements in the tree
func TestSubjectTreeIterOrdered(t *testing.T) {
	st := NewSubjectTree[int]()
	st.Insert(b("foo.bar.A"), 1)
	st.Insert(b("foo.bar.B"), 2)
	st.Insert(b("foo.bar.C"), 3)
	st.Insert(b("foo.baz.A"), 11)
	st.Insert(b("foo.baz.B"), 22)
	st.Insert(b("foo.baz.C"), 33)
	st.Insert(b("foo.bar"), 42)

	checkValMap := map[string]int{
		"foo.bar.A": 1,
		"foo.bar.B": 2,
		"foo.bar.C": 3,
		"foo.baz.A": 11,
		"foo.baz.B": 22,
		"foo.baz.C": 33,
		"foo.bar":   42,
	}
	checkOrder := []string{
		"foo.bar",
		"foo.bar.A",
		"foo.bar.B",
		"foo.bar.C",
		"foo.baz.A",
		"foo.baz.B",
		"foo.baz.C",
	}
	var received int
	walk := func(subject []byte, v *int) bool {
		if expected := checkOrder[received]; expected != string(subject) {
			t.Fatalf("Expected %q for %d item returned, got %q", expected, received, subject)
		}
		received++
		require_True(t, v != nil)
		if expected := checkValMap[string(subject)]; expected != *v {
			t.Fatalf("Expected %q to have value of %d, but got %d", subject, expected, *v)
		}
		return true
	}
	// Kick in the iter.
	st.IterOrdered(walk)
	require_Equal(t, received, len(checkOrder))

	// Make sure we can terminate properly.
	received = 0
	st.IterOrdered(func(subject []byte, v *int) bool {
		received++
		return received != 4
	})
	require_Equal(t, received, 4)
}

//-------------------
//  Test for Iteration Fast
//-------------------

// Test case to check fast iteration of elements in the tree
func TestSubjectTreeIterFast(t *testing.T) {
	st := NewSubjectTree[int]()
	st.Insert(b("foo.bar.A"), 1)
	st.Insert(b("foo.bar.B"), 2)
	st.Insert(b("foo.bar.C"), 3)
	st.Insert(b("foo.baz.A"), 11)
	st.Insert(b("foo.baz.B"), 22)
	st.Insert(b("foo.baz.C"), 33)
	st.Insert(b("foo.bar"), 42)

	checkValMap := map[string]int{
		"foo.bar.A": 1,
		"foo.bar.B": 2,
		"foo.bar.C": 3,
		"foo.baz.A": 11,
		"foo.baz.B": 22,
		"foo.baz.C": 33,
		"foo.bar":   42,
	}
	var received int
	walk := func(subject []byte, v *int) bool {
		received++
		require_True(t, v != nil)
		if expected := checkValMap[string(subject)]; expected != *v {
			t.Fatalf("Expected %q to have value of %d, but got %d", subject, expected, *v)
		}
		return true
	}
	// Kick in the iter.
	st.IterFast(walk)
	require_Equal(t, received, len(checkValMap))

	// Make sure we can terminate properly.
	received = 0
	st.IterFast(func(subject []byte, v *int) bool {
		received++
		return received != 4
	})
	require_Equal(t, received, 4)
}
