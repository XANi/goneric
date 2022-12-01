package goneric

func FilterMap[K comparable, V any](filterFunc func(k K, v V) (accept bool), in map[K]V) (out map[K]V) {
	out = make(map[K]V, 0)
	for k, v := range in {
		if filterFunc(k, v) {
			out[k] = v
		}
	}
	return out
}
func FilterSlice[V any](filterFunc func(idx int, v V) (accept bool), in []V) (out []V) {
	out = make([]V, 0)
	for idx, v := range in {
		if filterFunc(idx, v) {
			out = append(out, v)
		}
	}
	return out
}
