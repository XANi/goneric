package goneric

// Must takes any value and error, returns the value and panics if error happens.
func Must[T any](in T, err error) (out T) {
	if err != nil {
		panic(err)
	}
	return in
}

// IgnoreErr takes any value, and on error returns the default value
// You should probably not use it...
func IgnoreErr[T any](in T, err error) (out T) {
	if err != nil {
		var o T
		return o
	}
	return in
}
