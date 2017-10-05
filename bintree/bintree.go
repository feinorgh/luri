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

func (n *Node) Insert(number *big.Int) error {
	if n == nil {
		return errors.New("Cannot insert a value into a nil tree")
	}
	c := number.Cmp(&n.Number)
	switch {
	case c == 0:
		return nil
	case c < 0:
		if n.Left == nil {
			n.Left = &Node{Number: *number}
			return nil
		}
		return n.Left.Insert(number)
	case c > 0:
		if n.Right == nil {
			n.Right = &Node{Number: *number}
			return nil
		}
		return n.Right.Insert(number)
	}
	return nil
}

func (n *Node) Find(number *big.Int) bool {
	if n == nil {
		return false
	}
	c := number.Cmp(&n.Number)
	switch {
	case c == 0:
		return true
	case c < 0:
		return n.Left.Find(number)
	default:
		return n.Right.Find(number)
	}

}

func (t *Tree) Insert(number *big.Int) error {
	if t.Root == nil {
		t.Root = &Node{Number: *number}
		return nil
	}
	return t.Root.Insert(number)
}

func (t *Tree) Find(number *big.Int) bool {
	if t.Root == nil {
		return false
	}
	return t.Root.Find(number)
}

func (t *Tree) Traverse(n *Node, f func(*Node)) {
	if n == nil {
		return
	}
	t.Traverse(n.Left, f)
	f(n)
	t.Traverse(n.Right, f)
}
