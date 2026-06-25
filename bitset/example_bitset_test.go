package bitset_test

import (
	"fmt"
	"github.com/lock14/collections/bitset"
)

func ExampleBitSet_setOperations() {
	// BitSet provides a highly optimized set of non-negative integers
	// with support for fast bitwise set operations.
	b := bitset.New()

	// Set bits
	b.SetBit(1)
	b.SetBit(2)
	b.SetBit(3)

	fmt.Println("Contains 2:", b.GetBit(2))
	fmt.Println("Contains 5:", b.GetBit(5))

	// Flip a range of bits (inclusive start, exclusive end)
	// This will flip bits 2, 3, 4 (2 and 3 become false, 4 becomes true)
	b.FlipRange(2, 5)

	fmt.Println("After flip:")
	fmt.Println("Contains 2:", b.GetBit(2))
	fmt.Println("Contains 4:", b.GetBit(4))

	// Print all set bits
	fmt.Println("Set bits:")
	for bit := range b.All() {
		fmt.Println("-", bit)
	}

	// Output:
	// Contains 2: true
	// Contains 5: false
	// After flip:
	// Contains 2: false
	// Contains 4: true
	// Set bits:
	// - 1
	// - 4
}
