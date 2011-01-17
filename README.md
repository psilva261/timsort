# timsort

**timsort** is a Go implementation of Tim Peters's mergesort
sorting algorithm.

For many input types its 2-3 times faster than Go's built-in sorting.

This implementation was derived from Java's TimSort object by Josh Bloch,
which, in turn, was based on the original code by Tim Peters:

	http://svn.python.org/projects/python/trunk/Objects/listsort.txt


## Installation

1. git clone git://github.com/pgmmpk/timsort.git
2. cd timsort
3. make install

Alternatively, you can intall using `goinstall github.com/pgmmpk/timsort`, but
if you do this, the import statement in your programs will be `import github.com/pgmmpk/timsort` instead of just `import timsort`.

## Example

	package main

	import (
		"timsort"
		"fmt"
	)

	type Record struct {
		ssn int
		name string
	}

	func BySsn(a, b interface{}) bool {
		return a.(Record).ssn < b.(Record).ssn
	}

	func ByName(a, b interface{}) bool {
		return a.(Record).name < b.(Record).name
	}

	func main() {
		db := make([]interface{}, 3)
		db[0] = Record{123456789, "joe"}
		db[1] = Record{993456780, "sue"}
		db[2] = Record{345623452, "mary"}

		// sorts array by ssn (ascending)
		timsort.Sort(db, BySsn)
		fmt.Printf("sorted by ssn: %v\n", db)

		// now re-sort same array by name (ascending)
		timsort.Sort(db, ByName)
		fmt.Printf("sorted by name: %v\n", db)
	}

## Todo

* more testing
* benchmarking
* /done/ try different values for MIN_MERGE contstant (original Tim's code used 64, while java code uses 32)
* try replacing stackLen computation with simple constant (as in original Tim's code)

Mike K.

