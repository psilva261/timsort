# Benchmarking timsort

Tested on: Intel i7 8809G @3.10G running Debian 4.18.20-2 (2018-11-23) x86_64 GNU/Linux, go1.11.2 linux/amd64

Numbers are ns/op as reported by `go test -test.bench=.*` / `scripts/*.sh`. First number is for timsort, followed by standard sort.Sort (for ints) or sort.Stable otherwise, in parentheses.
`Xor100` means - sorting 100 elements generated using `Xor` method,
`Random1M` means - sorting 1 meg (1024*1024) records generated randomly. 
For more detail on data shapes see the source - [bench_test.go][bench_test.go]. 
Three columns represent three benchmark runs. 

### Comparing stable sort

benchmark_table.sh:
    
    Sorted100:    1421          (668)         1362          (696)         1400          (682)         
    RevSorted100: 1540          (7934)        1541          (7858)        1544          (8351)        
    Xor100:       4797          (5909)        4787          (6060)        4913          (5848)        
    Random100:    5211          (8144)        5251          (7878)        5259          (7979)        
    
    Sorted1K:     5528          (4856)        5604          (4958)        5575          (4994)        
    RevSorted1K:  6694          (86263)       6549          (89099)       6473          (89702)       
    Xor1K:        45905         (115574)      45706         (121319)      48774         (114212)      
    Random1K:     117190        (206253)      113427        (202286)      116871        (208261)      
    
    Sorted1M:     4238347       (8991637)     4125454       (8718271)     4199811       (8524249)     
    RevSorted1M:  6498133       (122685926)   6274050       (113631337)   6218917       (117388612)   
    Xor1M:        103830787     (226644420)   97834093      (232690160)   100051210     (228723419)   
    Random1M:     348764192     (663024283)   367589312     (690514504)   377960208     (733221145) ```

Timsort is crazy fast on patterned inputs, but even for random input it is much faster than stable sort in Go. Exception is a fully sorted input, where Go's stable sort greatly benefits on small inputs because it uses insertion sort in implementation.
Disclaimer: The above benchmark applies only to one specific type of data element (`record` structure as defined in [bench_test.go][bench_test.go]). For other data types results may vary. 


### Benchmarking primitives

benchmarkint_table.sh:

    Sorted100:    948           (1977)        942           (1926)        918           (1965)        
    RevSorted100: 1003          (2161)        1016          (2166)        1004          (2133)        
    Xor100:       3327          (2831)        3297          (2858)        3256          (2876)        
    Random100:    3640          (2895)        3580          (2900)        3546          (2892)        
    
    Sorted1K:     3664          (29693)       3757          (29639)       3568          (29545)       
    RevSorted1K:  4724          (31183)       4629          (31691)       4653          (31430)       
    Xor1K:        30678         (53533)       31136         (53791)       32052         (55790)       
    Random1K:     93291         (78021)       93153         (75886)       96178         (77012)       
    
    Sorted1M:     2612806       (63640085)    2657223       (66381050)    2782516       (65576825)    
    RevSorted1M:  3236216       (65236079)    3241050       (62996742)    3445819       (66550484)    
    Xor1M:        62054181      (61609233)    58021830      (59778514)    57094210      (65384476)    
    Random1M:     212171678     (185525353)   199493986     (184466052)   223314023     (188326317)```

If stable sort is not required then default sort in Go (quicksort) gives better results for random inputs.


[bench_test.go]: http://github.com/psilva261/timsort/blob/master/bench_test.go
