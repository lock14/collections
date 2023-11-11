package main

import (
	"fmt"

	"github.com/lock14/collections/bitset"
)

// Prints all primes less than 100
func main() {
	n := uint32(100)
	b := primesLessThan(n)
	fmt.Println(2)
	for i := uint32(3); i < n; i += 2 {
		if b.Get(i) {
			fmt.Println(i)
		}
	}
}

func primesLessThan(n uint32) *bitset.BitSet {
	b := bitset.New(bitset.NumBits(n))
	if n > 2 {
		b.Set(0)
		b.Set(1)
		for i := uint32(4); i < n; i += 2 {
			b.Set(i)
		}
		for i := uint32(3); (i*i) > i && (i*i) < n; i += 2 {
			if !b.Get(i) {
				// i is prime
				for j := i * i; j > i && j < n; j += i {
					b.Set(j)
				}
			}
		}
		b.Flip()
	}
	return b
}
