# Benchmarking timsort

Tested on: Intel Core 2 Duo 2.13 GHz running macOS 10.13.5 (64 bit), Go 1.10.1

Numbers are ns/op as reported by `go test -test.bench=.*` / `scripts/benchmark_table.sh`. First number is for timsort, followed by standard sort in parantheses.
`Xor100` means - sorting 100 elements generated using `Xor` method,
`Random1M` means - sorting 1 meg (1024*1024) records generated randomly. 
For more detail on data shapes see the source - [bench_test.go][bench_test.go]. 
Three columns represent three benchmark runs. 

    Sorted100:    4183          (6680)        4189          (6657)        4178          (6612)
    RevSorted100: 4341          (7677)        4621          (7631)        4481          (7812)
    Xor100:       18160         (11691)       17353         (11911)       17707         (11753)
    Random100:    22044         (13274)       22194         (13458)       21805         (13310)

    Sorted1K:     16768         (98000)       16750         (97830)       16277         (97651)
    RevSorted1K:  19432         (101731)      19594         (101747)      19557         (101764)
    Xor1K:        158106        (165457)      158474        (165644)      157338        (165257)
    Random1K:     298479        (199176)      309216        (198687)      296884        (200112)

    Sorted1M:     13835534      (209459201)   13906732      (209015568)   13836590      (210636112) 
    RevSorted1M:  22408878      (213692594)   20677165      (213109279)   20579740      (213315686) 
    Xor1M:        403355901     (174788669)   336437853     (175793278)   379731433     (174839708) 
    Random1M:     958138342     (405581976)   898752000     (416798092)   954800710     (415745244) 

Not surprisingly, timsort is crazy fast on sorted inputs. But even for random and quasi-random (xor) inputs, timsort is just 50% slower than built-in sort.

### Disclaimer

The above benchmark applies only to one specific type of data element (`record` structure as defined in [bench_test.go][bench_test.go]). For other data types results may vary.

[bench_test.go]: http://github.com/psilva261/timsort/blob/master/bench_test.go
