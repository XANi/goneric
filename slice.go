package goneric

// CompareSliceSet compares 2 slices and returns true if all elements of slice v1
// are the same as in v2
// order does not matter, the duplicates are ignored
func CompareSliceSet[T comparable](v1 []T, v2 []T) bool {
	m1 := map[T]bool{}
	m2 := map[T]bool{}
	for _, v := range v1 {
		m1[v] = true
	}
	for _, v := range v2 {
		if _, ok := m1[v]; !ok {
			return false
		}
		m2[v] = true
	}
	for k := range m1 {
		if _, ok := m2[k]; !ok {
			return false
		}
	}
	return true
}

// SliceMap turn slice into a map via extracting key from it using helper function
// and setting the map value to that slice
// `[]Struct{} -> map[func(Struct)]Struct`
func SliceMap[T any, M comparable](f func(T) M, a []T) map[M]T {
	n := make(map[M]T, len(a))
	for _, e := range a {
		n[f(e)] = e
	}
	return n
}

// SliceMapSkip works like `SliceMap` but
// allows slice->map function to skip elements via returning true to second argument
// `[]Struct{} -> map[func(Struct)]Struct`
func SliceMapSkip[T any, Z comparable](comparable func(T) (comparable Z, skip bool), slice []T) (m map[Z]T) {
	m = make(map[Z]T, len(slice))
	for _, e := range slice {
		k, skip := comparable(e)
		if !skip {
			m[k] = e
		}
	}
	return m
}

// SliceMapSet turns slice into map with key being slice elements and value being true boolean
// `[]Comparable -> map[Comparable]bool{true}`
func SliceMapSet[T comparable](a []T) (n map[T]bool) {
	n = make(map[T]bool, 0)
	for _, e := range a {
		n[e] = true
	}
	return n
}

// SliceMapSetFunc turns slice into map with key being slice elements passed thru specified function
// and value being true boolean
// `[]Any -> map[func(Any)Comparable]bool{true}`
func SliceMapSetFunc[T any, M comparable](mapFunc func(T) M, a []T) (n map[M]bool) {
	n = make(map[M]bool, 0)
	for _, e := range a {
		n[mapFunc(e)] = true
	}
	return n
}

// SliceDiff compares two slices of comparable values and
// returns slice of elements that are only in first/left element
// and ones that are only in right element
// Duplicates are ignored.
// `([]T, []T) -> (leftOnly []T, rightOnly []T)`
func SliceDiff[T comparable](v1 []T, v2 []T) (inLeft []T, inRight []T) {
	m1 := map[T]bool{}
	m2 := map[T]bool{}
	// we want to return empty slice, not nil slice
	inLeft = []T{}
	inRight = []T{}
	for _, v := range v1 {
		m1[v] = true
	}
	for _, v := range v2 {
		if _, ok := m1[v]; !ok {
			inRight = append(inRight, v)
		}
		m2[v] = true
	}
	for _, v := range v1 {
		if _, ok := m2[v]; !ok {
			inLeft = append(inLeft, v)
		}
	}
	return
}

// SliceDiffFunc compares two slices of any value
// using one conversion function per type to convert it into conmparables
// returns slice of elements that are only in first/left element
// and ones that are only in right element.
// Duplicates are ignored.
//([]DataT,[]ReturnT) -> (leftOnly []DataT, rightOnly []ReturnT)
func SliceDiffFunc[T1 any, T2 any, Z comparable](
	v1 []T1,
	v2 []T2,
	convertV1 func(T1) Z,
	convertV2 func(T2) Z,
) (inLeft []T1, inRight []T2) {
	m1 := map[Z]bool{}
	m2 := map[Z]bool{}
	// we want to return empty slice, not nil slice
	inLeft = []T1{}
	inRight = []T2{}
	for _, v := range v1 {
		m1[convertV1(v)] = true
	}
	for _, v := range v2 {
		converted := convertV2(v)
		if _, ok := m1[converted]; !ok {
			inRight = append(inRight, v)
		}
		m2[convertV2(v)] = true
	}
	for _, v := range v1 {
		converted := convertV1(v)
		if _, ok := m2[converted]; !ok {
			inLeft = append(inLeft, v)
		}
	}
	return
}

// SliceIn checks if slice contains a value. Value must be comparable
func SliceIn[T comparable](slice []T, contains T) bool {
	for _, v := range slice {
		if v == contains {
			return true
		}
	}
	return false
}
