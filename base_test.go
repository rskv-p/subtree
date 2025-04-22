package subtree

import (
	"fmt"
	"testing"
)

//-------------------
//  Test for SubjectTree Basics
//-------------------

// Test basic insertions, updates, and lookups on the subject tree.
func TestSubjectTreeBasics(t *testing.T) {
	st := NewSubjectTree[int]()
	require_Equal(t, st.Size(), 0)
	// Insert a single leaf node.
	old, updated := st.Insert(b("foo.bar.baz"), 22)
	require_True(t, old == nil)
	require_False(t, updated)
	require_Equal(t, st.Size(), 1)
	// Find should not work with a wildcard.
	_, found := st.Find(b("foo.bar.*"))
	require_False(t, found)
	// But it should work with a literal.
	v, found := st.Find(b("foo.bar.baz"))
	require_True(t, found)
	require_Equal(t, *v, 22)
	// Update the existing leaf.
	old, updated = st.Insert(b("foo.bar.baz"), 33)
	require_True(t, old != nil)
	require_Equal(t, *old, 22)
	require_True(t, updated)
	require_Equal(t, st.Size(), 1)
	// Split the tree by adding a new prefix.
	old, updated = st.Insert(b("foo.bar"), 22)
	require_True(t, old == nil)
	require_False(t, updated)
	require_Equal(t, st.Size(), 2)
	// Verify that both the prefix and the leaf can be retrieved.
	v, found = st.Find(b("foo.bar"))
	require_True(t, found)
	require_Equal(t, *v, 22)
	v, found = st.Find(b("foo.bar.baz"))
	require_True(t, found)
	require_Equal(t, *v, 33)
}

//-------------------
//  Test for Node Growth in SubjectTree
//-------------------

// Test how the tree grows from a small node4 to a larger node structure.
func TestSubjectTreeNodeGrow(t *testing.T) {
	st := NewSubjectTree[int]()
	for i := 0; i < 4; i++ {
		subj := b(fmt.Sprintf("foo.bar.%c", 'A'+i))
		old, updated := st.Insert(subj, 22)
		require_True(t, old == nil)
		require_False(t, updated)
	}
	// At this point, we should have filled a node4.
	_, ok := st.root.(*node4)
	require_True(t, ok)
	// Insert another subject to trigger growth to node10.
	old, updated := st.Insert(b("foo.bar.E"), 22)
	require_True(t, old == nil)
	require_False(t, updated)
	_, ok = st.root.(*node10)
	require_True(t, ok)
	// Insert additional subjects to fill a node10.
	for i := 5; i < 10; i++ {
		subj := b(fmt.Sprintf("foo.bar.%c", 'A'+i))
		old, updated := st.Insert(subj, 22)
		require_True(t, old == nil)
		require_False(t, updated)
	}
	// Trigger growth to node16.
	old, updated = st.Insert(b("foo.bar.K"), 22)
	require_True(t, old == nil)
	require_False(t, updated)
	_, ok = st.root.(*node16)
	require_True(t, ok)
	// Insert more subjects to fill a node16.
	for i := 11; i < 16; i++ {
		subj := b(fmt.Sprintf("foo.bar.%c", 'A'+i))
		old, updated := st.Insert(subj, 22)
		require_True(t, old == nil)
		require_False(t, updated)
	}
	// Trigger growth to node48.
	old, updated = st.Insert(b("foo.bar.Q"), 22)
	require_True(t, old == nil)
	require_False(t, updated)
	_, ok = st.root.(*node48)
	require_True(t, ok)
	// Fill the node48 with subjects.
	for i := 17; i < 48; i++ {
		subj := b(fmt.Sprintf("foo.bar.%c", 'A'+i))
		old, updated := st.Insert(subj, 22)
		require_True(t, old == nil)
		require_False(t, updated)
	}
	// Trigger growth to node256.
	subj := b(fmt.Sprintf("foo.bar.%c", 'A'+49))
	old, updated = st.Insert(subj, 22)
	require_True(t, old == nil)
	require_False(t, updated)
	_, ok = st.root.(*node256)
	require_True(t, ok)
}

//-------------------
//  Test for Node Deletion in SubjectTree
//-------------------

// Test case for deleting nodes in the subject tree and shrinking back to smaller nodes.
func TestSubjectTreeNodeDelete(t *testing.T) {
	st := NewSubjectTree[int]()
	st.Insert(b("foo.bar.A"), 22)
	v, found := st.Delete(b("foo.bar.A"))
	require_True(t, found)
	require_Equal(t, *v, 22)
	require_Equal(t, st.root, nil)
	v, found = st.Delete(b("foo.bar.A"))
	require_False(t, found)
	require_Equal(t, v, nil)
	v, found = st.Find(b("foo.foo.A"))
	require_False(t, found)
	require_Equal(t, v, nil)
	// Test node4 and shrink after deletions.
	st.Insert(b("foo.bar.A"), 11)
	st.Insert(b("foo.bar.B"), 22)
	st.Insert(b("foo.bar.C"), 33)
	// Delete and check shrinkage back to leaf.
	v, found = st.Delete(b("foo.bar.C"))
	require_True(t, found)
	require_Equal(t, *v, 33)
	v, found = st.Delete(b("foo.bar.B"))
	require_True(t, found)
	require_Equal(t, *v, 22)
	require_True(t, st.root.isLeaf())
	v, found = st.Delete(b("foo.bar.A"))
	require_True(t, found)
	require_Equal(t, *v, 11)
	require_Equal(t, st.root, nil)
	// Shrink back up to a node10.
	for i := 0; i < 5; i++ {
		subj := fmt.Sprintf("foo.bar.%c", 'A'+i)
		st.Insert(b(subj), 22)
	}
	_, ok := st.root.(*node10)
	require_True(t, ok)
	v, found = st.Delete(b("foo.bar.A"))
	require_True(t, found)
	require_Equal(t, *v, 22)
	_, ok = st.root.(*node4)
	require_True(t, ok)
	// Shrink to node16.
	for i := 0; i < 11; i++ {
		subj := fmt.Sprintf("foo.bar.%c", 'A'+i)
		st.Insert(b(subj), 22)
	}
	_, ok = st.root.(*node16)
	require_True(t, ok)
	v, found = st.Delete(b("foo.bar.A"))
	require_True(t, found)
	require_Equal(t, *v, 22)
	_, ok = st.root.(*node10)
	require_True(t, ok)
	v, found = st.Find(b("foo.bar.B"))
	require_True(t, found)
	require_Equal(t, *v, 22)
}

//-------------------
//  Test for Node48 Operations
//-------------------

// Test for operations on node48, adding and deleting children.
func TestSubjectTreeNode48(t *testing.T) {
	var a, b, c leaf[int]
	var n node48

	// Add children to node48 and check their correct placement.
	n.addChild('A', &a)
	require_Equal(t, n.key['A'], 1)
	require_True(t, n.child[0] != nil)
	require_Equal(t, n.child[0].(*leaf[int]), &a)
	require_Equal(t, len(n.children()), 1)

	// Find child 'A' and ensure correct leaf node is returned.
	child := n.findChild('A')
	require_True(t, child != nil)
	require_Equal(t, (*child).(*leaf[int]), &a)

	// Add more children to the node.
	n.addChild('B', &b)
	require_Equal(t, n.key['B'], 2)
	require_True(t, n.child[1] != nil)
	require_Equal(t, n.child[1].(*leaf[int]), &b)
	require_Equal(t, len(n.children()), 2)

	// Delete child 'A' and verify the node shrinks correctly.
	n.deleteChild('A')
	require_Equal(t, len(n.children()), 2)
	require_Equal(t, n.key['A'], 0) // Now deleted
	require_Equal(t, n.key['B'], 2) // Untouched
	require_Equal(t, n.key['C'], 1) // Where 'A' was

	// Ensure the proper children remain.
	child = n.findChild('A')
	require_Equal(t, child, nil)
	require_True(t, n.child[0] != nil)
	require_Equal(t, n.child[0].(*leaf[int]), &c)
}

//-------------------
//  Test for Insert with Longer Leaf Suffix and Trailing Nulls
//-------------------

// Test for inserting subjects with longer suffixes containing trailing nulls.
func TestSubjectTreeInsertLongerLeafSuffixWithTrailingNulls(t *testing.T) {
	st := NewSubjectTree[int]()
	subj := []byte("foo.bar.baz_")
	// Add in 10 nulls.
	for i := 0; i < 10; i++ {
		subj = append(subj, 0)
	}

	st.Insert(subj, 1)
	// Add in 10 more nulls.
	subj2 := subj
	for i := 0; i < 10; i++ {
		subj2 = append(subj, 0)
	}
	st.Insert(subj2, 2)

	// Ensure both subjects can be found.
	v, found := st.Find(subj)
	require_True(t, found)
	require_Equal(t, *v, 1)
	v, found = st.Find(subj2)
	require_True(t, found)
	require_Equal(t, *v, 2)
}

//-------------------
//  Test for Inserting with noPivot
//-------------------

// Test case to ensure no subjects with the noPivot (DEL) value are inserted.
func TestSubjectTreeInsertWithNoPivot(t *testing.T) {
	st := NewSubjectTree[int]()
	subj := []byte("foo.bar.baz.")
	subj = append(subj, noPivot)
	old, updated := st.Insert(subj, 22)
	require_True(t, old == nil)
	require_False(t, updated)
	require_Equal(t, st.Size(), 0)
}
