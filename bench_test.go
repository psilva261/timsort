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

func LessThanByKey(a, b interface{}) bool {
	return a.(*record).key < b.(*record).key
}

type RecordSlice []record

func (s *RecordSlice) Len() int {
	return len(*s)
}

func (s *RecordSlice) Swap(i, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}

func (s *RecordSlice) Less(i, j int) bool {
	return (*s)[i].key < (*s)[j].key
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
	v = make(records, size)
	switch shape {

	case "xor":
		for i := 0; i < size; i++ {
			v[i] = &record{0xff & (i ^ 0xab), i}
		}

	case "sorted":
		for i := 0; i < size; i++ {
			v[i] = &record{i, i}
		}

	case "revsorted":
		for i := 0; i < size; i++ {
			v[i] = &record{size - i, i}
		}

	case "random":
		rand.Seed(1)

		for i := 0; i < size; i++ {
			v[i] = &record{rand.Int(), i}
		}

	default:
		panic(shape)
	}

	return v
}

func makeRecords(size int, shape string) (v RecordSlice) {
	v = make(RecordSlice, size)
	switch shape {

	case "xor":
		for i := 0; i < size; i++ {
			v[i] = record{0xff & (i ^ 0xab), i}
		}

	case "sorted":
		for i := 0; i < size; i++ {
			v[i] = record{i, i}
		}

	case "revsorted":
		for i := 0; i < size; i++ {
			v[i] = record{size - i, i}
		}

	case "random":
		rand.Seed(1)

		for i := 0; i < size; i++ {
			v[i] = record{rand.Int(), i}
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

func benchmarkTimsortInterface(b *testing.B, size int, shape string) {
	b.StopTimer()

	for j := 0; j < b.N; j++ {
		v := makeRecords(size, shape)

		b.StartTimer()
		TimSort(&v)
		b.StopTimer()
	}
}

func benchmarkStandardSort(b *testing.B, size int, shape string) {
	b.StopTimer()

	for j := 0; j < b.N; j++ {
		v := makeRecords(size, shape)

		b.StartTimer()
		sort.Sort(&v)
		b.StopTimer()
	}
}

func BenchmarkTimsortXor100(b *testing.B) {
	benchmarkTimsort(b, 100, "xor")
}

func BenchmarkTimsortInterXor100(b *testing.B) {
	benchmarkTimsortInterface(b, 100, "xor")
}

func BenchmarkStandardSortXor100(b *testing.B) {
	benchmarkStandardSort(b, 100, "xor")
}

func BenchmarkTimsortSorted100(b *testing.B) {
	benchmarkTimsort(b, 100, "sorted")
}

func BenchmarkTimsortInterSorted100(b *testing.B) {
	benchmarkTimsortInterface(b, 100, "sorted")
}

func BenchmarkStandardSortSorted100(b *testing.B) {
	benchmarkStandardSort(b, 100, "sorted")
}

func BenchmarkTimsortRevSorted100(b *testing.B) {
	benchmarkTimsort(b, 100, "revsorted")
}

func BenchmarkTimsortInterRevSorted100(b *testing.B) {
	benchmarkTimsortInterface(b, 100, "revsorted")
}

func BenchmarkStandardSortRevSorted100(b *testing.B) {
	benchmarkStandardSort(b, 100, "revsorted")
}

func BenchmarkTimsortRandom100(b *testing.B) {
	benchmarkTimsort(b, 100, "random")
}

func BenchmarkTimsortInterRandom100(b *testing.B) {
	benchmarkTimsortInterface(b, 100, "random")
}

func BenchmarkStandardSortRandom100(b *testing.B) {
	benchmarkStandardSort(b, 100, "random")
}

func BenchmarkTimsortXor1K(b *testing.B) {
	benchmarkTimsort(b, 1024, "xor")
}

func BenchmarkTimsortInterXor1K(b *testing.B) {
	benchmarkTimsortInterface(b, 1024, "xor")
}

func BenchmarkStandardSortXor1K(b *testing.B) {
	benchmarkStandardSort(b, 1024, "xor")
}

func BenchmarkTimsortSorted1K(b *testing.B) {
	benchmarkTimsort(b, 1024, "sorted")
}

func BenchmarkTimsortInterSorted1K(b *testing.B) {
	benchmarkTimsortInterface(b, 1024, "sorted")
}

func BenchmarkStandardSortSorted1K(b *testing.B) {
	benchmarkStandardSort(b, 1024, "sorted")
}

func BenchmarkTimsortRevSorted1K(b *testing.B) {
	benchmarkTimsort(b, 1024, "revsorted")
}

func BenchmarkTimsortInterRevSorted1K(b *testing.B) {
	benchmarkTimsortInterface(b, 1024, "revsorted")
}

func BenchmarkStandardSortRevSorted1K(b *testing.B) {
	benchmarkStandardSort(b, 1024, "revsorted")
}

func BenchmarkTimsortRandom1K(b *testing.B) {
	benchmarkTimsort(b, 1024, "random")
}

func BenchmarkTimsortInterRandom1K(b *testing.B) {
	benchmarkTimsortInterface(b, 1024, "random")
}

func BenchmarkStandardSortRandom1K(b *testing.B) {
	benchmarkStandardSort(b, 1024, "random")
}

func BenchmarkTimsortXor1M(b *testing.B) {
	benchmarkTimsort(b, 1024*1024, "xor")
}

func BenchmarkTimsortInterXor1M(b *testing.B) {
	benchmarkTimsortInterface(b, 1024*1024, "xor")
}

func BenchmarkStandardSortXor1M(b *testing.B) {
	benchmarkStandardSort(b, 1024*1024, "xor")
}

func BenchmarkTimsortSorted1M(b *testing.B) {
	benchmarkTimsort(b, 1024*1024, "sorted")
}

func BenchmarkTimsortInterSorted1M(b *testing.B) {
	benchmarkTimsortInterface(b, 1024*1024, "sorted")
}

func BenchmarkStandardSortSorted1M(b *testing.B) {
	benchmarkStandardSort(b, 1024*1024, "sorted")
}

func BenchmarkTimsortRevSorted1M(b *testing.B) {
	benchmarkTimsort(b, 1024*1024, "revsorted")
}

func BenchmarkTimsortInterRevSorted1M(b *testing.B) {
	benchmarkTimsortInterface(b, 1024*1024, "revsorted")
}

func BenchmarkStandardSortRevSorted1M(b *testing.B) {
	benchmarkStandardSort(b, 1024*1024, "revsorted")
}

func BenchmarkTimsortRandom1M(b *testing.B) {
	benchmarkTimsort(b, 1024*1024, "random")
}

func BenchmarkTimsortInterRandom1M(b *testing.B) {
	benchmarkTimsortInterface(b, 1024*1024, "random")
}

func BenchmarkStandardSortRandom1M(b *testing.B) {
	benchmarkStandardSort(b, 1024*1024, "random")
}
