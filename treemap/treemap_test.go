package treemap

import (
	"cmp"
	"fmt"
	"iter"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/lock14/collections/comparator"
)

type dummyPayload struct {
	id int
}

// assertBTreeInvariants deterministically verifies B-Tree structural invariants and ensures
// spare capacity in keys, values, and children slices contains no stale/leaked references.
func assertBTreeInvariants[K comparable, V comparable](t *testing.T, m *TreeMap[K, V]) {
	t.Helper()
	if m.root == nil {
		if m.size != 0 {
			t.Fatalf("nil root with non-zero size %d", m.size)
		}
		return
	}

	var zeroK K
	var zeroV V
	leafDepth := -1

	var checkNode func(n *node[K, V], isRoot bool, depth int, hasLo bool, lo K, hasHi bool, hi K) int
	checkNode = func(n *node[K, V], isRoot bool, depth int, hasLo bool, lo K, hasHi bool, hi K) int {
		if len(n.keys) != len(n.values) {
			t.Fatalf("keys len %d != values len %d", len(n.keys), len(n.values))
		}
		if !isRoot && len(n.keys) < m.degree-1 {
			t.Fatalf("underfull non-root node: %d keys (degree %d)", len(n.keys), m.degree)
		}
		if len(n.keys) > 2*m.degree-1 {
			t.Fatalf("overfull node: %d keys (degree %d)", len(n.keys), m.degree)
		}

		// Verify spare capacity is zeroed (prevents GC memory leaks).
		fullKeys := n.keys[:cap(n.keys)]
		for idx := len(n.keys); idx < len(fullKeys); idx++ {
			if fullKeys[idx] != zeroK {
				t.Fatalf("stale key in spare capacity at index %d: %v", idx, fullKeys[idx])
			}
		}
		fullVals := n.values[:cap(n.values)]
		for idx := len(n.values); idx < len(fullVals); idx++ {
			if fullVals[idx] != zeroV {
				t.Fatalf("stale value in spare capacity at index %d: %v", idx, fullVals[idx])
			}
		}
		fullChildren := n.children[:cap(n.children)]
		for idx := len(n.children); idx < len(fullChildren); idx++ {
			if fullChildren[idx] != nil {
				t.Fatalf("stale child pointer in spare capacity at index %d", idx)
			}
		}

		for i, k := range n.keys {
			if i > 0 && m.comparator(n.keys[i-1], k) >= 0 {
				t.Fatalf("keys not strictly sorted: %v >= %v", n.keys[i-1], k)
			}
			if hasLo && m.comparator(k, lo) <= 0 {
				t.Fatalf("key %v <= lower bound %v", k, lo)
			}
			if hasHi && m.comparator(k, hi) >= 0 {
				t.Fatalf("key %v >= upper bound %v", k, hi)
			}
		}

		count := len(n.keys)
		if n.leaf {
			if len(n.children) != 0 {
				t.Fatalf("leaf node has %d children", len(n.children))
			}
			if leafDepth == -1 {
				leafDepth = depth
			} else if depth != leafDepth {
				t.Fatalf("inconsistent leaf depth: got %d, want %d", depth, leafDepth)
			}
			return count
		}

		if len(n.children) != len(n.keys)+1 {
			t.Fatalf("internal node has %d children for %d keys", len(n.children), len(n.keys))
		}
		for i, child := range n.children {
			if child == nil {
				t.Fatalf("nil active child at index %d", i)
			}
			cHasLo, cLo := hasLo, lo
			if i > 0 {
				cHasLo, cLo = true, n.keys[i-1]
			}
			cHasHi, cHi := hasHi, hi
			if i < len(n.keys) {
				cHasHi, cHi = true, n.keys[i]
			}
			count += checkNode(child, false, depth+1, cHasLo, cLo, cHasHi, cHi)
		}
		return count
	}

	counted := checkNode(m.root, true, 0, false, zeroK, false, zeroK)
	if counted != m.size {
		t.Fatalf("counted %d keys in B-Tree, but m.size is %d", counted, m.size)
	}
}

func TestTreeMap_Constructors(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		wantPanic bool
		check     func(*testing.T)
	}{
		{
			name:      "degree_less_than_2_panics",
			wantPanic: true,
			check: func(_ *testing.T) {
				NewOrdered[int, int](WithDegree[int](1))
			},
		},
		{
			name:      "negative_degree_panics",
			wantPanic: true,
			check: func(_ *testing.T) {
				NewOrdered[int, int](WithDegree[int](-2))
			},
		},
		{
			name:      "missing_comparator_panics",
			wantPanic: true,
			check: func(_ *testing.T) {
				New[int, int]()
			},
		},
		{
			name:      "nil_comparator_override_in_new_ordered_panics",
			wantPanic: true,
			check: func(_ *testing.T) {
				NewOrdered[int, int](WithComparator[int](nil))
			},
		},
		{
			name: "custom_comparator_and_degree",
			check: func(t *testing.T) {
				m := New[int, string](
					WithDegree[int](3),
					WithComparator(comparator.Reverse(comparator.NaturalOrder[int]())),
				)
				m.Put(1, "a")
				m.Put(3, "c")
				m.Put(2, "b")
				if got := slices.Collect(m.Keys()); !slices.Equal(got, []int{3, 2, 1}) {
					t.Errorf("expected reverse order [3 2 1], got %v", got)
				}
			},
		},
		{
			name: "new_ordered_option_precedence_and_slice_non_aliasing",
			check: func(t *testing.T) {
				opts := make([]Option[int], 1, 4)
				opts[0] = WithComparator(comparator.Reverse(comparator.NaturalOrder[int]()))
				sentinelCalled := false
				fullSlice := opts[:2]
				fullSlice[1] = func(*config[int]) { sentinelCalled = true }

				m := NewOrdered[int, int](opts...)
				if sentinelCalled {
					t.Fatalf("NewOrdered overwrote or read spare capacity of caller slice")
				}
				m.Put(1, 10)
				m.Put(2, 20)
				if k, _ := m.First(); k != 2 {
					t.Errorf("expected custom ReverseOrder comparator to win in NewOrdered, got First()=%d", k)
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				r := recover()
				if tc.wantPanic && r == nil {
					t.Errorf("expected panic but got none")
				}
				if !tc.wantPanic && r != nil {
					t.Errorf("expected no panic but got: %v", r)
				}
			}()
			tc.check(t)
		})
	}
}

func TestTreeMap_OperationsAndSinglePassBTree(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T)
	}{
		{
			name: "basic_put_get_contains_and_clear",
			check: func(t *testing.T) {
				m := NewOrdered[int, string](WithDegree[int](2))
				if v, ok := m.Get(10); ok || v != "" {
					t.Errorf("expected ('', false) on empty map, got (%q, %v)", v, ok)
				}
				for _, kv := range []struct {
					k int
					v string
				}{
					{10, "ten"}, {20, "twenty"}, {5, "five"}, {6, "six"},
					{12, "twelve"}, {30, "thirty"}, {7, "seven"}, {17, "seventeen"},
				} {
					m.Put(kv.k, kv.v)
				}
				assertBTreeInvariants(t, m)
				if m.Size() != 8 || m.Empty() {
					t.Errorf("expected size 8, got %d", m.Size())
				}
				if val, ok := m.Get(12); !ok || val != "twelve" {
					t.Errorf("expected (twelve, true), got (%q, %v)", val, ok)
				}
				if m.ContainsKey(100) || !m.ContainsKey(30) {
					t.Errorf("ContainsKey mismatch")
				}
				m.Clear()
				if !m.Empty() || m.Size() != 0 {
					t.Errorf("expected empty map after Clear")
				}
				assertBTreeInvariants(t, m)
			},
		},
		{
			name: "overwrite_in_full_root_full_child_median_internal_and_leaf",
			check: func(t *testing.T) {
				m := NewOrdered[int, int](WithDegree[int](2))
				// With degree=2, max keys per node is 2*2-1 = 3.
				// Fill root to capacity [10, 20, 30] and overwrite in full root!
				m.Put(10, 1)
				m.Put(20, 2)
				m.Put(30, 3)
				m.Put(20, 200) // Overwrites in full root without splitting!
				if m.Size() != 3 || !m.root.leaf {
					t.Fatalf("expected size 3 unsplit root, got size=%d leaf=%v", m.Size(), m.root.leaf)
				}
				if v, _ := m.Get(20); v != 200 {
					t.Errorf("expected 200, got %d", v)
				}

				// Split root by inserting 40, then fill right child [30, 40, 50].
				// Now root=[20], left=[10], right=[30, 40, 50] (full child, median=40!).
				m.Put(40, 4)
				m.Put(50, 5)
				if len(m.root.children[1].keys) != 3 {
					t.Fatalf("expected full right child of 3 keys")
				}
				// Overwrite median key 40 of full child before split!
				m.Put(40, 400)
				if m.Size() != 5 {
					t.Fatalf("expected size 5 after overwriting full child median, got %d", m.Size())
				}
				if v, _ := m.Get(40); v != 400 {
					t.Errorf("expected 400, got %d", v)
				}

				// Insert 60 to trigger splitChild with cmp > 0, and overwrite every key in deeper tree.
				for i := 1; i <= 40; i++ {
					m.Put(i*5, i)
				}
				prevSize := m.Size()
				for i := 1; i <= 40; i++ {
					m.Put(i*5, i*100)
				}
				if m.Size() != prevSize {
					t.Errorf("size changed during overwrites: got %d, want %d", m.Size(), prevSize)
				}
				assertBTreeInvariants(t, m)
			},
		},
		{
			name: "delete_predecessor_successor_merge_and_absent_keys",
			check: func(t *testing.T) {
				m := NewOrdered[int, int](WithDegree[int](2))
				// Remove on empty map is a no-op
				m.Remove(42)

				// Construct case for getSuccessor: root=[20, 40], children=[10], [30, 35], [50]
				for _, k := range []int{10, 20, 30, 40, 50, 35} {
					m.Put(k, k)
				}
				m.Remove(20) // Left child has 1 key (t-1), right child has 2 keys (>=t) -> getSuccessor
				if m.ContainsKey(20) || m.Size() != 5 {
					t.Fatalf("expected 20 removed and size 5, got size=%d", m.Size())
				}
				assertBTreeInvariants(t, m)

				// Build deep degree-2 tree so getPredecessor & getSuccessor descend non-leaf subtrees,
				// and test removing absent keys (which trigger proactive fill/merge without changing size).
				m2 := NewOrdered[int, *dummyPayload](WithDegree[int](2))
				const total = 150
				for i := 1; i <= total; i++ {
					m2.Put(i*2, &dummyPayload{id: i})
				}
				assertBTreeInvariants(t, m2)

				// Remove absent odd keys across the entire tree; size must stay 150!
				for i := 0; i <= total+1; i++ {
					m2.Remove(i*2 + 1)
				}
				if m2.Size() != total {
					t.Fatalf("removing absent keys changed size: got %d, want %d", m2.Size(), total)
				}
				assertBTreeInvariants(t, m2)

				// Remove internal keys from root repeatedly to hit deep getPredecessor/getSuccessor/merge,
				// checking invariants & zeroed spare pointer slots after every single deletion!
				for !m2.Empty() {
					rootKey := m2.root.keys[len(m2.root.keys)/2]
					m2.Remove(rootKey)
					assertBTreeInvariants(t, m2)
				}
			},
		},
		{
			name: "zero_value_defensive_branches",
			check: func(t *testing.T) {
				var m TreeMap[int, int]
				m.degree = 2
				m.comparator = cmp.Compare[int]

				if v, ok := m.Get(1); ok || v != 0 {
					t.Errorf("expected (0, false) on nil root")
				}
				m.Remove(1)
				for range m.All() {
					t.Fatal("expected empty All on nil root")
				}
				for range m.Backward() {
					t.Fatal("expected empty Backward on nil root")
				}
				for range m.From(0) {
					t.Fatal("expected empty From on nil root")
				}
				for range m.Between(0, 10) {
					t.Fatal("expected empty Between on nil root")
				}

				// Put on nil root initializes root node cleanly
				m.Put(5, 50)
				if v, ok := m.Get(5); !ok || v != 50 {
					t.Errorf("expected (50, true), got (%d, %v)", v, ok)
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.check(t)
		})
	}
}

// TestTreeMap_RandomizedStress exercises 1,000 shuffled insertions and deletions on a minimum-degree
// B-Tree using a deterministic PRNG and verifies structural + GC-zeroing invariants throughout.
func TestTreeMap_RandomizedStress(t *testing.T) {
	t.Parallel()
	rng := rand.New(rand.NewPCG(42, 1024))
	m := NewOrdered[int, *dummyPayload](WithDegree[int](2))

	const n = 1000
	keys := make([]int, n)
	for i := range keys {
		keys[i] = i
	}
	rng.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	for _, k := range keys {
		m.Put(k, &dummyPayload{id: k * 10})
	}
	if m.Size() != n {
		t.Fatalf("expected size %d, got %d", n, m.Size())
	}
	assertBTreeInvariants(t, m)

	rng.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	for idx, k := range keys {
		m.Remove(k)
		if idx%50 == 0 || m.Empty() {
			assertBTreeInvariants(t, m)
		}
	}
	if !m.Empty() {
		t.Fatalf("expected empty map after removing all %d keys", n)
	}
}

func TestTreeMap_SequencedAndNavigable(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		wantPanic bool
		check     func(*testing.T)
	}{
		{
			name:      "first_empty_panics",
			wantPanic: true,
			check: func(_ *testing.T) {
				NewOrdered[int, int]().First()
			},
		},
		{
			name:      "last_empty_panics",
			wantPanic: true,
			check: func(_ *testing.T) {
				NewOrdered[int, int]().Last()
			},
		},
		{
			name:      "poll_first_empty_panics",
			wantPanic: true,
			check: func(_ *testing.T) {
				NewOrdered[int, int]().PollFirst()
			},
		},
		{
			name:      "poll_last_empty_panics",
			wantPanic: true,
			check: func(_ *testing.T) {
				NewOrdered[int, int]().PollLast()
			},
		},
		{
			name: "multi_level_first_last_poll_first_poll_last",
			check: func(t *testing.T) {
				m := NewOrdered[int, int](WithDegree[int](2))
				for i := 1; i <= 25; i++ {
					m.Put(i*10, i*100)
				}
				if k, v := m.First(); k != 10 || v != 100 {
					t.Errorf("First() = (%d, %d), want (10, 100)", k, v)
				}
				if k, v := m.Last(); k != 250 || v != 2500 {
					t.Errorf("Last() = (%d, %d), want (250, 2500)", k, v)
				}
				if k, v := m.PollFirst(); k != 10 || v != 100 || m.ContainsKey(10) {
					t.Errorf("PollFirst() = (%d, %d)", k, v)
				}
				if k, v := m.PollLast(); k != 250 || v != 2500 || m.ContainsKey(250) {
					t.Errorf("PollLast() = (%d, %d)", k, v)
				}
			},
		},
		{
			name: "multi_level_lower_floor_ceiling_higher",
			check: func(t *testing.T) {
				empty := NewOrdered[int, int]()
				if _, _, ok := empty.Lower(10); ok {
					t.Errorf("expected Lower false on empty")
				}
				if _, _, ok := empty.Floor(10); ok {
					t.Errorf("expected Floor false on empty")
				}
				if _, _, ok := empty.Ceiling(10); ok {
					t.Errorf("expected Ceiling false on empty")
				}
				if _, _, ok := empty.Higher(10); ok {
					t.Errorf("expected Higher false on empty")
				}

				m := NewOrdered[int, int](WithDegree[int](2))
				for i := 1; i <= 30; i++ {
					m.Put(i*10, i*100)
				}

				if _, _, ok := m.Lower(10); ok {
					t.Errorf("Lower(10) should be false")
				}
				if _, _, ok := m.Floor(5); ok {
					t.Errorf("Floor(5) should be false")
				}
				if _, _, ok := m.Ceiling(305); ok {
					t.Errorf("Ceiling(305) should be false")
				}
				if _, _, ok := m.Higher(300); ok {
					t.Errorf("Higher(300) should be false")
				}

				for i := 1; i <= 30; i++ {
					key := i * 10
					// Exact key checks across all internal & leaf nodes
					if k, v, ok := m.Floor(key); !ok || k != key || v != i*100 {
						t.Fatalf("Floor(%d) = (%d, %d, %v)", key, k, v, ok)
					}
					if k, v, ok := m.Ceiling(key); !ok || k != key || v != i*100 {
						t.Fatalf("Ceiling(%d) = (%d, %d, %v)", key, k, v, ok)
					}
					if i > 1 {
						if k, _, ok := m.Lower(key); !ok || k != (i-1)*10 {
							t.Fatalf("Lower(%d) = (%d, %v)", key, k, ok)
						}
						if k, _, ok := m.Floor(key - 5); !ok || k != (i-1)*10 {
							t.Fatalf("Floor(%d) = (%d, %v)", key-5, k, ok)
						}
					}
					if i < 30 {
						if k, _, ok := m.Higher(key); !ok || k != (i+1)*10 {
							t.Fatalf("Higher(%d) = (%d, %v)", key, k, ok)
						}
						if k, _, ok := m.Ceiling(key + 5); !ok || k != (i+1)*10 {
							t.Fatalf("Ceiling(%d) = (%d, %v)", key+5, k, ok)
						}
					}
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				r := recover()
				if tc.wantPanic && r == nil {
					t.Errorf("expected panic, got none")
				}
				if !tc.wantPanic && r != nil {
					t.Errorf("unexpected panic: %v", r)
				}
			}()
			tc.check(t)
		})
	}
}

func TestTreeMap_IteratorsAndRangeQueries(t *testing.T) {
	t.Parallel()
	const n = 25
	m := NewOrdered[int, string](WithDegree[int](2))
	wantKeys := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		k := i * 10
		m.Put(k, fmt.Sprintf("v%d", k))
		wantKeys = append(wantKeys, k)
	}
	wantRevKeys := slices.Clone(wantKeys)
	slices.Reverse(wantRevKeys)

	collectSeq2Keys := func(seq iter.Seq2[int, string]) []int {
		var out []int
		for k := range seq {
			out = append(out, k)
		}
		return out
	}

	cases := []struct {
		name  string
		check func(*testing.T)
	}{
		{
			name: "forward_and_backward_full_and_every_early_exit_position",
			check: func(t *testing.T) {
				if got := slices.Collect(m.Keys()); !slices.Equal(got, wantKeys) {
					t.Fatalf("Keys() = %v, want %v", got, wantKeys)
				}
				if got := slices.Collect(m.BackwardKeys()); !slices.Equal(got, wantRevKeys) {
					t.Fatalf("BackwardKeys() = %v, want %v", got, wantRevKeys)
				}
				if len(slices.Collect(m.Values())) != n || len(slices.Collect(m.BackwardValues())) != n {
					t.Fatalf("unexpected Values/BackwardValues length")
				}

				// Early exit at every index 1..n exercises leaf, internal separator, and rightmost/leftmost child exits
				for stop := 1; stop <= n; stop++ {
					var fGot []int
					for k := range m.All() {
						fGot = append(fGot, k)
						if len(fGot) == stop {
							break
						}
					}
					if !slices.Equal(fGot, wantKeys[:stop]) {
						t.Fatalf("All() stop=%d: got %v, want %v", stop, fGot, wantKeys[:stop])
					}

					var bGot []int
					for k := range m.Backward() {
						bGot = append(bGot, k)
						if len(bGot) == stop {
							break
						}
					}
					if !slices.Equal(bGot, wantRevKeys[:stop]) {
						t.Fatalf("Backward() stop=%d: got %v, want %v", stop, bGot, wantRevKeys[:stop])
					}
				}

				// Early break on Keys, Values, BackwardKeys, BackwardValues wrappers
				for range m.Keys() {
					break
				}
				for range m.Values() {
					break
				}
				for range m.BackwardKeys() {
					break
				}
				for range m.BackwardValues() {
					break
				}
			},
		},
		{
			name: "from_to_between_exhaustive_bounds_and_early_breaks",
			check: func(t *testing.T) {
				empty := NewOrdered[int, string](WithDegree[int](2))
				if collectSeq2Keys(empty.From(10)) != nil {
					t.Errorf("expected empty From on empty map")
				}

				// Invalid/empty intervals in Between
				if collectSeq2Keys(m.Between(100, 100)) != nil {
					t.Errorf("expected empty Between(100, 100)")
				}
				if collectSeq2Keys(m.Between(200, 50)) != nil {
					t.Errorf("expected empty Between(200, 50)")
				}

				// Exhaustively verify From, To, and Between across boundary, exact-key, and mid-gap probes!
				probes := []int{0, 5, 10, 15, 40, 45, 80, 120, 125, 160, 200, 250, 260}
				for _, lo := range probes {
					var expFrom []int
					var expTo []int
					for _, k := range wantKeys {
						if k >= lo {
							expFrom = append(expFrom, k)
						}
						if k < lo {
							expTo = append(expTo, k)
						}
					}
					if got := collectSeq2Keys(m.From(lo)); !slices.Equal(got, expFrom) {
						t.Fatalf("From(%d) = %v, want %v", lo, got, expFrom)
					}
					if got := collectSeq2Keys(m.To(lo)); !slices.Equal(got, expTo) {
						t.Fatalf("To(%d) = %v, want %v", lo, got, expTo)
					}

					for _, hi := range probes {
						var expBetween []int
						for _, k := range wantKeys {
							if k >= lo && k < hi {
								expBetween = append(expBetween, k)
							}
						}
						got := collectSeq2Keys(m.Between(lo, hi))
						if !slices.Equal(got, expBetween) {
							t.Fatalf("Between(%d, %d) = %v, want %v", lo, hi, got, expBetween)
						}
					}
				}

				// Test early break at every stop position inside rangeInOrder
				expectedSub := wantKeys[2:22] // [30..220]
				for stop := 1; stop <= len(expectedSub); stop++ {
					var got []int
					for k := range m.Between(30, 225) {
						got = append(got, k)
						if len(got) == stop {
							break
						}
					}
					if !slices.Equal(got, expectedSub[:stop]) {
						t.Fatalf("Between early break stop=%d: got %v, want %v", stop, got, expectedSub[:stop])
					}
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.check(t)
		})
	}
}

func TestTreeMap_String(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		entries [][2]int
		want    string
	}{
		{
			name:    "empty",
			entries: nil,
			want:    "map[]",
		},
		{
			name:    "single",
			entries: [][2]int{{1, 10}},
			want:    "map[1:10]",
		},
		{
			name:    "multiple_sorted_order",
			entries: [][2]int{{3, 30}, {1, 10}, {2, 20}},
			want:    "map[1:10 2:20 3:30]",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			m := NewOrdered[int, int]()
			for _, entry := range tc.entries {
				m.Put(entry[0], entry[1])
			}
			if got := m.String(); got != tc.want {
				t.Errorf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}
