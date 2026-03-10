// Package bintree provides an add-only, unbalanced binary search tree for
// math/big.Int values.
//
// All operations (Insert, Find, Traverse) are implemented iteratively to avoid
// stack overflows on degenerate (fully skewed) trees that can arise when values
// are inserted in sorted order.
package bintree

import (
	"errors"
	"math/big"
)

type Node struct {
	Number big.Int
	Left   *Node
	Right  *Node
}

type Tree struct {
	Root *Node
}

// Insert inserts a number onto a Node that automatically gets placed into the right position in the tree.
func (n *Node) Insert(number *big.Int) error {
	if n == nil {
		return errors.New("cannot insert into a nil node")
	}
	curr := n
	for {
		c := number.Cmp(&curr.Number)
		switch {
		case c == 0:
			return nil
		case c < 0:
			if curr.Left == nil {
				curr.Left = &Node{Number: *number}
				return nil
			}
			curr = curr.Left
		default:
			if curr.Right == nil {
				curr.Right = &Node{Number: *number}
				return nil
			}
			curr = curr.Right
		}
	}
}

// Find returns true if `number' is already in the tree, false if it isn't.
func (n *Node) Find(number *big.Int) bool {
	curr := n
	for curr != nil {
		c := number.Cmp(&curr.Number)
		switch {
		case c == 0:
			return true
		case c < 0:
			curr = curr.Left
		default:
			curr = curr.Right
		}
	}
	return false
}

// Insert inserts a big.Int into the tree. If the tree does not have a root, this will become the root node.
// If there is already a root node, it delegates to the Node.Insert function.
func (t *Tree) Insert(number *big.Int) error {
	if t.Root == nil {
		t.Root = &Node{Number: *number}
		return nil
	}
	return t.Root.Insert(number)
}

// Find returns true if this number is in the tree, false if it isn't.
func (t *Tree) Find(number *big.Int) bool {
	if t.Root == nil {
		return false
	}
	return t.Root.Find(number)
}

// Traverse visits all nodes in the tree in ascending order, calling f() on each.
// Uses an iterative in-order traversal to avoid stack overflow on large or
// degenerate (linked-list) trees.
func (t *Tree) Traverse(n *Node, f func(*Node)) {
	stack := make([]*Node, 0)
	curr := n
	for curr != nil || len(stack) > 0 {
		for curr != nil {
			stack = append(stack, curr)
			curr = curr.Left
		}
		curr = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		f(curr)
		curr = curr.Right
	}
}
