package pair_test

import (
	"fmt"
	"github.com/lock14/collections/pair"
)

func ExamplePair() {
	// pair.Pair is a simple, generic composite type.
	// It is especially useful as a map key when you need a 2D coordinate or a composite key.

	type Coordinate pair.Pair[int, int]

	grid := make(map[Coordinate]string)

	// Create a new coordinate pair
	origin := pair.New(0, 0)
	target := pair.New(5, 10)

	grid[Coordinate(origin)] = "Start"
	grid[Coordinate(target)] = "End"

	fmt.Println("At origin:", grid[Coordinate(pair.New(0, 0))])
	fmt.Println("At target:", grid[Coordinate(pair.New(5, 10))])

	// Output:
	// At origin: Start
	// At target: End
}

func ExamplePair_Fst() {
	p := pair.New("hello", 42)
	fmt.Println(p.Fst())

	// Output:
	// hello
}

func ExamplePair_Snd() {
	p := pair.New("hello", 42)
	fmt.Println(p.Snd())

	// Output:
	// 42
}

func ExamplePair_Unwrap() {
	p := pair.New("hello", 42)
	fst, snd := p.Unwrap()
	fmt.Println(fst)
	fmt.Println(snd)

	// Output:
	// hello
	// 42
}

func ExamplePair_String() {
	p := pair.New("hello", 42)
	fmt.Println(p.String())

	// Output:
	// (hello, 42)
}
