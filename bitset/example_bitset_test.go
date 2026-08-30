package bitset_test

import (
	"fmt"
	"github.com/lock14/collections/bitset"
)

func Example_sieveOfEratosthenes() {
	// Find all prime numbers less than 30 using the Sieve of Eratosthenes.
	n := 30
	sieve := bitset.New(bitset.WithCapacity(n))

	// Mark 0 and 1 as non-prime
	sieve.SetBit(0)
	sieve.SetBit(1)

	// Mark even numbers > 2
	for i := 4; i < n; i += 2 {
		sieve.SetBit(i)
	}

	// Sieve odd composite numbers
	for i := 3; i*i < n; i += 2 {
		if !sieve.GetBit(i) {
			for j := i * i; j < n; j += i {
				sieve.SetBit(j)
			}
		}
	}

	// Flip bits in range [0, n) so set bits represent prime numbers
	sieve.FlipRange(0, n)

	// Iterate and print all discovered primes
	for p := range sieve.SetBits() {
		fmt.Printf("%d ", p)
	}
	fmt.Println()

	// Output:
	// 2 3 5 7 11 13 17 19 23 29
}

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

func ExampleBitSet_ClearBit() {
	b := bitset.New()
	b.SetBit(5)
	fmt.Println(b.GetBit(5))
	b.ClearBit(5)
	fmt.Println(b.GetBit(5))
	// Output:
	// true
	// false
}

func ExampleBitSet_SetBit() {
	b := bitset.New()
	b.SetBit(10)
	fmt.Println(b.GetBit(10))
	// Output:
	// true
}

func ExampleBitSet_GetBit() {
	b := bitset.New()
	b.SetBit(3)
	fmt.Println(b.GetBit(3))
	fmt.Println(b.GetBit(4))
	// Output:
	// true
	// false
}

func ExampleBitSet_Capacity() {
	b := bitset.New(bitset.WithCapacity(100))
	fmt.Println(b.Capacity() >= 100)
	// Output:
	// true
}

func ExampleBitSet_Size() {
	b := bitset.New()
	b.SetBit(1)
	b.SetBit(10)
	b.SetBit(100)
	fmt.Println(b.Size())
	// Output:
	// 3
}

func ExampleBitSet_Length() {
	b := bitset.New()
	b.SetBit(1)
	b.SetBit(10)
	// Length is the highest set bit plus one
	fmt.Println(b.Length())
	// Output:
	// 11
}

func ExampleBitSet_Flip() {
	b := bitset.New()
	b.SetBit(0)
	b.Flip()
	fmt.Println(b.GetBit(0))
	fmt.Println(b.GetBit(1))
	// Output:
	// false
	// true
}

func ExampleBitSet_FlipRange() {
	b := bitset.New()
	b.FlipRange(2, 5)
	fmt.Println(b.GetBit(1))
	fmt.Println(b.GetBit(2))
	fmt.Println(b.GetBit(4))
	fmt.Println(b.GetBit(5))
	// Output:
	// false
	// true
	// true
	// false
}

func ExampleBitSet_ToBytes() {
	b := bitset.New()
	b.SetBit(0) // Bit 0 -> Byte 0 has value 0x01
	b.SetBit(8) // Bit 8 -> Byte 1 has value 0x01

	// ToBytes returns a little-endian byte slice: []byte{0x01, 0x01}
	fmt.Printf("%v\n", b.ToBytes())
	// Output:
	// [1 1]
}

func ExampleFromBytes() {
	// Reconstruct a BitSet from a little-endian byte slice
	// 0x05 (binary 00000101) sets bits 0 and 2
	b := bitset.FromBytes([]byte{0x05})

	fmt.Println("Contains 0:", b.Contains(0))
	fmt.Println("Contains 1:", b.Contains(1))
	fmt.Println("Contains 2:", b.Contains(2))
	// Output:
	// Contains 0: true
	// Contains 1: false
	// Contains 2: true
}

func ExampleBitSet_String() {
	b := bitset.New()
	b.SetBit(0)
	b.SetBit(4)
	fmt.Println(b.String())
	// Output:
	// [0 4]
}

func ExampleBitSet_SetBits() {
	b := bitset.New()
	b.SetBit(1)
	b.SetBit(5)
	for bit := range b.SetBits() {
		fmt.Println(bit)
	}
	// Output:
	// 1
	// 5
}

func ExampleBitSet_Add() {
	b := bitset.New()
	b.Add(7)
	fmt.Println(b.Contains(7))
	// Output:
	// true
}

func ExampleBitSet_RemoveElement() {
	b := bitset.New()
	b.Add(7)
	b.RemoveElement(7)
	fmt.Println(b.Contains(7))
	// Output:
	// false
}

func ExampleBitSet_Contains() {
	b := bitset.New()
	b.Add(7)
	fmt.Println(b.Contains(7))
	fmt.Println(b.Contains(8))
	// Output:
	// true
	// false
}

func ExampleBitSet_Empty() {
	b := bitset.New()
	fmt.Println(b.Empty())
	b.Add(1)
	fmt.Println(b.Empty())
	// Output:
	// true
	// false
}

func ExampleBitSet_All() {
	b := bitset.New()
	b.Add(2)
	b.Add(4)
	for bit := range b.All() {
		fmt.Println(bit)
	}
	// Output:
	// 2
	// 4
}

func ExampleBitSet_Clear() {
	b := bitset.New()
	b.Add(1)
	b.Clear()
	fmt.Println(b.Empty())
	// Output:
	// true
}

func ExampleBitSet_Remove() {
	b := bitset.New()
	b.Add(10)
	b.Add(20)
	val := b.Remove()
	fmt.Println(val)
	fmt.Println(b.Size())
	// Output:
	// 10
	// 1
}

func ExampleBitSet_AddAll() {
	b1 := bitset.New()
	b2 := bitset.New()
	b2.Add(1)
	b2.Add(3)

	b1.AddAll(b2.All())
	fmt.Println(b1.Contains(1), b1.Contains(3))
	// Output:
	// true true
}

func ExampleBitSet_RemoveAll() {
	b1 := bitset.New()
	b1.Add(1)
	b1.Add(2)

	b2 := bitset.New()
	b2.Add(2)

	b1.RemoveAll(b2)
	fmt.Println(b1.Contains(1), b1.Contains(2))
	// Output:
	// true false
}

func ExampleBitSet_RetainAll() {
	b1 := bitset.New()
	b1.Add(1)
	b1.Add(2)

	b2 := bitset.New()
	b2.Add(2)
	b2.Add(3)

	b1.RetainAll(b2)
	fmt.Println(b1.Contains(1), b1.Contains(2), b1.Contains(3))
	// Output:
	// false true false
}

func ExampleBitSet_ContainsAll() {
	b1 := bitset.New()
	b1.Add(1)
	b1.Add(2)

	b2 := bitset.New()
	b2.Add(2)

	fmt.Println(b1.ContainsAll(b2))
	fmt.Println(b2.ContainsAll(b1))
	// Output:
	// true
	// false
}

func ExampleBitSet_First() {
	b := bitset.New()
	b.Add(5)
	b.Add(10)
	fmt.Println(b.First())
	// Output:
	// 5
}

func ExampleBitSet_Last() {
	b := bitset.New()
	b.Add(5)
	b.Add(10)
	fmt.Println(b.Last())
	// Output:
	// 10
}

func ExampleBitSet_PollFirst() {
	b := bitset.New()
	b.Add(5)
	b.Add(10)
	fmt.Println(b.PollFirst())
	fmt.Println(b.Contains(5))
	// Output:
	// 5
	// false
}

func ExampleBitSet_PollLast() {
	b := bitset.New()
	b.Add(5)
	b.Add(10)
	fmt.Println(b.PollLast())
	fmt.Println(b.Contains(10))
	// Output:
	// 10
	// false
}

func ExampleBitSet_AddFirst() {
	b := bitset.New()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	b.AddFirst(1)
	// Output:
	// AddFirst is not supported on SortedSet
}

func ExampleBitSet_AddLast() {
	b := bitset.New()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	b.AddLast(1)
	// Output:
	// AddLast is not supported on SortedSet
}

func ExampleBitSet_Lower() {
	b := bitset.New()
	b.Add(5)
	b.Add(10)
	val, found := b.Lower(10)
	fmt.Println(val, found)
	// Output:
	// 5 true
}

func ExampleBitSet_Floor() {
	b := bitset.New()
	b.Add(5)
	b.Add(10)
	val, found := b.Floor(9)
	fmt.Println(val, found)
	val2, found2 := b.Floor(10)
	fmt.Println(val2, found2)
	// Output:
	// 5 true
	// 10 true
}

func ExampleBitSet_Ceiling() {
	b := bitset.New()
	b.Add(5)
	b.Add(10)
	val, found := b.Ceiling(6)
	fmt.Println(val, found)
	val2, found2 := b.Ceiling(10)
	fmt.Println(val2, found2)
	// Output:
	// 10 true
	// 10 true
}

func ExampleBitSet_Higher() {
	b := bitset.New()
	b.Add(5)
	b.Add(10)
	val, found := b.Higher(5)
	fmt.Println(val, found)
	// Output:
	// 10 true
}

func ExampleBitSet_Backward() {
	b := bitset.New()
	b.Add(1)
	b.Add(5)
	b.Add(10)
	for bit := range b.Backward() {
		fmt.Println(bit)
	}
	// Output:
	// 10
	// 5
	// 1
}

func ExampleBitSet_From() {
	b := bitset.New()
	b.Add(1)
	b.Add(5)
	b.Add(10)
	for bit := range b.From(5) {
		fmt.Println(bit)
	}
	// Output:
	// 5
	// 10
}

func ExampleBitSet_To() {
	b := bitset.New()
	b.Add(1)
	b.Add(5)
	b.Add(10)
	for bit := range b.To(6) {
		fmt.Println(bit)
	}
	// Output:
	// 1
	// 5
}

func ExampleBitSet_Between() {
	b := bitset.New()
	b.Add(1)
	b.Add(5)
	b.Add(10)
	b.Add(15)
	for bit := range b.Between(5, 12) {
		fmt.Println(bit)
	}
	// Output:
	// 5
	// 10
}
