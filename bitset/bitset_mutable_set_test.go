package bitset

import (
	"github.com/lock14/collections/arraylist"
	"github.com/lock14/collections/hashset"
	"slices"
	"testing"
)

func TestBitSet_MutableSet(t *testing.T) {
	t.Parallel()
	b := New(NumBits(100))
	b.Add(10)
	b.Add(20)

	if !b.Contains(10) {
		t.Errorf("expected to contain 10")
	}
	if !b.Contains(20) {
		t.Errorf("expected to contain 20")
	}

	b.RemoveElement(10)
	if b.Contains(10) {
		t.Errorf("expected 10 to be removed")
	}

	if b.Empty() {
		t.Errorf("expected not to be empty")
	}

	b.RemoveElement(20)
	if !b.Empty() {
		t.Errorf("expected to be empty")
	}

	b.Add(5)
	b.Add(15)
	val := b.Remove()
	if val != 5 && val != 15 {
		t.Errorf("Remove returned unexpected value: %v", val)
	}

	b.Clear()
	if !b.Empty() {
		t.Errorf("expected empty after Clear")
	}

	b.AddAll(slices.Values([]int{1, 2, 3}))
	if !b.Contains(1) || !b.Contains(2) || !b.Contains(3) {
		t.Errorf("AddAll failed")
	}

	other := New(NumBits(100))
	other.Add(2)
	other.Add(3)
	other.Add(4)

	b.RemoveAll(other)
	if b.Contains(2) || b.Contains(3) || !b.Contains(1) {
		t.Errorf("RemoveAll failed")
	}

	b.Clear()
	b.AddAll(slices.Values([]int{1, 2, 3}))
	b.RetainAll(other) // retains 2, 3
	if b.Contains(1) || !b.Contains(2) || !b.Contains(3) {
		t.Errorf("RetainAll failed")
	}

	// test RetainAll when b has more words than other
	b.Clear()
	b.Add(100) // maxWordInUse = 2
	other.Clear()
	other.Add(2) // maxWordInUse = 1
	b.RetainAll(other)
	if b.Contains(100) {
		t.Errorf("RetainAll did not clear upper words")
	}

	b.Clear()
	b.AddAll(slices.Values([]int{2, 3}))
	other.Clear()
	other.AddAll(slices.Values([]int{2, 3, 4}))
	if b.ContainsAll(other) {
		t.Errorf("expected not to contain all")
	}

	// test ContainsAll when other has more words than b
	other.Clear()
	other.AddAll(slices.Values([]int{2, 3, 100}))
	if b.ContainsAll(other) {
		t.Errorf("expected not to contain all")
	}
	// test ContainsAll when other has more words than b but empty upper words
	other.Clear()
	other.Add(2)
	other.ensureSize(5)
	other.maxWordInUse = 6
	if !b.ContainsAll(other) {
		t.Errorf("expected to contain all even with empty upper words")
	}

	other.Clear()
	other.AddAll(slices.Values([]int{2, 3}))
	if !other.ContainsAll(b) {
		t.Errorf("ContainsAll failed")
	}

	collected := slices.Collect(b.All())
	if !slices.Equal(collected, []int{2, 3}) {
		t.Errorf("All iterator failed")
	}

	// Test generic Collection branch (using arraylist.Wrap)
	list := arraylist.Wrap([]int{2, 3, 4})
	b.Clear()
	b.AddAll(slices.Values([]int{1, 2, 3, 4, 5}))
	b.RemoveAll(list)
	if b.Contains(2) || b.Contains(3) || b.Contains(4) || !b.Contains(1) || !b.Contains(5) {
		t.Errorf("RemoveAll with generic collection failed")
	}

	b.Clear()
	b.AddAll(slices.Values([]int{1, 2, 3, 4, 5}))
	b.RetainAll(list)
	if b.Contains(1) || b.Contains(5) || !b.Contains(2) || !b.Contains(3) || !b.Contains(4) {
		t.Errorf("RetainAll with generic collection failed")
	}

	if !b.ContainsAll(arraylist.Wrap([]int{2, 3, 4})) {
		t.Errorf("ContainsAll with generic collection should be true")
	}
	if b.ContainsAll(arraylist.Wrap([]int{2, 3, 4, 99})) {
		t.Errorf("ContainsAll with generic collection should be false")
	}

	// Test generic Set branch (using hashset.New)
	hset := hashset.New[int]()
	hset.Add(2)
	hset.Add(3)
	hset.Add(4)

	b.Clear()
	b.AddAll(slices.Values([]int{1, 2, 3, 4, 5}))
	b.RemoveAll(hset)
	if b.Contains(2) || b.Contains(3) || b.Contains(4) || !b.Contains(1) || !b.Contains(5) {
		t.Errorf("RemoveAll with Set failed")
	}

	b.Clear()
	b.AddAll(slices.Values([]int{1, 2, 3, 4, 5}))
	b.RetainAll(hset)
	if b.Contains(1) || b.Contains(5) || !b.Contains(2) || !b.Contains(3) || !b.Contains(4) {
		t.Errorf("RetainAll with Set failed")
	}

	// Test Size() lazy evaluation
	b.Clear()
	b.AddAll(slices.Values([]int{1, 2, 3}))
	b.FlipRange(0, 5)  // sets size to -1
	if b.Size() != 2 { // triggers recomputeSize
		t.Errorf("Lazy size evaluation failed: expected 2, got %v", b.Size())
	}

	// Test Remove() panic
	b.Clear()
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Remove() did not panic on empty set")
			}
		}()
		b.Remove()
	}()

	// Test Remove() panic when size is broken
	b.Clear()
	b.size = 1 // break size
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Remove() did not panic on broken size")
			}
		}()
		b.Remove()
	}()
}
