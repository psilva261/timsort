package timsort

// Package timsort provides fast stable sort, uses external comparator.
//
// A stable, adaptive, iterative mergesort that requires far fewer than
// n lg(n) comparisons when running on partially sorted arrays, while
// offering performance comparable to a traditional mergesort when run
// on random arrays.  Like all proper mergesorts, this sort is stable and
// runs O(n log n) time (worst case).  In the worst case, this sort requires
// temporary storage space for n/2 object references; in the best case,
// it requires only a small constant amount of space.
//
// This implementation was derived from Java's TimSort object by Josh Bloch,
// which, in turn, was based on the original code by Tim Peters:
//
// http://svn.python.org/projects/python/trunk/Objects/listsort.txt
//
// Mike K.

const (
	/**
	 * This is the minimum sized sequence that will be merged.  Shorter
	 * sequences will be lengthened by calling binarySort.  If the entire
	 * array is less than this length, no merges will be performed.
	 *
	 * This constant should be a power of two.  It was 64 in Tim Peter's C
	 * implementation, but 32 was empirically determined to work better in
	 * this implementation.  In the unlikely event that you set this constant
	 * to be a number that's not a power of two, you'll need to change the
	 * {@link #minRunLength} computation.
	 *
	 * If you decrease this constant, you must change the stackLen
	 * computation in the TimSort constructor, or you risk an
	 * ArrayOutOfBounds exception.  See listsort.txt for a discussion
	 * of the minimum stack length required as a function of the length
	 * of the array being sorted and the minimum merge sequence length.
	 */
	minMerge = 32
	// mk: tried higher MIN_MERGE and got slower sorting (348->375)
	//	c_MIN_MERGE = 64

	/**
	 * When we get into galloping mode, we stay there until both runs win less
	 * often than cminGallop consecutive times.
	 */
	minGallop = 7

	/**
	 * Maximum initial size of tmp array, which is used for merging.  The array
	 * can grow to accommodate demand.
	 *
	 * Unlike Tim's original C version, we do not allocate this much storage
	 * when sorting smaller arrays.  This change was required for performance.
	 */
	initialTmpStorageLength = 256
)

// LessThan is Delegate type that sorting uses as a comparator
type LessThan[T any] func(a, b T) bool

type timSortHandler[T any] struct {

	/**
	 * The array being sorted.
	 */
	a []T

	/**
	 * The comparator for this sort.
	 */
	lt LessThan[T]

	/**
	 * This controls when we get *into* galloping mode.  It is initialized
	 * to cminGallop.  The mergeLo and mergeHi methods nudge it higher for
	 * random data, and lower for highly structured data.
	 */
	minGallop int

	/**
	 * Temp storage for merges.
	 */
	tmp []T // Actual runtime type will be Object[], regardless of T

	/**
	 * A stack of pending runs yet to be merged.  Run i starts at
	 * address base[i] and extends for len[i] elements.  It's always
	 * true (so long as the indices are in bounds) that:
	 *
	 *     runBase[i] + runLen[i] == runBase[i + 1]
	 *
	 * so we could cut the storage for this, but it's a minor amount,
	 * and keeping all the info explicit simplifies the code.
	 */
	stackSize int // Number of pending runs on stack
	runBase   []int
	runLen    []int
}

/**
 * Creates a TimSort instance to maintain the state of an ongoing sort.
 *
 * @param a the array to be sorted
 * @param c the comparator to determine the order of the sort
 */
func newTimSort[T any](a []T, lt LessThan[T]) (h *timSortHandler[T]) {
	h = &timSortHandler[T]{}

	h.a = a
	h.lt = lt
	h.minGallop = minGallop
	h.stackSize = 0

	// Allocate temp storage (which may be increased later if necessary)
	len := len(a)

	tmpSize := initialTmpStorageLength
	if len < 2*tmpSize {
		tmpSize = len / 2
	}

	h.tmp = make([]T, tmpSize)

	/*
	 * Allocate runs-to-be-merged stack (which cannot be expanded).  The
	 * stack length requirements are described in listsort.txt.  The C
	 * version always uses the same stack length (85), but this was
	 * measured to be too expensive when sorting "mid-sized" arrays (e.g.,
	 * 100 elements) in Java.  Therefore, we use smaller (but sufficiently
	 * large) stack lengths for smaller arrays.  The "magic numbers" in the
	 * computation below must be changed if c_MIN_MERGE is decreased.  See
	 * the c_MIN_MERGE declaration above for more information.
	 */
	// mk: confirmed that for small sorts this optimization gives measurable (albeit small)
	// performance enhancement
	stackLen := 40
	if len < 120 {
		stackLen = 5
	} else if len < 1542 {
		stackLen = 10
	} else if len < 119151 {
		stackLen = 19
	}

	h.runBase = make([]int, stackLen)
	h.runLen = make([]int, stackLen)

	return h
}

// Sort an array using the provided comparator
func Sort[T any](a []T, lt LessThan[T]) {
	lo := 0
	hi := len(a)
	nRemaining := hi

	if nRemaining < 2 {
		return // Arrays of size 0 and 1 are always sorted
	}

	// If array is small, do a "mini-TimSort" with no merges
	if nRemaining < minMerge {
		initRunLen := countRunAndMakeAscending(a, lo, hi, lt)

		binarySort(a, lo, hi, lo+initRunLen, lt)
		return
	}

	/**
	 * March over the array once, left to right, finding natural runs,
	 * extending short natural runs to minRun elements, and merging runs
	 * to maintain stack invariant.
	 */

	ts := newTimSort(a, lt)
	minRun := minRunLength(nRemaining)
	for {
		// Identify next run
		runLen := countRunAndMakeAscending(a, lo, hi, lt)

		// If run is short, extend to min(minRun, nRemaining)
		if runLen < minRun {
			force := minRun
			if nRemaining <= minRun {
				force = nRemaining
			}
			binarySort(a, lo, lo+force, lo+runLen, lt)
			runLen = force
		}

		// Push run onto pending-run stack, and maybe merge
		ts.pushRun(lo, runLen)
		ts.mergeCollapse()

		// Advance to find next run
		lo += runLen
		nRemaining -= runLen
		if nRemaining == 0 {
			break
		}
	}

	ts.mergeForceCollapse()
}

/**
 * Sorts the specified portion of the specified array using a binary
 * insertion sort.  This is the best method for sorting small numbers
 * of elements.  It requires O(n log n) compares, but O(n^2) data
 * movement (worst case).
 *
 * If the initial part of the specified range is already sorted,
 * this method can take advantage of it: the method assumes that the
 * elements from index {@code lo}, inclusive, to {@code start},
 * exclusive are already sorted.
 *
 * @param a the array in which a range is to be sorted
 * @param lo the index of the first element in the range to be sorted
 * @param hi the index after the last element in the range to be sorted
 * @param start the index of the first element in the range that is
 *        not already known to be sorted (@code lo <= start <= hi}
 * @param c comparator to used for the sort
 */
func binarySort[T any](a []T, lo, hi, start int, lt LessThan[T]) {
	if start == lo {
		start++
	}

	for ; start < hi; start++ {
		pivot := a[start]

		// Set left (and right) to the index where a[start] (pivot) belongs
		left := lo
		right := start

		/*
		 * Invariants:
		 *   pivot >= all in [lo, left).
		 *   pivot <  all in [right, start).
		 */
		for left < right {
			mid := int(uint(left+right) >> 1)
			if lt(pivot, a[mid]) {
				right = mid
			} else {
				left = mid + 1
			}
		}

		/*
		 * The invariants still hold: pivot >= all in [lo, left) and
		 * pivot < all in [left, start), so pivot belongs at left.  Note
		 * that if there are elements equal to pivot, left points to the
		 * first slot after them -- that's why this sort is stable.
		 * Slide elements over to make room to make room for pivot.
		 */
		n := start - left // The number of elements to move
		// just an optimization for copy in default case
		if n <= 2 {
			if n == 2 {
				a[left+2] = a[left+1]
			}
			if n > 0 {
				a[left+1] = a[left]
			}
		} else {
			copy(a[left+1:], a[left:left+n])
		}
		a[left] = pivot
	}
}

/*
*
  - Returns the length of the run beginning at the specified position in
  - the specified array and reverses the run if it is descending (ensuring
  - that the run will always be ascending when the method returns).
    *
  - A run is the longest ascending sequence with:
    *
  - a[lo] <= a[lo + 1] <= a[lo + 2] <= ...
    *
  - or the longest descending sequence with:
    *
  - a[lo] >  a[lo + 1] >  a[lo + 2] >  ...
    *
  - For its intended use in a stable mergesort, the strictness of the
  - definition of "descending" is needed so that the call can safely
  - reverse a descending sequence without violating stability.
    *
  - @param a the array in which a run is to be counted and possibly reversed
  - @param lo index of the first element in the run
  - @param hi index after the last element that may be contained in the run.
    It is required that @code{lo < hi}.
  - @param c the comparator to used for the sort
  - @return  the length of the run beginning at the specified position in
  - the specified array
*/
func countRunAndMakeAscending[T any](a []T, lo, hi int, lt LessThan[T]) int {
	runHi := lo + 1
	if runHi == hi {
		return 1
	}

	// Find end of run, and reverse range if descending
	if lt(a[runHi], a[lo]) { // Descending
		runHi++

		for runHi < hi && lt(a[runHi], a[runHi-1]) {
			runHi++
		}
		reverseRange(a, lo, runHi)
	} else { // Ascending
		for runHi < hi && !lt(a[runHi], a[runHi-1]) {
			runHi++
		}
	}

	return runHi - lo
}

/**
 * Reverse the specified range of the specified array.
 *
 * @param a the array in which a range is to be reversed
 * @param lo the index of the first element in the range to be reversed
 * @param hi the index after the last element in the range to be reversed
 */
func reverseRange[T any](a []T, lo, hi int) {
	hi--
	for lo < hi {
		a[lo], a[hi] = a[hi], a[lo]
		lo++
		hi--
	}
}

/**
 * Returns the minimum acceptable run length for an array of the specified
 * length. Natural runs shorter than this will be extended with
 * {@link #binarySort}.
 *
 * Roughly speaking, the computation is:
 *
 *  If n < c_MIN_MERGE, return n (it's too small to bother with fancy stuff).
 *  Else if n is an exact power of 2, return c_MIN_MERGE/2.
 *  Else return an int k, c_MIN_MERGE/2 <= k <= c_MIN_MERGE, such that n/k
 *   is close to, but strictly less than, an exact power of 2.
 *
 * For the rationale, see listsort.txt.
 *
 * @param n the length of the array to be sorted
 * @return the length of the minimum run to be merged
 */
func minRunLength(n int) int {
	r := 0 // Becomes 1 if any 1 bits are shifted off
	for n >= minMerge {
		r |= (n & 1)
		n >>= 1
	}
	return n + r
}

/**
 * Pushes the specified run onto the pending-run stack.
 *
 * @param runBase index of the first element in the run
 * @param runLen  the number of elements in the run
 */
func (h *timSortHandler[T]) pushRun(runBase, runLen int) {
	h.runBase[h.stackSize] = runBase
	h.runLen[h.stackSize] = runLen
	h.stackSize++
}

/**
 * Examines the stack of runs waiting to be merged and merges adjacent runs
 * until the stack invariants are reestablished:
 *
 *     1. runLen[i - 3] > runLen[i - 2] + runLen[i - 1]
 *     2. runLen[i - 2] > runLen[i - 1]
 *
 * This method is called each time a new run is pushed onto the stack,
 * so the invariants are guaranteed to hold for i < stackSize upon
 * entry to the method.
 */
func (h *timSortHandler[T]) mergeCollapse() {
	for h.stackSize > 1 {
		n := h.stackSize - 2
		if (n > 0 && h.runLen[n-1] <= h.runLen[n]+h.runLen[n+1]) ||
			(n > 1 && h.runLen[n-2] <= h.runLen[n-1]+h.runLen[n]) {
			if h.runLen[n-1] < h.runLen[n+1] {
				n--
			}
			h.mergeAt(n)
		} else if h.runLen[n] <= h.runLen[n+1] {
			h.mergeAt(n)
		} else {
			break // Invariant is established
		}
	}
}

/**
 * Merges all runs on the stack until only one remains.  This method is
 * called once, to complete the sort.
 */
func (h *timSortHandler[T]) mergeForceCollapse() {
	for h.stackSize > 1 {
		n := h.stackSize - 2
		if n > 0 && h.runLen[n-1] < h.runLen[n+1] {
			n--
		}
		h.mergeAt(n)
	}
}

/**
 * Merges the two runs at stack indices i and i+1.  Run i must be
 * the penultimate or antepenultimate run on the stack.  In other words,
 * i must be equal to stackSize-2 or stackSize-3.
 *
 * @param i stack index of the first of the two runs to merge
 */
func (h *timSortHandler[T]) mergeAt(i int) {
	base1 := h.runBase[i]
	len1 := h.runLen[i]
	base2 := h.runBase[i+1]
	len2 := h.runLen[i+1]

	/*
	 * Record the length of the combined runs; if i is the 3rd-last
	 * run now, also slide over the last run (which isn't involved
	 * in this merge).  The current run (i+1) goes away in any case.
	 */
	h.runLen[i] = len1 + len2
	if i == h.stackSize-3 {
		h.runBase[i+1] = h.runBase[i+2]
		h.runLen[i+1] = h.runLen[i+2]
	}
	h.stackSize--

	/*
	 * Find where the first element of run2 goes in run1. Prior elements
	 * in run1 can be ignored (because they're already in place).
	 */
	k := gallopRight(h.a[base2], h.a, base1, len1, 0, h.lt)
	base1 += k
	len1 -= k
	if len1 == 0 {
		return
	}

	/*
	 * Find where the last element of run1 goes in run2. Subsequent elements
	 * in run2 can be ignored (because they're already in place).
	 */
	len2 = gallopLeft(h.a[base1+len1-1], h.a, base2, len2, len2-1, h.lt)
	if len2 == 0 {
		return
	}

	// Merge remaining runs, using tmp array with min(len1, len2) elements
	if len1 <= len2 {
		h.mergeLo(base1, len1, base2, len2)
	} else {
		h.mergeHi(base1, len1, base2, len2)
	}
}

/**
 * Locates the position at which to insert the specified key into the
 * specified sorted range; if the range contains an element equal to key,
 * returns the index of the leftmost equal element.
 *
 * @param key the key whose insertion point to search for
 * @param a the array in which to search
 * @param base the index of the first element in the range
 * @param len the length of the range; must be > 0
 * @param hint the index at which to begin the search, 0 <= hint < n.
 *     The closer hint is to the result, the faster this method will run.
 * @param c the comparator used to order the range, and to search
 * @return the int k,  0 <= k <= n such that a[b + k - 1] < key <= a[b + k],
 *    pretending that a[b - 1] is minus infinity and a[b + n] is infinity.
 *    In other words, key belongs at index b + k; or in other words,
 *    the first k elements of a should precede key, and the last n - k
 *    should follow it.
 */
func gallopLeft[T any](key T, a []T, base, len, hint int, c LessThan[T]) int {
	lastOfs := 0
	ofs := 1

	if c(a[base+hint], key) {
		// Gallop right until a[base+hint+lastOfs] < key <= a[base+hint+ofs]
		maxOfs := len - hint
		for ofs < maxOfs && c(a[base+hint+ofs], key) {
			lastOfs = ofs
			ofs = (ofs << 1) + 1
			if ofs <= 0 { // int overflow
				ofs = maxOfs
			}
		}
		if ofs > maxOfs {
			ofs = maxOfs
		}

		// Make offsets relative to base
		lastOfs += hint
		ofs += hint
	} else { // key <= a[base + hint]
		// Gallop left until a[base+hint-ofs] < key <= a[base+hint-lastOfs]
		maxOfs := hint + 1
		for ofs < maxOfs && !c(a[base+hint-ofs], key) {
			lastOfs = ofs
			ofs = (ofs << 1) + 1
			if ofs <= 0 { // int overflow
				ofs = maxOfs
			}
		}
		if ofs > maxOfs {
			ofs = maxOfs
		}

		// Make offsets relative to base
		tmp := lastOfs
		lastOfs = hint - ofs
		ofs = hint - tmp
	}

	/*
	 * Now a[base+lastOfs] < key <= a[base+ofs], so key belongs somewhere
	 * to the right of lastOfs but no farther right than ofs.  Do a binary
	 * search, with invariant a[base + lastOfs - 1] < key <= a[base + ofs].
	 */
	lastOfs++
	for lastOfs < ofs {
		m := lastOfs + (ofs-lastOfs)/2

		if c(a[base+m], key) {
			lastOfs = m + 1 // a[base + m] < key
		} else {
			ofs = m // key <= a[base + m]
		}
	}

	return ofs
}

/**
 * Like gallopLeft, except that if the range contains an element equal to
 * key, gallopRight returns the index after the rightmost equal element.
 *
 * @param key the key whose insertion point to search for
 * @param a the array in which to search
 * @param base the index of the first element in the range
 * @param len the length of the range; must be > 0
 * @param hint the index at which to begin the search, 0 <= hint < n.
 *     The closer hint is to the result, the faster this method will run.
 * @param c the comparator used to order the range, and to search
 * @return the int k,  0 <= k <= n such that a[b + k - 1] <= key < a[b + k]
 */
func gallopRight[T any](key T, a []T, base, len, hint int, c LessThan[T]) int {
	ofs := 1
	lastOfs := 0
	if c(key, a[base+hint]) {
		// Gallop left until a[b+hint - ofs] <= key < a[b+hint - lastOfs]
		maxOfs := hint + 1
		for ofs < maxOfs && c(key, a[base+hint-ofs]) {
			lastOfs = ofs
			ofs = (ofs << 1) + 1
			if ofs <= 0 { // int overflow
				ofs = maxOfs
			}
		}
		if ofs > maxOfs {
			ofs = maxOfs
		}

		// Make offsets relative to b
		tmp := lastOfs
		lastOfs = hint - ofs
		ofs = hint - tmp
	} else { // a[b + hint] <= key
		// Gallop right until a[b+hint + lastOfs] <= key < a[b+hint + ofs]
		maxOfs := len - hint
		for ofs < maxOfs && !c(key, a[base+hint+ofs]) {
			lastOfs = ofs
			ofs = (ofs << 1) + 1
			if ofs <= 0 { // int overflow
				ofs = maxOfs
			}
		}
		if ofs > maxOfs {
			ofs = maxOfs
		}

		// Make offsets relative to b
		lastOfs += hint
		ofs += hint
	}

	/*
	 * Now a[b + lastOfs] <= key < a[b + ofs], so key belongs somewhere to
	 * the right of lastOfs but no farther right than ofs.  Do a binary
	 * search, with invariant a[b + lastOfs - 1] <= key < a[b + ofs].
	 */
	lastOfs++
	for lastOfs < ofs {
		m := lastOfs + (ofs-lastOfs)/2

		if c(key, a[base+m]) {
			ofs = m // key < a[b + m]
		} else {
			lastOfs = m + 1 // a[b + m] <= key
		}
	}
	return ofs
}

/**
 * Merges two adjacent runs in place, in a stable fashion.  The first
 * element of the first run must be greater than the first element of the
 * second run (a[base1] > a[base2]), and the last element of the first run
 * (a[base1 + len1-1]) must be greater than all elements of the second run.
 *
 * For performance, this method should be called only when len1 <= len2;
 * its twin, mergeHi should be called if len1 >= len2.  (Either method
 * may be called if len1 == len2.)
 *
 * @param base1 index of first element in first run to be merged
 * @param len1  length of first run to be merged (must be > 0)
 * @param base2 index of first element in second run to be merged
 *        (must be aBase + aLen)
 * @param len2  length of second run to be merged (must be > 0)
 */
func (h *timSortHandler[T]) mergeLo(base1, len1, base2, len2 int) {
	// Copy first run into temp array
	a := h.a // For performance
	tmp := h.ensureCapacity(len1)

	copy(tmp, a[base1:base1+len1])

	cursor1 := 0     // Indexes into tmp array
	cursor2 := base2 // Indexes int a
	dest := base1    // Indexes int a

	// Move first element of second run and deal with degenerate cases
	a[dest] = a[cursor2]
	dest++
	cursor2++
	len2--
	if len2 == 0 {
		copy(a[dest:dest+len1], tmp)
		return
	}
	if len1 == 1 {
		copy(a[dest:dest+len2], a[cursor2:cursor2+len2])
		a[dest+len2] = tmp[cursor1] // Last elt of run 1 to end of merge
		return
	}

	lt := h.lt               // Use local variable for performance
	minGallop := h.minGallop //  "    "       "     "      "

outer:
	for {
		count1 := 0 // Number of times in a row that first run won
		count2 := 0 // Number of times in a row that second run won

		/*
		 * Do the straightforward thing until (if ever) one run starts
		 * winning consistently.
		 */
		for {
			if lt(a[cursor2], tmp[cursor1]) {
				a[dest] = a[cursor2]
				dest++
				cursor2++
				count2++
				count1 = 0
				len2--
				if len2 == 0 {
					break outer
				}
			} else {
				a[dest] = tmp[cursor1]
				dest++
				cursor1++
				count1++
				count2 = 0
				len1--
				if len1 == 1 {
					break outer
				}
			}
			if (count1 | count2) >= minGallop {
				break
			}
		}

		/*
		 * One run is winning so consistently that galloping may be a
		 * huge win. So try that, and continue galloping until (if ever)
		 * neither run appears to be winning consistently anymore.
		 */
		for {
			count1 = gallopRight(a[cursor2], tmp, cursor1, len1, 0, lt)
			if count1 != 0 {
				copy(a[dest:dest+count1], tmp[cursor1:cursor1+count1])
				dest += count1
				cursor1 += count1
				len1 -= count1
				if len1 <= 1 { // len1 == 1 || len1 == 0
					break outer
				}
			}
			a[dest] = a[cursor2]
			dest++
			cursor2++
			len2--
			if len2 == 0 {
				break outer
			}

			count2 = gallopLeft(tmp[cursor1], a, cursor2, len2, 0, lt)
			if count2 != 0 {
				copy(a[dest:dest+count2], a[cursor2:cursor2+count2])
				dest += count2
				cursor2 += count2
				len2 -= count2
				if len2 == 0 {
					break outer
				}
			}
			a[dest] = tmp[cursor1]
			dest++
			cursor1++
			len1--
			if len1 == 1 {
				break outer
			}
			minGallop--
			if count1 < minGallop && count2 < minGallop {
				break
			}
		}
		if minGallop < 0 {
			minGallop = 0
		}
		minGallop += 2 // Penalize for leaving gallop mode
	} // End of "outer" loop

	if minGallop < 1 {
		minGallop = 1
	}
	h.minGallop = minGallop // Write back to field

	if len1 == 1 {

		copy(a[dest:dest+len2], a[cursor2:cursor2+len2])
		a[dest+len2] = tmp[cursor1] //  Last elt of run 1 to end of merge
	} else {
		copy(a[dest:dest+len1], tmp[cursor1:cursor1+len1])
	}
}

/**
 * Like mergeLo, except that this method should be called only if
 * len1 >= len2; mergeLo should be called if len1 <= len2.  (Either method
 * may be called if len1 == len2.)
 *
 * @param base1 index of first element in first run to be merged
 * @param len1  length of first run to be merged (must be > 0)
 * @param base2 index of first element in second run to be merged
 *        (must be aBase + aLen)
 * @param len2  length of second run to be merged (must be > 0)
 */
func (h *timSortHandler[T]) mergeHi(base1, len1, base2, len2 int) {
	// Copy second run into temp array
	a := h.a // For performance
	tmp := h.ensureCapacity(len2)

	copy(tmp, a[base2:base2+len2])

	cursor1 := base1 + len1 - 1 // Indexes into a
	cursor2 := len2 - 1         // Indexes into tmp array
	dest := base2 + len2 - 1    // Indexes into a

	// Move last element of first run and deal with degenerate cases
	a[dest] = a[cursor1]
	dest--
	cursor1--
	len1--
	if len1 == 0 {
		dest -= len2 - 1
		copy(a[dest:dest+len2], tmp)
		return
	}
	if len2 == 1 {
		dest -= len1 - 1
		cursor1 -= len1 - 1
		copy(a[dest:dest+len1], a[cursor1:cursor1+len1])
		a[dest-1] = tmp[cursor2]
		return
	}

	lt := h.lt               // Use local variable for performance
	minGallop := h.minGallop //  "    "       "     "      "

outer:
	for {
		count1 := 0 // Number of times in a row that first run won
		count2 := 0 // Number of times in a row that second run won

		/*
		 * Do the straightforward thing until (if ever) one run
		 * appears to win consistently.
		 */
		for {
			if lt(tmp[cursor2], a[cursor1]) {
				a[dest] = a[cursor1]
				dest--
				cursor1--
				count1++
				count2 = 0
				len1--
				if len1 == 0 {
					break outer
				}
			} else {
				a[dest] = tmp[cursor2]
				dest--
				cursor2--
				count2++
				count1 = 0
				len2--
				if len2 == 1 {
					break outer
				}
			}
			if (count1 | count2) >= minGallop {
				break
			}
		}

		/*
		 * One run is winning so consistently that galloping may be a
		 * huge win. So try that, and continue galloping until (if ever)
		 * neither run appears to be winning consistently anymore.
		 */
		for {
			gr := gallopRight(tmp[cursor2], a, base1, len1, len1-1, lt)
			count1 = len1 - gr
			if count1 != 0 {
				dest -= count1
				cursor1 -= count1
				len1 -= count1
				copy(a[dest+1:dest+1+count1], a[cursor1+1:cursor1+1+count1])
				if len1 == 0 {
					break outer
				}
			}
			a[dest] = tmp[cursor2]
			dest--
			cursor2--
			len2--
			if len2 == 1 {
				break outer
			}

			gl := gallopLeft(a[cursor1], tmp, 0, len2, len2-1, lt)
			count2 = len2 - gl
			if count2 != 0 {
				dest -= count2
				cursor2 -= count2
				len2 -= count2
				copy(a[dest+1:dest+1+count2], tmp[cursor2+1:cursor2+1+count2])
				if len2 <= 1 { // len2 == 1 || len2 == 0
					break outer
				}
			}
			a[dest] = a[cursor1]
			dest--
			cursor1--
			len1--
			if len1 == 0 {
				break outer
			}
			minGallop--

			if count1 < minGallop && count2 < minGallop {
				break
			}
		}
		if minGallop < 0 {
			minGallop = 0
		}
		minGallop += 2 // Penalize for leaving gallop mode
	} // End of "outer" loop

	if minGallop < 1 {
		minGallop = 1
	}

	h.minGallop = minGallop // Write back to field

	if len2 == 1 {
		dest -= len1
		cursor1 -= len1

		copy(a[dest+1:dest+1+len1], a[cursor1+1:cursor1+1+len1])
		a[dest] = tmp[cursor2] // Move first elt of run2 to front of merge
	} else {
		copy(a[dest-(len2-1):dest+1], tmp)
	}
}

/**
 * Ensures that the external array tmp has at least the specified
 * number of elements, increasing its size if necessary.  The size
 * increases exponentially to ensure amortized linear time complexity.
 *
 * @param minCapacity the minimum required capacity of the tmp array
 * @return tmp, whether or not it grew
 */
func (h *timSortHandler[T]) ensureCapacity(minCapacity int) []T {
	if len(h.tmp) < minCapacity {
		// Compute smallest power of 2 > minCapacity
		newSize := minCapacity
		newSize |= newSize >> 1
		newSize |= newSize >> 2
		newSize |= newSize >> 4
		newSize |= newSize >> 8
		newSize |= newSize >> 16
		newSize++

		if newSize < 0 { // Not bloody likely!
			newSize = minCapacity
		} else {
			ns := len(h.a) / 2
			if ns < newSize {
				newSize = ns
			}
		}

		h.tmp = make([]T, newSize)
	}

	return h.tmp
}
