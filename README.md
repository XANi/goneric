[![go report card](https://goreportcard.com/badge/github.com/XANi/goneric "go report card")](https://goreportcard.com/report/github.com/XANi/goneric)
[![test status](https://github.com/go-gorm/gorm/workflows/tests/badge.svg?branch=master "test status")](https://github.com/XANi/goneric)
[![codecov](https://codecov.io/gh/XANi/goneric/branch/master/graph/badge.svg?token=079HADYAJG)](https://codecov.io/gh/XANi/goneric)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/XANi/goneric?tab=doc)

# Goneric 

Collection of generics-related utility functions

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

### Map

* `Map` - Map variadic input thru function
* `MapSlice` - Map slice thru function
* `MapSliceKey` - Convert map to slice of its keys. `map[Comparable]V -> []Comparable` 
* `MapSliceValue` - Convert map to slice of its values. `map[Comparable]V -> []V`
* `MapErr` - Same as `Map` but function can return error that will stop the loop and propagate it out
* `MapSliceErr` - Same as `MapSlice` but function can return error that will stop the loop and propagate it out

## Worker

* `WorkerPool` - spawn x goroutines with workers return after input channel is closed and all requests are parsed

## Parallel

* `ParallelMap` - like `Map` but runs function in parallel up to specified number of goroutines. Ordered.
* `ParallelMapSlice` - like `MapSlice` but runs function in parallel up to specified number of goroutines. Ordered.
* `ParallelMapSliceChannel` - runs slice elements thru function and sends it to channel
* `ParallelMapSliceChannelFinisher` - runs slice elements thru function and sends it to channel. 
   Returns `finisher chan(bool){true}` that will return single `true` message when all workers finish

## Analytics

* Coverage map

  ![img](https://codecov.io/gh/XANi/goneric/branch/master/graphs/tree.svg?token=079HADYAJG)