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
func FilterChan[T any](filterFunc func(T) bool, inCh chan T) (outCh chan T) {
	outCh = make(chan T, 1)
	go func() {
		for v := range inCh {
			if filterFunc(v) {
				outCh <- v
			}
		}
		close(outCh)
	}()
	return outCh
}
