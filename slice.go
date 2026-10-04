package collections

import (
	"iter"
)

type Slice[E any] []E

type SliceTF[I, O any] = func(I) O
type SliceTFE[I, O any] = func(I) (O, error)

type Result[E any] struct {
	Val E
	Err error
}

// type ErrSeq[O any] = iter.Seq2[O, error]

func ToSlice[S ~[]E, E any](s S) Slice[E] {
	return Slice[E](s)
}

func (s Slice[E]) Map[O any](tf SliceTF[E, O]) iter.Seq[O] {
	return func(yield func(O) bool) {
		for _, k := range s {
			if !yield(tf(k)) {
				return
			}
		}
	}
}

func (s Slice[E]) MapErr[O any](tf SliceTFE[E, O]) ErrSeq[O] {
	return func(yield func(O, error) bool) {
		for _, k := range s {
			if !yield(tf(k)) {
				return
			}
		}
	}
}

func (s Slice[E]) MapToSlice[O any](tf SliceTF[E, O]) (res Slice[O]) {
	res = make(Slice[O], len(s))
	for i, v := range s {
		res[i] = tf(v)
	}
	return
}

func (s Slice[E]) MapToSliceErr[O any](tf SliceTFE[E, O]) (res Slice[O], err error) {
	res = make(Slice[O], len(s))
	for i, v := range s {
		res[i], err = tf(v)
		if err != nil {
			return
		}
	}
	return
}
