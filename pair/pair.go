package pair

import "fmt"

type Pair[T1 any, T2 any] struct {
	fst T1
	snd T2
}

func New[T1 any, T2 any](t1 T1, t2 T2) *Pair[T1, T2] {
	return &Pair[T1, T2]{
		fst: t1,
		snd: t2,
	}
}

func (p *Pair[T1, T2]) Fst() *T1 {
	return &p.fst
}

func (p *Pair[T1, T2]) Snd() *T2 {
	return &p.snd
}

func (p *Pair[T1, T2]) UnWrap() (T1, T2) {
	return p.fst, p.snd
}

func (p *Pair[T1, T2]) String() string {
	return fmt.Sprintf("(%v, %v)", p.fst, p.snd)
}
