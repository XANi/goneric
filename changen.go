package goneric

import "sync"

// ChanGen sends function output to provided channel
func ChanGen[T any](genFunc func() T, ch chan T) {
	// we return channel because this is likely to be used
	// as source of data for rest of the pipeline
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

// ChanGenCloser generates a channel that is fed from function results. Running closer func will stop it.
// Closing is asynchronous and number of events in flight depends on channel size so don't rely on stooping after exact number of calls,
// If you need synchronous close here you're probably doing something wrong
func ChanGenCloser[T any](genFunc func() T, ch chan T) (closer func(closeChannel ...bool)) {
	m := sync.RWMutex{}
	cl := false
	// TODO check performance, check whether there is a faster way
	closer = func(close ...bool) {
		if len(close) > 0 && close[0] {
			cl = true
		}
		m.Lock()
	}

	go func() {
		for {
			if m.TryRLock() {
				m.RUnlock()
				ch <- genFunc()
			} else {
				if cl {
					close(ch)
				}
				return
			}
		}
	}()
	return closer
}
