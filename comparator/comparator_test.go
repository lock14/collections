package comparator

import (
	"math"
	"testing"
)

func TestNaturalOrder(t *testing.T) {
	t.Parallel()
	nan := math.NaN()
	negZero := math.Copysign(0, -1)

	cases := []struct {
		name     string
		cmp      int
		expected int
	}{
		{
			name:     "int_less",
			cmp:      NaturalOrder[int]()(1, 2),
			expected: -1,
		},
		{
			name:     "int_equal",
			cmp:      NaturalOrder[int]()(5, 5),
			expected: 0,
		},
		{
			name:     "int_greater",
			cmp:      NaturalOrder[int]()(10, 2),
			expected: 1,
		},
		{
			name:     "string_less",
			cmp:      NaturalOrder[string]()("apple", "banana"),
			expected: -1,
		},
		{
			name:     "string_equal",
			cmp:      NaturalOrder[string]()("cat", "cat"),
			expected: 0,
		},
		{
			name:     "string_greater",
			cmp:      NaturalOrder[string]()("zebra", "yak"),
			expected: 1,
		},
		{
			name:     "float_less",
			cmp:      NaturalOrder[float64]()(1.5, 2.5),
			expected: -1,
		},
		{
			name:     "float_equal",
			cmp:      NaturalOrder[float64]()(3.14, 3.14),
			expected: 0,
		},
		{
			name:     "float_greater",
			cmp:      NaturalOrder[float64]()(9.9, 8.8),
			expected: 1,
		},
		{
			name:     "float_nan_less_than_number",
			cmp:      NaturalOrder[float64]()(nan, 1.0),
			expected: -1,
		},
		{
			name:     "float_nan_equal_nan",
			cmp:      NaturalOrder[float64]()(nan, nan),
			expected: 0,
		},
		{
			name:     "float_neg_zero_equal_pos_zero",
			cmp:      NaturalOrder[float64]()(negZero, 0.0),
			expected: 0,
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
	revInt := Reverse(NaturalOrder[int]())
	revFloat := Reverse(NaturalOrder[float64]())
	nan := math.NaN()

	minIntCmp := Reverse(Comparator[int](func(a, b int) int {
		if a < b {
			return math.MinInt
		}
		if a > b {
			return 1
		}
		return 0
	}))

	cases := []struct {
		name string
		got  int
		want func(int) bool
		desc string
	}{
		{
			name: "reverse_less_becomes_greater",
			got:  revInt(1, 2),
			want: func(v int) bool { return v == 1 },
			desc: "1",
		},
		{
			name: "reverse_equal_stays_equal",
			got:  revInt(5, 5),
			want: func(v int) bool { return v == 0 },
			desc: "0",
		},
		{
			name: "reverse_greater_becomes_less",
			got:  revInt(10, 2),
			want: func(v int) bool { return v == -1 },
			desc: "-1",
		},
		{
			name: "reverse_min_int_does_not_overflow",
			got:  minIntCmp(1, 2),
			want: func(v int) bool { return v > 0 },
			desc: "> 0",
		},
		{
			name: "reverse_min_int_greater_becomes_negative",
			got:  minIntCmp(2, 1),
			want: func(v int) bool { return v < 0 },
			desc: "< 0",
		},
		{
			name: "reverse_min_int_equal_stays_zero",
			got:  minIntCmp(2, 2),
			want: func(v int) bool { return v == 0 },
			desc: "0",
		},
		{
			name: "reverse_float_nan_greater_than_number",
			got:  revFloat(nan, 1.0),
			want: func(v int) bool { return v == 1 },
			desc: "1",
		},
		{
			name: "reverse_float_nan_equal_nan",
			got:  revFloat(nan, nan),
			want: func(v int) bool { return v == 0 },
			desc: "0",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if !tc.want(tc.got) {
				t.Errorf("expected %s, got %d", tc.desc, tc.got)
			}
		})
	}
}
