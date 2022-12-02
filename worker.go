package goneric

import (
	"sync"
)

// WorkerPool spawns `concurrency` goroutines eating from input channel and sending it to output channel
// caller should take care of closing input channel after it finished sending requests
// output channel will be closed after input is processed and closed
func WorkerPool[T1, T2 any](input chan T1, output chan T2, worker func(T1) T2, concurrency int) {
	if concurrency < 1 {
		panic("RTFM")
	}
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for w := range input {
				output <- worker(w)
			}
		}()
	}
	wg.Wait()
}

// WorkerPoolClose spawns `concurrency` goroutines eating from input channel and sending it to output channel
// caller should take care of closing input channel after it finished sending requests
// output channel will be closed after input is processed and closed
func WorkerPoolClose[T1, T2 any](input chan T1, output chan T2, worker func(T1) T2, concurrency int) {
	WorkerPool(input, output, worker, concurrency)
	close(output)
}

// WorkerPoolBackground spawns `concurrency` goroutines eating from input channel and returns output channel with results
func WorkerPoolBackground[T1, T2 any](input chan T1, worker func(T1) T2, concurrency int) (output chan T2) {
	if concurrency < 1 {
		panic("RTFM")
	}
	output = make(chan T2, concurrency/2+1)
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for w := range input {
				output <- worker(w)
			}
		}()
	}
	return output
}

// WorkerPoolBackgroundClose spawns `concurrency` goroutines eating from input channel and returns output channel with results
// output channel will be closed after input is procesed and closed
func WorkerPoolBackgroundClose[T1, T2 any](input chan T1, worker func(T1) T2, concurrency int) (output chan T2) {
	if concurrency < 1 {
		panic("RTFM")
	}
	output = make(chan T2, concurrency/2+1)
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for w := range input {
				output <- worker(w)
			}
		}()
	}
	go func() {
		wg.Wait()
		close(output)
	}()
	return output
}

// WorkerPoolFinisher runs WorkerPool in the background and
// returns channel that returns `true` then closes when workers finish
func WorkerPoolFinisher[T1, T2 any](input chan T1, output chan T2, worker func(T1) T2, concurrency int) chan bool {
	finisher := make(chan bool, 1)
	go func() {
		WorkerPoolClose(input, output, worker, concurrency)
		finisher <- true
		close(finisher)
	}()
	return finisher
}
