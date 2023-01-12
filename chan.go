package goneric

import (
	"time"
)

// ChanToSlice loads channel messages to slice until channel is closed
func ChanToSlice[T any](inCh chan T) []T {
	s := make([]T, 0)
	for v := range inCh {
		s = append(s, v)
	}
	return s
}

// ChanToSliceN loads up to n elements from to slice to channel
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

// ChanToSliceNTimeout loads up to n elements from to slice to channel or up until timeout expires
func ChanToSliceNTimeout[T any](inCh chan T, n int, timeout time.Duration) []T {
	s := make([]T, 0)
	idx := 0
	t := time.After(timeout)
	for {
		select {
		case <-t:
			return s
		case v := <-inCh:
			s = append(s, v)
			idx++
			if idx >= n {
				return s
			}
		}
	}

}

// SliceToChan feeds slice to channel
func SliceToChan[T any](in []T, out chan T, closeOutputChan ...bool) {
	go func() {
		for _, v := range in {
			out <- v
		}
		if len(closeOutputChan) > 0 && closeOutputChan[0] {
			close(out)
		}
	}()
	return
}
