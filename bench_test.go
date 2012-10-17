package timsort

import (
	"math/rand"
	"sort"
	"testing"
)

type record struct {
	key, order int
}

type records []interface{}

func (p records) Len() int {
	return len(p)
}

func (p records) Less(i, j int) bool {
	return p[i].(*record).Less(p[j])
}

func (p records) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (self *record) Less(o interface{}) bool {
	return self.key < o.(*record).key
}

func LessThanByKey(a, b interface{}) bool {
	return a.(*record).key < b.(*record).key
}

func LessThanByKeyByOrder(a, b interface{}) bool {
	aa := a.(*record)
	bb := b.(*record)

	if aa.key < bb.key {
		return true
	} else if aa.key == bb.key {
		return aa.order < bb.order
	}

	return false
}

func makeVector(size int, shape string) (v records) {
	switch shape {

	case "xor":
		for i := 0; i < size; i++ {
			v = append(v, &record{0xff & (i ^ 0xab), i})
		}

	case "sorted":
		for i := 0; i < size; i++ {
			v = append(v, &record{i, i})
		}

	case "revsorted":
		for i := 0; i < size; i++ {
			v = append(v, &record{size - i, i})
		}

	case "random":
		rand.Seed(1)

		for i := 0; i < size; i++ {
			v = append(v, &record{rand.Int(), i})
		}

	default:
		panic(shape)
	}

	return v
}

func benchmarkTimsort(b *testing.B, size int, shape string) {
	b.StopTimer()

	for j := 0; j < b.N; j++ {
		v := makeVector(size, shape)

		b.StartTimer()
		Sort(v, LessThanByKey)
		b.StopTimer()
	}
}

func benchmarkStandardSort(b *testing.B, size int, shape string) {
	b.StopTimer()

	for j := 0; j < b.N; j++ {
		v := makeVector(size, shape)

		b.StartTimer()
		sort.Sort(&v)
		b.StopTimer()
	}
}

func BenchmarkTimsortXor100(b *testing.B) {
	benchmarkTimsort(b, 100, "xor")
}

func BenchmarkStandardSortXor100(b *testing.B) {
	benchmarkStandardSort(b, 100, "xor")
}

func BenchmarkTimsortSorted100(b *testing.B) {
	benchmarkTimsort(b, 100, "sorted")
}

func BenchmarkStandardSortSorted100(b *testing.B) {
	benchmarkStandardSort(b, 100, "sorted")
}

func BenchmarkTimsortRevSorted100(b *testing.B) {
	benchmarkTimsort(b, 100, "revsorted")
}

func BenchmarkStandardSortRevSorted100(b *testing.B) {
	benchmarkStandardSort(b, 100, "revsorted")
}

func BenchmarkTimsortRandom100(b *testing.B) {
	benchmarkTimsort(b, 100, "random")
}

func BenchmarkStandardSortRandom100(b *testing.B) {
	benchmarkStandardSort(b, 100, "random")
}

func BenchmarkTimsortXor1K(b *testing.B) {
	benchmarkTimsort(b, 1024, "xor")
}

func BenchmarkStandardSortXor1K(b *testing.B) {
	benchmarkStandardSort(b, 1024, "xor")
}

func BenchmarkTimsortSorted1K(b *testing.B) {
	benchmarkTimsort(b, 1024, "sorted")
}

func BenchmarkStandardSortSorted1K(b *testing.B) {
	benchmarkStandardSort(b, 1024, "sorted")
}

func BenchmarkTimsortRevSorted1K(b *testing.B) {
	benchmarkTimsort(b, 1024, "revsorted")
}

func BenchmarkStandardSortRevSorted1K(b *testing.B) {
	benchmarkStandardSort(b, 1024, "revsorted")
}

func BenchmarkTimsortRandom1K(b *testing.B) {
	benchmarkTimsort(b, 1024, "random")
}

func BenchmarkStandardSortRandom1K(b *testing.B) {
	benchmarkStandardSort(b, 1024, "random")
}

func BenchmarkTimsortXor1M(b *testing.B) {
	benchmarkTimsort(b, 1024*1024, "xor")
}

func BenchmarkStandardSortXor1M(b *testing.B) {
	benchmarkStandardSort(b, 1024*1024, "xor")
}

func BenchmarkTimsortSorted1M(b *testing.B) {
	benchmarkTimsort(b, 1024*1024, "sorted")
}

func BenchmarkStandardSortSorted1M(b *testing.B) {
	benchmarkStandardSort(b, 1024*1024, "sorted")
}

func BenchmarkTimsortRevSorted1M(b *testing.B) {
	benchmarkTimsort(b, 1024*1024, "revsorted")
}

func BenchmarkStandardSortRevSorted1M(b *testing.B) {
	benchmarkStandardSort(b, 1024*1024, "revsorted")
}

func BenchmarkTimsortRandom1M(b *testing.B) {
	benchmarkTimsort(b, 1024*1024, "random")
}

func BenchmarkStandardSortRandom1M(b *testing.B) {
	benchmarkStandardSort(b, 1024*1024, "random")
}
