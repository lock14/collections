// Package bitset provides a vector of bits that grows as needed.
package bitset

import (
	"fmt"
	"iter"
	"math/bits"
	"strings"

	"github.com/lock14/collections"
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
	size         int
}

var _ collections.MutableSet[int] = (*BitSet)(nil)

// Config holds the values for configuring a BitSet.
type Config struct {
	numBits int
}

// Option configures a BitSet config
type Option func(*Config)

// NumBits provides the option to set the number of bits used in a BitSet.
func NumBits(n int) Option {
	return func(c *Config) {
		c.numBits = n
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
	ensureNonNegative(config.numBits)
	return &BitSet{
		bits:         make([]uint64, (config.numBits/wordSize)+min(1, config.numBits%wordSize)),
		maxWordInUse: 0,
		size:         0,
	}
}

// ClearBit sets the bit specified by the index to false.
func (b *BitSet) ClearBit(bit int) {
	index, shift := convert(bit)
	if index >= b.maxWordInUse {
		return
	}
	if b.bits[index]&(1<<shift) != 0 {
		if b.size != -1 {
			b.size--
		}
		b.bits[index] &= ^(1 << shift)
		if index == b.maxWordInUse-1 && b.bits[index] == 0 {
			b.maxWordInUse = b.lastNonZeroWord() + 1
		}
	}
}

// Set sets the bit at the specified index to true.
func (b *BitSet) SetBit(bit int) {
	index, shift := convert(bit)
	b.ensureSize(index)
	if b.bits[index]&(1<<shift) == 0 {
		if b.size != -1 {
			b.size++
		}
		b.bits[index] |= 1 << shift
		if index+1 > b.maxWordInUse {
			b.maxWordInUse = index + 1
		}
	}
}

// Get returns the value of the bit with the specified index.
func (b *BitSet) GetBit(bit int) bool {
	index, shift := convert(bit)
	if index >= b.maxWordInUse {
		return false
	}
	return (b.bits[index]>>shift)&1 == 1
}

// Capacity returns the maximum number of bits this bit set can currently hold without resizing.
func (b *BitSet) Capacity() int {
	return len(b.bits) * wordSize
}

// Size returns the number of set bits in the BitSet.
func (b *BitSet) Size() int {
	if b.size == -1 {
		b.recomputeSize()
	}
	return b.size
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
// equivalent to b.FlipRange(0, b.Capacity())
func (b *BitSet) Flip() {
	for i := 0; i < len(b.bits); i++ {
		b.bits[i] = ^b.bits[i]
	}
	b.maxWordInUse = b.lastNonZeroWord() + 1
	b.size = -1
}

// FlipRange sets each bit from the specified start bit (inclusive) to the
// specified end bit (exclusive) to the complement of its current value.
func (b *BitSet) FlipRange(start int, end int) {
	startIndex, startShift := convert(start)
	endIndex, endShift := convert(end)
	if end != b.Capacity() {
		b.ensureSize(endIndex)
	}

	startMask := ^(^uint64(0) << startShift)
	endMask := ^uint64(0) << endShift

	if startIndex == endIndex {
		// flip middle bits, keep upper and lower bits the same
		middleMask := ^(startMask | endMask)
		oldWord := b.bits[startIndex]
		lowerBits := oldWord & startMask
		middleBits := (^oldWord) & middleMask
		upperBits := oldWord & endMask
		b.bits[startIndex] = lowerBits | middleBits | upperBits
	} else {
		// flip upper bits, keep lower bits the same
		oldStart := b.bits[startIndex]
		lowerBits := oldStart & startMask
		upperBits := (^oldStart) & ^startMask
		b.bits[startIndex] = upperBits | lowerBits

		// flip all bits at each of the middles indices
		for i := startIndex + 1; i < endIndex; i++ {
			b.bits[i] = ^b.bits[i]
		}

		if end != b.Capacity() {
			// flip lower bits, keep upper bits the same
			oldEnd := b.bits[endIndex]
			lowerBits = (^oldEnd) & ^endMask
			upperBits = oldEnd & endMask
			b.bits[endIndex] = upperBits | lowerBits
		}
	}
	b.maxWordInUse = b.lastNonZeroWord() + 1
	b.size = -1
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
	b.recomputeSize()
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

// SetBits returns an iterator that iterates over the set bits of this BitSet.
// It uses word-level iteration with bits.TrailingZeros64 for efficiency.
func (b *BitSet) SetBits() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; i < b.maxWordInUse; i++ {
			word := b.bits[i]
			for word != 0 {
				tz := bits.TrailingZeros64(word)
				if !yield(i*wordSize + tz) {
					return
				}
				word &= word - 1 // clear lowest set bit
			}
		}
	}
}

func convert(bit int) (int, int) {
	ensureNonNegative(bit)
	return bit / wordSize, bit % wordSize
}

func defaultConfig() *Config {
	return &Config{
		numBits: DefaultNumBits,
	}
}

func (b *BitSet) ensureSize(index int) {
	if index >= len(b.bits) {
		newBits := make([]uint64, index+1)
		copy(newBits, b.bits)
		b.bits = newBits
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

func (b *BitSet) recomputeSize() {
	size := 0
	for i := 0; i < b.maxWordInUse; i++ {
		size += bits.OnesCount64(b.bits[i])
	}
	b.size = size
}

// Add inserts the specified element into the bit set.
func (b *BitSet) Add(t int) {
	b.SetBit(t)
}

// RemoveElement removes the specified element from the bit set.
func (b *BitSet) RemoveElement(t int) {
	b.ClearBit(t)
}

// Contains returns true if this bit set contains the specified element.
func (b *BitSet) Contains(t int) bool {
	return b.GetBit(t)
}

// Empty returns true if the collection contains no elements.
func (b *BitSet) Empty() bool {
	return b.size == 0
}

// All returns an iterator over all the elements in the bit set.
func (b *BitSet) All() iter.Seq[int] {
	return b.SetBits()
}

// Clear removes all elements from the bit set.
func (b *BitSet) Clear() {
	b.bits = make([]uint64, len(b.bits))
	b.maxWordInUse = 0
	b.size = 0
}

// Remove removes and returns a single element from the bit set.
func (b *BitSet) Remove() int {
	if b.size == 0 {
		panic("remove from empty set")
	}
	for i := 0; i < b.maxWordInUse; i++ {
		if b.bits[i] != 0 {
			tz := bits.TrailingZeros64(b.bits[i])
			val := i*wordSize + tz
			b.ClearBit(val)
			return val
		}
	}
	panic("size is not zero but no bits set")
}

// AddAll inserts all elements from the given sequence into the collection.
func (b *BitSet) AddAll(seq iter.Seq[int]) {
	for v := range seq {
		b.SetBit(v)
	}
}

// RemoveAll removes all elements of the specified collection from this set.
func (b *BitSet) RemoveAll(col collections.Collection[int]) {
	if other, ok := col.(*BitSet); ok {
		for i := 0; i < b.maxWordInUse && i < other.maxWordInUse; i++ {
			b.bits[i] &= ^other.bits[i]
		}
		b.maxWordInUse = b.lastNonZeroWord() + 1
		b.size = -1
		return
	}
	for v := range col.All() {
		b.ClearBit(v)
	}
}

// RetainAll retains only the elements in this set that are contained in the specified collection.
func (b *BitSet) RetainAll(col collections.Collection[int]) {
	if other, ok := col.(*BitSet); ok {
		for i := 0; i < b.maxWordInUse; i++ {
			if i < other.maxWordInUse {
				b.bits[i] &= other.bits[i]
			} else {
				b.bits[i] = 0
			}
		}
		b.maxWordInUse = b.lastNonZeroWord() + 1
		b.size = -1
		return
	}
	if set, ok := col.(collections.Set[int]); ok {
		for v := range b.All() {
			if !set.Contains(v) {
				b.ClearBit(v)
			}
		}
		return
	}

	// For generic non-set collections, build a temporary BitSet to avoid O(N*M).
	temp := New(NumBits(b.Capacity()))
	for v := range col.All() {
		if b.Contains(v) {
			temp.SetBit(v)
		}
	}
	b.bits = temp.bits
	b.maxWordInUse = temp.maxWordInUse
	b.size = temp.size
}

// ContainsAll returns true if this set contains all elements of the specified collection.
func (b *BitSet) ContainsAll(col collections.Collection[int]) bool {
	if other, ok := col.(*BitSet); ok {
		for i := 0; i < other.maxWordInUse; i++ {
			if i >= b.maxWordInUse {
				if other.bits[i] != 0 {
					return false
				}
			} else if (b.bits[i] & other.bits[i]) != other.bits[i] {
				return false
			}
		}
		return true
	}
	for v := range col.All() {
		if !b.Contains(v) {
			return false
		}
	}
	return true
}
