package result_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/lock14/collections/result"
)

func TestConstructors(t *testing.T) {
	t.Parallel()

	errExample := errors.New("sample error")

	cases := []struct {
		name    string
		res     result.Result[int, error]
		wantOk  bool
		wantVal int
		wantErr error
	}{
		{
			name:    "ok_constructor",
			res:     result.Ok[int, error](42),
			wantOk:  true,
			wantVal: 42,
			wantErr: nil,
		},
		{
			name:    "err_constructor",
			res:     result.Err[int, error](errExample),
			wantOk:  false,
			wantVal: 0,
			wantErr: errExample,
		},
		{
			name:    "of_with_nil_err",
			res:     result.Of(100, nil),
			wantOk:  true,
			wantVal: 100,
			wantErr: nil,
		},
		{
			name:    "of_with_non_nil_err",
			res:     result.Of(100, errExample),
			wantOk:  false,
			wantVal: 0,
			wantErr: errExample,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.res.IsOk(); got != tc.wantOk {
				t.Errorf("IsOk() = %v, want %v", got, tc.wantOk)
			}
			if got := tc.res.IsErr(); got == tc.wantOk {
				t.Errorf("IsErr() = %v, want %v", got, !tc.wantOk)
			}

			val, ok := tc.res.Ok()
			if ok != tc.wantOk {
				t.Errorf("Ok() ok = %v, want %v", ok, tc.wantOk)
			}
			if tc.wantOk && val != tc.wantVal {
				t.Errorf("Ok() val = %v, want %v", val, tc.wantVal)
			}

			errVal, isErr := tc.res.Err()
			if isErr == tc.wantOk {
				t.Errorf("Err() isErr = %v, want %v", isErr, !tc.wantOk)
			}
			if !tc.wantOk && !errors.Is(errVal, tc.wantErr) {
				t.Errorf("Err() err = %v, want %v", errVal, tc.wantErr)
			}

			uVal, uErr := tc.res.Unwrap()
			if tc.wantOk {
				if uVal != tc.wantVal || uErr != nil {
					t.Errorf("Unwrap() = (%v, %v), want (%v, nil)", uVal, uErr, tc.wantVal)
				}
			} else {
				if !errors.Is(uErr, tc.wantErr) {
					t.Errorf("Unwrap() err = %v, want %v", uErr, tc.wantErr)
				}
			}
		})
	}
}

func TestMustGet(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		res       result.Result[string, string]
		wantVal   string
		wantPanic bool
	}{
		{
			name:      "ok_returns_value",
			res:       result.Ok[string, string]("success"),
			wantVal:   "success",
			wantPanic: false,
		},
		{
			name:      "err_panics",
			res:       result.Err[string, string]("failure reason"),
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
						t.Errorf("expected MustGet to panic on Err")
					}
				}()
				_ = tc.res.MustGet()
			} else {
				if got := tc.res.MustGet(); got != tc.wantVal {
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
		res      result.Result[int, string]
		fallback int
		want     int
	}{
		{
			name:     "ok_returns_contained",
			res:      result.Ok[int, string](42),
			fallback: 0,
			want:     42,
		},
		{
			name:     "err_returns_fallback",
			res:      result.Err[int, string]("err"),
			fallback: 99,
			want:     99,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.res.OrElse(tc.fallback); got != tc.want {
				t.Errorf("OrElse() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestOrElseGet(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		res            result.Result[int, string]
		supplierVal    int
		expectSupplier bool
		want           int
	}{
		{
			name:           "ok_does_not_invoke_supplier",
			res:            result.Ok[int, string](10),
			supplierVal:    20,
			expectSupplier: false,
			want:           10,
		},
		{
			name:           "err_invokes_supplier",
			res:            result.Err[int, string]("err"),
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

			got := tc.res.OrElseGet(supplier)
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
		name    string
		res     result.Result[int, string]
		mapper  func(int) string
		wantOk  bool
		wantVal string
		wantErr string
	}{
		{
			name:    "ok_maps_value",
			res:     result.Ok[int, string](123),
			mapper:  strconv.Itoa,
			wantOk:  true,
			wantVal: "123",
			wantErr: "",
		},
		{
			name:    "err_preserves_error",
			res:     result.Err[int, string]("original error"),
			mapper:  strconv.Itoa,
			wantOk:  false,
			wantVal: "",
			wantErr: "original error",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mapped := tc.res.Map(tc.mapper)
			if got := mapped.IsOk(); got != tc.wantOk {
				t.Errorf("Map().IsOk() = %v, want %v", got, tc.wantOk)
			}
			if tc.wantOk {
				if val, ok := mapped.Ok(); !ok || val != tc.wantVal {
					t.Errorf("Map().Ok() = (%v, %v), want (%v, true)", val, ok, tc.wantVal)
				}
			} else {
				if errVal, isErr := mapped.Err(); !isErr || errVal != tc.wantErr {
					t.Errorf("Map().Err() = (%v, %v), want (%v, true)", errVal, isErr, tc.wantErr)
				}
			}
		})
	}
}

func TestMapErr(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		res     result.Result[int, int]
		mapper  func(int) string
		wantOk  bool
		wantVal int
		wantErr string
	}{
		{
			name:    "ok_preserves_value",
			res:     result.Ok[int, int](42),
			mapper:  strconv.Itoa,
			wantOk:  true,
			wantVal: 42,
			wantErr: "",
		},
		{
			name:    "err_maps_error",
			res:     result.Err[int, int](404),
			mapper:  func(code int) string { return fmt.Sprintf("HTTP %d", code) },
			wantOk:  false,
			wantVal: 0,
			wantErr: "HTTP 404",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mapped := tc.res.MapErr(tc.mapper)
			if got := mapped.IsOk(); got != tc.wantOk {
				t.Errorf("MapErr().IsOk() = %v, want %v", got, tc.wantOk)
			}
			if tc.wantOk {
				if val, ok := mapped.Ok(); !ok || val != tc.wantVal {
					t.Errorf("MapErr().Ok() = (%v, %v), want (%v, true)", val, ok, tc.wantVal)
				}
			} else {
				if errVal, isErr := mapped.Err(); !isErr || errVal != tc.wantErr {
					t.Errorf("MapErr().Err() = (%v, %v), want (%v, true)", errVal, isErr, tc.wantErr)
				}
			}
		})
	}
}

func TestFlatMap(t *testing.T) {
	t.Parallel()

	parseEven := func(s string) result.Result[int, string] {
		n, err := strconv.Atoi(s)
		if err != nil {
			return result.Err[int, string]("invalid integer")
		}
		if n%2 != 0 {
			return result.Err[int, string]("not even")
		}
		return result.Ok[int, string](n)
	}

	cases := []struct {
		name    string
		res     result.Result[string, string]
		wantOk  bool
		wantVal int
		wantErr string
	}{
		{
			name:    "ok_to_ok",
			res:     result.Ok[string, string]("42"),
			wantOk:  true,
			wantVal: 42,
			wantErr: "",
		},
		{
			name:    "ok_to_err_invalid",
			res:     result.Ok[string, string]("abc"),
			wantOk:  false,
			wantVal: 0,
			wantErr: "invalid integer",
		},
		{
			name:    "ok_to_err_odd",
			res:     result.Ok[string, string]("43"),
			wantOk:  false,
			wantVal: 0,
			wantErr: "not even",
		},
		{
			name:    "err_propagates",
			res:     result.Err[string, string]("initial failure"),
			wantOk:  false,
			wantVal: 0,
			wantErr: "initial failure",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			flatMapped := tc.res.FlatMap(parseEven)
			if got := flatMapped.IsOk(); got != tc.wantOk {
				t.Errorf("FlatMap().IsOk() = %v, want %v", got, tc.wantOk)
			}
			if tc.wantOk {
				if val, ok := flatMapped.Ok(); !ok || val != tc.wantVal {
					t.Errorf("FlatMap().Ok() = (%v, %v), want (%v, true)", val, ok, tc.wantVal)
				}
			} else {
				if errVal, isErr := flatMapped.Err(); !isErr || errVal != tc.wantErr {
					t.Errorf("FlatMap().Err() = (%v, %v), want (%v, true)", errVal, isErr, tc.wantErr)
				}
			}
		})
	}
}

func TestIfOk(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		res          result.Result[string, string]
		expectAction bool
		wantArg      string
	}{
		{
			name:         "ok_invokes_action",
			res:          result.Ok[string, string]("data"),
			expectAction: true,
			wantArg:      "data",
		},
		{
			name:         "err_does_not_invoke_action",
			res:          result.Err[string, string]("err"),
			expectAction: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			called := false
			var calledArg string
			tc.res.IfOk(func(v string) {
				called = true
				calledArg = v
			})

			if called != tc.expectAction {
				t.Errorf("action called = %v, want %v", called, tc.expectAction)
			}
			if tc.expectAction && calledArg != tc.wantArg {
				t.Errorf("action arg = %v, want %v", calledArg, tc.wantArg)
			}
		})
	}
}

func TestIfErr(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		res          result.Result[string, string]
		expectAction bool
		wantArg      string
	}{
		{
			name:         "ok_does_not_invoke_action",
			res:          result.Ok[string, string]("data"),
			expectAction: false,
		},
		{
			name:         "err_invokes_action",
			res:          result.Err[string, string]("err"),
			expectAction: true,
			wantArg:      "err",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			called := false
			var calledArg string
			tc.res.IfErr(func(e string) {
				called = true
				calledArg = e
			})

			if called != tc.expectAction {
				t.Errorf("action called = %v, want %v", called, tc.expectAction)
			}
			if tc.expectAction && calledArg != tc.wantArg {
				t.Errorf("action arg = %v, want %v", calledArg, tc.wantArg)
			}
		})
	}
}

func TestIfOkOrElse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name            string
		res             result.Result[string, string]
		expectOkAction  bool
		expectErrAction bool
		wantArg         string
	}{
		{
			name:            "ok_branch",
			res:             result.Ok[string, string]("ok_val"),
			expectOkAction:  true,
			expectErrAction: false,
			wantArg:         "ok_val",
		},
		{
			name:            "err_branch",
			res:             result.Err[string, string]("err_val"),
			expectOkAction:  false,
			expectErrAction: true,
			wantArg:         "err_val",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			okCalled := false
			errCalled := false
			var capturedArg string

			tc.res.IfOkOrElse(
				func(v string) {
					okCalled = true
					capturedArg = v
				},
				func(e string) {
					errCalled = true
					capturedArg = e
				},
			)

			if okCalled != tc.expectOkAction {
				t.Errorf("okCalled = %v, want %v", okCalled, tc.expectOkAction)
			}
			if errCalled != tc.expectErrAction {
				t.Errorf("errCalled = %v, want %v", errCalled, tc.expectErrAction)
			}
			if capturedArg != tc.wantArg {
				t.Errorf("capturedArg = %v, want %v", capturedArg, tc.wantArg)
			}
		})
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		res  fmt.Stringer
		want string
	}{
		{
			name: "ok_int",
			res:  result.Ok[int, string](42),
			want: "Ok(42)",
		},
		{
			name: "ok_string",
			res:  result.Ok[string, error]("success"),
			want: "Ok(success)",
		},
		{
			name: "err_string",
			res:  result.Err[int, string]("failed"),
			want: "Err(failed)",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.res.String(); got != tc.want {
				t.Errorf("String() = %v, want %v", got, tc.want)
			}
		})
	}
}
