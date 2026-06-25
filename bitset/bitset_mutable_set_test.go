package bitset

import (
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

	b.Clear()
	b.AddAll(slices.Values([]int{2, 3}))
	if b.ContainsAll(other) {
		t.Errorf("expected not to contain all")
	}
	if !other.ContainsAll(b) {
		t.Errorf("ContainsAll failed")
	}

	collected := slices.Collect(b.All())
	if !slices.Equal(collected, []int{2, 3}) {
		t.Errorf("All iterator failed")
	}
}
