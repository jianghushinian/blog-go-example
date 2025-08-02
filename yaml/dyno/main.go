package main

import (
	"fmt"

	"github.com/icza/dyno"
)

func main() {
	y := map[string]interface{}{
		"object": map[interface{}]interface{}{
			"a": 1,
			"array": []interface{}{
				map[string]interface{}{
					"null_value": interface{}(nil),
				},
				map[string]interface{}{
					"boolean": true,
				},
				map[string]interface{}{
					"integer": 1,
				},
			},
			"key": "value",
			1:     2,
		},
	}

	// 按路径获取值
	// 混合使用字符串键（字典）和整型索引（切片）
	get, err := dyno.Get(y, "object", "array", 2, "integer")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%T, %#v\n", get, get)

	get, err = dyno.GetInt(y, "object", "array", 2, "integer")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%T, %#v\n", get, get)

	m := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "江湖十年",
			"address": map[string]interface{}{
				"city": "Beijing",
				"zip":  10115,
			},
		},
	}

	// 按路径获取值
	get, err = dyno.SGet(m, "user", "address", "city")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%T, %#v\n", get, get)

	// 按路径设置值
	err = dyno.Set(m, "Hangzhou", "user", "address", "city")
	if err != nil {
		panic(err)
	}

	// 按路径获取值
	get, err = dyno.SGet(m, "user", "address", "city")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%T, %#v\n", get, get)
}
