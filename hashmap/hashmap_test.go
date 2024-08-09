package hashmap

import "testing"

type Foo struct {
	bar int
}

type Person struct {
	Name string
	age  int
	Foo  Foo
}

func TestPut(t *testing.T) {
	t.Parallel()
	hm := New[int, int]()
	hm.Put(2, 3)

	if got, ok := hm.Get(2); !ok || got != 3 {
		t.Errorf("got/want(%d, %t)/(%+v, %t)", got, ok, Person{Name: "bob", age: 5}, true)
	}
}

func TestHash(t *testing.T) {
	t.Parallel()

	p1 := &Person{
		Name: "bob",
		age:  5,
	}
	p2 := &Person{
		Name: "bob",
		age:  5,
	}
	m := make(map[Person]int)
	hm := New[Person, int]()

	m[*p1] = 2
	hm.Put(*p1, 2)

	if p1 == p2 {
		t.Errorf("BAD")
	}
	if _, ok := m[*p2]; !ok {
		t.Errorf("Did not get key")
	}
	if _, ok := hm.Get(*p2); !ok {
		t.Errorf("Did not get key")
	}

}
