package goneric

import (
	"sync"
)

// ParallelMapSlice takes variadic argument and runs all of them thru function in parallel, up to `concurrency` goroutines
// Order of elements in slice is kept
func ParallelMap[T1, T2 any](mapFunc func(T1) T2, concurrency int, slice ...T1) []T2 {
	return ParallelMapSlice(mapFunc, concurrency, slice)
}

//ParallelMapSlice takes slice and runs it thru function in parallel, up to `concurrency` goroutines
// Order of elements in slice is kept
func ParallelMapSlice[T1, T2 any](mapFunc func(T1) T2, concurrency int, slice []T1) []T2 {
	out := make([]T2, len(slice))
	inCh := make(chan struct {
		V   T1
		idx int
	}, concurrency/2+1) // TODO check whether size matters here for anything
	outCh := make(chan struct {
		V   T2
		idx int
	}, concurrency/2+1) // we're just guessing here that having open slot in channel is good for performance, test that

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for v := range outCh {
			out[v.idx] = v.V
		}
	}()
	// Go Generic Gymnastics.
	// We can't generate named types here (as of 1.19)
	// so we have to do with anonymous ones
	go func() {
		defer wg.Done()
		WorkerPool(
			inCh,
			outCh,
			func(i struct {
				V   T1
				idx int
			}) struct {
				V   T2
				idx int
			} {
				return struct {
					V   T2
					idx int
				}{
					V:   mapFunc(i.V),
					idx: i.idx,
				}

			},
			concurrency)
	}()
	for idx, v := range slice {
		inCh <- struct {
			V   T1
			idx int
		}{V: v, idx: idx}
	}
	close(inCh)
	wg.Wait()
	return out
}

// ParallelMapSliceChannel feeds slice to function in parallel and returns channels with function output
// channel is closed when function finishes. Caller should close input channel when it finishes sending
// or else it will leak goroutines
func ParallelMapSliceChannel[T1, T2 any](mapFunc func(T1) T2, concurrency int, slice []T1) chan T2 {
	in := make(chan T1, 1)
	out := make(chan T2, concurrency/2+1)
	go func() {
		for _, v := range slice {
			in <- v
		}
		close(in)
	}()
	go func() {
		WorkerPool(in, out, mapFunc, concurrency)
	}()

	return out
}

// ParallelMapSliceChannelFinisher feeds slice to function in parallel and returns channels with function output
// channel is closed when function finishes. Caller should close input channel when it finishes sending
// or else it will leak goroutines
// Second channel will return true (and then be closed) when the worker finishes parsing
func ParallelMapSliceChannelFinisher[T1, T2 any](mapFunc func(T1) T2, concurrency int, slice []T1) (chan T2, chan bool) {
	in := make(chan T1, 1)
	out := make(chan T2, concurrency/2+1)
	finisher := make(chan bool, 1)
	go func() {
		for _, v := range slice {
			in <- v
		}
		close(in)
	}()
	go func() {
		WorkerPool(in, out, mapFunc, concurrency)
		finisher <- true
		close(finisher)
	}()

	return out, finisher
}
