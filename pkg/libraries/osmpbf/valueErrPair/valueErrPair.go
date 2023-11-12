package valueErrPair

type Pair[T any] struct {
	Value T
	Err   error
}
