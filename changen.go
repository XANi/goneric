package goneric

import "sync"

// ChanGen sends function output to provided channel
// The generator goroutine runs forever, use ChanGenCloser if you need to be able to stop it
func ChanGen[T any](genFunc func() T, ch chan T) {
	go func() {
		for {
			ch <- genFunc()
		}
	}()
}

// ChanGenN runs function n times and sends the result to the provided channel
// resulting channel will be sent `true` when the function finishes last send
// Function gets id of element starting from 0.
// setting optional argument to true will close the channel after finishing
func ChanGenN[T any](count int, genFunc func(idx int) T, out chan T, closeOutputChan ...bool) (finished chan bool) {
	// we take channel as argument instead of returning channel for more flexibility
	finished = make(chan bool, 1)
	go func() {
		for i := 0; i < count; i++ {
			out <- genFunc(i)
		}
		if len(closeOutputChan) > 0 && closeOutputChan[0] {
			close(out)
		}
		finished <- true

	}()
	return finished
}

// ChanGenCloser feeds function results to the provided channel. Running closer func will stop it.
// Stopping is asynchronous and number of events in flight depends on channel size so don't rely on stopping after exact number of calls,
// If you need synchronous close here you're probably doing something wrong
// Calling closer with closer(true) will close the channel after the generator stops, by default it is left open.
// Calling closer more than once is safe, calls after the first one are ignored.
func ChanGenCloser[T any](genFunc func() T, ch chan T) (closer func(closeChannel ...bool)) {
	stop := make(chan struct{})
	// written only before close(stop), read only after stop is closed,
	// so channel close ordering makes it race-free without atomics
	closeOut := false
	stopOnce := sync.Once{}
	closer = func(closeChannel ...bool) {
		stopOnce.Do(func() {
			if len(closeChannel) > 0 && closeChannel[0] {
				closeOut = true
			}
			close(stop)
		})
	}

	go func() {
		for {
			select {
			case <-stop:
				if closeOut {
					close(ch)
				}
				return
			default:
				ch <- genFunc()
			}
		}
	}()
	return closer
}
