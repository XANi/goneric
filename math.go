package goneric

import "sort"

func Sum[T Number](n ...T) (sum T) {
	for _, v := range n {
		sum = sum + v
	}
	return sum
}

func SumF64[T Number](n ...T) (sum float64) {
	for _, v := range n {
		sum = sum + float64(v)
	}
	return sum
}

// Max returns biggest number
// will panic on empty
func Max[T Number](n ...T) (max T) {
	max = n[0]
	for _, v := range n[1:] {
		if v > max {
			max = v
		}
	}
	return
}

// Min returns smallest number
// will panic on empty
func Min[T Number](n ...T) (min T) {
	min = n[0]
	for _, v := range n[1:] {
		if v < min {
			min = v
		}
	}
	return
}

func Avg[T Number](n ...T) (avg T) {
	return Sum(n...) / T(len(n))
}

// AvgF64 calculates average with final division using float64 type
// int overflow can still happen
func AvgF64[T Number](n ...T) (avg float64) {
	return float64(Sum(n...)) / float64(len(n))
}

// AvgF64F64 calculates average after converting any input to float64
// to avoid integer overflows
func AvgF64F64[T Number](n ...T) (avg float64) {
	for _, v := range n {
		avg = avg + float64(v)
	}
	avg = avg / float64(len(n))
	return avg
}

func Median[T Number](n ...T) (median T) {
	// TODO probably want math-specific sort here, not to call function every time.
	sort.Slice(n, func(x, y int) bool { return n[x] < n[y] })
	if (len(n) % 2) == 0 {
		return ((n[len(n)/2-1]) + (n[(len(n) / 2)])) / 2
	} else {
		return n[len(n)/2]
	}
}

// MedianF64 calculates median with final division using float64 type
func MedianF64[T Number](n ...T) (median float64) {
	// TODO probably want math-specific sort here, not to call function every time.
	sort.Slice(n, func(x, y int) bool { return n[x] < n[y] })
	if (len(n) % 2) == 0 {
		return (float64(n[len(n)/2-1]) + float64(n[(len(n)/2)])) / 2.0
	} else {
		return float64(n[len(n)/2])
	}
}
