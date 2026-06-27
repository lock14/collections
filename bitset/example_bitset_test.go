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
	b := bitset.New(bitset.WithNumBits(100))
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
	b.SetBit(0)
	b.SetBit(8)
	fmt.Printf("%x\n", b.ToBytes())
	// Output:
	// 0101
}

func ExampleBitSet_String() {
	b := bitset.New()
	b.SetBit(0)
	b.SetBit(4)
	fmt.Println(b.String())
	// Output:
	// 0000000000000011
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
