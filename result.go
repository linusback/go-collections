package collections


type Result[E any] struct {
	Val E
	Err error
}

func (r Result[T]) IsOk() bool {
	return r.Err == nil
}

func (r Result[T]) IsErr() bool {
	return r.Err != nil
}

func (r Result[T]) Error() error {
	return r.Err
}

func (o Result[T]) Unwrap() T {
	if o.IsErr() {
		panic("can't unwrap err value")
	}
	return o.Val
}
