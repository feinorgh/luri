package bintree

import (
	"math/big"
	"testing"
)

// helpers

func bi(n int64) *big.Int { return big.NewInt(n) }

func mustInsert(t *testing.T, tree *Tree, n int64) {
	t.Helper()
	if err := tree.Insert(bi(n)); err != nil {
		t.Fatalf("Insert(%d) unexpected error: %v", n, err)
	}
}

func collectTraversal(t *Tree) []int64 {
	var out []int64
	t.Traverse(t.Root, func(n *Node) {
		out = append(out, n.Number.Int64())
	})
	return out
}

// --- Tree.Insert / Tree.Find ---

func TestInsertAndFindPresent(t *testing.T) {
	tree := &Tree{}
	mustInsert(t, tree, 42)
	if !tree.Find(bi(42)) {
		t.Error("Find(42) = false, want true")
	}
}

func TestFindAbsent(t *testing.T) {
	tree := &Tree{}
	mustInsert(t, tree, 10)
	if tree.Find(bi(99)) {
		t.Error("Find(99) = true, want false")
	}
}

func TestFindOnEmptyTree(t *testing.T) {
	tree := &Tree{}
	if tree.Find(bi(1)) {
		t.Error("Find on empty tree = true, want false")
	}
}

func TestInsertDuplicateIgnored(t *testing.T) {
	tree := &Tree{}
	mustInsert(t, tree, 7)
	mustInsert(t, tree, 7)
	got := collectTraversal(tree)
	if len(got) != 1 {
		t.Errorf("duplicate insert: got %d nodes, want 1", len(got))
	}
}

func TestInsertFirstBecomesRoot(t *testing.T) {
	tree := &Tree{}
	mustInsert(t, tree, 5)
	if tree.Root == nil {
		t.Fatal("Root is nil after first insert")
	}
	if tree.Root.Number.Int64() != 5 {
		t.Errorf("Root.Number = %d, want 5", tree.Root.Number.Int64())
	}
}

func TestNodeInsertNilReturnsError(t *testing.T) {
	var n *Node
	if err := n.Insert(bi(1)); err == nil {
		t.Error("Insert on nil Node: want error, got nil")
	}
}

// --- Tree.Traverse ---

func TestTraverseEmptyTree(t *testing.T) {
	tree := &Tree{}
	called := 0
	tree.Traverse(tree.Root, func(_ *Node) { called++ })
	if called != 0 {
		t.Errorf("Traverse on empty tree called f %d times, want 0", called)
	}
}

func TestTraverseAscendingOrder(t *testing.T) {
	tree := &Tree{}
	values := []int64{5, 3, 8, 1, 4, 7, 9, 2, 6}
	for _, v := range values {
		mustInsert(t, tree, v)
	}
	got := collectTraversal(tree)
	for i := 1; i < len(got); i++ {
		if got[i] <= got[i-1] {
			t.Errorf("traversal not ascending at index %d: %v", i, got)
			break
		}
	}
	if len(got) != len(values) {
		t.Errorf("got %d nodes, want %d", len(got), len(values))
	}
}

func TestTraverseSingleNode(t *testing.T) {
	tree := &Tree{}
	mustInsert(t, tree, 42)
	got := collectTraversal(tree)
	if len(got) != 1 || got[0] != 42 {
		t.Errorf("single node traversal = %v, want [42]", got)
	}
}

// TestTraverseDegenerateTree inserts values in sorted order, which produces a
// fully right-skewed (linked-list) tree. The recursive implementation would
// overflow the stack; the iterative implementation must handle it.
func TestTraverseDegenerateTree(t *testing.T) {
	const n = 10_000
	tree := &Tree{}
	for i := int64(0); i < n; i++ {
		mustInsert(t, tree, i)
	}
	got := collectTraversal(tree)
	if len(got) != n {
		t.Errorf("degenerate traversal: got %d nodes, want %d", len(got), n)
	}
	for i, v := range got {
		if v != int64(i) {
			t.Errorf("degenerate traversal[%d] = %d, want %d", i, v, i)
			break
		}
	}
}

// TestFindInDegenerateTree ensures iterative Find doesn't stack-overflow.
func TestFindInDegenerateTree(t *testing.T) {
	const n = 10_000
	tree := &Tree{}
	for i := int64(0); i < n; i++ {
		mustInsert(t, tree, i)
	}
	if !tree.Find(bi(n - 1)) {
		t.Errorf("Find(%d) = false in degenerate tree, want true", n-1)
	}
	if tree.Find(bi(n)) {
		t.Errorf("Find(%d) = true in degenerate tree, want false", n)
	}
}

// --- Negative numbers and big values ---

func TestNegativeNumbers(t *testing.T) {
	tree := &Tree{}
	mustInsert(t, tree, -5)
	mustInsert(t, tree, -1)
	mustInsert(t, tree, -10)
	got := collectTraversal(tree)
	want := []int64{-10, -5, -1}
	for i, v := range got {
		if v != want[i] {
			t.Errorf("negative[%d] = %d, want %d", i, v, want[i])
		}
	}
}

func TestBigIntValues(t *testing.T) {
	tree := &Tree{}
	large, _ := new(big.Int).SetString("999999999999999999999999999999", 10)
	small, _ := new(big.Int).SetString("1", 10)
	if err := tree.Insert(large); err != nil {
		t.Fatalf("Insert large: %v", err)
	}
	if err := tree.Insert(small); err != nil {
		t.Fatalf("Insert small: %v", err)
	}
	if !tree.Find(large) {
		t.Error("Find(large) = false, want true")
	}
	if !tree.Find(small) {
		t.Error("Find(small) = false, want true")
	}
}
