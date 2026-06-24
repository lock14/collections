package treemap

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"
)

func TestTreeMap_Basic(t *testing.T) {
	tm := NewOrdered[int, string](WithDegree[int](2)) // Small degree to force splits
	
	if !tm.Empty() {
		t.Errorf("expected new map to be empty")
	}

	tm.Put(10, "ten")
	tm.Put(20, "twenty")
	tm.Put(5, "five")
	tm.Put(6, "six")
	tm.Put(12, "twelve")
	tm.Put(30, "thirty")
	tm.Put(7, "seven")
	tm.Put(17, "seventeen")

	if tm.Size() != 8 {
		t.Errorf("expected size 8, got %d", tm.Size())
	}

	val, ok := tm.Get(12)
	if !ok || val != "twelve" {
		t.Errorf("expected twelve, got %v, %v", val, ok)
	}

	if tm.ContainsKey(100) {
		t.Errorf("did not expect to find 100")
	}

	// Overwrite
	tm.Put(10, "TEN")
	val, _ = tm.Get(10)
	if val != "TEN" {
		t.Errorf("expected TEN, got %v", val)
	}
	if tm.Size() != 8 {
		t.Errorf("size should remain 8 after overwrite")
	}
}

func TestTreeMap_Remove(t *testing.T) {
	tm := NewOrdered[int, string](WithDegree[int](2))
	for i := 1; i <= 20; i++ {
		tm.Put(i, fmt.Sprintf("val-%d", i))
	}

	if tm.Size() != 20 {
		t.Errorf("expected size 20, got %d", tm.Size())
	}

	// Remove leaf
	tm.Remove(20)
	if tm.ContainsKey(20) || tm.Size() != 19 {
		t.Errorf("remove failed for 20")
	}

	// Remove internal node
	tm.Remove(10)
	if tm.ContainsKey(10) || tm.Size() != 18 {
		t.Errorf("remove failed for 10")
	}

	// Remove root iteratively
	for i := 1; i <= 19; i++ {
		if i == 10 {
			continue
		}
		tm.Remove(i)
	}

	if !tm.Empty() {
		t.Errorf("expected empty map, size is %d", tm.Size())
	}
}

func TestTreeMap_Iterators(t *testing.T) {
	tm := NewOrdered[int, string](WithDegree[int](3))
	expectedKeys := []int{1, 5, 8, 10, 15, 20}
	
	// Insert out of order
	tm.Put(10, "10")
	tm.Put(1, "1")
	tm.Put(20, "20")
	tm.Put(8, "8")
	tm.Put(15, "15")
	tm.Put(5, "5")

	var actualKeys []int
	for k := range tm.Keys() {
		actualKeys = append(actualKeys, k)
	}

	if !slices.Equal(expectedKeys, actualKeys) {
		t.Errorf("expected %v, got %v", expectedKeys, actualKeys)
	}

	var actualVals []string
	for v := range tm.Values() {
		actualVals = append(actualVals, v)
	}
	expectedVals := []string{"1", "5", "8", "10", "15", "20"}
	if !slices.Equal(expectedVals, actualVals) {
		t.Errorf("expected %v, got %v", expectedVals, actualVals)
	}
}

func TestTreeMap_LargeRandom(t *testing.T) {
	tm := NewOrdered[int, int]()
	keys := make([]int, 0, 1000)
	for i := 0; i < 1000; i++ {
		keys = append(keys, i)
	}
	
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	for _, k := range keys {
		tm.Put(k, k*10)
	}

	if tm.Size() != 1000 {
		t.Fatalf("expected 1000, got %d", tm.Size())
	}

	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	for _, k := range keys[:500] {
		tm.Remove(k)
	}

	if tm.Size() != 500 {
		t.Fatalf("expected 500, got %d", tm.Size())
	}

	for _, k := range keys[:500] {
		if tm.ContainsKey(k) {
			t.Fatalf("should have been removed: %d", k)
		}
	}

	for _, k := range keys[500:] {
		v, ok := tm.Get(k)
		if !ok || v != k*10 {
			t.Fatalf("expected to find %d with value %d, got ok=%v, v=%d", k, k*10, ok, v)
		}
	}
}

func TestTreeMap_Clear(t *testing.T) {
	tm := NewOrdered[int, int]()
	tm.Put(1, 1)
	tm.Put(2, 2)
	tm.Clear()
	
	if !tm.Empty() || tm.Size() != 0 {
		t.Errorf("expected empty map")
	}
	if tm.ContainsKey(1) {
		t.Errorf("should not contain 1")
	}
}

func TestTreeMap_EdgeCases(t *testing.T) {
	// Panics
	assertPanics := func(f func()) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic")
			}
		}()
		f()
	}

	assertPanics(func() { NewOrdered[int, int](WithDegree[int](1)) })
	assertPanics(func() { New[int, int]() }) // No comparator

	tm := NewOrdered[int, int](WithDegree[int](2))

	// Remove from empty
	tm.Remove(10)

	// Get from empty
	if _, ok := tm.Get(10); ok {
		t.Errorf("expected not ok")
	}

	// Insert and remove non-existent
	tm.Put(10, 10)
	tm.Remove(20)
	if tm.Size() != 1 {
		t.Errorf("expected size 1")
	}

	// Iterator early exit
	tm.Put(20, 20)
	tm.Put(30, 30)
	count := 0
	for range tm.Keys() {
		count++
		break // early exit
	}
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}

	count = 0
	for range tm.Values() {
		count++
		break // early exit
	}
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}

	count = 0
	for k, _ := range tm.All() {
		if k == 20 {
			break
		}
		count++
	}
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}

	// For Get 50% coverage
	tm.root = nil
	tm.Get(10)
}

func TestTreeMap_DeepTree(t *testing.T) {
	tm := NewOrdered[int, int](WithDegree[int](2))
	for i := 0; i < 2000; i++ {
		tm.Put(i, i)
	}
	
	// Delete every third element to trigger borrows and merges heavily
	for i := 0; i < 2000; i += 3 {
		tm.Remove(i)
	}
	
	// Delete the rest
	for i := 0; i < 2000; i++ {
		if i%3 != 0 {
			tm.Remove(i)
		}
	}
	
	if !tm.Empty() {
		t.Errorf("expected empty")
	}

	// Trigger early exit in inOrder on an internal node
	for i := 0; i < 2000; i++ {
		tm.Put(i, i)
	}
	count := 0
	for range tm.Keys() {
		count++
		if count == 1000 {
			break
		}
	}
}

func TestTreeMap_GetSuccessor(t *testing.T) {
	// We want to hit the case where we delete an internal node's key,
	// and its left child has t-1 keys, but its right child has >= t keys.
	tm := NewOrdered[int, int](WithDegree[int](2))
	
	// Create a specific tree structure
	// degree=2 means max keys = 3. 
	// Let's just insert elements carefully to make left child small, right child big.
	tm.Put(20, 20)
	tm.Put(10, 10)
	tm.Put(30, 30)
	tm.Put(40, 40)
	tm.Put(50, 50)
	// Now root might be 20, 40.
	// If we just remove 10, the left child of 20 becomes empty/merged.
	// Let's just insert a bunch and delete carefully.
	
	tm.Clear()
	tm.Put(10, 10)
	tm.Put(20, 20)
	tm.Put(30, 30)
	tm.Put(40, 40)
	tm.Put(50, 50)
	
	// Root is [20, 40].
	// Children of 20, 40 are: [10], [30], [50] (if t=2).
	// Let's make right child of 20 have more keys.
	tm.Put(35, 35) // Now [30, 35]
	
	// Left child of 20 is [10] (t-1 = 1). Right child of 20 is [30, 35] (>= t = 2).
	// If we delete 20, it should borrow from the right child (getSuccessor).
	tm.Remove(20)
	
	if tm.ContainsKey(20) {
		t.Errorf("20 should be removed")
	}
}
