package bitset

import (
	"fmt"
	"github.com/lock14/collections/iterator"
	"math/bits"
	"strings"
)

const (
	DefaultNumBits = 64
	wordSize       = 64
	wordFmt        = "%016X"
)

// BitSet represents a vector of bits that grows as needed.
type BitSet struct {
	bits         []uint64
	maxWordInUse int
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
		bits:         make([]uint64, (config.NumBits/wordSize)+min(1, config.NumBits%wordSize)),
		maxWordInUse: 0,
	}
}

// Clear sets the bit specified by the index to false.
func (b *BitSet) Clear(bit int) {
	index, shift := convert(bit)
	b.ensureSize(index)
	b.bits[index] &= ^(1 << shift)
	if index == b.maxWordInUse && b.bits[index] == 0 {
		b.maxWordInUse = b.lastNonZeroWord() + 1
	}
}

// Set sets the bit at the specified index to true.
func (b *BitSet) Set(bit int) {
	index, shift := convert(bit)
	b.ensureSize(index)
	b.bits[index] |= 1 << shift
	if index > b.maxWordInUse {
		b.maxWordInUse = index
	}
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

// Length returns the 'logical size' of this BitSet.
// The 'logical size' is the highest set bit in the BitSet plus one.
// Returns zero if no bits are set.
func (b *BitSet) Length() int {
	if b.maxWordInUse == 0 {
		return 0
	}
	return wordSize*(b.maxWordInUse-1) + (wordSize - bits.LeadingZeros64(b.bits[b.maxWordInUse-1]))
}

// Flip sets each bit to the complement of its current value. This call is
// equivalent to b.FlipRange(0, b.Size())
func (b *BitSet) Flip() {
	for i := 0; i < len(b.bits); i++ {
		b.bits[i] = ^b.bits[i]
	}
	b.maxWordInUse = b.lastNonZeroWord() + 1
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
	b.maxWordInUse = b.lastNonZeroWord() + 1
}

// FromBytes returns new BitSet containing all the bits in the given byte array.
func FromBytes(bytes []byte) *BitSet {
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
	b.maxWordInUse = b.lastNonZeroWord() + 1
	return b
}

// ToBytes returns a byte array containing all the set bits in this BitSet.
func (b *BitSet) ToBytes() []byte {
	n := b.maxWordInUse
	if n == 0 {
		return []byte{}
	}
	length := 8 * (n - 1)
	for word := b.bits[n-1]; word != 0; word >>= 8 {
		length++
	}
	bytes := make([]byte, length)
	k := 0
	for i := 0; i < n-1; i++ {
		for j := 0; j < 8; j++ {
			bytes[k] = byte(0xFF & (b.bits[i] >> (j * 8)))
			k++
		}
	}
	for word := b.bits[n-1]; word != 0; word >>= 8 {
		bytes[k] = byte(word & 0xFF)
		k++
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
func (b *BitSet) Iterator() iterator.Iterator[int] {
	return b.SetBitIterator()
}

func (b *BitSet) Elements() chan *int {
	return iterator.Elements(b.Iterator())
}

// SetBitIterator returns an iterator that iterates over the set bits of this BitSet
func (b *BitSet) SetBitIterator() iterator.Iterator[int] {
	bi := &setBitIterator{
		bitSet: b,
	}
	bi.bitIndex = bi.getNextSetIndex(0)
	return bi
}

// UnsetBitIterator returns an iterator that iterates over the unset bits of this BitSet
func (b *BitSet) UnsetBitIterator() iterator.Iterator[int] {
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

func (bi *setBitIterator) Next() (*int, error) {
	if bi.Empty() {
		return nil, fmt.Errorf("cannot get front of an empty iterator")
	}
	v := bi.bitIndex
	bi.bitIndex = bi.getNextSetIndex(bi.bitIndex + 1)
	return &v, nil
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

func (bi *unSetBitIterator) Next() (*int, error) {
	if bi.Empty() {
		return nil, fmt.Errorf("cannot get front of an empty iterator")
	}
	v := bi.bitIndex
	bi.bitIndex = bi.getNextUnSetIndex(bi.bitIndex + 1)
	return &v, nil
}

func (bi *unSetBitIterator) getNextUnSetIndex(start int) int {
	for start < bi.bitSet.Size() && bi.bitSet.Get(start) {
		start++
	}
	return start
}
