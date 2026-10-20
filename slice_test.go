package collections

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"
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

var sizes = []int{
	10,
	1_000,
	1_000_000,
	10_000_000,
	50_000_000,
}

var intCache = createIntCache()

func createIntCache() (res [][]int) {
	for _, size := range sizes {
		res = append(res, createSignedSlice[int](size))
	}
	return res
}

var strCache = createStrCache()

func createStrCache() (res [][]string) {
	for _, size := range sizes {
		res = append(res, createStringSlice(size))
	}
	return res
}

func createSignedSliceCached(size int) (s []int) {
	if index := slices.Index(sizes, size); index >= 0 {
		return intCache[index]
	}
	return createSignedSlice[int](size)
}

func createSignedSlice[S constraints.Signed](size int) (s []S) {
	s = make([]S, size)
	for v := range size {
		s[v] = S(v)
	}
	return
}

func createStringSliceCached(size int) (s []string) {
	if index := slices.Index(sizes, size); index >= 0 {
		return strCache[index]
	}
	return createStringSlice(size)
}

func createStringSlice(size int) (s []string) {
	s = make([]string, size)
	for v := range size {
		s[v] = strconv.Itoa(v)
	}
	return
}

type iterBench struct {
	size      int
	name      string
	benchFunc func(b *testing.B)
}

type testF[S ~[]E, S2 ~[]E2, E, E2 any] = func(S) S2
type createF[S ~[]E, E any] = func(int) S
type correctF[S ~[]E, E any] = func(S) bool
type correctCreateF[S ~[]E, E any] = func(int) correctF[S, E]

func createIterBench[S ~[]E, S2 ~[]E2, E, E2 any](name string, test testF[S, S2, E, E2], create createF[S, E], correctCreate correctCreateF[S2, E2], sizes ...int) (r []*iterBench) {
	maxSize := len(strconv.Itoa(slices.Max(sizes)))
	for _, size := range sizes {
		res := new(iterBench)
		r = append(r, res)
		res.name = fmt.Sprintf("%s%s", name, sixePadding(size, maxSize))
		res.size = size
		data := create(size)
		correct := correctCreate(size)
		res.benchFunc = func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			benchRes := test(data)
			if !correct(benchRes) {
				b.Errorf("failed bench %s", res.name)
			}
		}
	}
	return r
}

func sixePadding(size, maxSize int) string {
	sizeStr := strconv.Itoa(size)
	padding := maxSize - len(sizeStr)
	if padding == 0 {
		return "_" + sizeStr
	}
	return strings.Repeat("_", padding+1) + sizeStr
}

func iterIntToStringFilter(ints []int) []string {
	return SliceToSeq(ints).Filter(evenNumbers).Map(strconv.Itoa).ToSlice()
}

func iterIntToStringFilterAlloc(ints []int) []string {
	return SliceToLenSeq(ints).FilterEstimate(evenNumbers, len(ints)/2).Map(strconv.Itoa).ToSlice()
}

func loopIntToStringFilter(ints []int) (s []string) {
	for _, v := range ints {
		if evenNumbers(v) {
			s = append(s, strconv.Itoa(v))
		}
	}
	return s
}

func toIntNoErr(s string) (i int) {
	i, _ = strconv.Atoi(s)
	return i
}

func iterStringToIntFilter(strs []string) []int {
	return SliceToSeq(strs).Map(toIntNoErr).Filter(evenNumbers).ToSlice()
}

func iterStringToIntFilterAlloc(strs []string) []int {
	return SliceToLenSeq(strs).Map(toIntNoErr).FilterEstimate(evenNumbers, len(strs)/2).ToSlice()
}

func loopStringToIntFilter(strs []string) (s []int) {
	for _, v := range strs {
		e := toIntNoErr(v)
		if evenNumbers(e) {
			s = append(s, e)
		}
	}
	return s
}

func PreloopStringToIntFilter(strs []string) (s []int) {
	s = make([]int, 0, len(strs)/2)
	for _, v := range strs {
		e := toIntNoErr(v)
		if evenNumbers(e) {
			s = append(s, e)
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

func checkStringRes(size int) correctF[[]string, string] {
	const expectedFirst = "0"
	expectedLenght := size / 2
	lastElement := expectedLenght - 1
	expectedLast := strconv.Itoa(size - 2)
	return func(res []string) bool {
		lenght := len(res)
		return lenght == expectedLenght && res[0] == expectedFirst && res[lastElement] == expectedLast
	}
}

func checkIntRes(size int) correctF[[]int, int] {
	const expectedFirst = 0
	expectedLenght := size / 2
	lastElement := expectedLenght - 1
	expectedLast := size - 2
	return func(res []int) bool {
		lenght := len(res)
		return lenght == expectedLenght && res[0] == expectedFirst && res[lastElement] == expectedLast
	}
}

func createIntegerIterBench(name string, test testF[[]int, []string, int, string]) []*iterBench {
	return createIterBench(name, test, createSignedSliceCached, checkStringRes, sizes...)
}

func createStringIterBench(name string, test testF[[]string, []int, string, int]) []*iterBench {
	return createIterBench(name, test, createStringSliceCached, checkIntRes, sizes...)
}

func Benchmark_Iter_ItoA(b *testing.B) {
	var benchmarks []*iterBench
	benchmarks = append(benchmarks, createIntegerIterBench("iterat", iterIntToStringFilter)...)
	benchmarks = append(benchmarks, createIntegerIterBench("itallo", iterIntToStringFilterAlloc)...)
	benchmarks = append(benchmarks, createIntegerIterBench("normal", loopIntToStringFilter)...)
	benchmarks = append(benchmarks, createIntegerIterBench("palloc", PreLoopIntToStringFilter)...)
	b.Logf("size of cache %d, (largest is %d)\n", len(intCache), len(intCache[(len(intCache)-1)]))
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

func Benchmark_Iter_AtoI(b *testing.B) {
	var benchmarks []*iterBench
	benchmarks = append(benchmarks, createStringIterBench("iterat", iterStringToIntFilter)...)
	benchmarks = append(benchmarks, createStringIterBench("itallo", iterStringToIntFilterAlloc)...)
	benchmarks = append(benchmarks, createStringIterBench("normal", loopStringToIntFilter)...)
	benchmarks = append(benchmarks, createStringIterBench("palloc", PreloopStringToIntFilter)...)
	b.Logf("size of cache %d, (largest is %d)\n", len(strCache), len(strCache[(len(strCache)-1)]))
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
