package timsort

import (
	"math/rand"
	"sort"
	"testing"
)

type ints []int

func (p *ints) Len() int           { return len(*p) }
func (p *ints) Less(i, j int) bool { return (*p)[i] < (*p)[j] }
func (p *ints) Swap(i, j int)      { (*p)[i], (*p)[j] = (*p)[j], (*p)[i] }


func LessThanInt(a, b int) bool {
	return a < b
}

func makeInts(size int, shape string) (v ints) {
	v = make(ints, 0, size)
	switch shape {

	case "xor":
		for i := 0; i < size; i++ {
			v = append(v, 0xff&(i^0xab))
		}

	case "sorted":
		for i := 0; i < size; i++ {
			v = append(v, i)
		}

	case "revsorted":
		for i := 0; i < size; i++ {
			v = append(v, size-i)
		}

	case "random":
		rand.Seed(1)

		for i := 0; i < size; i++ {
			v = append(v, rand.Int())
		}

	default:
		panic(shape)
	}

	return v
}

func benchmarkTimsortI(b *testing.B, size int, shape string) {
	b.StopTimer()

	for j := 0; j < b.N; j++ {
		v := makeInts(size, shape)

		b.StartTimer()
		Ints(v, LessThanInt)
		b.StopTimer()
	}
}

func benchmarkStandardSortI(b *testing.B, size int, shape string) {
	b.StopTimer()

	for j := 0; j < b.N; j++ {
		v := makeInts(size, shape)

		b.StartTimer()
		sort.Sort(&v)
		b.StopTimer()
	}
}

func BenchmarkTimsortIXor100(b *testing.B) {
	benchmarkTimsortI(b, 100, "xor")
}

func BenchmarkStandardSortIXor100(b *testing.B) {
	benchmarkStandardSortI(b, 100, "xor")
}

func BenchmarkTimsortISorted100(b *testing.B) {
	benchmarkTimsortI(b, 100, "sorted")
}

func BenchmarkStandardSortISorted100(b *testing.B) {
	benchmarkStandardSortI(b, 100, "sorted")
}

func BenchmarkTimsortIRevSorted100(b *testing.B) {
	benchmarkTimsortI(b, 100, "revsorted")
}

func BenchmarkStandardSortIRevSorted100(b *testing.B) {
	benchmarkStandardSortI(b, 100, "revsorted")
}

func BenchmarkTimsortIRandom100(b *testing.B) {
	benchmarkTimsortI(b, 100, "random")
}

func BenchmarkStandardSortIRandom100(b *testing.B) {
	benchmarkStandardSortI(b, 100, "random")
}

func BenchmarkTimsortIXor1K(b *testing.B) {
	benchmarkTimsortI(b, 1024, "xor")
}

func BenchmarkStandardSortIXor1K(b *testing.B) {
	benchmarkStandardSortI(b, 1024, "xor")
}

func BenchmarkTimsortISorted1K(b *testing.B) {
	benchmarkTimsortI(b, 1024, "sorted")
}

func BenchmarkStandardSortISorted1K(b *testing.B) {
	benchmarkStandardSortI(b, 1024, "sorted")
}

func BenchmarkTimsortIRevSorted1K(b *testing.B) {
	benchmarkTimsortI(b, 1024, "revsorted")
}

func BenchmarkStandardSortIRevSorted1K(b *testing.B) {
	benchmarkStandardSortI(b, 1024, "revsorted")
}

func BenchmarkTimsortIRandom1K(b *testing.B) {
	benchmarkTimsortI(b, 1024, "random")
}

func BenchmarkStandardSortIRandom1K(b *testing.B) {
	benchmarkStandardSortI(b, 1024, "random")
}

func BenchmarkTimsortIXor1M(b *testing.B) {
	benchmarkTimsortI(b, 1024*1024, "xor")
}

func BenchmarkStandardSortIXor1M(b *testing.B) {
	benchmarkStandardSortI(b, 1024*1024, "xor")
}

func BenchmarkTimsortISorted1M(b *testing.B) {
	benchmarkTimsortI(b, 1024*1024, "sorted")
}

func BenchmarkStandardSortISorted1M(b *testing.B) {
	benchmarkStandardSortI(b, 1024*1024, "sorted")
}

func BenchmarkTimsortIRevSorted1M(b *testing.B) {
	benchmarkTimsortI(b, 1024*1024, "revsorted")
}

func BenchmarkStandardSortIRevSorted1M(b *testing.B) {
	benchmarkStandardSortI(b, 1024*1024, "revsorted")
}

func BenchmarkTimsortIRandom1M(b *testing.B) {
	benchmarkTimsortI(b, 1024*1024, "random")
}

func BenchmarkStandardSortIRandom1M(b *testing.B) {
	benchmarkStandardSortI(b, 1024*1024, "random")
}
