# luri

**List Unique Random Integers**

`luri` generates a sorted list of unique random integers sampled from an
inclusive range `[lower, upper]`. Both bounds are
[`math/big.Int`](https://pkg.go.dev/math/big) values, so the range can be
arbitrarily large — well beyond the limits of a 64-bit integer.

## Algorithm

The sampling strategy is chosen automatically based on the size of the
requested range:

### Small ranges (≤ 1 000 000 values) — partial Fisher-Yates shuffle

An index array `[0, 1, …, rangeSize-1]` is allocated once. A partial
[Fisher-Yates shuffle](https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle)
is then performed for exactly `count` steps:

```
for i in 0 .. count-1:
    j = random integer in [i, rangeSize)
    swap(indices[i], indices[j])
```

The first `count` elements of the shuffled array, offset by `lowerBound`,
become the result. This is **O(count)** time with no possibility of
collision, making it optimal when the range fits comfortably in memory.

The threshold of 1 000 000 is a practical balance: allocating a slice of
that many `int64` values requires roughly 8 MB of memory, which is
acceptable on virtually any modern system, while keeping the algorithm
fast for typical use-cases.

### Large ranges (> 1 000 000 values) — hash-map rejection sampling

A `map[string]struct{}` tracks already-selected values. Random integers are
drawn from `[0, rangeSize)` using
[`(*big.Int).Rand`](https://pkg.go.dev/math/big#Int.Rand) until `count`
distinct values have been collected, then each is offset by `lowerBound`.

This approach handles ranges of arbitrary magnitude (including ranges that
exceed `int64`) and is efficient when `count` is much smaller than
`rangeSize`. For each of the `count` unique values to collect, the expected
draws required for the k-th value (with k−1 already seen) is
`rangeSize / (rangeSize − k + 1)`. Summing over all k gives:

```
E[draws] = Σ_{k=1}^{count} rangeSize/(rangeSize − k + 1)
         ≈ count + count·(count−1) / (2·rangeSize)
         ≈ count   (when count << rangeSize)
```

This is the standard
[rejection-sampling](https://en.wikipedia.org/wiki/Rejection_sampling)
argument. When `count` is close to `rangeSize`, the expected draws grow
significantly and the Fisher-Yates path should be preferred instead, which
is why the threshold exists.

### Output

Results are sorted in ascending order before printing, regardless of which
sampling path was taken.

## Code structure

| Path | Description |
|------|-------------|
| `luri.go` | Entry point: flag parsing, algorithm dispatch, sorting, output |
| `luri_test.go` | Unit tests for `fisherYatesSample` and `rejectionSample` |
| `bintree/bintree.go` | Add-only, iterative unbalanced [BST](https://en.wikipedia.org/wiki/Binary_search_tree) for `*big.Int` values |
| `bintree/bintree_test.go` | Unit tests for the `bintree` package |

The `bintree` package provides an iterative (stack-overflow-safe) binary
search tree keyed on `*big.Int`. All tree operations — `Insert`, `Find`, and
`Traverse` — are implemented without recursion to handle arbitrarily large or
degenerate (fully skewed) trees safely.

## Building

Requires Go 1.21 or later.

```
git clone https://github.com/feinorgh/luri
cd luri
go build .
```

## Usage

```
Usage of ./luri:
  -c int
        size of set (shorthand) (default 1)
  -count int
        size of set (default 1)
  -l big.Int
        the lower bound big.Int (shorthand) (default "1")
  -lower big.Int
        the lower bound big.Int (default "1")
  -u big.Int
        the upper bound big.Int (inclusive) (default "100")
  -upper big.Int
        the upper bound big.Int (inclusive) (default "100")
  -v    be verbose (shorthand)
  -verbose
        be verbose
```

Both bounds are **inclusive**. The program exits with an error if:

- either bound is not a valid integer,
- `lower > upper`, or
- `count > upper - lower + 1` (the requested set is larger than the range).

### Examples

Pick 5 unique integers from 1 to 10 (inclusive):

```
$ luri -l 1 -u 10 -c 5
2
4
6
7
9
```

Pick 3 values from an astronomically large range:

```
$ luri -l 100000000000000000000 -u 999999999999999999999 -c 3
317423091827364509128
541890234761092834710
876123409812374650912
```

Show bounds and count before output:

```
$ luri -l 1 -u 20 -c 5 -v
Lower Bound:  1
Upper Bound:  20
Count:  5
3
8
11
14
19
```

## Running the tests

```
go test ./...
```

## Copyright & License

This program is Copyright (C) 2017 Pär Karlsson

It is released under the GNU General Public License. See the file COPYING for details.
