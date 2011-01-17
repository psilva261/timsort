package timsort

import (
	"testing"
	"rand"
	"fmt"
	"container/vector"
)

type val struct {
	key, order, payload int
}

func makeTestArray(size int) []interface{} {
	a := make([]interface{}, size)

	for i := 0; i < size; i++ {
		a[i] = val{i & 0xeeeeee, i, 0}
	}

	return a
}

func IsSorted(a []interface{}, lessThan LessThan) bool {
	len := len(a)

	if len < 2 {
		return true
	}

	prev := a[0]
	for i := 1; i < len; i++ {
		if lessThan(a[1], prev) {
			return false
		}
	}

	return true
}

func KeyLessThan(a, b interface{}) bool {
	return a.(val).key < b.(val).key
}

func KeyOrderLessThan(a, b interface{}) bool {
	if a.(val).key < b.(val).key {
		return true
	} else if a.(val).key == b.(val).key {
		return a.(val).order < b.(val).order
	}

	return false
}

func OrderLessThan(a, b interface{}) bool {
	return a.(val).order < b.(val).order
}

func TestSmoke(t *testing.T) {
	a := make([]interface{}, 3)
	a[0] = val{3, 0, 0}
	a[1] = val{1, 1, 0}
	a[2] = val{2, 2, 0}

	Sort(a, KeyLessThan)

	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}
}


func TestSmokeStability(t *testing.T) {
	a := make([]interface{}, 3)
	a[0] = val{3, 0, 0}
	a[1] = val{2, 1, 0}
	a[2] = val{2, 2, 0}

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

func Test1M(t *testing.T) {
	a := makeTestArray(1024 * 1024)

	Sort(a, KeyLessThan)
	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}
}

func makeRandomArray(size int) []interface{} {
	a := make([]interface{}, size)

	for i := 0; i < size; i++ {
		a[i] = val{rand.Intn(100), i, 0}
	}

	return a
}

func Equals(a, b interface{}) bool {
	return a.(val).key == b.(val).key && a.(val).order == b.(val).order && a.(val).payload == b.(val).payload
}

func TestRandom1M(t *testing.T) {
	size := 1024 * 1024

	a := makeRandomArray(size)
	b := make([]interface{}, size)
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

type person struct {
	ssn  int
	name string
}

func (self *person) Less(o interface{}) bool {
	return self.ssn < o.(*person).ssn
}

func TestSortVector(t *testing.T) {
	var v vector.Vector

	v.Push(&person{123, "jack"})
	v.Push(&person{124, "mary"})
	v.Push(&person{120, "bill"})

	SortVector(v)

	if v[0].(*person).name != "bill" {
		t.Error("bill")
	}
	if v[1].(*person).name != "jack" {
		t.Error("jack")
	}
	if v[2].(*person).name != "mary" {
		t.Error("mary")
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
	tmp1 := make([]interface{}, 1025*1025)
	tmp2 := make([]interface{}, 1025*1025)
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
						data[i] = val{i % m, i, 0}
					case _Rand:
						data[i] = val{rand.Intn(m), i, 0}
					case _Stagger:
						data[i] = val{(i*m + i) % n, i, 0}
					case _Plateau:
						if i < m {
							data[i] = val{i, i, 0}
						} else {
							data[i] = val{m, i, 0}
						}
					case _Shuffle:
						if rand.Intn(m) != 0 {
							j += 2
							data[i] = val{j, i, 0}
						} else {
							k += 2
							data[i] = val{k, i, 0}
						}
					}
				}

				mdata := tmp2[0:n]
				for mode := 0; mode < _NMode; mode++ {
					switch mode {
					case _Copy:
						for i := 0; i < n; i++ {
							mdata[i] = val{data[i].(val).key, i, 0}
						}
					case _Reverse:
						for i := 0; i < n; i++ {
							mdata[i] = val{data[n-i-1].(val).key, i, 0}
						}
					case _ReverseFirstHalf:
						for i := 0; i < n/2; i++ {
							mdata[i] = val{data[n/2-i-1].(val).key, i, 0}
						}
						for i := n / 2; i < n; i++ {
							mdata[i] = val{data[i].(val).key, i, 0}
						}
					case _ReverseSecondHalf:
						for i := 0; i < n/2; i++ {
							mdata[i] = val{data[i].(val).key, i, 0}
						}
						for i := n / 2; i < n; i++ {
							mdata[i] = val{data[n-(i-n/2)-1].(val).key, i, 0}
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
							mdata[i] = val{data[i].(val).key + i%5, i, 0}
						}
					}

					desc := fmt.Sprintf("n=%d m=%d dist=%s mode=%s", n, m, dists[dist], modes[mode])

					for i := 0; i < len(mdata); i++ {
						mdata[i] = val{mdata[i].(val).key, i, 0}
					}

					gdata := make([]interface{}, len(mdata))
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
