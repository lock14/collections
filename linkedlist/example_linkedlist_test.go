package linkedlist_test

import (
	"fmt"
	"github.com/lock14/collections/linkedlist"
	"slices"
)

func ExampleLinkedList() {
	// LinkedList is a doubly-linked list, providing fast O(1) insertions
	// and deletions at the front and back.
	list := linkedlist.New[string]()

	list.AddBack("Middle")
	list.AddFront("Start")
	list.AddBack("End")

	// Set elements at a specific index (requires O(N) traversal to index)
	list.Set(1, "Second")

	// Iteration
	for val := range list.All() {
		fmt.Println(val)
	}

	// Output:
	// Start
	// Second
	// End
}

func ExampleLinkedList_AddFront() {
	list := linkedlist.New[int]()
	list.AddFront(1)
	list.AddFront(2)
	fmt.Println(list.String())
	// Output:
	// [2 1]
}

func ExampleLinkedList_RemoveFront() {
	list := linkedlist.New[int]()
	list.AddBack(1)
	list.AddBack(2)
	val := list.RemoveFront()
	fmt.Println(val)
	fmt.Println(list.String())
	// Output:
	// 1
	// [2]
}

func ExampleLinkedList_AddBack() {
	list := linkedlist.New[int]()
	list.AddBack(1)
	list.AddBack(2)
	fmt.Println(list.String())
	// Output:
	// [1 2]
}

func ExampleLinkedList_RemoveBack() {
	list := linkedlist.New[int]()
	list.AddBack(1)
	list.AddBack(2)
	val := list.RemoveBack()
	fmt.Println(val)
	fmt.Println(list.String())
	// Output:
	// 2
	// [1]
}

func ExampleLinkedList_Peek() {
	list := linkedlist.New[int]()
	list.AddBack(1)
	list.AddBack(2)
	fmt.Println(list.Peek())
	// Output:
	// 1
}

func ExampleLinkedList_PeekFront() {
	list := linkedlist.New[int]()
	list.AddBack(1)
	list.AddBack(2)
	fmt.Println(list.PeekFront())
	// Output:
	// 1
}

func ExampleLinkedList_PeekBack() {
	list := linkedlist.New[int]()
	list.AddBack(1)
	list.AddBack(2)
	fmt.Println(list.PeekBack())
	// Output:
	// 2
}

func ExampleLinkedList_Add() {
	list := linkedlist.New[int]()
	list.Add(1)
	list.Add(2)
	fmt.Println(list.String())
	// Output:
	// [1 2]
}

func ExampleLinkedList_Remove() {
	list := linkedlist.New[int]()
	list.Add(1)
	list.Add(2)
	val := list.Remove()
	fmt.Println(val)
	fmt.Println(list.String())
	// Output:
	// 1
	// [2]
}

func ExampleLinkedList_Push() {
	list := linkedlist.New[int]()
	list.Push(1)
	list.Push(2)
	fmt.Println(list.String())
	// Output:
	// [2 1]
}

func ExampleLinkedList_Get() {
	list := linkedlist.New[int]()
	list.Add(10)
	list.Add(20)
	list.Add(30)
	fmt.Println(list.Get(1))
	// Output:
	// 20
}

func ExampleLinkedList_Set() {
	list := linkedlist.New[int]()
	list.Add(10)
	list.Add(20)
	list.Add(30)
	list.Set(1, 99)
	fmt.Println(list.String())
	// Output:
	// [10 99 30]
}

func ExampleLinkedList_Pop() {
	list := linkedlist.New[int]()
	list.Push(1)
	list.Push(2)
	val := list.Pop()
	fmt.Println(val)
	fmt.Println(list.String())
	// Output:
	// 2
	// [1]
}

func ExampleLinkedList_Size() {
	list := linkedlist.New[int]()
	list.Add(1)
	list.Add(2)
	fmt.Println(list.Size())
	// Output:
	// 2
}

func ExampleLinkedList_AddAll() {
	list := linkedlist.New[int]()
	list.AddAll(slices.Values([]int{1, 2, 3}))
	fmt.Println(list.String())
	// Output:
	// [1 2 3]
}

func ExampleLinkedList_Empty() {
	list := linkedlist.New[int]()
	fmt.Println(list.Empty())
	list.Add(1)
	fmt.Println(list.Empty())
	// Output:
	// true
	// false
}

func ExampleLinkedList_Clear() {
	list := linkedlist.New[int]()
	list.Add(1)
	list.Add(2)
	list.Clear()
	fmt.Println(list.Size())
	fmt.Println(list.String())
	// Output:
	// 0
	// []
}

func ExampleLinkedList_String() {
	list := linkedlist.New[int]()
	list.Add(1)
	list.Add(2)
	fmt.Println(list.String())
	// Output:
	// [1 2]
}

func ExampleLinkedList_All() {
	list := linkedlist.New[int]()
	list.Add(1)
	list.Add(2)
	for val := range list.All() {
		fmt.Println(val)
	}
	// Output:
	// 1
	// 2
}

func ExampleLinkedList_Backward() {
	list := linkedlist.New[int]()
	list.Add(1)
	list.Add(2)
	for val := range list.Backward() {
		fmt.Println(val)
	}
	// Output:
	// 2
	// 1
}
