package main

import (
	"fmt"
	"math"

	"github.com/lock14/collections/bitset"
)

// Prints all primes less than 100
func main() {
	n := uint32(100)
	b := bitset.New(bitset.NumBits(n))
	b.Set(0)
	b.Set(1)
	bound := uint32(math.Sqrt(float64(n))) + 1
	for i := uint32(4); i < bound; i += 2 {
		b.Set(i)
	}
	for i := uint32(3); i < bound; i += 2 {
		if !b.Get(i) {
			// i is prime
			for j := i * i; j < n; j += i {
				b.Set(j)
			}
		}
	}
	b.Flip()
	fmt.Println(2)
	for i := uint32(3); i < n; i += 2 {
		if b.Get(i) {
			fmt.Println(i)
		}
	}
}
