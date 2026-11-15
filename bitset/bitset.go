// Package bitset provides a vector of bits that grows as needed.
package bitset

import (
	"fmt"
	"iter"
	"math"
	"math/bits"
	"slices"
	"strconv"
	"strings"

	"github.com/lock14/collections"
)

var (
	_ collections.MutableNavigableSet[int] = (*BitSet)(nil)
	_ fmt.Stringer                         = (*BitSet)(nil)
)

const (
	// DefaultCapacity is the default initial bit capacity of a BitSet.
	DefaultCapacity = 64
	wordSize        = 64
)

// config holds configuration options for constructing a BitSet.
type config struct {
	capacity int
}

// Option configures a BitSet.
type Option func(*config)

// WithCapacity sets the initial bit capacity of a BitSet.
// Passing a negative capacity causes New to panic.
func WithCapacity(capacity int) Option {
	return func(c *config) {
		c.capacity = capacity
	}
}

// BitSet represents a dynamically growing vector of bits implementing MutableNavigableSet[int].
type BitSet struct {
	bits         []uint64
	maxWordInUse int
	size         int
}

// New creates a BitSet whose initial capacity is large enough to represent bits
// with indices in the range [0, capacity) without resizing. If no options are provided,
// DefaultCapacity is used. All bits are initially false. Panics if capacity is negative.
func New(opts ...Option) *BitSet {
	cfg := defaultConfig()
	for _, option := range opts {
		option(cfg)
	}
	if cfg.capacity < 0 {
		panic(fmt.Sprintf("bitset: capacity cannot be negative: %d", cfg.capacity))
	}
	return &BitSet{
		bits:         make([]uint64, (cfg.capacity/wordSize)+min(1, cfg.capacity%wordSize)),
		maxWordInUse: 0,
		size:         0,
	}
}

// FromBytes creates a new BitSet containing all set bits from data interpreted in little-endian byte order.
func FromBytes(data []byte) *BitSet {
	s := New(WithCapacity(len(data) * 8))
	k := 0
	for i := 0; i < len(data); i += 8 {
		var word uint64
		for j := 0; i+j < len(data) && j < 8; j++ {
			word |= uint64(data[i+j]) << (8 * j)
		}
		s.bits[k] = word
		k++
	}
	s.maxWordInUse = s.lastNonZeroWord() + 1
	s.recomputeSize()
	return s
}

// SetBit sets the bit at the specified non-negative index to true.
// Panics if bit is negative.
func (s *BitSet) SetBit(bit int) {
	index, shift := convert(bit)
	s.ensureSize(index)
	if s.bits[index]&(1<<shift) == 0 {
		if s.size != -1 {
			s.size++
		}
		s.bits[index] |= 1 << shift
		if index+1 > s.maxWordInUse {
			s.maxWordInUse = index + 1
		}
	}
}

// ClearBit sets the bit at the specified non-negative index to false.
// Panics if bit is negative.
func (s *BitSet) ClearBit(bit int) {
	index, shift := convert(bit)
	if index >= s.maxWordInUse {
		return
	}
	if s.bits[index]&(1<<shift) != 0 {
		if s.size != -1 {
			s.size--
		}
		s.bits[index] &= ^(1 << shift)
		if index == s.maxWordInUse-1 && s.bits[index] == 0 {
			s.maxWordInUse = s.lastNonZeroWord() + 1
			if s.maxWordInUse == 0 {
				s.size = 0
			}
		}
	}
}

// Flip sets each bit in [0, s.Capacity()) to the complement of its current value.
func (s *BitSet) Flip() {
	for i := 0; i < len(s.bits); i++ {
		s.bits[i] = ^s.bits[i]
	}
	s.invalidateSize()
}

// FlipRange sets each bit in the half-open range [start, end) to the complement of its current value.
// Panics if start is negative or if end < start.
func (s *BitSet) FlipRange(start, end int) {
	ensureNonNegative(start)
	if end < start {
		panic(fmt.Sprintf("bitset: start index %d greater than end index %d", start, end))
	}
	if start == end {
		return
	}
	startIndex, startShift := convert(start)
	lastIndex, lastShift := convert(end - 1)
	s.ensureSize(lastIndex)

	firstMask := ^uint64(0) << startShift
	lastMask := ^uint64(0) >> (wordSize - 1 - lastShift)
	if startIndex == lastIndex {
		s.bits[startIndex] ^= firstMask & lastMask
	} else {
		s.bits[startIndex] ^= firstMask
		for i := startIndex + 1; i < lastIndex; i++ {
			s.bits[i] = ^s.bits[i]
		}
		s.bits[lastIndex] ^= lastMask
	}
	s.invalidateSize()
}

// Add inserts the specified non-negative integer into the set.
// Panics if t is negative.
func (s *BitSet) Add(t int) {
	s.SetBit(t)
}

// AddAll inserts all elements from the given sequence into the set.
// Panics if any element is negative.
func (s *BitSet) AddAll(seq iter.Seq[int]) {
	for v := range seq {
		s.SetBit(v)
	}
}

// Remove removes and returns the smallest set bit in the set.
// Panics if the set is empty.
func (s *BitSet) Remove() int {
	return s.PollFirst()
}

// RemoveElement removes the specified non-negative integer from the set if present.
// Panics if t is negative.
func (s *BitSet) RemoveElement(t int) {
	s.ClearBit(t)
}

// RemoveAll removes all elements in col from this set.
func (s *BitSet) RemoveAll(col collections.Collection[int]) {
	if other, ok := col.(*BitSet); ok {
		for i := 0; i < s.maxWordInUse && i < other.maxWordInUse; i++ {
			s.bits[i] &= ^other.bits[i]
		}
		s.invalidateSize()
		return
	}
	for v := range col.All() {
		s.ClearBit(v)
	}
}

// RetainAll retains only the elements in this set that are also contained in col.
func (s *BitSet) RetainAll(col collections.Collection[int]) {
	if other, ok := col.(*BitSet); ok {
		for i := 0; i < s.maxWordInUse; i++ {
			if i < other.maxWordInUse {
				s.bits[i] &= other.bits[i]
			} else {
				s.bits[i] = 0
			}
		}
		s.invalidateSize()
		return
	}
	if set, ok := col.(collections.Set[int]); ok {
		for v := range s.All() {
			if !set.Contains(v) {
				s.ClearBit(v)
			}
		}
		return
	}

	temp := New(WithCapacity(s.Capacity()))
	for v := range col.All() {
		if s.Contains(v) {
			temp.SetBit(v)
		}
	}
	s.bits = temp.bits
	s.maxWordInUse = temp.maxWordInUse
	s.size = temp.size
}

// PollFirst removes and returns the smallest element in the set.
// Panics if the set is empty.
func (s *BitSet) PollFirst() int {
	if s.Empty() {
		panic("bitset: PollFirst called on empty set")
	}
	val := s.First()
	s.ClearBit(val)
	return val
}

// PollLast removes and returns the largest element in the set.
// Panics if the set is empty.
func (s *BitSet) PollLast() int {
	if s.Empty() {
		panic("bitset: PollLast called on empty set")
	}
	val := s.Last()
	s.ClearBit(val)
	return val
}

// Clear sets all bits in the set to false while retaining allocated capacity.
func (s *BitSet) Clear() {
	clear(s.bits)
	s.maxWordInUse = 0
	s.size = 0
}

// GetBit returns the value of the bit at the specified non-negative index.
// Panics if bit is negative.
func (s *BitSet) GetBit(bit int) bool {
	index, shift := convert(bit)
	if index >= s.maxWordInUse {
		return false
	}
	return (s.bits[index]>>shift)&1 == 1
}

// Contains returns true if this set contains the specified non-negative integer.
// Panics if t is negative.
func (s *BitSet) Contains(t int) bool {
	return s.GetBit(t)
}

// ContainsAll returns true if this set contains every element of col.
func (s *BitSet) ContainsAll(col collections.Collection[int]) bool {
	if other, ok := col.(*BitSet); ok {
		for i := 0; i < other.maxWordInUse; i++ {
			if i >= s.maxWordInUse {
				if other.bits[i] != 0 {
					return false
				}
			} else if (s.bits[i] & other.bits[i]) != other.bits[i] {
				return false
			}
		}
		return true
	}
	for v := range col.All() {
		if !s.Contains(v) {
			return false
		}
	}
	return true
}

// Capacity returns the number of bits this BitSet can hold in its active word slice without growing.
func (s *BitSet) Capacity() int {
	return len(s.bits) * wordSize
}

// Size returns the number of set bits (population count) in the BitSet.
func (s *BitSet) Size() int {
	if s.size == -1 {
		s.recomputeSize()
	}
	return s.size
}

// Length returns the logical length of the BitSet: the index of the highest set bit plus one,
// or zero if the BitSet is empty.
func (s *BitSet) Length() int {
	if s.maxWordInUse == 0 {
		return 0
	}
	return wordSize*(s.maxWordInUse-1) + (wordSize - bits.LeadingZeros64(s.bits[s.maxWordInUse-1]))
}

// Empty returns true if the BitSet contains no set bits.
func (s *BitSet) Empty() bool {
	return s.maxWordInUse == 0
}

// First returns the smallest element in the set.
// Panics if the set is empty.
func (s *BitSet) First() int {
	if s.Empty() {
		panic("bitset: First called on empty set")
	}
	i := 0
	for s.bits[i] == 0 {
		i++
	}
	return i*wordSize + bits.TrailingZeros64(s.bits[i])
}

// Last returns the largest element in the set.
// Panics if the set is empty.
func (s *BitSet) Last() int {
	if s.Empty() {
		panic("bitset: Last called on empty set")
	}
	i := s.maxWordInUse - 1
	return i*wordSize + (wordSize - 1 - bits.LeadingZeros64(s.bits[i]))
}

// Lower returns the greatest element strictly less than t, or (0, false) if no such element exists.
func (s *BitSet) Lower(t int) (int, bool) {
	if t <= 0 {
		return 0, false
	}
	return s.Floor(t - 1)
}

// Floor returns the greatest element less than or equal to t, or (0, false) if no such element exists.
func (s *BitSet) Floor(t int) (int, bool) {
	if t < 0 || s.Empty() {
		return 0, false
	}
	index, shift := convert(t)
	if index >= s.maxWordInUse {
		return s.Last(), true
	}

	mask := ^uint64(0) >> (wordSize - 1 - shift)
	if w := s.bits[index] & mask; w != 0 {
		return index*wordSize + (wordSize - 1 - bits.LeadingZeros64(w)), true
	}

	for i := index - 1; i >= 0; i-- {
		if s.bits[i] != 0 {
			return i*wordSize + (wordSize - 1 - bits.LeadingZeros64(s.bits[i])), true
		}
	}
	return 0, false
}

// Ceiling returns the least element greater than or equal to t, or (0, false) if no such element exists.
func (s *BitSet) Ceiling(t int) (int, bool) {
	if s.Empty() {
		return 0, false
	}
	if t <= 0 {
		return s.First(), true
	}
	index, shift := convert(t)
	if index >= s.maxWordInUse {
		return 0, false
	}

	if w := s.bits[index] & (^uint64(0) << shift); w != 0 {
		return index*wordSize + bits.TrailingZeros64(w), true
	}

	for i := index + 1; i < s.maxWordInUse; i++ {
		if s.bits[i] != 0 {
			return i*wordSize + bits.TrailingZeros64(s.bits[i]), true
		}
	}
	return 0, false
}

// Higher returns the least element strictly greater than t, or (0, false) if no such element exists.
func (s *BitSet) Higher(t int) (int, bool) {
	if t == math.MaxInt {
		return 0, false
	}
	if t < 0 {
		return s.Ceiling(0)
	}
	return s.Ceiling(t + 1)
}

// ToBytes returns a minimal little-endian byte slice representing all set bits in the BitSet.
func (s *BitSet) ToBytes() []byte {
	n := s.maxWordInUse
	if n == 0 {
		return []byte{}
	}
	length := 8 * (n - 1)
	for word := s.bits[n-1]; word != 0; word >>= 8 {
		length++
	}
	out := make([]byte, length)
	k := 0
	for i := 0; i < n-1; i++ {
		for j := 0; j < 8; j++ {
			out[k] = byte(0xFF & (s.bits[i] >> (j * 8)))
			k++
		}
	}
	for word := s.bits[n-1]; word != 0; word >>= 8 {
		out[k] = byte(word & 0xFF)
		k++
	}
	return out
}

// SetBits returns an iterator over all set bit indices in ascending order.
func (s *BitSet) SetBits() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; i < s.maxWordInUse; i++ {
			word := s.bits[i]
			for word != 0 {
				tz := bits.TrailingZeros64(word)
				if !yield(i*wordSize + tz) {
					return
				}
				word &= word - 1
			}
		}
	}
}

// All returns an iterator over all set bit indices in ascending order.
func (s *BitSet) All() iter.Seq[int] {
	return s.SetBits()
}

// Backward returns an iterator over all set bit indices in descending order.
func (s *BitSet) Backward() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := s.maxWordInUse - 1; i >= 0; i-- {
			w := s.bits[i]
			for w != 0 {
				bit := wordSize - 1 - bits.LeadingZeros64(w)
				if !yield(i*wordSize + bit) {
					return
				}
				w &= ^(1 << bit)
			}
		}
	}
}

// From returns an iterator over all elements greater than or equal to from in ascending order.
func (s *BitSet) From(from int) iter.Seq[int] {
	return func(yield func(int) bool) {
		if from < 0 {
			from = 0
		}
		index, shift := convert(from)
		if index >= s.maxWordInUse {
			return
		}

		w := s.bits[index] & (^uint64(0) << shift)
		for w != 0 {
			tz := bits.TrailingZeros64(w)
			if !yield(index*wordSize + tz) {
				return
			}
			w &= w - 1
		}

		for i := index + 1; i < s.maxWordInUse; i++ {
			w = s.bits[i]
			for w != 0 {
				tz := bits.TrailingZeros64(w)
				if !yield(i*wordSize + tz) {
					return
				}
				w &= w - 1
			}
		}
	}
}

// To returns an iterator over all elements strictly less than to in ascending order.
func (s *BitSet) To(to int) iter.Seq[int] {
	return func(yield func(int) bool) {
		if to <= 0 {
			return
		}
		lastIndex, lastShift := convert(to - 1)
		for i := 0; i < s.maxWordInUse && i <= lastIndex; i++ {
			w := s.bits[i]
			if i == lastIndex {
				w &= ^uint64(0) >> (wordSize - 1 - lastShift)
			}
			for w != 0 {
				tz := bits.TrailingZeros64(w)
				if !yield(i*wordSize + tz) {
					return
				}
				w &= w - 1
			}
		}
	}
}

// Between returns an iterator over all elements in the half-open range [from, to) in ascending order.
func (s *BitSet) Between(from, to int) iter.Seq[int] {
	return func(yield func(int) bool) {
		if from >= to || to <= 0 {
			return
		}
		if from < 0 {
			from = 0
		}
		index, shift := convert(from)
		if index >= s.maxWordInUse {
			return
		}
		lastIndex, lastShift := convert(to - 1)

		for i := index; i < s.maxWordInUse && i <= lastIndex; i++ {
			w := s.bits[i]
			if i == index {
				w &= ^uint64(0) << shift
			}
			if i == lastIndex {
				w &= ^uint64(0) >> (wordSize - 1 - lastShift)
			}
			for w != 0 {
				tz := bits.TrailingZeros64(w)
				if !yield(i*wordSize + tz) {
					return
				}
				w &= w - 1
			}
		}
	}
}

// String returns a string representation of the BitSet matching Go slice formatting.
func (s *BitSet) String() string {
	if s.Empty() {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	first := true
	var buf [20]byte
	for bit := range s.All() {
		if !first {
			sb.WriteByte(' ')
		}
		first = false
		sb.Write(strconv.AppendInt(buf[:0], int64(bit), 10))
	}
	sb.WriteByte(']')
	return sb.String()
}

func defaultConfig() *config {
	return &config{
		capacity: DefaultCapacity,
	}
}

func ensureNonNegative(i int) {
	if i < 0 {
		panic(fmt.Sprintf("runtime error: index out of range [%d]", i))
	}
}

func convert(bit int) (int, int) {
	ensureNonNegative(bit)
	return bit / wordSize, bit % wordSize
}

func (s *BitSet) ensureSize(index int) {
	if index >= len(s.bits) {
		oldLen := len(s.bits)
		s.bits = slices.Grow(s.bits, index+1-oldLen)
		s.bits = s.bits[:index+1]
		clear(s.bits[oldLen:])
	}
}

func (s *BitSet) lastNonZeroWord() int {
	for i := len(s.bits) - 1; i >= 0; i-- {
		if s.bits[i] != 0 {
			return i
		}
	}
	return -1
}

func (s *BitSet) invalidateSize() {
	s.maxWordInUse = s.lastNonZeroWord() + 1
	if s.maxWordInUse == 0 {
		s.size = 0
	} else {
		s.size = -1
	}
}

func (s *BitSet) recomputeSize() {
	size := 0
	for i := 0; i < s.maxWordInUse; i++ {
		size += bits.OnesCount64(s.bits[i])
	}
	s.size = size
}
