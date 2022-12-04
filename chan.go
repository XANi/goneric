package goneric

import "sync"

//ChanGen generates a channel that is fed from function results
func ChanGen[T any](genFunc func() T) chan T {
	ch := make(chan T, 1)
	go func() {
		for {
			ch <- genFunc()
		}
	}()
	return ch
}

// ChanGenCloser generates a channel that is fed from function results. Running closer func will stop it
// the channel is of size one which means new value will be immediately generated without waiting
// for consumer
// That also causes few messages to be generated (theoretically at most 2) after channel is closed
// If you need synchronous close here you're probably doing something wrong
func ChanGenCloser[T any](genFunc func() T) (out chan T, closer func()) {
	ch := make(chan T, 1)
	m := sync.RWMutex{}
	// TODO check performance, check whether there is a faster way
	closer = func() {
		m.Lock()
	}

	go func() {
		for {
			if m.TryRLock() {
				m.RUnlock()
				ch <- genFunc()
			} else {
				close(ch)
				return
			}
		}
	}()
	return ch, closer
}

//ChanToSlice loads channel messages to slice until channel is closed
func ChanToSlice[T any](inCh chan T) []T {
	s := make([]T, 0)
	for v := range inCh {
		s = append(s, v)
	}
	return s
}

//ChanToSliceN loads up to n elements from to slice to channel
func ChanToSliceN[T any](inCh chan T, n int) []T {
	s := make([]T, 0)
	idx := 0
	for v := range inCh {
		s = append(s, v)
		idx++
		if idx >= n {
			return s
		}
	}
	return s
}

// SliceToChan returns channel with background goroutine feeding it data from slice
func SliceToChan[T any](in []T) chan T {
	out := make(chan T, 1)
	go func() {
		for _, v := range in {
			out <- v
		}
	}()
	return out
}

// SliceToChan returns channel with background goroutine feeding it data from slice
// then closes the channel
func SliceToChanClose[T any](in []T) chan T {
	outCh := make(chan T, 1)
	go func() {
		for _, v := range in {
			outCh <- v
		}
		close(outCh)
	}()
	return outCh
}
