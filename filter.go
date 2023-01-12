package goneric

// FilterMap runs function on every element of map and adds it to result map if it returned true
func FilterMap[K comparable, V any](filterFunc func(k K, v V) (accept bool), in map[K]V) (out map[K]V) {
	out = make(map[K]V, 0)
	for k, v := range in {
		if filterFunc(k, v) {
			out[k] = v
		}
	}
	return out
}

// FilterSlice runs function on every element of slice and adds it to result slice if it returned true
func FilterSlice[V any](filterFunc func(idx int, v V) (accept bool), in []V) (out []V) {
	out = make([]V, 0)
	for idx, v := range in {
		if filterFunc(idx, v) {
			out = append(out, v)
		}
	}
	return out
}

//FilterChan filters elements going thru a channel
// close is propagated
func FilterChan[T any](filterFunc func(T) bool, in chan T) (out chan T) {
	out = make(chan T, 1)
	go func() {
		for v := range in {
			if filterFunc(v) {
				out <- v
			}
		}
		close(out)
	}()
	return out
}

//FilterChanErr filters elements going thru channel, redirecting errors to separate channel
// both channels need to be read or else it will stall
// close is propagated
func FilterChanErr[T any](filterFunc func(T) (bool, error), in chan T) (out chan T, errCh chan error) {
	out = make(chan T, 1)
	errCh = make(chan error, 1)
	go func() {
		for v := range in {
			ok, err := filterFunc(v)
			if err != nil {
				errCh <- err
			} else if ok {
				out <- v
			}
		}
		close(out)
		close(errCh)
	}()
	return out, errCh
}
