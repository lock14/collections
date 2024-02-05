package bitset

import (
	"fmt"
	"github.com/lock14/collections/iterator"
	"github.com/lock14/collections/util"
	"strings"
)

const (
	DefaultNumBits = 64
	wordSize       = 64
	wordFmt        = "%016X"
)

// BitSet represents a vector of bits that grows as needed.
type BitSet struct {
	bits []uint64
}

// Config holds the values for configuring a BitSet.
type Config struct {
	NumBits int
}

// Option configures a BitSet config
type Option func(*Config)

// NumBits provides the option to set the number of bits used in a BitSet.
func NumBits(n int) Option {
	return func(c *Config) {
		c.NumBits = n
	}
}

// iterator over the set bits
type setBitIterator struct {
	bitSet   *BitSet
	bitIndex int
}

// iterator over the unset bits
type unSetBitIterator struct {
	bitSet   *BitSet
	bitIndex int
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
	ensureNonNegative(config.NumBits)
	return &BitSet{
		bits: make([]uint64, (config.NumBits/wordSize)+min(1, config.NumBits%wordSize)),
	}
}

// Clear sets the bit specified by the index to false.
func (b *BitSet) Clear(bit int) {
	index, shift := convert(bit)
	b.ensureSize(index)
	b.bits[index] &= ^(1 << shift)
}

// Set sets the bit at the specified index to true.
func (b *BitSet) Set(bit int) {
	index, shift := convert(bit)
	b.ensureSize(index)
	b.bits[index] |= 1 << shift
}

// Get returns the value of the bit with the specified index.
func (b *BitSet) Get(bit int) bool {
	index, shift := convert(bit)
	b.ensureSize(index)
	return (b.bits[index]>>shift)&1 == 1
}

// Size returns the number of bits in this bit set.
func (b *BitSet) Size() int {
	return len(b.bits) * wordSize
}

// Flip sets each bit to the complement of its current value. This call is
// equivalent to b.FlipRange(0, b.Size())
func (b *BitSet) Flip() {
	for i := 0; i < len(b.bits); i++ {
		b.bits[i] = ^b.bits[i]
	}
}

// FlipRange sets each bit from the specified start bit (inclusive) to the
// specified end bit (exclusive) to the complement of its current value.
func (b *BitSet) FlipRange(start int, end int) {
	startIndex, startShift := convert(start)
	endIndex, endShift := convert(end)
	if end != b.Size() {
		b.ensureSize(endIndex)
	}

	startMask := ^(^uint64(0) << startShift)
	endMask := ^uint64(0) << endShift

	if startIndex == endIndex {
		// flip middle bits, keep upper and lower bits the same
		middleMask := ^(startMask | endMask)
		lowerBits := b.bits[startIndex] & startMask
		middleBits := (^b.bits[startIndex]) & middleMask
		upperBits := b.bits[endIndex] & endMask
		b.bits[startIndex] = lowerBits | middleBits | upperBits
	} else {
		// flip upper bits, keep lower bits the same
		lowerBits := b.bits[startIndex] & startMask
		upperBits := (^b.bits[startIndex]) & ^startMask
		b.bits[startIndex] = upperBits | lowerBits

		// flip all bits at each of the middles indices
		for i := startIndex + 1; i < endIndex; i++ {
			b.bits[i] = ^b.bits[i]
		}

		if end != b.Size() {
			// flip lower bits, keep upper bits the same
			lowerBits = (^b.bits[endIndex]) & ^endMask
			upperBits = b.bits[endIndex] & endMask
			b.bits[endIndex] = upperBits | lowerBits
		}
	}
}

// TODO: expose this as a public function once its ready
func fromBytes(bytes []byte) *BitSet {
	b := New(NumBits(len(bytes) * 8))
	k := 0
	for i := 0; i < len(bytes); i += 8 {
		word := uint64(0)
		for j := 0; i+j < len(bytes) && j < 8; j++ {
			b := uint64(bytes[i+j])
			bShift := b << (8 * j)
			word |= bShift
		}
		b.bits[k] = word
		k++
	}
	return b
}

// TODO: expose this as a public function once its ready
func (b *BitSet) toBytes() []byte {
	bytes := make([]byte, len(b.bits)*8)
	k := 0
	for i := 0; i < len(b.bits); i++ {
		for j := 0; j < 8; j++ {
			bytes[k] = byte(0xFF & (b.bits[i] >> (j * 8)))
			k++
		}
	}
	return bytes
}

// String returns a hexadecimal representation of the bits in this BitSet
func (b *BitSet) String() string {
	s := make([]string, len(b.bits))
	for i := 0; i < len(s); i++ {
		s[i] = fmt.Sprintf(wordFmt, b.bits[len(b.bits)-1-i])
	}
	return strings.Join(s, "")
}

// Iterator is an alias for SetBitIterator
func (b *BitSet) Iterator() iterator.ForwardIterator[int] {
	return b.SetBitIterator()
}

func (b *BitSet) Elements() chan *int {
	return iterator.Range[int](b)
}

// SetBitIterator returns an iterator that iterates over the set bits of this BitSet
func (b *BitSet) SetBitIterator() iterator.ForwardIterator[int] {
	bi := &setBitIterator{
		bitSet: b,
	}
	bi.bitIndex = bi.getNextSetIndex(0)
	return bi
}

// UnsetBitIterator returns an iterator that iterates over the unset bits of this BitSet
func (b *BitSet) UnsetBitIterator() iterator.ForwardIterator[int] {
	bi := &unSetBitIterator{
		bitSet: b,
	}
	bi.bitIndex = bi.getNextUnSetIndex(0)
	return bi
}

func convert(bit int) (int, int) {
	ensureNonNegative(bit)
	return bit / wordSize, bit % wordSize
}

func defaultConfig() *Config {
	return &Config{
		NumBits: DefaultNumBits,
	}
}

func (b *BitSet) ensureSize(index int) {
	for index >= len(b.bits) {
		b.bits = append(b.bits, 0)
	}
}

// TODO: work on incorporating this
func (b *BitSet) lastNonZeroWord() int {
	for i := len(b.bits) - 1; i >= 0; i-- {
		if b.bits[i] != 0 {
			return i
		}
	}
	return -1
}

func ensureNonNegative(i int) {
	if i < 0 {
		panic(fmt.Sprintf("runtime error: index out of range [%d]", i))
	}
}

// Iterator stuff

func (bi *setBitIterator) Empty() bool {
	return bi.bitIndex >= len(bi.bitSet.bits)*wordSize
}

func (bi *setBitIterator) Increment() error {
	if bi.Empty() {
		return fmt.Errorf("cannot pop front of an empty iterator")
	}
	bi.bitIndex = bi.getNextSetIndex(bi.bitIndex + 1)
	return nil
}

func (bi *setBitIterator) GetFront() (*int, error) {
	if bi.Empty() {
		return nil, fmt.Errorf("cannot get front of an empty iterator")
	}
	v := bi.bitIndex
	return &v, nil
}

func (bi *setBitIterator) MustIncrement() {
	util.MustDo(bi.Increment())
}

func (bi *setBitIterator) MustGetFront() *int {
	return util.MustGet(bi.GetFront())
}

func (bi *setBitIterator) getNextSetIndex(start int) int {
	for start < bi.bitSet.Size() && !bi.bitSet.Get(start) {
		start++
	}
	return start
}

func (bi *unSetBitIterator) Empty() bool {
	return bi.bitIndex >= len(bi.bitSet.bits)*wordSize
}

func (bi *unSetBitIterator) Increment() error {
	if bi.Empty() {
		return fmt.Errorf("cannot pop front of an empty iterator")
	}
	bi.bitIndex = bi.getNextUnSetIndex(bi.bitIndex + 1)
	return nil
}

func (bi *unSetBitIterator) GetFront() (*int, error) {
	if bi.Empty() {
		return nil, fmt.Errorf("cannot get front of an empty iterator")
	}
	v := bi.bitIndex
	return &v, nil
}

func (bi *unSetBitIterator) getNextUnSetIndex(start int) int {
	for start < bi.bitSet.Size() && bi.bitSet.Get(start) {
		start++
	}
	return start
}

func (bi *unSetBitIterator) MustIncrement() {
	util.MustDo(bi.Increment())
}

func (bi *unSetBitIterator) MustGetFront() *int {
	return util.MustGet(bi.GetFront())
}
