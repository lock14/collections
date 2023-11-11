package bitset

import "fmt"

const (
	DefaultNumBits uint32 = 64
	wordSize              = 64
	panicFmt              = "bit %d out of range [0, %d)"
)

type BitSet struct {
	bits []uint64
}

type Config struct {
	NumBits uint32
}

type Option func(*Config)

func NumBits(n uint32) Option {
	return func(c *Config) {
		c.NumBits = n
	}
}

func New(opts ...Option) *BitSet {
	config := defaultConfig()
	for _, option := range opts {
		option(config)
	}
	return &BitSet{
		bits: make([]uint64, (config.NumBits/wordSize)+min(1, config.NumBits%wordSize)),
	}
}

func (b *BitSet) Clear(bit uint32) {
	index, shift := convert(bit)
	if index >= uint32(len(b.bits)) {
		panic(fmt.Sprintf(panicFmt, bit, len(b.bits)*wordSize))
	}
	b.bits[index] &= ^(1 << shift)
}

func (b *BitSet) Set(bit uint32) {
	index, shift := convert(bit)
	if index >= uint32(len(b.bits)) {
		panic(fmt.Sprintf(panicFmt, bit, len(b.bits)*wordSize))
	}
	b.bits[index] |= 1 << shift
}

func (b *BitSet) Get(bit uint32) bool {
	index, shift := convert(bit)
	if index >= uint32(len(b.bits)) {
		panic(fmt.Sprintf(panicFmt, bit, len(b.bits)*wordSize))
	}
	return (b.bits[index]>>shift)&1 == 1
}

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

func min(a uint32, b uint32) uint32 {
	if a <= b {
		return a
	}
	return b
}
