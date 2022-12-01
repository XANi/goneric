package goneric

import (
	"sync"
)

// WorkerPool spawns `concurrency` goroutines eating from input channel and sending it to output channel
// caller should take care of closing input channel after it finished sending requests
// output channel will be closed after input is closed and all workers finish
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
	close(output)
}
