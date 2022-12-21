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

// WorkerPoolDrain runs function per input without returning anything. Goroutines close on channel close.
// returns finish channel that returns single boolean true after goroutines finish
func WorkerPoolDrain[T1 any](worker func(T1), concurrency int, input chan T1) (finish chan bool) {
	finish = make(chan bool, 1)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func() {
				defer wg.Done()
				for w := range input {
					worker(w)
				}
			}()
		}
		wg.Wait()
		finish <- true
	}()
	return finish
}

// WorkerPoolAsync returns a function that adds new job to queue and returns a channel with result, and function to stop worker
func WorkerPoolAsync[T1, T2 any](worker func(T1) T2, concurrency int) (async func(T1) chan T2, stop func()) {
	inCh := make(chan Response[T1, T2], concurrency/2+1)
	finish := WorkerPoolDrain(func(in Response[T1, T2]) {
		in.ReturnCh <- worker(in.Data)
	}, concurrency, inCh)

	return func(in T1) (out chan T2) {
			ch := make(chan T2, 1)
			inCh <- Response[T1, T2]{
				ReturnCh: ch,
				Data:     in,
			}
			return ch
		}, func() {
			close(inCh)
			<-finish
		}
}
