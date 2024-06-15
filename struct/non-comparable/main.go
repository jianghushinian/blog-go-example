package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type Normal struct {
	a string
	B int
}

type NoCompare struct {
	a string
	B map[string]int
}

type DisallowCompare struct {
	_ [0]func()
}

// Value ref: https://github.com/golang/go/blob/master/src/log/slog/value.go#L21
type Value struct {
	_ [0]func() // disallow ==
	// num holds the value for Kinds Int64, Uint64, Float64, Bool and Duration,
	// the string length for KindString, and nanoseconds since the epoch for KindTime.
	num uint64
	// If any is of type Kind, then the value is in num as described above.
	// If any is of type *time.Location, then the Kind is Time and time.Time value
	// can be constructed from the Unix nanos in num and the location (monotonic time
	// is not preserved).
	// If any is of type stringptr, then the Kind is String and the string value
	// consists of the length in num and the pointer in any.
	// Otherwise, the Kind is Any and any is the value.
	// (This implies that Attrs cannot store values of type Kind, *time.Location
	// or stringptr.)
	any any
}

func main() {
	// 可比较
	{
		n1 := Normal{
			a: "a",
			B: 10,
		}
		n2 := Normal{
			a: "a",
			B: 10,
		}
		n3 := Normal{
			a: "b",
			B: 20,
		}

		fmt.Println(n1 == n2) // true
		fmt.Println(n1 == n3) // false
	}

	// 不可比较
	{
		n1 := NoCompare{
			a: "a",
			B: map[string]int{
				"a": 10,
			},
		}
		n2 := NoCompare{
			a: "a",
			B: map[string]int{
				"a": 10,
			},
		}

		// invalid operation: n1 == n2 (struct containing map[string]int cannot be compared)
		// fmt.Println(n1 == n2)

		fmt.Println(reflect.DeepEqual(n1, n2)) // true

		fmt.Println(unsafe.Sizeof(n1), unsafe.Sizeof(n1)) // 24 24
	}

	// 零成本不可比较
	{
		d1 := DisallowCompare{}
		d2 := DisallowCompare{}

		// invalid operation: d1 == d2 (struct containing [0]func() cannot be compared)
		// fmt.Println(d1 == d2)

		fmt.Println(unsafe.Sizeof(d1), unsafe.Sizeof(d2)) // 0 0

		fmt.Println(reflect.DeepEqual(d1, d2)) // true
	}

	// slog.Value 不可比较
	{
		v1 := Value{
			num: 1,
			any: 2,
		}
		v2 := Value{
			num: 1,
			any: 2,
		}

		// invalid operation: v1 == v2 (struct containing [0]func() cannot be compared)
		// fmt.Println(v1 == v2)

		fmt.Println(unsafe.Sizeof(v1), unsafe.Sizeof(v2)) // 24 24
		fmt.Println(reflect.DeepEqual(v1, v2))            // true
	}
}
