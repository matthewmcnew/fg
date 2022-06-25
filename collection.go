package fg

import "sort"

type Collection[E any] []E

func CollectionOf[E any](c []E) Collection[E] {
	return c
}

func CollectionFrom[E any](c ...E) Collection[E] {
	return c
}

func (c Collection[E]) Filter(f func(e E) bool) Collection[E] {
	var filtered Collection[E]
	for _, e := range c {
		if f(e) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

func (c Collection[E]) Contains(needle E) bool {
	return c.AnyMatch(func(e E) bool {
		return any(e) == any(needle)
	})
}

func (c Collection[E]) AllMatch(f func(e E) bool) bool {
	for _, e := range c {
		if !f(e) {
			return false
		}
	}
	return true
}

func (c Collection[E]) AnyMatch(f func(e E) bool) bool {
	for _, e := range c {
		if f(e) {
			return true
		}
	}
	return false
}

func (c Collection[E]) Map(f func(e E) E) Collection[E] {
	return Map(c, f)
}

func (c Collection[E]) MapE(f func(e E) (E, error)) (Collection[E], error) {
	return MapE(c, f)
}

func (c Collection[E]) MapString(f func(e E) string) Collection[string] {
	return Map(c, f)
}

func (c Collection[E]) MapStringE(f func(e E) (string, error)) (Collection[string], error) {
	return MapE(c, f)
}

func (c Collection[E]) FlatMap(f func(e E) []E) Collection[E] {
	return FlatMap(c, f)
}

func (c Collection[E]) FlatMapE(f func(e E) ([]E, error)) (Collection[E], error) {
	return FlatMapE(c, f)
}

func (c Collection[E]) Reduce(initial E, f func(sub E, element E) E) E {
	for _, e := range c {
		initial = f(initial, e)
	}

	return initial
}

func (c Collection[E]) Sort(compare func(i E, j E) bool) Collection[E] {
	sort.Slice(c, func(i, j int) bool {
		return compare(c[i], c[j])
	})
	return c
}

func (c Collection[E]) Concat(b Collection[E]) Collection[E] {
	var concat Collection[E]
	for _, e := range c {
		concat = append(concat, e)
	}
	for _, e := range b {
		concat = append(concat, e)
	}
	return concat
}

func (c Collection[E]) Intersect(b Collection[E]) Collection[E] {
	var intersect Collection[E]
	for _, x := range c {
		for _, y := range b {
			if any(x) == any(y) {
				intersect = append(intersect, x)
			}
		}
	}
	return intersect
}

func (c Collection[E]) ToStringMap(f func(E) string) map[string]E {
	return ToMap(c, f, Identity[E])
}

func (c Collection[E]) Unwrap() []E {
	return c
}

func (c Collection[E]) Distinct() Collection[E] {
	distinct := Collection[E]{}
	seen := map[any]struct{}{}
	for _, e := range c {
		_, ok := seen[e]
		if ok {
			continue
		}

		seen[e] = struct{}{}
		distinct = append(distinct, e)
	}
	return distinct
}

func MapE[E any, B any](collection Collection[E], o func(E) (B, error)) (Collection[B], error) {
	var mapped = Collection[B]{}
	for _, e := range collection {
		b, err := o(e)
		if err != nil {
			return nil, err
		}
		mapped = append(mapped, b)
	}

	return mapped, nil
}

func Map[E any, B any](collection Collection[E], o func(E) B) Collection[B] {
	results, err := MapE(collection, addError(o))
	if err != nil {
		panic("unexpected err")
	}
	return results
}

func FlatMapE[E any, B any](collection Collection[E], o func(E) ([]B, error)) (Collection[B], error) {
	var mapped = Collection[B]{}
	for _, e := range collection {
		toFlatten, err := o(e)
		if err != nil {
			return nil, err
		}
		for _, i := range toFlatten {
			mapped = append(mapped, i)
		}
	}

	return mapped, nil
}

func FlatMap[E any, B any](collection Collection[E], o func(E) []B) Collection[B] {
	results, err := FlatMapE(collection, addError(o))
	if err != nil {
		panic("unexpected error")
	}
	return results
}

func ToMap[E any, K comparable, V any](c Collection[E], keyMapper func(E) K, valueMapper func(E) V) map[K]V {
	m := map[K]V{}
	for _, e := range c {
		m[keyMapper(e)] = valueMapper(e)
	}

	return m
}

func addError[E any, B any](o func(E) B) func(e E) (B, error) {
	return func(e E) (B, error) {
		return o(e), nil
	}
}
