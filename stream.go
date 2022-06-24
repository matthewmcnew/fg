package fg

import "sort"

type Stream[E any] []E

func StreamOf[E any](s []E) Stream[E] {
	return s
}

func StreamFrom[E any](e ...E) Stream[E] {
	return e
}

func (s Stream[E]) Filter(f func(e E) bool) Stream[E] {
	var filtered Stream[E]
	for _, e := range s {
		if f(e) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

func (s Stream[E]) AllMatch(f func(e E) bool) bool {
	for _, e := range s {
		if !f(e) {
			return false
		}
	}
	return true
}

func (s Stream[E]) AnyMatch(f func(e E) bool) bool {
	for _, e := range s {
		if f(e) {
			return true
		}
	}
	return false
}

func (s Stream[E]) Map(f func(e E) E) Stream[E] {
	return Map(s, f)
}

func (s Stream[E]) MapString(f func(e E) string) Stream[string] {
	return Map(s, f)
}

func (s Stream[E]) Reduce(initial E, f func(sub E, element E) E) E {
	for _, e := range s {
		initial = f(initial, e)
	}

	return initial
}

func (s Stream[E]) Sort(compare func(i E, j E) bool) Stream[E] {
	sort.Slice(s, func(i, j int) bool {
		return compare(s[i], s[j])
	})
	return s
}

func (s Stream[E]) Concat(b Stream[E]) Stream[E] {
	var concat Stream[E]
	for _, e := range s {
		concat = append(concat, e)
	}
	for _, e := range b {
		concat = append(concat, e)
	}
	return concat
}

func (s Stream[E]) Unwrap() []E {
	return s
}

func (s Stream[E]) Distinct() Stream[E] {
	distinct := Stream[E]{}
	seen := map[any]struct{}{}
	for _, e := range s {
		_, ok := seen[e]
		if ok {
			continue
		}

		seen[e] = struct{}{}
		distinct = append(distinct, e)
	}
	return distinct
}

func Map[E any, B any](stream Stream[E], o func(E) B) Stream[B] {
	var mapped = Stream[B]{}
	for _, e := range stream {
		mapped = append(mapped, o(e))
	}

	return mapped
}
