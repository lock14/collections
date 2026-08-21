package optional_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/lock14/collections/optional"
)

func TestConstructors(t *testing.T) {
	t.Parallel()

	var anyInt any = 42
	var anyStr any = "hello"

	cases := []struct {
		name        string
		opt         optional.Option[any]
		wantPresent bool
		wantVal     any
	}{
		{
			name:        "zero_value",
			opt:         optional.Option[any]{},
			wantPresent: false,
			wantVal:     nil,
		},
		{
			name:        "empty",
			opt:         optional.Empty[any](),
			wantPresent: false,
			wantVal:     nil,
		},
		{
			name:        "of_int",
			opt:         optional.Of[any](42),
			wantPresent: true,
			wantVal:     42,
		},
		{
			name:        "of_string",
			opt:         optional.Of[any]("hello"),
			wantPresent: true,
			wantVal:     "hello",
		},
		{
			name:        "of_ptr_nil",
			opt:         optional.OfPtr[any](nil),
			wantPresent: false,
			wantVal:     nil,
		},
		{
			name:        "of_ptr_non_nil_int",
			opt:         optional.OfPtr[any](&anyInt),
			wantPresent: true,
			wantVal:     42,
		},
		{
			name:        "of_ptr_non_nil_string",
			opt:         optional.OfPtr[any](&anyStr),
			wantPresent: true,
			wantVal:     "hello",
		},
		{
			name:        "of_ok_true",
			opt:         optional.OfOk[any]("ok_val", true),
			wantPresent: true,
			wantVal:     "ok_val",
		},
		{
			name:        "of_ok_false",
			opt:         optional.OfOk[any]("ignored", false),
			wantPresent: false,
			wantVal:     nil,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.opt.IsPresent(); got != tc.wantPresent {
				t.Errorf("IsPresent() = %v, want %v", got, tc.wantPresent)
			}
			if got := tc.opt.IsEmpty(); got == tc.wantPresent {
				t.Errorf("IsEmpty() = %v, want %v", got, !tc.wantPresent)
			}

			val, ok := tc.opt.Get()
			if ok != tc.wantPresent {
				t.Errorf("Get() ok = %v, want %v", ok, tc.wantPresent)
			}
			if tc.wantPresent && val != tc.wantVal {
				t.Errorf("Get() val = %v, want %v", val, tc.wantVal)
			}
		})
	}
}

func TestMustGet(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		opt       optional.Option[int]
		wantVal   int
		wantPanic bool
	}{
		{
			name:      "present_value",
			opt:       optional.Of(123),
			wantVal:   123,
			wantPanic: false,
		},
		{
			name:      "empty_option",
			opt:       optional.Empty[int](),
			wantPanic: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected MustGet() to panic on empty Option")
					}
				}()
				_ = tc.opt.MustGet()
			} else {
				if got := tc.opt.MustGet(); got != tc.wantVal {
					t.Errorf("MustGet() = %v, want %v", got, tc.wantVal)
				}
			}
		})
	}
}

func TestOrElse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		opt      optional.Option[string]
		fallback string
		want     string
	}{
		{
			name:     "present_returns_contained",
			opt:      optional.Of("present"),
			fallback: "fallback",
			want:     "present",
		},
		{
			name:     "empty_returns_fallback",
			opt:      optional.Empty[string](),
			fallback: "fallback",
			want:     "fallback",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.opt.OrElse(tc.fallback); got != tc.want {
				t.Errorf("OrElse() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestOrElseGet(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		opt            optional.Option[int]
		supplierVal    int
		expectSupplier bool
		want           int
	}{
		{
			name:           "present_does_not_invoke_supplier",
			opt:            optional.Of(10),
			supplierVal:    20,
			expectSupplier: false,
			want:           10,
		},
		{
			name:           "empty_invokes_supplier",
			opt:            optional.Empty[int](),
			supplierVal:    20,
			expectSupplier: true,
			want:           20,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			supplierCalled := false
			supplier := func() int {
				supplierCalled = true
				return tc.supplierVal
			}

			got := tc.opt.OrElseGet(supplier)
			if got != tc.want {
				t.Errorf("OrElseGet() = %v, want %v", got, tc.want)
			}
			if supplierCalled != tc.expectSupplier {
				t.Errorf("supplierCalled = %v, want %v", supplierCalled, tc.expectSupplier)
			}
		})
	}
}

func TestMap(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		opt         optional.Option[int]
		mapper      func(int) string
		wantPresent bool
		wantVal     string
	}{
		{
			name:        "present_transforms_type",
			opt:         optional.Of(123),
			mapper:      strconv.Itoa,
			wantPresent: true,
			wantVal:     "123",
		},
		{
			name:        "empty_stays_empty",
			opt:         optional.Empty[int](),
			mapper:      strconv.Itoa,
			wantPresent: false,
			wantVal:     "",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mapped := tc.opt.Map(tc.mapper)
			if got := mapped.IsPresent(); got != tc.wantPresent {
				t.Errorf("Map().IsPresent() = %v, want %v", got, tc.wantPresent)
			}
			if val, ok := mapped.Get(); ok != tc.wantPresent || (ok && val != tc.wantVal) {
				t.Errorf("Map().Get() = (%v, %v), want (%v, %v)", val, ok, tc.wantVal, tc.wantPresent)
			}
		})
	}
}

func TestFlatMap(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		opt         optional.Option[string]
		mapper      func(string) optional.Option[int]
		wantPresent bool
		wantVal     int
	}{
		{
			name: "present_to_present",
			opt:  optional.Of("42"),
			mapper: func(s string) optional.Option[int] {
				n, err := strconv.Atoi(s)
				if err != nil {
					return optional.Empty[int]()
				}
				return optional.Of(n)
			},
			wantPresent: true,
			wantVal:     42,
		},
		{
			name: "present_to_empty",
			opt:  optional.Of("not_a_number"),
			mapper: func(s string) optional.Option[int] {
				n, err := strconv.Atoi(s)
				if err != nil {
					return optional.Empty[int]()
				}
				return optional.Of(n)
			},
			wantPresent: false,
			wantVal:     0,
		},
		{
			name: "empty_to_empty",
			opt:  optional.Empty[string](),
			mapper: func(s string) optional.Option[int] {
				return optional.Of(999)
			},
			wantPresent: false,
			wantVal:     0,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			flatMapped := tc.opt.FlatMap(tc.mapper)
			if got := flatMapped.IsPresent(); got != tc.wantPresent {
				t.Errorf("FlatMap().IsPresent() = %v, want %v", got, tc.wantPresent)
			}
			if val, ok := flatMapped.Get(); ok != tc.wantPresent || (ok && val != tc.wantVal) {
				t.Errorf("FlatMap().Get() = (%v, %v), want (%v, %v)", val, ok, tc.wantVal, tc.wantPresent)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	t.Parallel()

	isEven := func(n int) bool {
		return n%2 == 0
	}

	cases := []struct {
		name        string
		opt         optional.Option[int]
		predicate   func(int) bool
		wantPresent bool
		wantVal     int
	}{
		{
			name:        "present_matching_predicate",
			opt:         optional.Of(4),
			predicate:   isEven,
			wantPresent: true,
			wantVal:     4,
		},
		{
			name:        "present_not_matching_predicate",
			opt:         optional.Of(5),
			predicate:   isEven,
			wantPresent: false,
			wantVal:     0,
		},
		{
			name:        "empty_stays_empty",
			opt:         optional.Empty[int](),
			predicate:   isEven,
			wantPresent: false,
			wantVal:     0,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			filtered := tc.opt.Filter(tc.predicate)
			if got := filtered.IsPresent(); got != tc.wantPresent {
				t.Errorf("Filter().IsPresent() = %v, want %v", got, tc.wantPresent)
			}
			if val, ok := filtered.Get(); ok != tc.wantPresent || (ok && val != tc.wantVal) {
				t.Errorf("Filter().Get() = (%v, %v), want (%v, %v)", val, ok, tc.wantVal, tc.wantPresent)
			}
		})
	}
}

func TestIfPresent(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		opt           optional.Option[string]
		expectAction  bool
		wantActionArg string
	}{
		{
			name:          "present_calls_action",
			opt:           optional.Of("payload"),
			expectAction:  true,
			wantActionArg: "payload",
		},
		{
			name:         "empty_does_not_call_action",
			opt:          optional.Empty[string](),
			expectAction: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			called := false
			var calledArg string
			tc.opt.IfPresent(func(v string) {
				called = true
				calledArg = v
			})

			if called != tc.expectAction {
				t.Errorf("action called = %v, want %v", called, tc.expectAction)
			}
			if tc.expectAction && calledArg != tc.wantActionArg {
				t.Errorf("action arg = %v, want %v", calledArg, tc.wantActionArg)
			}
		})
	}
}

func TestIfPresentOrElse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name              string
		opt               optional.Option[string]
		expectAction      bool
		expectEmptyAction bool
		wantActionArg     string
	}{
		{
			name:              "present_calls_action",
			opt:               optional.Of("hello"),
			expectAction:      true,
			expectEmptyAction: false,
			wantActionArg:     "hello",
		},
		{
			name:              "empty_calls_empty_action",
			opt:               optional.Empty[string](),
			expectAction:      false,
			expectEmptyAction: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actionCalled := false
			emptyActionCalled := false
			var actionArg string

			tc.opt.IfPresentOrElse(
				func(v string) {
					actionCalled = true
					actionArg = v
				},
				func() {
					emptyActionCalled = true
				},
			)

			if actionCalled != tc.expectAction {
				t.Errorf("actionCalled = %v, want %v", actionCalled, tc.expectAction)
			}
			if emptyActionCalled != tc.expectEmptyAction {
				t.Errorf("emptyActionCalled = %v, want %v", emptyActionCalled, tc.expectEmptyAction)
			}
			if tc.expectAction && actionArg != tc.wantActionArg {
				t.Errorf("actionArg = %v, want %v", actionArg, tc.wantActionArg)
			}
		})
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		opt  fmt.Stringer
		want string
	}{
		{
			name: "present_int",
			opt:  optional.Of(42),
			want: "Some(42)",
		},
		{
			name: "present_string",
			opt:  optional.Of("gopher"),
			want: "Some(gopher)",
		},
		{
			name: "empty_int",
			opt:  optional.Empty[int](),
			want: "None",
		},
		{
			name: "empty_string",
			opt:  optional.Empty[string](),
			want: "None",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.opt.String(); got != tc.want {
				t.Errorf("String() = %v, want %v", got, tc.want)
			}
		})
	}
}
