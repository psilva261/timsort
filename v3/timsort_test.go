package timsort

import (
	"fmt"
	"math/rand"
	"testing"
)

type val struct {
	key, order int
}

func makeTestArray(size int) []val {
	a := make([]val, size)

	for i := 0; i < size; i++ {
		a[i] = val{i & 0xeeeeee, i}
	}

	return a
}

func IsSorted[T any](a []T, lessThan LessThan[T]) bool {
	len := len(a)

	if len < 2 {
		return true
	}

	prev := a[0]
	for i := 1; i < len; i++ {
		if lessThan(a[i], prev) {
			return false
		}
		prev = a[i]
	}

	return true
}

func TestIsSorted(t *testing.T) {
	a := []val{
		{3, 1},
		{1, 5},
		{2, 3},
		{3, 4},
		{4, 5},
	}

	if IsSorted(a, OrderLessThan) {
		t.Error("Sorted")
	}

}

// use this comparator for sorting
func KeyLessThan(a, b val) bool {
	return a.key < b.key
}

type KeyLessThanSlice []val

func (s KeyLessThanSlice) Len() int {
	return len(s)
}

func (s KeyLessThanSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s KeyLessThanSlice) Less(i, j int) bool {
	return s[i].key < s[j].key
}

// use this comparator to validate sorted data (and prove its stable)
func KeyOrderLessThan(a, b val) bool {
	if a.key < b.key {
		return true
	} else if a.key == b.key {
		return a.order < b.order
	}

	return false
}

// use this comparator to restore the original order of elements (by sorting on order field)
func OrderLessThan(a, b val) bool {
	return a.order < b.order
}

type OrderLessThanSlice []val

func (s OrderLessThanSlice) Len() int {
	return len(s)
}

func (s OrderLessThanSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s OrderLessThanSlice) Less(i, j int) bool {
	return s[i].order < s[j].order
}

func TestSmoke(t *testing.T) {
	a := []val{
		{3, 0},
		{1, 1},
		{2, 2},
	}

	Sort(a, KeyLessThan)

	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}
}

func TestSmokeStability(t *testing.T) {
	a := []val{
		{3, 0},
		{2, 1},
		{2, 2},
	}

	Sort(a, KeyLessThan)

	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}
}

func Test0(t *testing.T) {
	a := []val{}

	Sort(a, KeyLessThan)
	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}
}

func Test1(t *testing.T) {
	a := makeTestArray(1)

	Sort(a, KeyLessThan)
	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}
}

func Test1K(t *testing.T) {
	a := makeTestArray(1024)

	Sort(a, KeyLessThan)
	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}
}

func Test100K(t *testing.T) {
	a := makeTestArray(100 * 1024)

	Sort(a, KeyLessThan)
	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}
}

func Test1M(t *testing.T) {
	a := makeTestArray(1024 * 1024)

	Sort(a, KeyLessThan)
	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}
}

func TestConst(t *testing.T) {
	a := []val{
		{1, 1},
		{1, 1},
		{1, 1},
	}

	Sort(a, KeyLessThan)
	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}
}

func makeRandomArray(size int) []val {
	a := make([]val, size)

	for i := 0; i < size; i++ {
		a[i] = val{rand.Intn(100), i}
	}

	return a
}

func Equals(a, b val) bool {
	return a.key == b.key && a.order == b.order
}

func TestRandom1M(t *testing.T) {
	size := 1024 * 1024

	a := makeRandomArray(size)
	b := make([]val, size)
	copy(b, a)

	Sort(a, KeyLessThan)
	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}

	// sort by order
	Sort(a, OrderLessThan)
	for i := 0; i < len(b); i++ {
		if !Equals(b[i], a[i]) {
			t.Error("oops")
		}
	}
}

const (
	_Sawtooth = iota
	_Rand
	_Stagger
	_Plateau
	_Shuffle
	_NDist
)

const (
	_Copy = iota
	_Reverse
	_ReverseFirstHalf
	_ReverseSecondHalf
	_Sorted
	_Dither
	_NMode
)

func TestBentleyMcIlroy(t *testing.T) {
	//	sizes := []int{100, 1023, 1024, 1025, 1023 * 1023, 1024 * 1024, 1025 * 1025}
	sizes := []int{100, 1023, 1024, 1025}
	dists := []string{"sawtooth", "rand", "stagger", "plateau", "shuffle"}
	modes := []string{"copy", "reverse", "reverse1", "reverse2", "sort", "dither"}
	tmp1 := make([]val, 1025*1025)
	tmp2 := make([]val, 1025*1025)
	for ni := 0; ni < len(sizes); ni++ {
		n := sizes[ni]
		for m := 1; m < 2*n; m *= 2 {
			for dist := 0; dist < _NDist; dist++ {
				j := 0
				k := 1
				data := tmp1[0:n]
				for i := 0; i < n; i++ {
					switch dist {
					case _Sawtooth:
						data[i] = val{i % m, i}
					case _Rand:
						data[i] = val{rand.Intn(m), i}
					case _Stagger:
						data[i] = val{(i*m + i) % n, i}
					case _Plateau:
						if i < m {
							data[i] = val{i, i}
						} else {
							data[i] = val{m, i}
						}
					case _Shuffle:
						if rand.Intn(m) != 0 {
							j += 2
							data[i] = val{j, i}
						} else {
							k += 2
							data[i] = val{k, i}
						}
					}
				}

				mdata := tmp2[0:n]
				for mode := 0; mode < _NMode; mode++ {
					switch mode {
					case _Copy:
						for i := 0; i < n; i++ {
							mdata[i] = val{data[i].key, i}
						}
					case _Reverse:
						for i := 0; i < n; i++ {
							mdata[i] = val{data[n-i-1].key, i}
						}
					case _ReverseFirstHalf:
						for i := 0; i < n/2; i++ {
							mdata[i] = val{data[n/2-i-1].key, i}
						}
						for i := n / 2; i < n; i++ {
							mdata[i] = val{data[i].key, i}
						}
					case _ReverseSecondHalf:
						for i := 0; i < n/2; i++ {
							mdata[i] = val{data[i].key, i}
						}
						for i := n / 2; i < n; i++ {
							mdata[i] = val{data[n-(i-n/2)-1].key, i}
						}
					case _Sorted:
						for i := 0; i < n; i++ {
							mdata[i] = data[i]
						}
						// SortInts is known to be correct
						// because mode Sort runs after mode _Copy.
						Sort(mdata, KeyLessThan)
					case _Dither:
						for i := 0; i < n; i++ {
							mdata[i] = val{data[i].key + i%5, i}
						}
					}

					desc := fmt.Sprintf("n=%d m=%d dist=%s mode=%s", n, m, dists[dist], modes[mode])

					for i := 0; i < len(mdata); i++ {
						mdata[i] = val{mdata[i].key, i}
					}

					gdata := make([]val, len(mdata))
					copy(gdata, mdata)

					Sort(mdata, KeyLessThan)

					// If we were testing C qsort, we'd have to make a copy
					// of the array and sort it ourselves and then compare
					// x against it, to ensure that qsort was only permuting
					// the data, not (for example) overwriting it with zeros.
					//
					// In go, we don't have to be so paranoid: since the only
					// mutating method Sort can call is TestingData.swap,
					// it suffices here just to check that the final array is sorted.
					if !IsSorted(mdata, KeyOrderLessThan) {
						t.Errorf("%s: ints not sorted", desc)
						t.Errorf("\t%v", mdata)
						t.FailNow()
					}

					Sort(mdata, OrderLessThan)
					for i := 0; i < len(data); i++ {
						if !Equals(gdata[i], mdata[i]) {
							t.Error("restore sort failed")
							t.Errorf("gdata=%v", gdata)
							t.Errorf("mdata=%v", mdata)
							t.Errorf("bad index: %v\n", i)
							t.FailNow()
						}
					}
				}
			}
		}
	}
}
