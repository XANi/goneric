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
Functions where it is sensible to have option to close channel should have that as last optional argument.
There are few exceptions for convenience, like variadic `Map`.

If possible, sensible, functions that take function parameter should have function as first parameter.

Channel-operating function *in general* should accept channel as parameter; the ones returning a channel should be under `Gen*` hierarchy


## Examples

### Get the list of elements that differ between slices of different types

```go
data1 := []string{"1", "2", "3", "4", "5"}
data2 := []float32{1, 7, 3, 4}
stringToInt := func(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
floatToInt := func(f float32) int {
	return int(f)
}
left, right := SliceDiffFunc(data1, data2, stringToInt, floatToInt)
fmt.Printf("left: %T%+v right: %T%+v", left, left, right, right)
// left: []string[2 5] right: []float32[7]
```



### Convert every string to float, exit on first error

```go
in := []string{"1", "2.2", "3", "cat","5"}
out, err := goneric.MapSliceErr(func(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}, in)
fmt.Printf("%+v, err: %s", out, err)
// [1 2.2 3], err: strconv.ParseFloat: parsing "cat": invalid syntax
```

### Convert every string to float, skip entries with error


```go
in := []string{"1", "2.2", "3", "cat", "5"}
out := goneric.MapSliceSkip(func(s string) (float64, bool) {
	f, err := strconv.ParseFloat(s, 64)
	return f, err != nil
}, in)
fmt.Printf("%+v", out)
// [1 2.2 3 5]
```

### Run every element of map thru function in parallel, creating new map, at max concurrency of 2 goroutines

```go
data := map[string]int{
	"a": 99,
	"b": 250,
	"c": 30,
	"d": 9,
}
mappedData := goneric.ParallelMapMap(func(k string, v int) (string, string) {
	time.Sleep(time.Millisecond * time.Duration(v))
	return k, fmt.Sprintf("0x%02x", v)
}, 2, data)
fmt.Printf("%+v", mappedData)
// map[a:0x63 b:0xfa c:0x1e d:0x09]
```


## Functions


### Slice 

* `CompareSliceSet` - check whether slices have same elements regardless of order
* `SliceMap` - Convert slice to map using function to get the key. `[]Struct{} -> f(s)K -> map[Struct.key]Struct` being the common usage
* `SliceMapFunc` - Convert slice to map using function to return key and value. `[]Struct{} -> f(s)(K,V) -> map[K]V`
* `SliceMapSet` - Convert slice to set-like map. `[]K -> map[K]bool{true}`
* `SliceMapSetFunc` - Convert slice to set-like map via helper function. `[]T -> map[func(T)Comparable]bool{true}`
* `SliceDiff` - return difference between 2 slices of same comparable type, in form of 2 variables where first have elements 
   that are only in first set, and second elements only in second set. `([]T, []T) -> (leftOnly []T, rightOnly []T)`
* `SliceDiffFunc`  - As `SliceDiff` but type of slice is irrelevant, via use of conversion function that converts it
   into comparable type. `([]T1,[]T2) -> (leftOnly []T1, rightOnly []T2)`
* `SliceIn` - Check whether value is in slice
* `SliceDedupe` - remove duplicates from `comparable` slice. `[]T -> []T`
* `SliceDedupeFunc` - remove duplicates from `any` slice via conversion function. `[]T -> []T`
* `SliceReverse` - reverses the order of elements in slice and returns reversed copy, `[]T -> []T`
* `SliceReverseInplace` - reverses the order of elements in slice in-place.
* `FirstOrEmpty` - return first element or empty value. `[]T -> T`
* `LastOrEmpty` - return last element or empty value. `[]T -> T`
* `FirstOrEmpty` - return first element or passed "default" value. `[]T -> T`
* `LastOrEmpty` - return last element or passed "default" value. `[]T -> T`


### Map

* `MapMap` - Map one map to another using a function. `map[K1]V1 -> map[K2]V2`
* `MapMapInplace` - Map one map to another using a function, filling existing passed map. `(in map[K1]V1, out map[K2]V2)`
* `Map` - Map variadic input thru function. `T1... -> []T2`
* `MapSlice` - Map slice thru function. `[]T1 -> []T2`
* `MapSliceKey` - Convert map to slice of its keys. `map[K]V -> []K` 
* `MapSliceValue` - Convert map to slice of its values. `map[K]V -> []V`
* `MapErr` - Same as `Map` but function can return error that will stop the loop and propagate it out.  `T1... -> ([]T2,err)`
* `MapSliceErr` - Same as `MapSlice` but function can return error that will stop the loop and propagate it out.  `T1... -> ([]T2,err)`
* `MapSliceSkip` - Same as `MapSlice` but function can return true in second argument to skip the entry. `[]T1 -> []T2`
* `MapSliceErrSkip` - Same as `MapSliceErr` but `ErrSkip` error type can be used to skip entry instead of erroring out.  `[]T1 -> ([]T2,err)`


### Filter

* `FilterMap` - Filter thru a map using a function. `map[K]V -> map[K]V`
* `FilterSlice` - Filter thru a slice using a function. `[]T -> []T`
* `FilterChan` - Filter thru a channel using a function. `in chan T -> out chan T`
* `FilterChanErr` - Filter thru a channel using a function, with separate output channel for that function errors. `in chan T -> (out chan T,err chan error)`


### Channel tools

* `ChanGen` - Feed function output to passed channel in a loop. `(f()T, chan T)`
* `ChanGenN` - Feed function output to passed channel in a loop N times, optionally close it. `(f()T, count, chan T)`
* `ChanGenCloser` - Use function to pass generated messages to channel, stop when closer function is called,  `(f()T, chan T) -> chan closeChannel`
* `ChanToSlice` - Loads data to slice from channel until channel is closed. `chan T -> []T`
* `ChanToSliceN` - Loads data to slice from channel to at most N elements. `(chan T,count) -> []T`
* `ChanToSliceNTimeout` - Loads data to slice from channel to at most N elements or until timeout passes. `(chan T,count,timeout) -> []T`
* `SliceToChan` - Sends slice to passed channel in background, optionally closes it. `[]T -> chan T`


### Worker

* `WorkerPool` - spawn x goroutines with workers and return after input channel is closed and all requests are parsed. Optionally close output
* `WorkerPoolBackground` - spawn x goroutines with workers in background and returns output channel. Optionally close output.
* `WorkerPoolFinisher` - spawn x goroutines with workers in background, returns finisher channel that signals with `bool{true}` when the processing ends.
* `WorkerPoolDrain` - spawn x goroutines that will run a function on the channel element without returning anything
* `WorkerPoolAsync` - function will run x goroutines for worker in the background and return a function that enqueues job and returns channel with result of that job, allowing to queue stuff to run in background conveniently


### Parallel

* `ParallelMap` - like `Map` but runs function in parallel up to specified number of goroutines. Ordered.
* `ParallelMap` - like `MapMap` but runs function in parallel up to specified number of goroutines. Ordered.
* `ParallelMapSlice` - like `MapSlice` but runs function in parallel up to specified number of goroutines. Ordered.
* `ParallelMapSliceChan` - runs slice elements thru function and sends it to channel
* `ParallelMapSliceChanFinisher` - runs slice elements thru function and sends it to channel. 
   Returns `finisher chan(bool){true}` that will return single `true` message when all workers finish and close it


### Async

* `Async` - run function in background goroutine and return result as a channel. `func()T -> chan T`
* `AsyncV` - run functions in background goroutine and return result as a channel, then close it. `funcList... -> chan T`
* `AsyncVUnpanic` - run function in background goroutine and return result as a channel, then close it, ignoring every panic. `funcList... -> chan T`
* `AsyncPipe` - run function in background, taking and returning values to pipe. Designed to be chained. `(in chan T1,  func(T1)T2) -> chan T2`
* `AsyncOut` - as `AsyncPipe` but takes output channel as argument .`(in chan T1, func(T1)T2, chan T2)`
* `AsyncIn` - converts value into channel with that value. `T -> chan T`


### (Re)Try

* `Retry` - retry function X times
* `RetryAfter` - retry with timeout, minimal, and maximal interval between retries.
* `Try` - tries each function in slice till first success


### Generators

Generator functions always return generated values. 
For ones that operate on passed on types look at `Type*Gen*` functions like `SliceGen`

* `GenSlice` - generate slice of given length via function. `func(idx int) T -> []T`
* `GenMap` - generate Map of given length via function. `func(idx int) K,V -> map[K]V`
* `GenChan` - returns channel fed from generator function ad infinitum. `func() K -> chan T`
* `GenChanN` - returns channel fed from generator function N times. `func() K -> chan T`
* `GenChanNCloser` - returns channel fed from generator that returns closer() function that will stop generator from running.`func() K -> (chan T,func closer())` 
* `GenSliceToChan` - returns channel fed from slice, optionally closes it, `[]K -> chan T`


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
* `ValueIndex` - represents slice element with index
* `KeyValue` - represents map key/value pair

## Miscellaneous 

* `Must` - Turn error into panic, returning non-err arguments
* `IgnoreErr` - Ignores error, returning default type value if error is passed

## Current project state

No API changes to existing functions AP will be made in 1.x.x releases


## Analytics

* Coverage map

  ![img](https://codecov.io/gh/XANi/goneric/branch/master/graphs/tree.svg?token=079HADYAJG)