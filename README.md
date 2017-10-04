# luri
List Unique Random Integers

This program generates an ordered list of integers of arbitrary size,
guaranteed to contain only unique integers.

It is a reimplementation of a C program that does the same thing,
with the purpose of learning about Go and the equivalent packages
in the Go core libraries (math/rand, math/big) etcetera.

Generating 1 000 000 unique integers with an upper bound of 2 000 000
takes roughly 5.2 s on my 2 GHz Intel Core i7-2630QM.

`
luri -c 1000000 -u 2000000
`

The equivalent C implementation (with the GMP library) is about 1 s
quicker on the same processor. Some optimizations could be done
in the Go implementation, making generation and binary tree searches
into goroutines to provide some concurrency.

Also, it would be possible to optimize an algorithm that generates
only the excluded numbers and prints the rest when the size of the
set is more than half the numeric interval.

## Usage

luri --help gives you the following menu.

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
        the upper bound big.Int (default "100")
  -upper big.Int
        the upper bound big.Int (default "100")
  -v    be verbose (shorthand)
  -verbose
        be verbose
```

## Copyright & License

This program is Copyright (C) 2017 Pär Karlsson

It is released under the GNU General Public License. See the file COPYING for details.
