package timsort

import (
	"sort"
)

// TimSort sorts the data defined by sort.Interface.
func TimSort(a sort.Interface) {
	indexes := make([]int, a.Len())
	for i := 0; i < len(indexes); i++ {
		indexes[i] = i
	}

	Ints(indexes, func(i, j int) bool {
		return a.Less(i, j)
	})

	for i := 0; i < len(indexes); i++ {
		j := indexes[i]
		if j == 0 {
			continue
		}
		for k := i; j != i; {
			a.Swap(j, k)
			k, j, indexes[j] = j, indexes[j], 0
		}
	}
}
