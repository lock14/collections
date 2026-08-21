package optional_test

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lock14/collections/optional"
)

func ExampleOption() {
	present := optional.Of("hello")
	empty := optional.Empty[string]()

	fmt.Println("present.IsPresent():", present.IsPresent())
	fmt.Println("present.MustGet():", present.MustGet())
	fmt.Println("empty.IsEmpty():", empty.IsEmpty())
	fmt.Println("empty.OrElse(\"default\"):", empty.OrElse("default"))

	// Output:
	// present.IsPresent(): true
	// present.MustGet(): hello
	// empty.IsEmpty(): true
	// empty.OrElse("default"): default
}

func ExampleOfOk() {
	scores := map[string]int{
		"alice": 95,
		"bob":   80,
	}

	aliceVal, aliceOk := scores["alice"]
	aliceScore := optional.OfOk(aliceVal, aliceOk)

	charlieVal, charlieOk := scores["charlie"]
	charlieScore := optional.OfOk(charlieVal, charlieOk)

	fmt.Println("Alice:", aliceScore.OrElse(0))
	fmt.Println("Charlie:", charlieScore.OrElse(0))

	// Output:
	// Alice: 95
	// Charlie: 0
}

func ExampleOfPtr() {
	val := "found"
	var presentPtr *string = &val
	var nilPtr *string = nil

	optPresent := optional.OfPtr(presentPtr)
	optEmpty := optional.OfPtr(nilPtr)

	fmt.Println(optPresent)
	fmt.Println(optEmpty)

	// Output:
	// Some(found)
	// None
}

func ExampleOption_Map() {
	intOpt := optional.Of(42)
	strOpt := intOpt.Map(func(value int) string {
		return fmt.Sprintf("value is %d", value)
	})

	emptyIntOpt := optional.Empty[int]()
	emptyStrOpt := emptyIntOpt.Map(func(value int) string {
		return fmt.Sprintf("value is %d", value)
	})

	fmt.Println(strOpt)
	fmt.Println(emptyStrOpt)

	// Output:
	// Some(value is 42)
	// None
}

func ExampleOption_FlatMap() {
	parsePositiveInt := func(value string) optional.Option[int] {
		n, err := strconv.Atoi(value)
		if err != nil || n <= 0 {
			return optional.Empty[int]()
		}
		return optional.Of(n)
	}

	valid := optional.Of("100").FlatMap(parsePositiveInt)
	invalid := optional.Of("abc").FlatMap(parsePositiveInt)
	negative := optional.Of("-5").FlatMap(parsePositiveInt)

	fmt.Println("valid:", valid)
	fmt.Println("invalid:", invalid)
	fmt.Println("negative:", negative)

	// Output:
	// valid: Some(100)
	// invalid: None
	// negative: None
}

func ExampleOption_Filter() {
	opt := optional.Of("hello world")

	longString := opt.Filter(func(value string) bool {
		return len(value) > 5
	})
	startsWithZ := opt.Filter(func(value string) bool {
		return strings.HasPrefix(value, "z")
	})

	fmt.Println("longString:", longString)
	fmt.Println("startsWithZ:", startsWithZ)

	// Output:
	// longString: Some(hello world)
	// startsWithZ: None
}

func ExampleOption_OrElseGet() {
	expensiveCompute := func() string {
		return "computed-fallback"
	}

	present := optional.Of("cached-value")
	empty := optional.Empty[string]()

	fmt.Println("present:", present.OrElseGet(expensiveCompute))
	fmt.Println("empty:", empty.OrElseGet(expensiveCompute))

	// Output:
	// present: cached-value
	// empty: computed-fallback
}

func ExampleOption_IfPresent() {
	opt := optional.Of("send email")
	opt.IfPresent(func(value string) {
		fmt.Println("Action:", value)
	})

	empty := optional.Empty[string]()
	empty.IfPresent(func(value string) {
		fmt.Println("Action:", value)
	})

	// Output:
	// Action: send email
}

func ExampleOption_IfPresentOrElse() {
	present := optional.Of("admin")
	empty := optional.Empty[string]()

	present.IfPresentOrElse(
		func(value string) {
			fmt.Printf("User role: %s\n", value)
		},
		func() {
			fmt.Println("User role: guest")
		},
	)

	empty.IfPresentOrElse(
		func(value string) {
			fmt.Printf("User role: %s\n", value)
		},
		func() {
			fmt.Println("User role: guest")
		},
	)

	// Output:
	// User role: admin
	// User role: guest
}
