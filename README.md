# Append-only set of uint64 with tiny overhead

This is just append-only, fixed size set of uint64. It is faster than map and has lower memory footprint,
but only for uniformly distributed uint64 values. *It will panic if there are no slots left.*
This is essentially hash table with open addressing, but without hash and with dead-simple addressing logic (lookup in nearest slots).
Not concurrent-safe.

## Benchmarks (see test code for details):

```
BenchmarkSet	20000000	       115 ns/op	       9 B/op	       0 allocs/op
BenchmarkMap	10000000	       202 ns/op	      17 B/op	       0 allocs/op
```

You should always allocate more slots than values you want to store. More slots - faster.
But with len(items)/8 extra slots (one extra byte per value) it gives pretty good results (see benchmark).