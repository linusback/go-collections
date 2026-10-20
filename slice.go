package collections

import (
	"iter"
	"slices"
	"strconv"
)

type Seq[E any] iter.Seq[E]
type FilterF[I any] = func(I) bool
type MapF[I, O any] = func(I) O

func SliceToSeq[S ~[]E, E any](slice S) Seq[E] {
	return func(yield func(E) bool) {
		for _,e := range slice {
			if !yield(e) {
				return
			}
		}
	}
}

func (s Seq[E]) Filter(ff FilterF[E]) Seq[E] {
	return func(yield func(E) bool) {
		for k := range s {
			if !ff(k) {
				continue
			}
			if !yield(k) {
				return
			}
		}
	}
}

func (s Seq[E]) Map[O any](tf MapF[E, O]) Seq[O] {
	return func(yield func(O) bool) {
		for k := range s {
			if !yield(tf(k)) {
				return
			}
		}
	}
}

func (s Seq[E]) ToSlice() []E {
	return slices.Collect(iter.Seq[E](s))
}



func (s Seq[E]) Map2[O any](tf MapF[E, O]) Seq[O] {
	return func(yield func(O) bool) {
		for k := range s {
			if !yield(tf(k)) {
				return
			}
		}
	}
}





type Slice[E any] []E

func (s Slice[E]) MapToSlice[O any](tf MapF[E, O]) (result Slice[O]) {
	result = make(Slice[O], len(s))
	for i, v := range s {
		result[i] = tf(v)
	}
	return
}



func Filter[E any](s iter.Seq[E], ff FilterF[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		for k := range s {
			if !ff(k) {
				continue
			}
			if !yield(k) {
				return
			}
		}
	}
}

func Map[E, O any](s iter.Seq[E], tf MapF[E, O]) iter.Seq[O] {
	return func(yield func(O) bool) {
		for k := range s {
			if !yield(tf(k)) {
				return
			}
		}
	}
}

func isEven(i int) bool {
	return i % 2 == 0
}

func Normal(integers []int) (res []string) {
	for _, v := range integers {
		if v % 2 == 0 {
			res = append(res, strconv.Itoa(v))
		}
	}
	return res
}

func UsingIter(integers []int) (res []string) {
	return slices.Collect(Map(Filter(slices.Values(integers), isEven), strconv.Itoa)) 
}