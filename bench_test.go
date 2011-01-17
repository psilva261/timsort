package timsort

import (
	"testing"
	"sort"
	"rand"
	//	"fmt"
	"container/vector"
)

type record struct {
	key, order int
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

func makeVector(size int, shape string) vector.Vector {

	var v vector.Vector

	switch shape {

	case "xor":
		for i := 0; i < size; i++ {
			v.Push(&record{0xff & (i ^ 0xab), i})
		}

	case "sorted":
		for i := 0; i < size; i++ {
			v.Push(&record{i, i})
		}

	case "revsorted":
		for i := 0; i < size; i++ {
			v.Push(&record{size - i, i})
		}

	case "random":
		rand.Seed(1)

		for i := 0; i < size; i++ {
			v.Push(&record{rand.Int(), i})
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
