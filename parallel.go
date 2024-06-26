package goneric

import (
	"sync"
)

// ParallelMapSlice takes variadic argument and runs all of them thru function in parallel, up to `concurrency` goroutines
// Order of elements in slice is kept
func ParallelMap[T1, T2 any](mapFunc func(T1) T2, concurrency int, slice ...T1) []T2 {
	return ParallelMapSlice(mapFunc, concurrency, slice)
}

// ParallelMapSlice takes slice and runs it thru function in parallel, up to `concurrency` goroutines
// Order of elements in slice is kept
func ParallelMapSlice[T1, T2 any](mapFunc func(T1) T2, concurrency int, slice []T1) []T2 {
	out := make([]T2, len(slice))
	inCh := make(chan ValueIndex[T1], concurrency/2+1)  // TODO check whether size matters here for anything
	outCh := make(chan ValueIndex[T2], concurrency/2+1) // we're just guessing here that having open slot in channel is good for performance, test that

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for v := range outCh {
			out[v.IDX] = v.V
		}
	}()
	go func() {
		defer wg.Done()
		WorkerPool(
			inCh,
			outCh,
			func(i ValueIndex[T1]) ValueIndex[T2] {
				return ValueIndex[T2]{
					V:   mapFunc(i.V),
					IDX: i.IDX,
				}

			},
			concurrency, true)
	}()
	for idx, v := range slice {
		inCh <- ValueIndex[T1]{V: v, IDX: idx}
	}
	close(inCh)
	wg.Wait()
	return out
}

// ParallelMapMap takes map and runs each element thru function in parallel, storing result in a map
func ParallelMapMap[K1, K2 comparable, V1, V2 any](
	mapFunc func(k K1, v V1) (K2, V2),
	concurrency int,
	in map[K1]V1,
) map[K2]V2 {
	out := make(map[K2]V2, len(in))
	wg := sync.WaitGroup{}
	inCh := make(chan KeyValue[K1, V1], concurrency/2+1)  // TODO check whether size matters here for anything
	outCh := make(chan KeyValue[K2, V2], concurrency/2+1) // we're just guessing here that having open slot in channel is good for performance, test that
	wg.Add(2)
	go func() {
		defer wg.Done()
		for v := range outCh {
			out[v.K] = v.V
		}
	}()
	go func() {
		defer wg.Done()
		WorkerPool(
			inCh,
			outCh,
			func(i KeyValue[K1, V1]) KeyValue[K2, V2] {
				k, v := mapFunc(i.K, i.V)
				return KeyValue[K2, V2]{
					K: k,
					V: v,
				}

			},
			concurrency, true)
	}()
	for k, v := range in {
		inCh <- KeyValue[K1, V1]{K: k, V: v}
	}
	close(inCh)
	wg.Wait()
	return out
}

// ParallelMapSliceChan feeds slice to function in parallel and returns channels with function output
// channel is closed when function finishes. Caller should close input channel when it finishes sending
// or else it will leak goroutines
func ParallelMapSliceChan[T1, T2 any](mapFunc func(T1) T2, concurrency int, slice []T1) chan T2 {
	in := make(chan T1, 1)
	out := make(chan T2, concurrency/2+1)
	go func() {
		for _, v := range slice {
			in <- v
		}
		close(in)
	}()
	go func() {
		WorkerPool(in, out, mapFunc, concurrency, true)
	}()

	return out
}

// ParallelMapSliceChanFinisher feeds slice to function in parallel and returns channels with function output
// channel is closed when function finishes. Caller should close input channel when it finishes sending
// or else it will leak goroutines
// Second channel will return true (and then be closed) when the worker finishes parsing
func ParallelMapSliceChanFinisher[T1, T2 any](mapFunc func(T1) T2, concurrency int, slice []T1) (chan T2, chan bool) {
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
		WorkerPool(in, out, mapFunc, concurrency, true)
		finisher <- true
		close(finisher)
	}()

	return out, finisher
}
