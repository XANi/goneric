package goneric

type ErrSkip struct{}

func (v ErrSkip) Error() string {
	return "fake error marking entries to skip"
}
