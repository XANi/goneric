[![go report card](https://goreportcard.com/badge/github.com/XANi/goneric "go report card")](https://goreportcard.com/report/github.com/XANi/goneric)
[![test status](https://github.com/go-gorm/gorm/workflows/tests/badge.svg?branch=master "test status")](https://github.com/XANi/goneric)
[![codecov](https://codecov.io/gh/XANi/goneric/branch/master/graph/badge.svg?token=079HADYAJG)](https://codecov.io/gh/XANi/goneric)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/XANI/goneric?tab=doc)

# Goneric 

Collection of generics-related utility functions

## Conventions

If not specified the returning value will be at least initiated to be empty, not nil.

Functions returning stuff like `map[key]boolean` (using map as set) set the boolean value to `true` 
because `m["nonexistent_key"] == false`


## Functions

### Slice 

* `CompareSliceSet` - check whether slices have same elements regardless of order
* `SliceMap` - Convert slice to map using function to get the key. `[]Struct{} -> map[Struct.key]Struct` being the common usage
* `SliceMapSet` - Convert slice to set-like map. `[]Comparable -> map[Comparable]bool{true}`
* `SliceMapSetFunc` - Convert slice to set-like map via helper function. `[]Any -> map[func(Any)Comparable]bool{true}`
* `SliceDiff` - return difference between 2 slices of same comparable type, in form of 2 variables where first have elements 
   that are only in first set, and second elements only in second set
* `SliceDiffFunc`  - As `SliceDiff` but type of slice is irrelevant, via use of conversion function that converts it
   into comparable type  

### Map

* `Map` - Map variadic input thru function
* `MapSlice` - Map slice thru function
* `MapSliceKey` - Convert map to slice of its keys. `map[Comparable]V -> []Comparable` 
* `MapSliceValue` - Convert map to slice of its values. `map[Comparable]V -> []V`
* `MapErr` - Same as `Map` but function can return error that will stop the loop and propagate it out
* `MapSliceErr` - Same as `MapSlice` but function can return error that will stop the loop and propagate it out


## Analytics

* Coverage map

  ![img](https://codecov.io/gh/XANi/goneric/branch/master/graphs/tree.svg?token=079HADYAJG)