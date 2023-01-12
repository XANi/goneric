package goneric

import "sync"

// All "Gen*" functions are supposed to return generated data, for ones using existing structures
// there are Type* functions

// GenSlice generates a slice of given length based on passed function.
// Function gets id of element starting from 0.
func GenSlice[T any](count int, f func(idx int) T) (out []T) {
	out = make([]T, count)
	for i := 0; i < count; i++ {
		out[i] = f(i)
	}
	return out
}

// GenMap generates a map of given size based on passed function.
// Function gets id of element starting from 0.
func GenMap[K comparable, V any](count int, f func(idx int) (K, V)) (out map[K]V) {
	out = make(map[K]V, count)
	for i := 0; i < count; i++ {
		k, v := f(i)
		out[k] = v
	}
	return out
}

// GenChan generates a channel that is fed from function results
func GenChan[T any](genFunc func() T) chan T {
	// we return channel because this is likely to be used
	// as source of data for rest of the pipeline
	ch := make(chan T, 1)
	go func() {
		for {
			ch <- genFunc()
		}
	}()
	return ch
}

// GenChanN generates channel that will run function n times and send result to channel, then close it
// Function gets id of element starting from 0.
// setting optional argument to true will close the channel after finishing
func GenChanN[T any](genFunc func(idx int) T, count int, closeOutputChan ...bool) (ch chan T) {
	// we take channel as argument instead of returning channel for more flexibility
	ch = make(chan T, 1)
	go func() {
		for i := 0; i < count; i++ {
			ch <- genFunc(i)
		}
		if len(closeOutputChan) > 0 && closeOutputChan[0] {
			close(ch)
		}
	}()
	return ch
}

// GenChanCloser generates a channel that is fed from function results. Running closer func will stop it
// the channel is of size one which means new value will be immediately generated without waiting
// for consumer
// That also causes few messages to be generated (theoretically at most 2) after channel is closed
// If you need synchronous close here you're probably doing something wrong
// calling closer with closer(true) will close the channel, by default it is not getting closed

func GenChanCloser[T any](genFunc func() T) (out chan T, closer func(closeChannel ...bool)) {
	ch := make(chan T, 1)
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
	return ch, closer
}

// GenSliceToChan returns channel with background goroutine feeding it data from slice
func GenSliceToChan[T any](in []T, closeOutputChan ...bool) (out chan T) {
	out = make(chan T, 1)
	go func() {
		for _, v := range in {
			out <- v
		}
		if len(closeOutputChan) > 0 && closeOutputChan[0] {
			close(out)
		}
	}()
	return out
}
