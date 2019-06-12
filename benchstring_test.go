package timsort

import (
	"math/rand"
	"sort"
	"strconv"
	"testing"
)

func makeStrings(size int, shape string) (v sort.StringSlice) {
	v = make(sort.StringSlice, 0, size)
	switch shape {

	case "xor":
		for i := 0; i < size; i++ {
			v = append(v, strconv.Itoa(0xff&(i^0xab)))
		}

	case "sorted":
		for i := 0; i < size; i++ {
			v = append(v, strconv.Itoa(i))
		}

	case "revsorted":
		for i := 0; i < size; i++ {
			v = append(v, strconv.Itoa(size-i))
		}

	case "random":
		rand.Seed(1)

		for i := 0; i < size; i++ {
			v = append(v, strconv.Itoa(rand.Int()))
		}

	default:
		panic(shape)
	}

	return v
}

func benchmarkTimsortStr(b *testing.B, size int, shape string) {
	b.StopTimer()

	for j := 0; j < b.N; j++ {
		v := makeStrings(size, shape)

		b.StartTimer()
		err := TimSort(v)
		b.StopTimer()
		if err != nil {
			b.Fatalf("TimSort: %v", err)
		}
	}
}

func benchmarkStandardSortStr(b *testing.B, size int, shape string) {
	b.StopTimer()

	for j := 0; j < b.N; j++ {
		v := makeStrings(size, shape)

		b.StartTimer()
		sort.Sort(&v)
		b.StopTimer()
	}
}

func BenchmarkTimsortStrXor100(b *testing.B) {
	benchmarkTimsortStr(b, 100, "xor")
}

func BenchmarkStandardSortStrXor100(b *testing.B) {
	benchmarkStandardSortStr(b, 100, "xor")
}

func BenchmarkTimsortStrSorted100(b *testing.B) {
	benchmarkTimsortStr(b, 100, "sorted")
}

func BenchmarkStandardSortStrSorted100(b *testing.B) {
	benchmarkStandardSortStr(b, 100, "sorted")
}

func BenchmarkTimsortStrRevSorted100(b *testing.B) {
	benchmarkTimsortStr(b, 100, "revsorted")
}

func BenchmarkStandardSortStrRevSorted100(b *testing.B) {
	benchmarkStandardSortStr(b, 100, "revsorted")
}

func BenchmarkTimsortStrRandom100(b *testing.B) {
	benchmarkTimsortStr(b, 100, "random")
}

func BenchmarkStandardSortStrRandom100(b *testing.B) {
	benchmarkStandardSortStr(b, 100, "random")
}

func BenchmarkTimsortStrXor1K(b *testing.B) {
	benchmarkTimsortStr(b, 1024, "xor")
}

func BenchmarkStandardSortStrXor1K(b *testing.B) {
	benchmarkStandardSortStr(b, 1024, "xor")
}

func BenchmarkTimsortStrSorted1K(b *testing.B) {
	benchmarkTimsortStr(b, 1024, "sorted")
}

func BenchmarkStandardSortStrSorted1K(b *testing.B) {
	benchmarkStandardSortStr(b, 1024, "sorted")
}

func BenchmarkTimsortStrRevSorted1K(b *testing.B) {
	benchmarkTimsortStr(b, 1024, "revsorted")
}

func BenchmarkStandardSortStrRevSorted1K(b *testing.B) {
	benchmarkStandardSortStr(b, 1024, "revsorted")
}

func BenchmarkTimsortStrRandom1K(b *testing.B) {
	benchmarkTimsortStr(b, 1024, "random")
}

func BenchmarkStandardSortStrRandom1K(b *testing.B) {
	benchmarkStandardSortStr(b, 1024, "random")
}

func BenchmarkTimsortStrXor1M(b *testing.B) {
	benchmarkTimsortStr(b, 1024*1024, "xor")
}

func BenchmarkStandardSortStrXor1M(b *testing.B) {
	benchmarkStandardSortStr(b, 1024*1024, "xor")
}

func BenchmarkTimsortStrSorted1M(b *testing.B) {
	benchmarkTimsortStr(b, 1024*1024, "sorted")
}

func BenchmarkStandardSortStrSorted1M(b *testing.B) {
	benchmarkStandardSortStr(b, 1024*1024, "sorted")
}

func BenchmarkTimsortStrRevSorted1M(b *testing.B) {
	benchmarkTimsortStr(b, 1024*1024, "revsorted")
}

func BenchmarkStandardSortStrRevSorted1M(b *testing.B) {
	benchmarkStandardSortStr(b, 1024*1024, "revsorted")
}

func BenchmarkTimsortStrRandom1M(b *testing.B) {
	benchmarkTimsortStr(b, 1024*1024, "random")
}

func BenchmarkStandardSortStrRandom1M(b *testing.B) {
	benchmarkStandardSortStr(b, 1024*1024, "random")
}
