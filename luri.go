/*
luri - list unique random integers

Copyright (C) 2017  Pär Karlsson <feinorgh@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"errors"
	"flag"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"time"
)

type Node struct {
	Number big.Int
	Left   *Node
	Right  *Node
}

type Tree struct {
	Root *Node
}

type Options struct {
	lowerBound *big.Int
	upperBound *big.Int
	count      int
	verbose    bool
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

func setOptions(opt *Options) {
	var lower, upper string
	var count int
	var verbose bool
	const (
		defaultLowerBound = "1"
		lowerUsage        = "the lower bound `big.Int`"
		defaultUpperBound = "100"
		upperUsage        = "the upper bound `big.Int`"
		defaultCount      = 1
		countUsage        = "size of set"
		defaultVerbose    = false
		verboseUsage      = "be verbose"
	)
	flag.StringVar(&lower, "lower", defaultLowerBound, lowerUsage)
	flag.StringVar(&lower, "l", defaultLowerBound, lowerUsage+" (shorthand)")
	flag.StringVar(&upper, "upper", defaultUpperBound, upperUsage)
	flag.StringVar(&upper, "u", defaultUpperBound, upperUsage)
	flag.IntVar(&count, "count", defaultCount, countUsage)
	flag.IntVar(&count, "c", defaultCount, countUsage+" (shorthand)")
	flag.BoolVar(&verbose, "verbose", defaultVerbose, verboseUsage)
	flag.BoolVar(&verbose, "v", defaultVerbose, verboseUsage+" (shorthand)")
	flag.Parse()
	opt.lowerBound = new(big.Int)
	opt.lowerBound.SetString(lower, 10)
	opt.upperBound = new(big.Int)
	opt.upperBound.SetString(upper, 10)
	opt.count = count
	opt.verbose = verbose
}

func main() {
	var opt Options
	setOptions(&opt)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if opt.verbose {
		fmt.Println("Lower Bound: ", opt.lowerBound)
		fmt.Println("Upper Bound: ", opt.upperBound)
		fmt.Println("Count: ", opt.count)
		fmt.Printf("%+v\n", opt)
	}

	i := new(big.Int)
	i.Set(i.Sub(opt.upperBound, opt.lowerBound))
	if c := i.Cmp(big.NewInt(int64(opt.count))); c < 0 {
		fmt.Fprintf(os.Stderr, "Error: the interval (%s) is less than requested size of set (%d)\n", i, opt.count)
		os.Exit(1)
	}
	tree := &Tree{}
	for n := 0; n < opt.count; n++ {
		for {
			x := new(big.Int)
			x.Rand(r, i)
			x.Add(x, opt.lowerBound)
			if tree.Find(x) == false {
				tree.Insert(x)
				break
			}
		}
	}
	tree.Traverse(tree.Root, func(n *Node) {
		fmt.Println(n.Number.String())
	})
}
