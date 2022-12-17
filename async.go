package goneric

// Async runs a function in goroutine and returns pipe with result
func Async[T1 any](f func() T1) chan T1 {
	out := make(chan T1, 1)
	go func() {
		out <- f()
	}()
	return out
}

// AsyncPipe runs a function in goroutine from input channel and returns pipe with result
func AsyncPipe[T1, T2 any](in chan T1, f func(T1) T2) chan T2 {
	out := make(chan T2, 1)
	go func() {
		out <- f(<-in)
	}()
	return out
}

// AsyncOut takes value and feeds it to function returning asynchronously to channel
func AsyncOut[T1, T2 any](in chan T1, f func(T1) T2, out chan T2) {
	go func() {
		out <- f(<-in)
	}()
}

//AsyncIn turns value into channel with value
func AsyncIn[T1 any](in T1) (out chan T1) {
	out = make(chan T1, 1)
	out <- in
	return out
}
