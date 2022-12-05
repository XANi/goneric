[![go report card](https://goreportcard.com/badge/github.com/XANi/goneric "go report card")](https://goreportcard.com/report/github.com/XANi/goneric)
[![test status](https://github.com/go-gorm/gorm/workflows/tests/badge.svg?branch=master "test status")](https://github.com/XANi/goneric)
[![codecov](https://codecov.io/gh/XANi/goneric/branch/master/graph/badge.svg?token=079HADYAJG)](https://codecov.io/gh/XANi/goneric)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/XANi/goneric?tab=doc)

# Goneric 

Collection of generics-related utility functions, slice/map/channel manipulation and some parallel processing.

## Conventions

If not specified the returning value will be at least initiated to be empty, not nil.

Functions returning stuff like `map[key]boolean` (using map as set) set the boolean value to `true` 
because `m["nonexistent_key"] == false`

Naming function goes in order of `operation_group`, `input type`, `modifier`/`output` with ones irrelevant skipped.
There are few exceptions for convenience, like variadic `Map`.


## Functions

### Slice 

* `CompareSliceSet` - check whether slices have same elements regardless of order
* `SliceMap` - Convert slice to map using function to get the key. `[]Struct{} -> map[Struct.key]Struct` being the common usage
* `SliceMapSet` - Convert slice to set-like map. `[]Comparable -> map[Comparable]bool{true}`
* `SliceMapSetFunc` - Convert slice to set-like map via helper function. `[]Any -> map[func(Any)Comparable]bool{true}`
* `SliceDiff` - return difference between 2 slices of same comparable type, in form of 2 variables where first have elements 
   that are only in first set, and second elements only in second set. `([]T, []T) -> (leftOnly []T, rightOnly []T)`
* `SliceDiffFunc`  - As `SliceDiff` but type of slice is irrelevant, via use of conversion function that converts it
   into comparable type. `([]T1,[]T2) -> (leftOnly []T1, rightOnly []T2)`
* `SliceIn` - Check whether value is in slice

### Map

* `Map` - Map variadic input thru function
* `MapSlice` - Map slice thru function
* `MapSliceKey` - Convert map to slice of its keys. `map[Comparable]V -> []Comparable` 
* `MapSliceValue` - Convert map to slice of its values. `map[Comparable]V -> []V`
* `MapErr` - Same as `Map` but function can return error that will stop the loop and propagate it out
* `MapSliceErr` - Same as `MapSlice` but function can return error that will stop the loop and propagate it out
* `MapSliceSkip` - Same as `MapSlice` but function can return true in second argument to skip the entry
* `MapSliceErrSkip` - Same as `MapSliceErr` but `ErrSkip` error type can be used to skip entry instead of erroring out

### Filter

* `FilterMap` - Filter thru a map using a function
* `FilterSlice` - Filter thru a slice using a function
* `FilterChan` - Filter thru a channel using a function
* `FilterChanErr` - Filter thru a channel using a function, with separate output channel for that function errors

### Channel tools

* `ChanGen` - Uses function to generate channel messages
* `ChanGenCloser` - Use function to generate channel messages, stop when closer function is called
* `ChanToSlice` - Loads data to slice from channel until channel is closed
* `ChanToSliceN` - Loads data to slice from channel to at most N elements
* `SliceToChan` - Sends slice to channel
* `SliceToChanClose` - Sends slice to channel, then closes it.

### Worker

* `WorkerPool` - spawn x goroutines with workers and return after input channel is closed and all requests are parsed
* `WorkerPoolClose` - spawn x goroutines with workers and return after input channel is closed and all requests are parsed. Also close output channel
* `WorkerPoolBackground` - spawn x goroutines with workers in background and returns output channel
* `WorkerPoolBackgroundClose` - spawn x goroutines with workers in background and returns output channel. 
   Close output channel if input channel is closed, after processing all messages
* `WorkerPoolFinisher` - spawn x goroutines with workers in background, returns finisher channel that signals with `bool{true}` when the processing ends.

### Parallel

* `ParallelMap` - like `Map` but runs function in parallel up to specified number of goroutines. Ordered.
* `ParallelMapSlice` - like `MapSlice` but runs function in parallel up to specified number of goroutines. Ordered.
* `ParallelMapSliceChan` - runs slice elements thru function and sends it to channel
* `ParallelMapSliceChanFinisher` - runs slice elements thru function and sends it to channel. 
   Returns `finisher chan(bool){true}` that will return single `true` message when all workers finish and close it

### Math

Not equivalent of `math` library, NaN math is ignored, zero length inputs might, sanitize your inputs.

Results unless specified otherwise will follow math of type, so median of `[]int{8,9}` will be `int{8}` coz of rounding.


* `Sum`
* `SumF64` - sum returning float64, for summing up small ints
* `Min`
* `Max`
* `Avg`
* `AvgF64` - average with final input calculated as float64. Addition is still in source 
* `AvgF64F64` - average with float64 accumulator. Use if you want to avoid overflow on small int type
* `Median`
* `MedianF64` - median with final division using float64 type to avoid overflows


## Types

* `Number` - any basic numeric types

## Current project state

Unstable interface - the function names might change, the new ones will still be committed with full test coverage, no benchmarks yet.
Take care with `get -u` until 1.0.0 hits.


## Analytics

* Coverage map

  ![img](https://codecov.io/gh/XANi/goneric/branch/master/graphs/tree.svg?token=079HADYAJG)