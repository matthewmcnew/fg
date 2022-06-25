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

func (c Collection[E]) FlatMap(f func(e E) []E) Collection[E] {
	return FlatMap(c, f)
}

func (c Collection[E]) MapString(f func(e E) string) Collection[string] {
	return Map(c, f)
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

func Map[E any, B any](stream Collection[E], o func(E) B) Collection[B] {
	var mapped = Collection[B]{}
	for _, e := range stream {
		mapped = append(mapped, o(e))
	}

	return mapped
}

func FlatMap[E any, B any](stream Collection[E], o func(E) []B) Collection[B] {
	var mapped = Collection[B]{}
	for _, e := range stream {
		toFlatten := o(e)
		for _, i := range toFlatten {
			mapped = append(mapped, i)
		}
	}

	return mapped
}
