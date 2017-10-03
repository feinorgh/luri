package main

import (
	"flag"
	"fmt"
	"math/big"
	// "os"
)

type Options struct {
	lowerBound *big.Int
	upperBound *big.Int
	count      int
	verbose    bool
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
	args := flag.Args()
	if opt.verbose {
		fmt.Println("Lower Bound: ", opt.lowerBound)
		fmt.Println("Upper Bound: ", opt.upperBound)
		fmt.Printf("%+v\n", opt)
		fmt.Println(opt)
		fmt.Println(args)
	}
}
