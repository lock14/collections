package arraylist_test

import (
	"fmt"
	"github.com/lock14/collections/arraylist"
	"slices"
)

func ExampleSliceWrapper() {
	// ArrayList is a dynamically resizing array, providing O(1) random access.
	list := arraylist.Wrap(make([]string, 0))

	list.Add("Apple")
	list.Add("Banana")
	list.Add("Cherry")

	// Set element at a specific index
	list.Set(1, "Blueberry")

	// Random access by index
	fmt.Println("Element at index 2:", list.Get(2))

	// Iteration
	fmt.Println("All fruits:")
	for val := range list.All() {
		fmt.Println("-", val)
	}

	// Output:
	// Element at index 2: Cherry
	// All fruits:
	// - Apple
	// - Blueberry
	// - Cherry
}

func ExampleSliceWrapper_Add() {
	list := arraylist.Wrap([]int{})
	list.Add(1)
	list.Add(2)
	fmt.Println(list.String())
	// Output:
	// [1, 2]
}

func ExampleSliceWrapper_Remove() {
	list := arraylist.Wrap([]int{1, 2, 3})
	val := list.Remove()
	fmt.Println(val)
	fmt.Println(list.String())
	// Output:
	// 3
	// [1, 2]
}

func ExampleSliceWrapper_Push() {
	stack := arraylist.Wrap([]int{})
	stack.Push(10)
	stack.Push(20)
	fmt.Println(stack.String())
	// Output:
	// [10, 20]
}

func ExampleSliceWrapper_Pop() {
	stack := arraylist.Wrap([]int{10, 20})
	val := stack.Pop()
	fmt.Println(val)
	fmt.Println(stack.String())
	// Output:
	// 20
	// [10]
}

func ExampleSliceWrapper_Peek() {
	stack := arraylist.Wrap([]int{10, 20})
	val := stack.Peek()
	fmt.Println(val)
	fmt.Println(stack.String())
	// Output:
	// 20
	// [10, 20]
}

func ExampleSliceWrapper_Clear() {
	list := arraylist.Wrap([]int{1, 2, 3})
	list.Clear()
	fmt.Println(list.Size())
	fmt.Println(list.String())
	// Output:
	// 0
	// []
}

func ExampleSliceWrapper_AddAll() {
	list := arraylist.Wrap([]int{1, 2})
	list.AddAll(slices.Values([]int{3, 4}))
	fmt.Println(list.String())
	// Output:
	// [1, 2, 3, 4]
}

func ExampleSliceWrapper_Size() {
	list := arraylist.Wrap([]int{1, 2, 3})
	fmt.Println(list.Size())
	list.Remove()
	fmt.Println(list.Size())
	// Output:
	// 3
	// 2
}

func ExampleSliceWrapper_Empty() {
	list := arraylist.Wrap([]int{})
	fmt.Println(list.Empty())
	list.Add(1)
	fmt.Println(list.Empty())
	// Output:
	// true
	// false
}

func ExampleSliceWrapper_Get() {
	list := arraylist.Wrap([]int{10, 20, 30})
	fmt.Println(list.Get(1))
	// Output:
	// 20
}

func ExampleSliceWrapper_Set() {
	list := arraylist.Wrap([]int{10, 20, 30})
	list.Set(1, 99)
	fmt.Println(list.String())
	// Output:
	// [10, 99, 30]
}

func ExampleSliceWrapper_String() {
	list := arraylist.Wrap([]string{"A", "B", "C"})
	fmt.Println(list.String())
	// Output:
	// [A, B, C]
}

func ExampleSliceWrapper_All() {
	list := arraylist.Wrap([]int{10, 20, 30})
	for val := range list.All() {
		fmt.Println(val)
	}
	// Output:
	// 10
	// 20
	// 30
}
