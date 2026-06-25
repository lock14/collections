package treemap

import (
	"fmt"
	"iter"
	"math/rand"
	"slices"
	"testing"
)

func TestTreeMap_Operations(t *testing.T) {
	cases := []struct {
		name     string
		degree   int
		ops      func(tm *TreeMap[int, string])
		validate func(t *testing.T, tm *TreeMap[int, string])
	}{
		{
			name:   "basic put and get",
			degree: 2,
			ops: func(tm *TreeMap[int, string]) {
				tm.Put(10, "ten")
				tm.Put(20, "twenty")
				tm.Put(5, "five")
				tm.Put(6, "six")
				tm.Put(12, "twelve")
				tm.Put(30, "thirty")
				tm.Put(7, "seven")
				tm.Put(17, "seventeen")
			},
			validate: func(t *testing.T, tm *TreeMap[int, string]) {
				if tm.Size() != 8 {
					t.Errorf("expected size 8, got %d", tm.Size())
				}
				if val, ok := tm.Get(12); !ok || val != "twelve" {
					t.Errorf("expected twelve, got %v, %v", val, ok)
				}
				if tm.ContainsKey(100) {
					t.Errorf("did not expect to find 100")
				}
			},
		},
		{
			name:   "overwrite existing key",
			degree: 2,
			ops: func(tm *TreeMap[int, string]) {
				tm.Put(10, "ten")
				tm.Put(10, "TEN")
			},
			validate: func(t *testing.T, tm *TreeMap[int, string]) {
				if val, _ := tm.Get(10); val != "TEN" {
					t.Errorf("expected TEN, got %v", val)
				}
				if tm.Size() != 1 {
					t.Errorf("expected size 1, got %d", tm.Size())
				}
			},
		},
		{
			name:   "remove leaf and internal nodes",
			degree: 2,
			ops: func(tm *TreeMap[int, string]) {
				for i := 1; i <= 20; i++ {
					tm.Put(i, fmt.Sprintf("val-%d", i))
				}
				tm.Remove(20) // leaf
				tm.Remove(10) // internal
			},
			validate: func(t *testing.T, tm *TreeMap[int, string]) {
				if tm.ContainsKey(20) || tm.ContainsKey(10) {
					t.Errorf("nodes were not removed")
				}
				if tm.Size() != 18 {
					t.Errorf("expected size 18, got %d", tm.Size())
				}
			},
		},
		{
			name:   "remove all nodes",
			degree: 2,
			ops: func(tm *TreeMap[int, string]) {
				for i := 1; i <= 20; i++ {
					tm.Put(i, fmt.Sprintf("val-%d", i))
				}
				for i := 1; i <= 20; i++ {
					tm.Remove(i)
				}
			},
			validate: func(t *testing.T, tm *TreeMap[int, string]) {
				if !tm.Empty() || tm.Size() != 0 {
					t.Errorf("expected empty map")
				}
			},
		},
		{
			name:   "clear map",
			degree: 2,
			ops: func(tm *TreeMap[int, string]) {
				tm.Put(1, "1")
				tm.Put(2, "2")
				tm.Clear()
			},
			validate: func(t *testing.T, tm *TreeMap[int, string]) {
				if !tm.Empty() || tm.Size() != 0 {
					t.Errorf("expected empty map")
				}
			},
		},
		{
			name:   "iterators return sorted order",
			degree: 3,
			ops: func(tm *TreeMap[int, string]) {
				tm.Put(10, "10")
				tm.Put(1, "1")
				tm.Put(20, "20")
				tm.Put(8, "8")
				tm.Put(15, "15")
				tm.Put(5, "5")
			},
			validate: func(t *testing.T, tm *TreeMap[int, string]) {
				expectedKeys := []int{1, 5, 8, 10, 15, 20}
				expectedVals := []string{"1", "5", "8", "10", "15", "20"}

				var actualKeys []int
				for k := range tm.Keys() {
					actualKeys = append(actualKeys, k)
				}
				if !slices.Equal(expectedKeys, actualKeys) {
					t.Errorf("expected keys %v, got %v", expectedKeys, actualKeys)
				}

				var actualVals []string
				for v := range tm.Values() {
					actualVals = append(actualVals, v)
				}
				if !slices.Equal(expectedVals, actualVals) {
					t.Errorf("expected vals %v, got %v", expectedVals, actualVals)
				}
			},
		},
		{
			name:   "iterator early exit",
			degree: 2,
			ops: func(tm *TreeMap[int, string]) {
				tm.Put(20, "20")
				tm.Put(30, "30")
			},
			validate: func(t *testing.T, tm *TreeMap[int, string]) {
				count := 0
				for range tm.Keys() {
					count++
					break
				}
				if count != 1 {
					t.Errorf("expected early exit, count %d", count)
				}
			},
		},
		{
			name:   "get from empty map",
			degree: 2,
			ops:    func(tm *TreeMap[int, string]) {},
			validate: func(t *testing.T, tm *TreeMap[int, string]) {
				if _, ok := tm.Get(10); ok {
					t.Errorf("expected not ok")
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tm := NewOrdered[int, string](WithDegree[int](tc.degree))
			tc.ops(tm)
			tc.validate(t, tm)
		})
	}
}

func TestTreeMap_Constructors(t *testing.T) {
	cases := []struct {
		name        string
		constructor func()
		expectPanic bool
	}{
		{
			name: "invalid degree panics",
			constructor: func() {
				NewOrdered[int, int](WithDegree[int](1))
			},
			expectPanic: true,
		},
		{
			name: "missing comparator panics",
			constructor: func() {
				New[int, int]()
			},
			expectPanic: true,
		},
		{
			name: "valid degree succeeds",
			constructor: func() {
				NewOrdered[int, int](WithDegree[int](32))
			},
			expectPanic: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				r := recover()
				if tc.expectPanic && r == nil {
					t.Errorf("expected panic but got none")
				}
				if !tc.expectPanic && r != nil {
					t.Errorf("expected no panic but got: %v", r)
				}
			}()
			tc.constructor()
		})
	}
}

// Procedural Stress Tests
// The following tests construct very specific B-Tree states or perform
// randomized large-scale operations which do not fit cleanly into a table format.

func TestTreeMap_LargeRandom(t *testing.T) {
	t.Parallel()
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

func TestTreeMap_DeepTree(t *testing.T) {
	t.Parallel()
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
}

func TestTreeMap_GetSuccessor(t *testing.T) {
	t.Parallel()
	// We want to hit the case where we delete an internal node's key,
	// and its left child has t-1 keys, but its right child has >= t keys.
	tm := NewOrdered[int, int](WithDegree[int](2))

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

func assertKV(t *testing.T, k int, v int, ok bool, expectedK int, expectedV int, expectedOk bool, name string) {
	t.Helper()
	if ok != expectedOk {
		t.Errorf("%s: expected ok=%v, got %v", name, expectedOk, ok)
	} else if ok && (k != expectedK || v != expectedV) {
		t.Errorf("%s: expected (%v, %v), got (%v, %v)", name, expectedK, expectedV, k, v)
	}
}

func TestTreeMap_Navigable(t *testing.T) {
	t.Parallel()
	tm := NewOrdered[int, int]()

	_, _, ok := tm.Lower(10)
	if ok {
		t.Errorf("expected empty map to return false for Lower")
	}

	tm.Put(10, 10)
	tm.Put(20, 20)
	tm.Put(30, 30)
	tm.Put(40, 40)
	tm.Put(50, 50)

	k, v, ok := tm.Lower(30)
	assertKV(t, k, v, ok, 20, 20, true, "Lower(30)")
	k, v, ok = tm.Lower(10)
	assertKV(t, k, v, ok, 0, 0, false, "Lower(10)")

	k, v, ok = tm.Floor(30)
	assertKV(t, k, v, ok, 30, 30, true, "Floor(30)")
	k, v, ok = tm.Floor(25)
	assertKV(t, k, v, ok, 20, 20, true, "Floor(25)")
	k, v, ok = tm.Floor(5)
	assertKV(t, k, v, ok, 0, 0, false, "Floor(5)")

	k, v, ok = tm.Ceiling(30)
	assertKV(t, k, v, ok, 30, 30, true, "Ceiling(30)")
	k, v, ok = tm.Ceiling(35)
	assertKV(t, k, v, ok, 40, 40, true, "Ceiling(35)")
	k, v, ok = tm.Ceiling(60)
	assertKV(t, k, v, ok, 0, 0, false, "Ceiling(60)")

	k, v, ok = tm.Higher(30)
	assertKV(t, k, v, ok, 40, 40, true, "Higher(30)")
	k, v, ok = tm.Higher(50)
	assertKV(t, k, v, ok, 0, 0, false, "Higher(50)")
}

func TestTreeMap_Sequenced(t *testing.T) {
	t.Parallel()

	// Test panics
	assertPanic := func(name string, f func()) {
		t.Helper()
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s did not panic", name)
			}
		}()
		f()
	}

	tm := NewOrdered[int, int]()
	assertPanic("First", func() { tm.First() })
	assertPanic("Last", func() { tm.Last() })
	assertPanic("PollFirst", func() { tm.PollFirst() })
	assertPanic("PollLast", func() { tm.PollLast() })

	tm.Put(10, 10)
	tm.Put(20, 20)
	tm.Put(30, 30)

	k, v := tm.First()
	if k != 10 || v != 10 {
		t.Errorf("First: expected (10, 10), got (%v, %v)", k, v)
	}

	k, v = tm.Last()
	if k != 30 || v != 30 {
		t.Errorf("Last: expected (30, 30), got (%v, %v)", k, v)
	}

	k, v = tm.PollFirst()
	if k != 10 || v != 10 {
		t.Errorf("PollFirst: expected (10, 10), got (%v, %v)", k, v)
	}
	if tm.ContainsKey(10) {
		t.Errorf("PollFirst didn't remove")
	}

	k, v = tm.PollLast()
	if k != 30 || v != 30 {
		t.Errorf("PollLast: expected (30, 30), got (%v, %v)", k, v)
	}
	if tm.ContainsKey(30) {
		t.Errorf("PollLast didn't remove")
	}

	assertPanic("PutFirst", func() { tm.PutFirst(1, 1) })
	assertPanic("PutLast", func() { tm.PutLast(1, 1) })
}

func TestTreeMap_Reversed(t *testing.T) {
	t.Parallel()
	tm := NewOrdered[int, int]()
	for i := 1; i <= 5; i++ {
		tm.Put(i*10, i*10)
	}

	var keys []int
	for k := range tm.ReversedKeys() {
		keys = append(keys, k)
	}
	if !slices.Equal(keys, []int{50, 40, 30, 20, 10}) {
		t.Errorf("ReversedKeys: %v", keys)
	}

	var vals []int
	for _, v := range tm.ReversedAll() {
		vals = append(vals, v)
	}
	if !slices.Equal(vals, []int{50, 40, 30, 20, 10}) {
		t.Errorf("ReversedAll: %v", vals)
	}

	// coverage for ReversedKeys early exit
	for k := range tm.ReversedKeys() {
		_ = k
		break
	}
	// coverage for ReversedValues early exit
	for v := range tm.ReversedValues() {
		_ = v
		break
	}
}

func TestTreeMap_Iterators_Bounds(t *testing.T) {
	t.Parallel()
	tm := NewOrdered[int, int]()
	for i := 1; i <= 10; i++ {
		tm.Put(i, i)
	}

	collect := func(seq iter.Seq2[int, int]) []int {
		var res []int
		for k := range seq {
			res = append(res, k)
		}
		return res
	}

	if got := collect(tm.AllFrom(5)); !slices.Equal(got, []int{5, 6, 7, 8, 9, 10}) {
		t.Errorf("AllFrom(5): %v", got)
	}
	if got := collect(tm.AllTo(5)); !slices.Equal(got, []int{1, 2, 3, 4}) {
		t.Errorf("AllTo(5): %v", got)
	}
	if got := collect(tm.AllBetween(3, 7)); !slices.Equal(got, []int{3, 4, 5, 6}) {
		t.Errorf("AllBetween(3, 7): %v", got)
	}

	// Early exit coverage
	for k := range tm.AllBetween(1, 10) {
		_ = k
		break
	}
}
