package bitset

import (
	"testing"
)

func TestSetBit(t *testing.T) {
	n := 128
	bitSet := New(NumBits(uint32(n)))
	for i := 0; i < n; i++ {
		bitSet.Set(uint32(i))
		if !bitSet.Get(uint32(i)) {
			t.Fatalf("excepted bit %d to be set, but it was not", i)
		}
	}
}
