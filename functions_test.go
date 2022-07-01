package fg_test

import (
	"errors"
	"github.com/matthewmcnew/fg"
	"strings"
	"testing"
)

func TestPredicate(t *testing.T) {
	assertEqual(t, true, fg.True[string]()(""))
	assertEqual(t, false, fg.False[string]()(""))

	assertEqual(t, false, fg.True[string]().Negate()(""))
	assertEqual(t, true, fg.False[string]().Negate()(""))

	var contains = func(subString string) fg.Predicate[string] {
		return func(s string) bool {
			return strings.Contains(s, subString)
		}
	}

	assertEqual(t, false, contains("z")("abc"))
	assertEqual(t, true, contains("a")("abc"))

	assertEqual(t, false, contains("z").Or(contains("q"))("abc"))
	assertEqual(t, true, contains("z").Or(contains("c"))("abc"))

	assertEqual(t, false, contains("a").And(contains("z"))("abc"))
	assertEqual(t, true, contains("a").And(contains("b"))("abc"))

	assertEqual(t, false, contains("a").Xor(contains("b"))("abc"))
	assertEqual(t, true, contains("a").Xor(contains("z"))("abc"))

	assertEqual(t, true, contains("a").Xor(contains("z"))("abc"))
}

func TestFunction(t *testing.T) {
	composedFunc := fg.Function[string, int](func(e string) int {
		return len(e)
	}).Compose(func(e string) string {
		return "my" + e
	})
	assertEqual(t, 3, composedFunc("a"))

	andThenFunc := fg.Function[string, int](func(e string) int {
		return len(e)
	}).AndThen(func(e int) int {
		return e + 1
	})
	assertEqual(t, 2, andThenFunc("a"))

	assertEqual(t, 3, fg.Identity[int]()(3))
	assertEqual(t, "3", fg.Identity[string]()("3"))
}

func TestCompose(t *testing.T) {
	lenOfFirstElement := fg.Compose[[]string, string, int](func(l []string) string {
		return l[0]
	},
		func(s string) int {
			return len(s)
		})

	assertEqual(t, 3, lenOfFirstElement([]string{"abc", "bd"}))
}

func TestComposeE(t *testing.T) {
	firstElement := func(l []string) (string, error) {
		if len(l) == 0 {
			return "", errors.New("issue")
		}
		return l[0], nil
	}
	lenOfElement := func(s string) (int, error) {
		return len(s), nil
	}

	lenOfFirstElement := fg.ComposeE(firstElement, lenOfElement, 0)

	element, err := lenOfFirstElement([]string{"abc", "bd"})
	assertNil(t, err)
	assertEqual(t, 3, element)
}
