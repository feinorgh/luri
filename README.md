# luri

**List Unique Random Integers**

`luri` generates a sorted list of unique random integers sampled from an
inclusive range `[lower, upper]`. Both bounds are `math/big.Int` values, so
the range can be arbitrarily large — well beyond the limits of a 64-bit
integer.

## Algorithm

The sampling strategy is chosen automatically based on the size of the
requested range:

### Small ranges (≤ 1 000 000 values) — partial Fisher-Yates shuffle

An index array `[0, 1, …, rangeSize-1]` is allocated once. A partial
Fisher-Yates shuffle is then performed for exactly `count` steps:

```
for i in 0 .. count-1:
    j = random integer in [i, rangeSize)
    swap(indices[i], indices[j])
```

The first `count` elements of the shuffled array, offset by `lowerBound`,
become the result. This is **O(count)** time with no possibility of
collision, making it optimal when the range fits comfortably in memory.

### Large ranges (> 1 000 000 values) — hash-map rejection sampling

A `map[string]struct{}` tracks already-selected values. Random integers are
drawn from `[0, rangeSize)` using `(*big.Int).Rand` until `count` distinct
values have been collected, then each is offset by `lowerBound`.

This approach handles ranges of arbitrary magnitude (including ranges that
exceed `int64`) and is efficient when `count` is much smaller than
`rangeSize` — selection probability of a collision approaches zero in that
regime, so the expected number of draws is approximately `count`.

### Output

Results are sorted in ascending order before printing, regardless of which
sampling path was taken.

## Building

Requires Go 1.18 or later.

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
