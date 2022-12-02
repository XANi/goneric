package goneric

type ErrSkip struct{}

func (v ErrSkip) Error() string {
	return "fake error marking entries to skip"
}

type Number interface {
	~int | ~int64 | ~int32 | ~int16 | ~int8 |
		~uint64 | ~uint32 | ~uint16 | ~uint8 |
		~float64 | ~float32
}
