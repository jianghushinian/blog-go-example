package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string
	// string
}

func main() {
	// 内置类型
	{
		age := 20

		val := reflect.ValueOf(age)
		typ := reflect.TypeOf(age)
		fmt.Println(val, typ) // 输出：20 int
	}

	// 自定义结构体类型
	{
		user := User{
			Name:  "江湖十年",
			Age:   20,
			Email: "jianghushinian007@outlook.com",
		}

		val := reflect.ValueOf(user)
		typ := reflect.TypeOf(user)
		fmt.Println(val, typ) // 输出：{江湖十年 20 jianghushinian007@outlook.com} main.User
	}

	// reflect.Value 常用方法
	{
		// 实例化 User 结构体指针
		user := &User{
			Name:  "江湖十年",
			Age:   20,
			Email: "jianghushinian007@outlook.com",
		}

		// 注意这里传递的是指针类型
		kind := reflect.ValueOf(user).Kind()
		fmt.Println(kind) // 输出：ptr
		kind = reflect.ValueOf(*user).Kind()
		fmt.Println(kind)                          // 输出：struct
		kind = reflect.ValueOf(user).Elem().Kind() // 输出：指针类型需要使用 Elem 方法获取指针指向的值
		fmt.Println(kind)                          // 输出：struct

		// 以下二者等价
		tpy := reflect.ValueOf(user).Type()
		fmt.Println(tpy) // 输出：*main.User
		tpy1 := reflect.TypeOf(user)
		fmt.Println(tpy1)                         // 输出：*main.User
		fmt.Println(reflect.DeepEqual(tpy, tpy1)) // 输出：true

		// 获取结构体值字段
		nameField := reflect.ValueOf(user).Elem().FieldByName("Name")
		ageField := reflect.ValueOf(user).Elem().FieldByIndex([]int{1}) // FieldByIndex 内部调用的也是 Field 方法
		emailField := reflect.ValueOf(user).Elem().Field(2)             // 获取结构体第 3 个字段

		// 结构体字段总个数
		numField := reflect.ValueOf(*user).NumField()
		fmt.Println(numField) // 输出：3

		// 获取结构体字段值
		fmt.Println(nameField.String())  // 输出："江湖十年"
		fmt.Println(ageField.Int())      // 输出：20
		fmt.Println(emailField.String()) // 输出：jianghushinian007@outlook.com

		// 设置结构体字段值，只有传递给 reflect.ValueOf 的值是指针时才可以这样做
		nameField.SetString("jianghushinian")             // 设置 Name 字段的值
		ageField.SetInt(18)                               // 设置 Age 字段的值
		emailField.SetString("jianghushinian007@163.com") // 设置 Email 字段的值
		fmt.Println(user)                                 // 输出：&{jianghushinian 18 jianghushinian007@163.com}
	}

	// reflect.Type 常用方法
	{
		// 实例化 User 结构体指针
		user := &User{
			Name:  "江湖十年",
			Age:   20,
			Email: "jianghushinian007@outlook.com",
		}

		// reflect.Type 也有 Kind 方法，与 reflect.Value 对应
		kind := reflect.TypeOf(user).Kind()
		fmt.Println(kind) // 输出：ptr
		kind = reflect.TypeOf(*user).Kind()
		fmt.Println(kind)                         // 输出：struct
		kind = reflect.TypeOf(user).Elem().Kind() // 输出：指针类型需要使用 Elem 方法获取指针指向的值
		fmt.Println(kind)                         // 输出：struct

		// 获取结构体类型字段
		// reflect.Type 也有这几个方法，与 reflect.Value 对应
		nameField, _ := reflect.TypeOf(user).Elem().FieldByName("Name")
		ageField := reflect.TypeOf(user).Elem().FieldByIndex([]int{1})
		emailField := reflect.TypeOf(user).Elem().Field(2)
		fmt.Printf("%+v\n", nameField)  // 输出：{Name:Name PkgPath: Type:string Tag:json:"name" Offset:0 Index:[0] Anonymous:false}
		fmt.Printf("%+v\n", ageField)   // 输出：{Name:Age PkgPath: Type:int Tag:json:"age" Offset:16 Index:[1] Anonymous:false}
		fmt.Printf("%+v\n", emailField) // 输出：{Name:Email PkgPath: Type:string Tag: Offset:24 Index:[2] Anonymous:false}

		// 获取字段标签
		tag := nameField.Tag
		fmt.Printf("%+v\n", tag)             // 输出：json:"name"
		fmt.Printf("%+v\n", tag.Get("json")) // 输出：name
	}
}
