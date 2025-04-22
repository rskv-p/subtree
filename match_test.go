package subtree

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

//-------------------
//  Test for Matching Leaf Only
//-------------------

// Test case to match leaf-only subjects with different wildcard placements.
func TestSubjectTreeMatchLeafOnly(t *testing.T) {
	st := NewSubjectTree[int]()
	st.Insert(b("foo.bar.baz.A"), 1)

	// Check all placements of pwc (partial wildcard) in token space.
	match(t, st, "foo.bar.*.A", 1)
	match(t, st, "foo.*.baz.A", 1)
	match(t, st, "foo.*.*.A", 1)
	match(t, st, "foo.*.*.*", 1)
	match(t, st, "*.*.*.*", 1)

	// Check fwc (full wildcard).
	match(t, st, ">", 1)
	match(t, st, "foo.>", 1)
	match(t, st, "foo.*.>", 1)
	match(t, st, "foo.bar.>", 1)
	match(t, st, "foo.bar.*.>", 1)

	// Check partials to ensure they do not match on leaf nodes.
	match(t, st, "foo.bar.baz", 0)
}

//-------------------
//  Test for Matching Nodes
//-------------------

// Test case to match nodes and check internal and terminal wildcards.
func TestSubjectTreeMatchNodes(t *testing.T) {
	st := NewSubjectTree[int]()
	st.Insert(b("foo.bar.A"), 1)
	st.Insert(b("foo.bar.B"), 2)
	st.Insert(b("foo.bar.C"), 3)
	st.Insert(b("foo.baz.A"), 11)
	st.Insert(b("foo.baz.B"), 22)
	st.Insert(b("foo.baz.C"), 33)

	// Test literals.
	match(t, st, "foo.bar.A", 1)
	match(t, st, "foo.baz.A", 1)
	match(t, st, "foo.bar", 0)

	// Test internal pwc (partial wildcard).
	match(t, st, "foo.*.A", 2)

	// Test terminal pwc (partial wildcard at the end).
	match(t, st, "foo.bar.*", 3)
	match(t, st, "foo.baz.*", 3)

	// Check fwc (full wildcard).
	match(t, st, ">", 6)
	match(t, st, "foo.>", 6)
	match(t, st, "foo.bar.>", 3)
	match(t, st, "foo.baz.>", 3)

	// Make sure prefix matches don't cause false positives.
	match(t, st, "foo.ba", 0)

	// Add "foo.bar" to make a more complex tree construction and re-test.
	st.Insert(b("foo.bar"), 42)

	// Test literals again.
	match(t, st, "foo.bar.A", 1)
	match(t, st, "foo.baz.A", 1)
	match(t, st, "foo.bar", 1)

	// Test internal pwc (partial wildcard).
	match(t, st, "foo.*.A", 2)

	// Test terminal pwc.
	match(t, st, "foo.bar.*", 3)
	match(t, st, "foo.baz.*", 3)

	// Check fwc.
	match(t, st, ">", 7)
	match(t, st, "foo.>", 7)
	match(t, st, "foo.bar.>", 3)
	match(t, st, "foo.baz.>", 3)
}

//-------------------
//  Test for Matching Subject Parameters
//-------------------

// Test case to match and check subject parameters.
func TestSubjectTreeMatchSubjectParam(t *testing.T) {
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

	// Make sure the subject parameter matches the correct value.
	st.Match([]byte(">"), func(subject []byte, v *int) {
		if expected, ok := checkValMap[string(subject)]; !ok {
			t.Fatalf("Unexpected subject parameter: %q", subject)
		} else if expected != *v {
			t.Fatalf("Expected %q to have value of %d, but got %d", subject, expected, *v)
		}
	})
}

//-------------------
//  Test for Matching Random Double PWC (Partial Wildcard)
//-------------------

// Test case for matching using random double pwc and checking for correctness.
func TestSubjectTreeMatchRandomDoublePWC(t *testing.T) {
	st := NewSubjectTree[int]()
	for i := 1; i <= 10_000; i++ {
		subj := fmt.Sprintf("foo.%d.%d", rand.Intn(20)+1, i)
		st.Insert(b(subj), 42)
	}
	match(t, st, "foo.*.*", 10_000)

	// Check with pwc and short interior token.
	seen, verified := 0, 0
	st.Match(b("*.2.*"), func(_ []byte, _ *int) {
		seen++
	})
	// Now check via walk to ensure accuracy.
	st.IterOrdered(func(subject []byte, v *int) bool {
		tokens := strings.Split(string(subject), ".")
		require_Equal(t, len(tokens), 3)
		if tokens[1] == "2" {
			verified++
		}
		return true
	})
	require_Equal(t, seen, verified)

	// Check with another pattern for matching.
	seen, verified = 0, 0
	st.Match(b("*.*.222"), func(_ []byte, _ *int) {
		seen++
	})
	st.IterOrdered(func(subject []byte, v *int) bool {
		tokens := strings.Split(string(subject), ".")
		require_Equal(t, len(tokens), 3)
		if tokens[2] == "222" {
			verified++
		}
		return true
	})
	require_Equal(t, seen, verified)
}

//-------------------
//  Test for Matching Invalid Wildcards
//-------------------

// Test case to check for invalid wildcard usage in subject matching.
func TestSubjectTreeMatchInvalidWildcard(t *testing.T) {
	st := NewSubjectTree[int]()
	st.Insert(b("foo.123"), 22)
	st.Insert(b("one.two.three.four.five"), 22)
	st.Insert(b("'*.123"), 22)
	match(t, st, "invalid.>", 0)
	match(t, st, ">", 3)
	match(t, st, `'*.*`, 1)
	match(t, st, `'*.*.*'`, 0)
	// None of these should match.
	match(t, st, "`>`", 0)
	match(t, st, `">"`, 0)
	match(t, st, `'>'`, 0)
	match(t, st, `'*.>'`, 0)
	match(t, st, `'*.>.`, 0)
	match(t, st, "`invalid.>`", 0)
	match(t, st, `'*.*'`, 0)
}

//-------------------
//  Test for Multiple Wildcard Match Basic
//-------------------

// Test case for basic matching with multiple wildcards.
func TestSubjectTreeMatchMultipleWildcardBasic(t *testing.T) {
	st := NewSubjectTree[int]()
	st.Insert(b("A.B.C.D.0.G.H.I.0"), 22)
	st.Insert(b("A.B.C.D.1.G.H.I.0"), 22)
	match(t, st, "A.B.*.D.1.*.*.I.0", 1)
}

//-------------------
//  Test for Long Tokens in SubjectTree
//-------------------

// Test case for inserting and deleting subjects with long tokens in the tree.
func TestSubjectTreeLongTokens(t *testing.T) {
	st := NewSubjectTree[int]()
	st.Insert(b("a1.aaaaaaaaaaaaaaaaaaaaaa0"), 1)
	st.Insert(b("a2.0"), 2)
	st.Insert(b("a1.aaaaaaaaaaaaaaaaaaaaaa1"), 3)
	st.Insert(b("a2.1"), 4)
	// Simulate purging of "a2.>"
	// This required to show bug.
	st.Delete(b("a2.0"))
	st.Delete(b("a2.1"))
	require_Equal(t, st.Size(), 2)
	v, found := st.Find(b("a1.aaaaaaaaaaaaaaaaaaaaaa0"))
	require_True(t, found)
	require_Equal(t, *v, 1)
	v, found = st.Find(b("a1.aaaaaaaaaaaaaaaaaaaaaa1"))
	require_True(t, found)
	require_Equal(t, *v, 3)
}

//-------------------
// Test: Performance of Iteration Over Subject Tree
//-------------------

// Test case to measure the performance of iterating over 1 million entries in the SubjectTree.
func TestSubjectTreeIterPerf(t *testing.T) {
	// Skip the test if results are not enabled in the flags
	if !*runResults {
		t.Skip()
	}

	// Create a new SubjectTree instance.
	st := NewSubjectTree[int]()

	// Insert 1 million random subjects into the tree.
	for i := 0; i < 1_000_000; i++ {
		subj := fmt.Sprintf("subj.%d.%d", rand.Intn(100)+1, i)
		st.Insert(b(subj), 22)
	}

	// Measure the time taken for iteration.
	start := time.Now()
	count := 0

	// Iterate over the tree and count the entries.
	st.IterOrdered(func(_ []byte, _ *int) bool {
		count++
		return true
	})

	// Log the time taken and the number of matched entries.
	t.Logf("Iter took %s and matched %d entries", time.Since(start), count)
}

//-------------------
//  Test Helper Functions
//-------------------

// require_True is a helper function that asserts the given boolean value is true.
// If the value is false, the test will fail with an error message.
func require_True(t *testing.T, b bool) {
	t.Helper() // Marks this function as a helper function for clearer stack traces.
	if !b {
		// Fail the test if the condition is false.
		t.Fatalf("require true, but got false")
	}
}

// require_False is a helper function that asserts the given boolean value is false.
// If the value is true, the test will fail with an error message.
func require_False(t *testing.T, b bool) {
	t.Helper() // Marks this function as a helper function for clearer stack traces.
	if b {
		// Fail the test if the condition is true.
		t.Fatalf("require false, but got true")
	}
}

// require_Equal is a generic helper function that asserts two values of any comparable type are equal.
// If they are not equal, the test will fail with an error message.
func require_Equal[T comparable](t *testing.T, a, b T) {
	t.Helper() // Marks this function as a helper function for clearer stack traces.
	if a != b {
		// Fail the test if the values are not equal, and include their types and values in the error message.
		t.Fatalf("require %T equal, but got: %v != %v", a, a, b)
	}
}

// b is a simple helper function to convert a string to a byte slice.
func b(s string) []byte {
	return []byte(s)
}
