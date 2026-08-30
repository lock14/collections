package comparator_test

import (
	"github.com/lock14/collections/comparator"
	"testing"
)

func TestNaturalOrder(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		cmp      int
		expected int
	}{
		{
			name:     "int_less",
			cmp:      comparator.NaturalOrder[int]()(1, 2),
			expected: -1,
		},
		{
			name:     "int_equal",
			cmp:      comparator.NaturalOrder[int]()(5, 5),
			expected: 0,
		},
		{
			name:     "int_greater",
			cmp:      comparator.NaturalOrder[int]()(10, 2),
			expected: 1,
		},
		{
			name:     "string_less",
			cmp:      comparator.NaturalOrder[string]()("apple", "banana"),
			expected: -1,
		},
		{
			name:     "string_equal",
			cmp:      comparator.NaturalOrder[string]()("cat", "cat"),
			expected: 0,
		},
		{
			name:     "string_greater",
			cmp:      comparator.NaturalOrder[string]()("zebra", "yak"),
			expected: 1,
		},
		{
			name:     "float_less",
			cmp:      comparator.NaturalOrder[float64]()(1.5, 2.5),
			expected: -1,
		},
		{
			name:     "float_equal",
			cmp:      comparator.NaturalOrder[float64]()(3.14, 3.14),
			expected: 0,
		},
		{
			name:     "float_greater",
			cmp:      comparator.NaturalOrder[float64]()(9.9, 8.8),
			expected: 1,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.cmp != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, tc.cmp)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	t.Parallel()
	revInt := comparator.Reverse(comparator.NaturalOrder[int]())

	cases := []struct {
		name     string
		cmp      int
		expected int
	}{
		{
			name:     "reverse_less_becomes_greater",
			cmp:      revInt(1, 2),
			expected: 1,
		},
		{
			name:     "reverse_equal_stays_equal",
			cmp:      revInt(5, 5),
			expected: 0,
		},
		{
			name:     "reverse_greater_becomes_less",
			cmp:      revInt(10, 2),
			expected: -1,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.cmp != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, tc.cmp)
			}
		})
	}
}
