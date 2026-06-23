package pair

import (
	"testing"
)

func TestPair(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		check func(*testing.T)
	}{
		{
			name: "basic_types",
			check: func(t *testing.T) {
				p := New(1, "hello")
				if p.Fst() != 1 {
					t.Errorf("expected 1")
				}
				if p.Snd() != "hello" {
					t.Errorf("expected hello")
				}
				f, s := p.Unwrap()
				if f != 1 || s != "hello" {
					t.Errorf("unwrap failed")
				}
			},
		},
		{
			name: "string_format",
			check: func(t *testing.T) {
				p := New(10, 20)
				if str := p.String(); str != "(10, 20)" {
					t.Errorf("wrong string format: %s", str)
				}
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.check(t)
		})
	}
}
