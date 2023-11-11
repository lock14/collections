package bitset

import "fmt"

const (
    DefaultNumBits uint32 = 64
    wordSize              = 64
    panicFmt              = "bit %d out of range [0, %d)"
)

// BitSet represents a vector of bits.
type BitSet struct {
    bits []uint64
}

// Config holds the values for configuring a BitSet.
type Config struct {
    NumBits uint32
}

// Option configures a BitSet config
type Option func(*Config)

// NumBits provides the option to set the number of bits used in a BitSet.
func NumBits(n uint32) Option {
    return func(c *Config) {
        c.NumBits = n
    }
}

// New creates a BitSet whose initial size is large enough to explicitly
// represent bits with indices in the range 0 through NumBits-1. If no
// configuration is used the DefaultNumBits is used as the number of bits.
// All bits are initially false.
func New(opts ...Option) *BitSet {
    config := defaultConfig()
    for _, option := range opts {
        option(config)
    }
    return &BitSet{
        bits: make([]uint64, (config.NumBits/wordSize)+min(1, config.NumBits%wordSize)),
    }
}

// Clear sets the bit specified by the index to false.
func (b *BitSet) Clear(bit uint32) {
    index, shift := convert(bit)
    if index >= uint32(len(b.bits)) {
        panic(fmt.Sprintf(panicFmt, bit, len(b.bits)*wordSize))
    }
    b.bits[index] &= ^(1 << shift)
}

// Set sets the bit at the specified index to true.
func (b *BitSet) Set(bit uint32) {
    index, shift := convert(bit)
    if index >= uint32(len(b.bits)) {
        panic(fmt.Sprintf(panicFmt, bit, len(b.bits)*wordSize))
    }
    b.bits[index] |= 1 << shift
}

// Get returns the value of the bit with the specified index.
func (b *BitSet) Get(bit uint32) bool {
    index, shift := convert(bit)
    if index >= uint32(len(b.bits)) {
        panic(fmt.Sprintf(panicFmt, bit, len(b.bits)*wordSize))
    }
    return (b.bits[index]>>shift)&1 == 1
}

// Size returns the number of bits in this bit set.
func (b *BitSet) Size() uint32 {
    return uint32(len(b.bits)) * wordSize
}

// Flip sets each bit to the complement of its current value.
func (b *BitSet) Flip() {
    for i := 0; i < len(b.bits); i++ {
        b.bits[i] = ^b.bits[i]
    }
}

func convert(bit uint32) (uint32, uint32) {
    return bit / wordSize, bit % wordSize
}

func defaultConfig() *Config {
    return &Config{
        NumBits: DefaultNumBits,
    }
}
