package goneric

// Map maps the list of variadic(...) values via function.
// It is provided as convenience, MapSlice() should be used when you have incoming slice
func Map[T1, T2 any](mapFunc func(v T1) T2, slice ...T1) []T2 {
	out := make([]T2, len(slice))
	for idx, v := range slice {
		out[idx] = mapFunc(v)
	}
	return out
}

// MapSlice maps slice using provided function
func MapSlice[T1, T2 any](mapFunc func(v T1) T2, slice []T1) []T2 {
	out := make([]T2, len(slice))
	for idx, v := range slice {
		out[idx] = mapFunc(v)
	}
	return out
}

// MapErr maps the list of variadic(...) values via function. and returns on first error
// Returns slice with elements that didn't return error before the failure
// so index of the first element in error is essentially `slice[len(out)]`
func MapErr[T1, T2 any](mapFunc func(v T1) (T2, error), slice ...T1) (out []T2, err error) {
	out = make([]T2, len(slice))
	for idx, v := range slice {
		out[idx], err = mapFunc(v)
		if err != nil {
			return out[:idx], err
		}
	}
	return out, nil
}

// MapSliceErr maps slice using provided function, returning first error
// Returns slice with elements that didn't return error
func MapSliceErr[T1, T2 any](mapFunc func(v T1) (T2, error), slice []T1) (out []T2, err error) {
	out = make([]T2, len(slice))
	for idx, v := range slice {
		out[idx], err = mapFunc(v)
		if err != nil {
			return out[:idx], err
		}
	}
	return out, nil
}

// MapSliceSkip maps slice using provided function, allowing to skip entries by returning true
// Returns slice with elements that didn't return true
func MapSliceSkip[T1, T2 any](mapFunc func(v T1) (T2, bool), slice []T1) (out []T2) {
	out = make([]T2, 0)
	for _, v := range slice {
		r, skip := mapFunc(v)
		if !skip {
			out = append(out, r)
		}
	}
	return out
}

// MapSliceErrSkip maps slice using provided function, allowing to skip entries by returning ErrSkip
// Returns on first non `ErrSkip` error
func MapSliceErrSkip[T1, T2 any](mapFunc func(v T1) (T2, error), slice []T1) (out []T2, err error) {
	out = make([]T2, 0)
	for _, v := range slice {
		r, err := mapFunc(v)
		switch err.(type) {
		case ErrSkip:
			continue
		case nil:
			out = append(out, r)
		default:
			return out, err
		}
	}
	return out, nil
}

// MapSliceKey returns keys of a map in slice
func MapSliceKey[K comparable, V any](in map[K]V) (out []K) {
	out = make([]K, len(in))
	i := 0
	for k := range in {
		out[i] = k
		i++
	}
	return out
}

// MapSliceValue returns values of a map in slice
func MapSliceValue[K comparable, V any](in map[K]V) (out []V) {
	out = make([]V, len(in))
	i := 0
	for _, v := range in {
		out[i] = v
		i++
	}
	return out
}

// MapToSlice converts map into slice via specified function
func MapToSlice[K1 comparable, V any, V2 any](f func(k K1, v V) V2, in map[K1]V) (out []V2) {
	out = make([]V2, len(in))
	i := 0
	for k, v := range in {
		out[i] = f(k, v)
		i++
	}
	return out
}

// MapMap runs every map element thru function that returns new key and value, and returns that in another map. Types can vary between in and out.
func MapMap[K1, K2 comparable, V1, V2 any](f func(k K1, v V1) (K2, V2), in map[K1]V1) (out map[K2]V2) {
	out = make(map[K2]V2, len(in))
	for kin, vin := range in {
		kout, vout := f(kin, vin)
		out[kout] = vout
	}
	return out
}

// MapMapInplace runs every map element thru function that returns new key and value, and puts it into existing map. Types can vary between in and out.
func MapMapInplace[K1, K2 comparable, V1, V2 any](f func(k K1, v V1) (K2, V2), in map[K1]V1, out map[K2]V2) {
	for kin, vin := range in {
		kout, vout := f(kin, vin)
		out[kout] = vout
	}
}
