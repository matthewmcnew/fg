# fg

Simple Type Safe Functional Programming in Golang with Generics! 

### Collections

```golang
type Collection[E any] []E
```

The Collection type is a slice of elements which allows functional operations to be performed. Many operations can be done in series with method chaining. 

Some helpers functions are provided to construct collections

```golang
var collection fg.Collection[string] = fg.CollectionOf([]string{"1", "2"})
var collection fg.Collection[string] = fg.CollectionFrom("1", "2")
```

#### Filter

Filter lets you filter a collection with a provided predicate function

```go
    func greaterThan2(e int) bool{
        if e > 2 {
            return true
        }
        return false
    }

    fg.CollectionFrom(1, 2, 3).Filter(greaterThan2) // []int{3}
```

#### FindFirst

FindFirst lets you find the first element that matches 

```go
    func greaterThan2(e int) bool{
        if e > 2 {
        return true
        }
        return false
    }

    fg.CollectionFrom(1, 2, 3).FindFirst(greaterThan2) // 3
```

#### Contains

Contains returns true if a list contains an element

```go
    func greaterThan2(e int) bool{
        if e > 2 {
            return true
        }
        return false
    }

    fg.CollectionFrom(1, 2, 3).Contains(greaterThan2) // true
    fg.CollectionFrom(1, 2).Contains(greaterThan2) // false
```

#### AllMatch

Contains returns true if all elements match

```go
    func greaterThan2(e int) bool{
        if e > 2 {
            return true
        }
        return false
    }

    fg.CollectionFrom(3, 4, 5).AllMatch(greaterThan2) // true
    fg.CollectionFrom(3, 2).AllMatch(greaterThan2) // false
```

#### Map

Map allows you to map all elements in a collection with a mapping function

```go
    func addPrefix(e string) string{
        return "prefix-" + e
    }

    fg.CollectionFrom("a", "b").Map(addPrefix) // []string{"prefix-a", "prefix-b"}
```

Golang does not allow methods to have type parameters. Mapping a collection to another type is possible with the fg.Map method. 

```go
    func addPrefix(e int) string{
        return fmt.Sprintf("prefix-%d", e)
    }

    fg.Map([]int{1, 2}, addPrefix) // []string{"prefix-1", "prefix-2"}
```

#### MapE

MapE allows map to be called with a mapping function that returns an error.

```go
    func addPrefix(e string) (string, error) {
		if strings.Contains(e, "badword") {
			return "", errors.New("we don't accept bad words around here")
        }   
		
        return "prefix-" + e
    }

    fg.CollectionFrom("a", "b", "badword").MapE(addPrefix) // error
```

#### FlatMap

FlatMap allows you to map all elements in a collection with a mapping function that returns a slice

```go
    func split(e string) []string{
		return []strings.Split(e, ".")
    }

    fg.CollectionFrom("a.b", "c.d").FlatMap(split) // []string{"a", "b", "c", "d"}
```

#### Reduce

Reduce allows you to reduce a collection by running each result through the provided combiner function.

```go
    func combine(sub string, e string) string {
        return sub + e
    }

    fg.CollectionFrom("b", "c").Reduce("a", combine) // abc
```

#### ToMap

To Map allows you to convert a collection to a Map with a key mapping function. 

```go
    func intToString(e int) string {
        return fmt.Sprintf("%d", e)
    }

    fg.CollectionFrom(1, 2).ToStringMap(intToString) // map[string]int{"1": 1, "2": 2}
```

#### Chaining

Collection Methods can be chained to support multiple operations in a sequence

```go
    positiveSum := fg.CollectionFrom(-1, 2, 4).
                Filter(func(e int) bool {
                    return e >= 0 
                }).Reduce(0, func(sum, e int) int {
                    return sum + e
                })
                

```

### Function Types

Function types are provided to help compose and combine functions

#### Predicate

A predicate is a function that returns true or false from a given type. The methods `negate`, `or`, `and`, and `xor` are available.    

```go
    var containsA fg.Predicate[string] = func(s string) bool {
	    	return strings.Contains(s, "A")
    }

    var containsB fg.Predicate[string] = func(s string) bool {
        return strings.Contains(s, "B")
    }

    containsA("ABC") //true
	containsA.Negate()("ABC") // false
    containsA.Or(containsB)("BCD") // true
    containsA.And(containsB)("BCD") // false
    containsA.Xor(containsB)("ABC") // false
```

#### Compose

Compose allows functions to be composed 

```go
    func firstElement(s []string) string {
        return s[0]
    }

    func lengthOfString(s string) int {
        return len(s)
    }

	lenOfFirstElement := fg.Compose(firstElement, lengthOfString)

    lenOfFirstElement([]string{"ab", "bcd"}) //2
```


