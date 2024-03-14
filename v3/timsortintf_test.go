package timsort

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
)

func TestSmokeTS(t *testing.T) {
	a := make([]interface{}, 3)
	a[0] = val{3, 0}
	a[1] = val{1, 1}
	a[2] = val{2, 2}

	TimSort(KeyLessThanSlice(a))

	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}
}

func TestSmokeStabilityTS(t *testing.T) {
	a := make([]interface{}, 3)
	a[0] = val{3, 0}
	a[1] = val{2, 1}
	a[2] = val{2, 2}

	TimSort(KeyLessThanSlice(a))

	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}
}

func Test1KTS(t *testing.T) {
	a := makeTestArray(1024)

	TimSort(KeyLessThanSlice(a))
	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}
}

func Test1MTS(t *testing.T) {
	a := makeTestArray(1024 * 1024)

	TimSort(KeyLessThanSlice(a))
	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}
}

func TestRandom1MTS(t *testing.T) {
	size := 1024 * 1024

	a := makeRandomArray(size)
	b := make([]interface{}, size)
	copy(b, a)

	TimSort(KeyLessThanSlice(a))
	if !IsSorted(a, KeyOrderLessThan) {
		t.Error("not sorted")
	}

	// sort by order
	TimSort(OrderLessThanSlice(a))
	for i := 0; i < len(b); i++ {
		if !Equals(b[i], a[i]) {
			t.Error("oops")
		}
	}
}

func TestBentleyMcIlroyTS(t *testing.T) {
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
							mdata[i] = val{data[i].(val).key, i}
						}
					case _Reverse:
						for i := 0; i < n; i++ {
							mdata[i] = val{data[n-i-1].(val).key, i}
						}
					case _ReverseFirstHalf:
						for i := 0; i < n/2; i++ {
							mdata[i] = val{data[n/2-i-1].(val).key, i}
						}
						for i := n / 2; i < n; i++ {
							mdata[i] = val{data[i].(val).key, i}
						}
					case _ReverseSecondHalf:
						for i := 0; i < n/2; i++ {
							mdata[i] = val{data[i].(val).key, i}
						}
						for i := n / 2; i < n; i++ {
							mdata[i] = val{data[n-(i-n/2)-1].(val).key, i}
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
							mdata[i] = val{data[i].(val).key + i%5, i}
						}
					}

					desc := fmt.Sprintf("n=%d m=%d dist=%s mode=%s", n, m, dists[dist], modes[mode])

					for i := 0; i < len(mdata); i++ {
						mdata[i] = val{mdata[i].(val).key, i}
					}

					gdata := make([]interface{}, len(mdata))
					copy(gdata, mdata)

					TimSort(KeyLessThanSlice(mdata))

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

					TimSort(OrderLessThanSlice(mdata))
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

func TestStrings1TS(t *testing.T) {
	a := []string{
		"ef",
		"abc",
		"aaa",
		"de",
		"ed",
	}
	TimSort(sort.StringSlice(a))
}
