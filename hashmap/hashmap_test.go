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

func GetAndPutSameAsBuiltInTestCase[K comparable, V comparable](entries map[K]V) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		built := make(map[K]V)
		hashMap := New[K, V]()

		for k, v := range entries {
			built[k] = v
			hashMap.Put(k, v)
		}
		for k := range built {
			gotV, gotOk := hashMap.Get(k)
			wantV, wantOK := built[k]
			if gotOk != wantOK {
				t.Errorf("got %v, want %v", gotOk, wantOK)
			}
			if gotV != wantV {
				t.Errorf("got %v, want %v", gotV, wantV)
			}
		}
	}
}

func TestGetAndPutSameAsBuiltInMapInt(t *testing.T) {
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
		t.Run(tc.name, GetAndPutSameAsBuiltInTestCase(tc.entries))
	}
}

func TestGetAndPutSameAsBuiltInMapFloat(t *testing.T) {
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
		t.Run(tc.name, GetAndPutSameAsBuiltInTestCase(tc.entries))
	}
}

func TestGetAndPutSameAsBuiltInMapString(t *testing.T) {
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
		t.Run(tc.name, GetAndPutSameAsBuiltInTestCase(tc.entries))
	}
}

func TestGetAndPutSameAsBuiltInMapStruct(t *testing.T) {
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
		t.Run(tc.name, GetAndPutSameAsBuiltInTestCase(tc.entries))
	}
}

func TestGetAndPutSameAsBuiltInMapStructPointers(t *testing.T) {
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
		t.Run(tc.name, GetAndPutSameAsBuiltInTestCase(tc.entries))
	}
}
