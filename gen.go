package goneric

func GenSlice[T any](count int, f func(idx int) T) (out []T) {
	out = make([]T, count)
	for i := 0; i < count; i++ {
		out[i] = f(i)
	}
	return out
}
func GenMap[K comparable, V any](count int, f func(idx int) (K, V)) (out map[K]V) {
	out = make(map[K]V, count)
	for i := 0; i < count; i++ {
		k, v := f(i)
		out[k] = v
	}
	return out
}
