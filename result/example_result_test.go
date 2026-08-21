package result_test

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/lock14/collections/result"
)

func ExampleResult() {
	resOk := result.Ok[int, string](42)
	resErr := result.Err[int, string]("out of bounds")

	fmt.Println("resOk.IsOk():", resOk.IsOk())
	fmt.Println("resOk.MustGet():", resOk.MustGet())
	fmt.Println("resErr.IsErr():", resErr.IsErr())
	fmt.Println("resErr.OrElse(0):", resErr.OrElse(0))

	// Output:
	// resOk.IsOk(): true
	// resOk.MustGet(): 42
	// resErr.IsErr(): true
	// resErr.OrElse(0): 0
}

func ExampleOf() {
	parse := func(s string) (int, error) {
		return strconv.Atoi(s)
	}

	valid := result.Of(parse("123"))
	invalid := result.Of(parse("not_a_number"))

	fmt.Println("valid:", valid.OrElse(0))
	fmt.Println("invalid:", invalid.OrElse(0))

	// Output:
	// valid: 123
	// invalid: 0
}

func ExampleResult_Map() {
	res := result.Ok[int, string](42)
	mapped := res.Map(func(value int) string {
		return fmt.Sprintf("answer=%d", value)
	})

	errRes := result.Err[int, string]("failure")
	errMapped := errRes.Map(func(value int) string {
		return fmt.Sprintf("answer=%d", value)
	})

	fmt.Println(mapped)
	fmt.Println(errMapped)

	// Output:
	// Ok(answer=42)
	// Err(failure)
}

func ExampleResult_MapErr() {
	res := result.Err[string, int](404)
	mapped := res.MapErr(func(err int) string {
		return fmt.Sprintf("HTTP error %d", err)
	})

	fmt.Println(mapped)

	// Output:
	// Err(HTTP error 404)
}

func ExampleResult_FlatMap() {
	divide := func(a, b int) result.Result[int, string] {
		if b == 0 {
			return result.Err[int, string]("division by zero")
		}
		return result.Ok[int, string](a / b)
	}

	res := result.Ok[int, string](100).
		FlatMap(func(value int) result.Result[int, string] {
			return divide(value, 2)
		}).
		FlatMap(func(value int) result.Result[int, string] {
			return divide(value, 5)
		})

	errRes := result.Ok[int, string](100).
		FlatMap(func(value int) result.Result[int, string] {
			return divide(value, 0)
		})

	fmt.Println("res:", res)
	fmt.Println("errRes:", errRes)

	// Output:
	// res: Ok(10)
	// errRes: Err(division by zero)
}

func ExampleResult_OrElseGet() {
	expensiveFallback := func() string {
		return "computed-default"
	}

	okRes := result.Ok[string, error]("cached")
	errRes := result.Err[string, error](errors.New("miss"))

	fmt.Println(okRes.OrElseGet(expensiveFallback))
	fmt.Println(errRes.OrElseGet(expensiveFallback))

	// Output:
	// cached
	// computed-default
}

func ExampleResult_IfOkOrElse() {
	handle := func(res result.Result[string, string]) {
		res.IfOkOrElse(
			func(value string) {
				fmt.Printf("Success: %s\n", value)
			},
			func(err string) {
				fmt.Printf("Failure: %s\n", err)
			},
		)
	}

	handle(result.Ok[string, string]("data loaded"))
	handle(result.Err[string, string]("connection timeout"))

	// Output:
	// Success: data loaded
	// Failure: connection timeout
}
