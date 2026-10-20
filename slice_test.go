package collections

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"testing"

	"golang.org/x/exp/constraints"
)

func Test_Map(t *testing.T) {
	integers := make(Slice[int64], 0, 10)
	for i := int64(0); i < 10; i++ {
		integers = append(integers, i)
		t.Logf("adding %d\n", i)
	}
	tr := func(i int64) string {
		return strconv.FormatInt(i, 10)
	}
	stringArr := integers.MapToSlice(tr)
	for i, v := range stringArr {
		t.Logf("strings %d: %s\n", i, v)
	}
}

func Test_Iter(t *testing.T) {
	integers := make([]int, 0, 10)
	for i := range 10 {
		integers = append(integers, i)
		t.Logf("adding %d\n", i)
	}

	stringArr := SliceToSeq(integers).
		Filter(evenNumbers).
		Map(strconv.Itoa).
		ToSlice()

	for i, v := range stringArr {
		t.Logf("strings %d: %s\n", i, v)
	}
}

type iterBench struct {
	size      int
	name      string
	benchFunc func(b *testing.B)
}

type testF[S ~[]E, S2 ~[]E2, E, E2 any] = func(S) S2
type createF[S ~[]E, E any] = func(int) S
type correctF[S ~[]E, E any] = func(S, int) bool

func createIterBench[S ~[]E, S2 ~[]E2, E, E2 any](name string, test testF[S, S2, E, E2], create createF[S, E], correct correctF[S2, E2], sizes ...int) (r []*iterBench) {

	for _, size := range sizes {
		res := new(iterBench)
		r = append(r, res)
		res.name = fmt.Sprintf("%s_%d", name, size)
		res.size = size
		data := create(size)
		res.benchFunc = func(b *testing.B) {
			benchRes := test(data)
			b.ReportAllocs()
			b.ResetTimer()
			if !correct(benchRes, size) {
				b.Errorf("failed bench %s", res.name)
			}
		}
	}
	return r
}

func createSignedSlice[S constraints.Signed](size int) (s []S) {
	s = make([]S, size)
	for v := range size {
		s[v] = S(v)
	}
	return
}

func iterIntToStringFilter(ints []int) []string {
	return SliceToSeq(ints).Filter(evenNumbers).Map(strconv.Itoa).ToSlice()
}

func loopIntToStringFilter(ints []int) (s []string) {
	for _, v := range ints {
		if evenNumbers(v) {
			s = append(s, strconv.Itoa(v))
		}
	}
	return s
}

func PreLoopIntToStringFilter(ints []int) (s []string) {
	s = make([]string, 0, len(ints)/2)
	for _, v := range ints {
		if evenNumbers(v) {
			s = append(s, strconv.Itoa(v))
		}
	}
	return s
}

func checkStringRes(res []string, size int) bool {
	return len(res) == size/2 && res[0] == "0" && res[len(res)-1] == strconv.Itoa(size-2)
}

func createIntegerIterBench(name string, test testF[[]int, []string, int, string]) []*iterBench {
	sizes := []int{
		10,
		1_000,
		1_000_000,
		10_000_000,
		50_000_000,
	}
	return createIterBench(name, test, createSignedSlice, checkStringRes, sizes...)
}

func Benchmark_Iter(b *testing.B) {
	var benchmarks []*iterBench
	benchmarks = append(benchmarks, createIntegerIterBench("iterat", iterIntToStringFilter)...)
	benchmarks = append(benchmarks, createIntegerIterBench("normal", loopIntToStringFilter)...)
	benchmarks = append(benchmarks, createIntegerIterBench("palloc", PreLoopIntToStringFilter)...)

	slices.SortFunc(benchmarks, func(a, b *iterBench) int {
		if a.size == b.size {
			return cmp.Compare(a.name, b.name)
		}
		return a.size - b.size
	})

	for _, v := range benchmarks {
		b.Run(v.name, v.benchFunc)
	}
}

func evenNumbers[S constraints.Signed](s S) bool {
	return s%2 == 0
}
