package fg_test

import (
	"fmt"
	"github.com/matthewmcnew/fg"
	"reflect"
	"testing"
)

func TestFilter(t *testing.T) {
	results := fg.CollectionOf([]int{1, 2, 3}).
		Filter(func(e int) bool {
			if e >= 2 {
				return true
			}
			return false
		}).Unwrap()

	assertEqual(t, []int{2, 3}, results)
}

func TestMap(t *testing.T) {
	results := fg.CollectionOf([]int{1, 2, 3}).
		MapString(func(e int) string {
			return fmt.Sprintf("%d", e)
		}).Unwrap()

	assertEqual(t, []string{"1", "2", "3"}, results)

	results = fg.Map([]int{1, 2, 3}, func(e int) string {
		return fmt.Sprintf("%d", e)
	})

	assertEqual(t, []string{"1", "2", "3"}, results)

	intMapResults := fg.CollectionOf([]int{1, 2, 3}).
		Map(func(e int) int {
			return e + 1
		}).Unwrap()

	assertEqual(t, []int{2, 3, 4}, intMapResults)
}

func TestFlatMap(t *testing.T) {
	results := fg.FlatMap([]int{1, 2, 3}, func(e int) []string {
		return []string{"1", "2", "3"}[:e]
	}).Unwrap()

	assertEqual(t, []string{"1", "1", "2", "1", "2", "3"}, results)

	intMapResults := fg.CollectionOf([]int{1, 2, 3}).
		FlatMap(func(e int) []int {
			return []int{1, 2, 3}[:e]
		}).Unwrap()

	assertEqual(t, []int{1, 1, 2, 1, 2, 3}, intMapResults)
}

func TestReduce(t *testing.T) {
	result := fg.CollectionFrom("a", "b", "c").
		Reduce("", func(sub string, e string) string {
			return sub + e
		})

	assertEqual(t, "abc", result)
}

func TestSort(t *testing.T) {
	result := fg.CollectionOf([]int{3, 1, 2}).
		Sort(func(i int, j int) bool {
			return i < j
		}).Unwrap()

	assertEqual(t, []int{1, 2, 3}, result)
}

func TestAllMatch(t *testing.T) {
	positive := fg.CollectionFrom(3, 1, 2).
		AllMatch(func(e int) bool { return e > 0 })

	negative := fg.CollectionFrom(3, 1, 2).
		AllMatch(func(e int) bool { return e < 0 })

	assertEqual(t, true, positive)
	assertEqual(t, false, negative)
}

func TestAnyMatch(t *testing.T) {
	positive := fg.CollectionFrom(3, 1, 2).
		AnyMatch(func(e int) bool { return e == 3 })

	negative := fg.CollectionFrom(3, 1, 2).
		AnyMatch(func(e int) bool { return e == -3 })

	assertEqual(t, true, positive)
	assertEqual(t, false, negative)
}

func TestConcat(t *testing.T) {
	results := fg.CollectionFrom(1, 2).
		Concat([]int{3, 4}).
		Unwrap()

	assertEqual(t, []int{1, 2, 3, 4}, results)
}

func TestIntersect(t *testing.T) {
	results := fg.CollectionFrom(1, 2, 3).
		Intersect([]int{2, 3, 4}).
		Unwrap()

	assertEqual(t, []int{2, 3}, results)
}

func TestDistinct(t *testing.T) {
	results := fg.CollectionFrom(1, 1, 2, 2, 3, 4, 4).
		Distinct().
		Unwrap()

	assertEqual(t, []int{1, 2, 3, 4}, results)
}

func assertEqual(t *testing.T, expected any, actual any) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("not equal %s %s", expected, actual)
	}
}
