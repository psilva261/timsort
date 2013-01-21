package timsort

import (
	"fmt"
	"math/rand"
	"testing"
)

func makeTestArrayI(size int) []int {
	a := make([]int, size)

	for i := 0; i < size; i++ {
		a[i] = i & 0xeeeeee
	}

	return a
}

func IsSortedI(a []int, lessThan IntLessThan) bool {
	len := len(a)

	if len < 2 {
		return true
	}

	prev := a[0]
	for i := 1; i < len; i++ {
		if lessThan(a[i], prev) {
			fmt.Println("false")
			return false
		}
	}

	return true
}

// use this comparator for sorting
func intLessThan(a, b int) bool {
	return a < b
}

func TestSmokeI(t *testing.T) {
	a := []int{3, 1, 2}

	err := Ints(a, intLessThan)
	if err != nil {
		t.Fatalf("sort: %v", err)
	}

	if !IsSortedI(a, intLessThan) {
		t.Error("not sorted")
	}
}

func Test1KI(t *testing.T) {
	a := makeTestArrayI(1024)

	err := Ints(a, intLessThan)
	if err != nil {
		fmt.Printf("sort: %v\n", err)
		t.Fatalf("sort: %v", err)
	}
	if !IsSortedI(a, intLessThan) {
		t.Error("not sorted")
	}
}

func Test1MI(t *testing.T) {
	a := makeTestArrayI(1024 * 1024)

	err := Ints(a, intLessThan)
	if err != nil {
		t.Fatalf("sort: %v", err)
	}
	if !IsSortedI(a, intLessThan) {
		t.Error("not sorted")
	}
}

func makeRandomArrayI(size int) []int {
	a := make([]int, size)

	for i := 0; i < size; i++ {
		a[i] = rand.Intn(100)
	}

	return a
}

func TestRandom1MI(t *testing.T) {
	size := 1024 * 1024

	a := makeRandomArrayI(size)
	b := make([]int, size)
	copy(b, a)

	err := Ints(a, intLessThan)
	if err != nil {
		t.Fatalf("sort: %v", err)
	}
	if !IsSortedI(a, intLessThan) {
		t.Error("not sorted")
	}
}
