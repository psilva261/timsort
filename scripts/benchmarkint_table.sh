#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset
# set -o xtrace


sortings=( Sorted RevSorted Xor Random )
amounts=( 100 1K 1M )

for amount in "${amounts[@]}"
do
	for sorting in "${sortings[@]}"
	do
		printf '%-14s' $sorting$amount:
		for n in {1..3}
		do
			printf '%-14s' $(go test -test.bench=TimsortI$sorting$amount | grep Benchmark | awk '{print $3}')
			printf '%-14s' '('$(go test -test.bench=StandardSortI$sorting$amount | grep Benchmark | awk '{print $3}')')'
		done
		echo
	done
	echo
done

