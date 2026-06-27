package bitset

import (
	"slices"
	"testing"

	linkedlist "github.com/lock14/collections/linkedlist"
)

func TestBitSet_Size(t *testing.T) {
	b := New()
	if b.Size() != 0 {
		t.Fatalf("expected 0 got %d", b.Size())
	}
	b.Add(5)
	b.Add(10)
	if b.Size() != 2 {
		t.Fatalf("expected 2 got %d", b.Size())
	}
}

func TestBitSet_Remove(t *testing.T) {
	b := New()
	b.Add(5)
	val := b.Remove()
	if val != 5 {
		t.Fatalf("expected 5, got %d", val)
	}
}

func TestBitSet_RetainAll(t *testing.T) {
	b := New()
	b.Add(1)
	b.Add(2)
	b.Add(3)

	other := linkedlist.New[int]()
	other.Add(2)
	other.Add(4)

	b.RetainAll(other)
	if b.Size() != 1 || !b.Contains(2) {
		t.Fatalf("expected only 2 to be retained, got size %d", b.Size())
	}
}

func TestBitSet_ContainsAll(t *testing.T) {
	b := New()
	b.Add(1)
	b.Add(2)

	other := linkedlist.New[int]()
	other.Add(1)
	other.Add(2)

	if !b.ContainsAll(other) {
		t.Fatal("expected contains all to be true")
	}
	other.Add(3)
	if b.ContainsAll(other) {
		t.Fatal("expected contains all to be false")
	}
}

func TestBitSet_FirstLast(t *testing.T) {
	b := New()
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic on First() when empty")
			}
		}()
		b.First()
	}()
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic on Last() when empty")
			}
		}()
		b.Last()
	}()
	b.Add(5)
	b.Add(10)
	if v := b.First(); v != 5 {
		t.Fatalf("expected 5, got %d", v)
	}
	if v := b.Last(); v != 10 {
		t.Fatalf("expected 10, got %d", v)
	}
}

func TestBitSet_PollFirstLast(t *testing.T) {
	b := New()
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic on PollFirst() when empty")
			}
		}()
		b.PollFirst()
	}()
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic on PollLast() when empty")
			}
		}()
		b.PollLast()
	}()
	b.Add(5)
	b.Add(10)
	if v := b.PollFirst(); v != 5 {
		t.Fatalf("expected 5, got %d", v)
	}
	if b.Contains(5) {
		t.Fatal("expected 5 to be removed")
	}
	if v := b.PollLast(); v != 10 {
		t.Fatalf("expected 10, got %d", v)
	}
	if b.Contains(10) {
		t.Fatal("expected 10 to be removed")
	}
}

func TestBitSet_AddFirstLast(t *testing.T) {
	b := New()
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic on AddFirst")
			}
		}()
		b.AddFirst(1)
	}()
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic on AddLast")
			}
		}()
		b.AddLast(1)
	}()
}

func TestBitSet_NavigableQueries(t *testing.T) {
	b := New()

	// Test empty conditions
	if _, ok := b.Lower(100); ok {
		t.Fatal("expected lower to be false on empty")
	}
	if _, ok := b.Floor(100); ok {
		t.Fatal("expected floor to be false on empty")
	}
	if _, ok := b.Ceiling(0); ok {
		t.Fatal("expected ceiling to be false on empty")
	}
	if _, ok := b.Higher(0); ok {
		t.Fatal("expected higher to be false on empty")
	}

	b.Add(5)
	b.Add(10)
	b.Add(15)

	if v, ok := b.Lower(10); !ok || v != 5 {
		t.Fatalf("expected 5, got %d", v)
	}
	if _, ok := b.Lower(5); ok {
		t.Fatal("expected no lower than 5")
	}
	if v, ok := b.Lower(100); !ok || v != 15 {
		t.Fatalf("expected 15, got %d", v)
	}

	if v, ok := b.Floor(10); !ok || v != 10 {
		t.Fatalf("expected 10, got %d", v)
	}
	if v, ok := b.Floor(12); !ok || v != 10 {
		t.Fatalf("expected 10, got %d", v)
	}
	if v, ok := b.Floor(100); !ok || v != 15 {
		t.Fatalf("expected 15, got %d", v)
	}

	if v, ok := b.Ceiling(10); !ok || v != 10 {
		t.Fatalf("expected 10, got %d", v)
	}
	if v, ok := b.Ceiling(12); !ok || v != 15 {
		t.Fatalf("expected 15, got %d", v)
	}
	if v, ok := b.Ceiling(-10); !ok || v != 5 {
		t.Fatalf("expected 5, got %d", v)
	}

	if v, ok := b.Higher(10); !ok || v != 15 {
		t.Fatalf("expected 15, got %d", v)
	}
	if _, ok := b.Higher(15); ok {
		t.Fatal("expected no higher than 15")
	}
	if v, ok := b.Higher(-10); !ok || v != 5 {
		t.Fatalf("expected 5, got %d", v)
	}

	b.Add(63)
	b.Add(64)
	b.Add(127)
	b.Add(128)
	b.Lower(64)
	b.Lower(128)
	b.Floor(63)
	b.Floor(127)
	b.Ceiling(64)
	b.Ceiling(128)
	b.Higher(63)
	b.Higher(127)
}

func TestBitSet_Iterators(t *testing.T) {
	b := New()

	// Test empty iterators
	if len(slices.Collect(b.Backward())) != 0 {
		t.Fatal("expected empty backward")
	}
	if len(slices.Collect(b.From(0))) != 0 {
		t.Fatal("expected empty from")
	}
	if len(slices.Collect(b.To(100))) != 0 {
		t.Fatal("expected empty to")
	}
	if len(slices.Collect(b.Between(0, 100))) != 0 {
		t.Fatal("expected empty between")
	}

	b.Add(5)
	b.Add(10)
	b.Add(15)

	back := slices.Collect(b.Backward())
	if !slices.Equal(back, []int{15, 10, 5}) {
		t.Fatalf("expected [15 10 5], got %v", back)
	}

	from := slices.Collect(b.From(10))
	if !slices.Equal(from, []int{10, 15}) {
		t.Fatalf("expected [10 15], got %v", from)
	}

	to := slices.Collect(b.To(15))
	if !slices.Equal(to, []int{5, 10}) {
		t.Fatalf("expected [5 10], got %v", to)
	}

	between := slices.Collect(b.Between(5, 15))
	if !slices.Equal(between, []int{5, 10}) {
		t.Fatalf("expected [5 10], got %v", between)
	}

	b.Add(63)
	b.Add(64)
	slices.Collect(b.From(64))
	slices.Collect(b.To(64))
	slices.Collect(b.Between(63, 65))
}
