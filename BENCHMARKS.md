# Benchmarking timsort

Tested on: AMD Athlon(tm) 64 X2 Dual Core Processor 5200+ running 64-bit ubuntu 10.10.

Numbers are ns/op as reported by `make bench`. First number is for timsort, followed by standard sort in parantheses.
`Xor100` means - sorting 100 elements generated using `Xor` method,
`Random1M` means - sorting 1 meg (1024*1024) records generated randomly. 
For more detail on data shapes see sorce - `bench_test.go`. 
Three columns represent three benchmark runs. 

    Sorted100:          7042(70578)            6933(68945)            6890(68955)
    RevSorted100:       7328(71239)            7283(70233)            7252(69968)
    Xor100:            30106(79482)           29673(78166)           29890(79522)
    Random100:         36738(86734)           36200(85511)           36390(85456)

    Sorted1K:	       49769(1160276)         48998(1144771)         49008(1146563)
    RevSorted1K:       53824(1148334)         53070(1134186)         52968(1136656)
    Xor1K:            351850(1232022)        340447(1217338)        338730(1215607)
    Random1K:         562683(1373354)        556980(1357697)        557346(1357735)

    Sorted1M:       49643000(2723631000)   49026980(2711728000)   49049180(2684859000)
    RevSorted1M:    58038420(2708541000)   57416740(2703899000)   57463300(2679528000)
    Xor1M:         688159400(1393372000)  676922600(1376313000)  677548800(1375267000)
    Random1M:     1618038000(3600431000) 1589860000(3548680000) 1591280000(3557897000)

Not surprisingly, timsort is crazy fast on sorted inputs. But even for random and quasi-random (xor) inputs, timsort is at least 2x faster than built-in sort.

### Disclaimer

The above benchmark applies only to one specific type of data element (`record` structure as defined in `bench_test.go`). For other data types results may vary.
