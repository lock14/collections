package hashmap

import (
	"github.com/lock14/collections"
	"testing"
)

type Foo struct {
	bar int
}

type Person struct {
	Name string
	age  int
	Foo  Foo
}

func GetPutAndRemoveSameAsBuiltInTestCase[K comparable, V comparable](entries map[K]V) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		builtIn := make(map[K]V)
		hashMap := New[K, V]()

		for k, v := range entries {
			builtIn[k] = v
			hashMap.Put(k, v)
			if got, want := hashMap.Size(), len(builtIn); got != want {
				t.Errorf("unexpected number of entries: got %d, want %d", got, want)
			}
			biEntries := make([]struct {
				k K
				v V
			}, 0, len(builtIn))
			for k, v := range builtIn {
				biEntries = append(biEntries, struct {
					k K
					v V
				}{k: k, v: v})
			}
			hmEntries := make([]struct {
				k K
				v V
			}, 0, len(builtIn))
			for k, v := range hashMap.All() {
				hmEntries = append(hmEntries, struct {
					k K
					v V
				}{k: k, v: v})
			}
			if got, want := len(hmEntries), len(biEntries); got != want {
				t.Errorf("unexpected number of entries: got %d, want %d", got, want)
			}
		}
		if got, want := hashMap.Size(), len(builtIn); got != want {
			t.Errorf("unexpected number of entries: got %d, want %d", got, want)
		}
		for k := range entries {
			gotV, gotOk := hashMap.Get(k)
			wantV, wantOK := builtIn[k]
			if gotOk != wantOK {
				t.Errorf("got %v, want %v", gotOk, wantOK)
			}
			if gotV != wantV {
				t.Errorf("got %v, want %v", gotV, wantV)
			}
			if got, want := hashMap.Size(), len(builtIn); got != want {
				t.Errorf("unexpected number of entries: got %d, want %d", got, want)
			}
			biEntries := make([]struct {
				k K
				v V
			}, 0, len(builtIn))
			for k, v := range builtIn {
				biEntries = append(biEntries, struct {
					k K
					v V
				}{k: k, v: v})
			}
			hmEntries := make([]struct {
				k K
				v V
			}, 0, hashMap.Size())
			for k, v := range hashMap.All() {
				hmEntries = append(hmEntries, struct {
					k K
					v V
				}{k: k, v: v})
			}
			if got, want := len(hmEntries), len(biEntries); got != want {
				t.Errorf("unexpected number of entries: got %d, want %d", got, want)
			}
		}
		if got, want := hashMap.Size(), len(builtIn); got != want {
			t.Errorf("unexpected number of entries: got %d, want %d", got, want)
		}
		for k := range entries {
			delete(builtIn, k)
			hashMap.Remove(k)
			if got, want := hashMap.Size(), len(builtIn); got != want {
				t.Errorf("unexpected number of entries: got %d, want %d", got, want)
			}
			biEntries := make([]struct {
				k K
				v V
			}, 0, len(builtIn))
			for k, v := range builtIn {
				biEntries = append(biEntries, struct {
					k K
					v V
				}{k: k, v: v})
			}
			hmEntries := make([]struct {
				k K
				v V
			}, 0, hashMap.Size())
			for k, v := range hashMap.All() {
				hmEntries = append(hmEntries, struct {
					k K
					v V
				}{k: k, v: v})
			}
			if got, want := len(hmEntries), len(biEntries); got != want {
				t.Errorf("unexpected number of entries: got %d, want %d", got, want)
			}
		}
		if got, want := hashMap.Size(), len(builtIn); got != want {
			t.Errorf("unexpected number of entries: got %d, want %d", got, want)
		}
	}
}

func TestGetPutAndRemoveSameAsBuiltInMapInt(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		entries map[int]int
	}{
		{
			name:    "empty",
			entries: map[int]int{},
		},
		{
			name:    "one",
			entries: map[int]int{45: 2},
		},
		{
			name: "many",
			entries: map[int]int{
				531279: 235585,
				774428: 293791,
				508129: 869760,
				867759: 339288,
				517856: 441415,
				203188: 565133,
				109294: 270937,
				392301: 97964,
				933428: 183517,
				107478: 721549,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, GetPutAndRemoveSameAsBuiltInTestCase(tc.entries))
	}
}

func TestGetPutAndRemoveSameAsBuiltInMapFloat(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		entries map[float64]float64
	}{
		{
			name:    "empty",
			entries: map[float64]float64{},
		},
		{
			name:    "one",
			entries: map[float64]float64{45: 2},
		},
		{
			name: "many",
			entries: map[float64]float64{
				531279: 235585,
				774428: 293791,
				508129: 869760,
				867759: 339288,
				517856: 441415,
				203188: 565133,
				109294: 270937,
				392301: 97964,
				933428: 183517,
				107478: 721549,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, GetPutAndRemoveSameAsBuiltInTestCase(tc.entries))
	}
}

func TestGetPutAndRemoveSameAsBuiltInMapString(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		entries map[string]string
	}{
		{
			name:    "empty",
			entries: map[string]string{},
		},
		{
			name:    "one",
			entries: map[string]string{"45": "2"},
		},
		{
			name: "many",
			entries: map[string]string{
				"531279": "235585",
				"774428": "293791",
				"508129": "869760",
				"867759": "339288",
				"517856": "441415",
				"203188": "565133",
				"109294": "270937",
				"392301": "97964",
				"933428": "183517",
				"107478": "721549",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, GetPutAndRemoveSameAsBuiltInTestCase(tc.entries))
	}
}

func TestGetPutAndRemoveSameAsBuiltInMapStruct(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		entries map[Person]Person
	}{
		{
			name:    "empty",
			entries: map[Person]Person{},
		},
		{
			name: "one",
			entries: map[Person]Person{
				{Name: "bob", age: 21}: {Name: "mary", age: 19},
			},
		},
		{
			name: "many",
			entries: map[Person]Person{
				{Name: "bob", age: 21}:   {Name: "mary", age: 26},
				{Name: "andy", age: 45}:  {Name: "rachel", age: 73},
				{Name: "susan", age: 67}: {Name: "anne", age: 62},
				{Name: "john", age: 34}:  {Name: "krystal", age: 55},
				{Name: "lilly", age: 35}: {Name: "penelope", age: 50},
				{Name: "hank", age: 89}:  {Name: "fred", age: 41},
				{Name: "mark", age: 12}:  {Name: "rick", age: 44},
				{Name: "bill", age: 23}:  {Name: "nancy", age: 48},
				{Name: "ben", age: 48}:   {Name: "colin", age: 26},
				{Name: "paul", age: 39}:  {Name: "wendy", age: 24},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, GetPutAndRemoveSameAsBuiltInTestCase(tc.entries))
	}
}

func TestGetPutAndRemoveSameAsBuiltInMapStructPointers(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		entries map[*Person]*Person
	}{
		{
			name:    "empty",
			entries: map[*Person]*Person{},
		},
		{
			name: "one",
			entries: map[*Person]*Person{
				{Name: "bob", age: 21}: {Name: "mary", age: 19},
			},
		},
		{
			name: "many",
			entries: map[*Person]*Person{
				{Name: "bob", age: 21}:   {Name: "mary", age: 26},
				{Name: "andy", age: 45}:  {Name: "rachel", age: 73},
				{Name: "susan", age: 67}: {Name: "anne", age: 62},
				{Name: "john", age: 34}:  {Name: "krystal", age: 55},
				{Name: "lilly", age: 35}: {Name: "penelope", age: 50},
				{Name: "hank", age: 89}:  {Name: "fred", age: 41},
				{Name: "mark", age: 12}:  {Name: "rick", age: 44},
				{Name: "bill", age: 23}:  {Name: "nancy", age: 48},
				{Name: "ben", age: 48}:   {Name: "colin", age: 26},
				{Name: "paul", age: 39}:  {Name: "wendy", age: 24},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, GetPutAndRemoveSameAsBuiltInTestCase(tc.entries))
	}
}

func TestType(t *testing.T) {
	t.Parallel()

	mapType(New[int, int]())
}

func mapType[K, V any](_ collections.Map[K, V]) {}
