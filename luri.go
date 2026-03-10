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
	"flag"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"sort"
	"time"
)

type Options struct {
	lowerBound *big.Int
	upperBound *big.Int
	count      int
	verbose    bool
}

// setOptions gets the options from the command line environment.
func setOptions(opt *Options) {
	var lower, upper string
	var count int
	var verbose bool
	const (
		defaultLowerBound = "1"
		lowerUsage        = "the lower bound `big.Int`"
		defaultUpperBound = "100"
		upperUsage        = "the upper bound `big.Int` (inclusive)"
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
	if _, ok := opt.lowerBound.SetString(lower, 10); !ok {
		fmt.Fprintf(os.Stderr, "Error: invalid lower bound: %q\n", lower)
		os.Exit(1)
	}
	opt.upperBound = new(big.Int)
	if _, ok := opt.upperBound.SetString(upper, 10); !ok {
		fmt.Fprintf(os.Stderr, "Error: invalid upper bound: %q\n", upper)
		os.Exit(1)
	}
	opt.count = count
	opt.verbose = verbose
}

// maxFisherYates is the maximum range size for which a partial Fisher-Yates
// shuffle is used. Above this threshold, map-based rejection sampling is used.
const maxFisherYates = 1_000_000

func main() {
	var opt Options
	setOptions(&opt)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// rangeSize = upperBound - lowerBound + 1 (upper bound is inclusive)
	rangeSize := new(big.Int).Sub(opt.upperBound, opt.lowerBound)
	rangeSize.Add(rangeSize, big.NewInt(1))
	if rangeSize.Sign() <= 0 {
		fmt.Fprintf(os.Stderr, "Error: upper bound must be >= lower bound\n")
		os.Exit(1)
	}
	if rangeSize.Cmp(big.NewInt(int64(opt.count))) < 0 {
		fmt.Fprintf(os.Stderr, "Error: the interval (%s) is less than requested size of set (%d)\n", rangeSize, opt.count)
		os.Exit(1)
	}

	if opt.verbose {
		fmt.Println("Lower Bound: ", opt.lowerBound)
		fmt.Println("Upper Bound: ", opt.upperBound)
		fmt.Println("Count: ", opt.count)
	}

	var numbers []*big.Int
	if rangeSize.IsInt64() && rangeSize.Int64() <= maxFisherYates {
		numbers = fisherYatesSample(opt.lowerBound, rangeSize.Int64(), opt.count, r)
	} else {
		numbers = rejectionSample(opt.lowerBound, rangeSize, opt.count, r)
	}

	sort.Slice(numbers, func(i, j int) bool {
		return numbers[i].Cmp(numbers[j]) < 0
	})

	for _, n := range numbers {
		fmt.Println(n.String())
	}
}

// fisherYatesSample performs a partial Fisher-Yates shuffle over the range
// [lowerBound, lowerBound+rangeSize) and returns count unique values.
func fisherYatesSample(lowerBound *big.Int, rangeSize int64, count int, r *rand.Rand) []*big.Int {
	indices := make([]int64, rangeSize)
	for i := int64(0); i < rangeSize; i++ {
		indices[i] = i
	}
	for i := 0; i < count; i++ {
		j := int64(i) + r.Int63n(rangeSize-int64(i))
		indices[i], indices[j] = indices[j], indices[i]
	}
	result := make([]*big.Int, count)
	for i := 0; i < count; i++ {
		v := new(big.Int).SetInt64(indices[i])
		v.Add(v, lowerBound)
		result[i] = v
	}
	return result
}

// rejectionSample uses a hash map to collect count unique random integers from
// [lowerBound, lowerBound+rangeSize). Suitable for large ranges where
// count << rangeSize.
func rejectionSample(lowerBound *big.Int, rangeSize *big.Int, count int, r *rand.Rand) []*big.Int {
	seen := make(map[string]struct{}, count)
	result := make([]*big.Int, 0, count)
	x := new(big.Int)
	for len(result) < count {
		x.Rand(r, rangeSize)
		key := x.String()
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, new(big.Int).Add(x, lowerBound))
		}
	}
	return result
}
