package fg

import (
	"fmt"
	"sort"
)

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

func (c Collection[E]) FindFirst(f func(e E) bool, empty E) (E, error) {
	for _, e := range c {
		if f(e) {
			return e, nil
		}
	}

	return empty, fmt.Errorf("could not find element")
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
	return Reduce(c, initial, f)
}

func (c Collection[E]) Sort(compare func(i E, j E) bool) Collection[E] {
	copy := make([]E, len(c))
	for i, e := range c {
		copy[i] = e
	}

	sort.Slice(copy, func(i, j int) bool {
		return compare(copy[i], copy[j])
	})
	return copy
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
	return ToMap(c, f, Identity[E]())
}

func (c Collection[E]) Unwrap() []E {
	return c
}

func (c Collection[E]) Distinct() Collection[E] {
	distinct := make([]E, 0, len(c))
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
	var mapped = make([]B, 0, len(collection))
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
	results, err := MapE(collection, Function[E, B](o).WithError())
	if err != nil {
		panic("unexpected err")
	}
	return results
}

func FlatMapE[E any, B any](collection Collection[E], o func(E) ([]B, error)) (Collection[B], error) {
	var mapped = make([][]B, 0, len(collection))
	for _, e := range collection {
		b, err := o(e)
		if err != nil {
			return nil, err
		}
		mapped = append(mapped, b)
	}

	return Flatten(mapped), nil
}

func Reduce[E any, B any](collection Collection[E], initial B, f func(sub B, element E) B) B {
	for _, e := range collection {
		initial = f(initial, e)
	}

	return initial
}

func FlatMap[E any, B any](collection Collection[E], o func(E) []B) Collection[B] {
	results, err := FlatMapE(collection, Function[E, []B](o).WithError())
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

func Flatten[E any](s ...[][]E) []E {
	size := 0
	for _, e := range s {
		size += len(e)
	}
	flattened := make([]E, 0, size)
	for _, t := range s {
		for _, e := range t {
			flattened = append(flattened, e...)
		}
	}
	return flattened
}
