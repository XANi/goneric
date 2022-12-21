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

// ValueIndex contains value with source index for ordered operation
type ValueIndex[V any] struct {
	V   V
	IDX int
}

// Response is used to pass data and channel to return value to worker
type Response[DataT, ReturnT any] struct {
	ReturnCh chan ReturnT
	Data     DataT
}
