# microspace [![Build Status](https://travis-ci.org/WatchBeam/microspace.svg?branch=master)](https://travis-ci.org/WatchBeam/microspace) [![godoc reference](https://godoc.org/github.com/WatchBeam/microspace?status.png)](https://godoc.org/github.com/WatchBeam/microspace)

Microspace is a spatial index that's optimized for very fast building and nearest-n lookups.

### Performance Characteristics

 - Inserting `n` points ultimately runs in `O(n log n)` time.
 - In the worse case, querying for the nearest neighbor runs in `O(n)` time.
 - But in practice it runs much faster. In the benchmark random distribution, querying for nearest neighbors in a set size of 10000 took ~10 nanoseconds longer than querying in a set size of 10.

```
# Generating a tree with `n` elements:
BenchmarkIndexCreate10                 200000              6010 ns/op
BenchmarkIndexCreate100                 20000             64149 ns/op
BenchmarkIndexCreate1000                 2000            839679 ns/op
BenchmarkIndexCreate10000                 200           8848370 ns/op

# Querying for 3 nearest neighbors in a random set of `n` elements:
BenchmarkIndexNearest10               5000000               294 ns/op
BenchmarkIndexNearest100              5000000               287 ns/op
BenchmarkIndexNearest1000             5000000               286 ns/op
BenchmarkIndexNearest10000            5000000               305 ns/op

# Querying for 3 nearest neighbors in the worst-case set of `n` elements:
BenchmarkIndexNearestWorstCase10      3000000               485 ns/op
BenchmarkIndexNearestWorstCase100      500000              2642 ns/op
BenchmarkIndexNearestWorstCase1000     100000             23939 ns/op
BenchmarkIndexNearestWorstCase10000      5000            240544 ns/op
```
