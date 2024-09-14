package hashmap

import (
	"iter"
	"maps"
)

type MapWrapper[K comparable, V any] struct {
	m map[K]V
}

func Wrap[K comparable, V any](m map[K]V) *MapWrapper[K, V] {
	return &MapWrapper[K, V]{
		m: m,
	}
}

func (hm *MapWrapper[K, V]) Get(key K) (V, bool) {
	v, ok := hm.m[key]
	return v, ok
}

func (hm *MapWrapper[K, V]) Put(key K, value V) {
	hm.m[key] = value
}

func (hm *MapWrapper[K, V]) Remove(key K) {
	delete(hm.m, key)
}

func (hm *MapWrapper[K, V]) Size() int {
	return len(hm.m)
}

func (hm *MapWrapper[K, V]) Empty() bool {
	return hm.Size() == 0
}

func (hm *MapWrapper[K, V]) All() iter.Seq2[K, V] {
	return maps.All(hm.m)
}

func (hm *MapWrapper[K, V]) Keys() iter.Seq[K] {
	return maps.Keys(hm.m)
}

func (hm *MapWrapper[K, V]) Values() iter.Seq[V] {
	return maps.Values(hm.m)
}
