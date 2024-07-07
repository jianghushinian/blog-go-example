package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type T1 struct {
	a int8
	b string
	c bool
}

type T2 struct {
	b string
	a int8
	c bool
}

type T3 struct {
	a int8
	b string
	c struct{}
}

type T4 struct {
	a int8
	c struct{}
	b string
}

// T5 推荐写法
type T5 struct {
	c struct{}
	a int8
	b string
}

type TaskResource1 struct {
	CPU     uint8  `json:"cpu"`
	GPU     uint8  `json:"gpu"`
	GPUType string `json:"gpuType"`
	Memory  uint16 `json:"memory"`
	Storage uint64 `json:"storage"`
}

type TaskResource2 struct {
	GPUType string `json:"gpuType"`
	Storage uint64 `json:"storage"`
	Memory  uint16 `json:"memory"`
	CPU     uint8  `json:"cpu"`
	GPU     uint8  `json:"gpu"`
}

func main() {
	{
		fmt.Printf("T1 size: %d\n", unsafe.Sizeof(T1{}))
		fmt.Printf("T2 size: %d\n", unsafe.Sizeof(T2{}))
	}

	{
		fmt.Printf("int8 align: %d\n", unsafe.Alignof(int8(1)))
		fmt.Printf("bool align: %d\n", unsafe.Alignof(true))
		fmt.Printf("string align: %d\n", unsafe.Alignof("Hello World"))
		fmt.Printf("T1 align: %d\n", unsafe.Alignof(T1{}))
		fmt.Printf("T2 align: %d\n", unsafe.Alignof(T1{}))
		fmt.Printf("empty struct align: %d\n", unsafe.Alignof(struct{}{}))
		fmt.Printf("int align: %d\n", unsafe.Alignof(int(3)))
		fmt.Printf("int array align: %d\n", unsafe.Alignof([3]int{1, 2, 3}))
	}

	{
		t1 := T1{}
		fmt.Println("# T1")
		fmt.Printf("T1.a: size=%d, offset=%v, align=%d\n", unsafe.Sizeof(t1.a), unsafe.Offsetof(t1.a), unsafe.Alignof(t1.a))
		fmt.Printf("T1.b: size=%d, offset=%v, align=%d\n", unsafe.Sizeof(t1.b), unsafe.Offsetof(t1.b), unsafe.Alignof(t1.b))
		fmt.Printf("T1.c: size=%d, offset=%v, align=%d\n", unsafe.Sizeof(t1.c), unsafe.Offsetof(t1.c), unsafe.Alignof(t1.c))

		t2 := T2{}
		fmt.Println("# T2")
		fmt.Printf("T2.b: size=%d, offset=%v, align=%d\n", unsafe.Sizeof(t2.b), unsafe.Offsetof(t2.b), unsafe.Alignof(t2.b))
		fmt.Printf("T2.a: size=%d, offset=%v, align=%d\n", unsafe.Sizeof(t2.a), unsafe.Offsetof(t2.a), unsafe.Alignof(t2.a))
		fmt.Printf("T2.c: size=%d, offset=%v, align=%d\n", unsafe.Sizeof(t2.c), unsafe.Offsetof(t2.c), unsafe.Alignof(t2.c))

		// panic
		// fmt.Println(unsafe.Offsetof(int8(1)))

		t3 := T3{}
		fmt.Println("# T3")
		fmt.Printf("T3.a: size=%d, offset=%v, align=%d\n", unsafe.Sizeof(t3.a), unsafe.Offsetof(t3.a), unsafe.Alignof(t3.a))
		fmt.Printf("T3.b: size=%d, offset=%v, align=%d\n", unsafe.Sizeof(t3.b), unsafe.Offsetof(t3.b), unsafe.Alignof(t3.b))
		fmt.Printf("T3.c: size=%d, offset=%v, align=%d\n", unsafe.Sizeof(t3.c), unsafe.Offsetof(t3.c), unsafe.Alignof(t3.c))
	}

	{
		for _, T := range []any{T1{}, T2{}, T3{}, T4{}, T5{}} {
			typ := reflect.TypeOf(T)
			fmt.Printf("%s size: %d\n", typ.Name(), typ.Size())

			n := typ.NumField()
			for i := 0; i < n; i++ {
				field := typ.Field(i)
				fmt.Printf("%s.%s: size=%d, offset=%v, align=%d\n",
					typ.Name(),
					field.Name,
					field.Type.Size(),
					field.Offset,
					field.Type.Align(),
				)
			}
		}
	}

	{
		t1 := T1{}
		fmt.Printf("T1: %+v\n", t1)
		b := (*string)(unsafe.Pointer(uintptr(unsafe.Pointer(&t1)) + unsafe.Offsetof(t1.b)))
		*b = "江湖十年"
		fmt.Printf("T1: %+v\n", t1)
	}

	{
		fmt.Printf("T3 size: %d\n", unsafe.Sizeof(T3{}))
		fmt.Printf("T4 size: %d\n", unsafe.Sizeof(T4{}))
		fmt.Printf("T5 size: %d\n", unsafe.Sizeof(T5{}))
	}

	{
		fmt.Printf("TaskResource1 size: %d\n", unsafe.Sizeof(TaskResource1{}))
		fmt.Printf("TaskResource2 size: %d\n", unsafe.Sizeof(TaskResource2{}))
	}
}
