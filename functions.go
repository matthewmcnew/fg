package fg

type Predicate[E any] func(e E) bool

func (f Predicate[E]) Or(orF Predicate[E]) Predicate[E] {
	return func(e E) bool {
		return f(e) || orF(e)
	}
}

func (f Predicate[E]) And(andF Predicate[E]) Predicate[E] {
	return func(e E) bool {
		return f(e) && andF(e)
	}
}

func (f Predicate[E]) Xor(xorF Predicate[E]) Predicate[E] {
	return func(e E) bool {
		return f(e) != xorF(e)
	}
}

func (f Predicate[E]) Negate() Predicate[E] {
	return func(e E) bool {
		return !f(e)
	}
}

func False[E any]() Predicate[E] {
	return func(_ E) bool { return false }
}

func True[E any]() Predicate[E] {
	return func(_ E) bool { return true }
}

type Function[E any, B any] func(e E) B

func (f Function[E, B]) Compose(c Function[E, E]) Function[E, B] {
	return func(e E) B {
		return f(c(e))
	}
}

func (f Function[E, B]) AndThen(c Function[B, B]) Function[E, B] {
	return func(e E) B {
		return c(f(e))
	}
}

func (f Function[E, B]) WithError() FunctionE[E, B] {
	return func(e E) (B, error) {
		return f(e), nil
	}
}

func Identity[E any]() Function[E, E] {
	return func(e E) E { return e }
}

func Compose[E any, B any, C any](a Function[E, B], b Function[B, C]) Function[E, C] {
	return func(e E) C {
		return b(a(e))
	}
}

type FunctionE[E any, B any] func(e E) (B, error)

func ComposeE[E any, B any, C any](a FunctionE[E, B], b FunctionE[B, C], zero C) FunctionE[E, C] {
	return func(e E) (C, error) {
		z, err := a(e)
		if err != nil {
			return zero, err
		}

		return b(z)
	}
}
