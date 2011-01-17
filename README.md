Go implementation of Tim Peter's mergesort sorting algorithm.

This implementation was derived from Java's TimSort object by Josh Bloch,
which, in turn, was based on the original code by Tim Peters:

 http://svn.python.org/projects/python/trunk/Objects/listsort.txt

TODO: more testing
TODO: benchmarking
DONE: try different values for MIN_MERGE contstant (original Tim's code used 64,
       while java code uses 32)
TODO: try replacing stackLen computation with simple constant (as in original Tim's code)

Mike K.

