package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/icza/dyno" // 用于递归转换 map[interface{}] => map[string]
	"gopkg.in/yaml.v3"
)

func main() {
	yamlExample := `
object:
  a: 1
  1: 2
  # "1": 3
  key: value
  array:
  - null_value: 
  - boolean: true
  - integer: 1
`

	// 解析 YAML
	var data interface{}
	if err := yaml.Unmarshal([]byte(yamlExample), &data); err != nil {
		log.Fatalf("YAML parse error: %v", err)
	}
	fmt.Printf("Type: %T\nValue: %#v\n", data, data)

	fmt.Println("--------------------------")

	// 关键：递归转换 map 键类型（interface{} → string）
	convertedData := dyno.ConvertMapI2MapS(data)

	// 转换为 JSON
	jsonData, err := json.Marshal(convertedData)
	if err != nil {
		log.Fatalf("JSON convert error: %v", err)
	}
	fmt.Printf("Type: %T\nValue: %s\n", jsonData, jsonData)
}

// ref: dyno.ConvertMapI2MapS

// ConvertMapI2MapS walks the given dynamic object recursively, and
// converts maps with interface{} key type to maps with string key type.
// This function comes handy if you want to marshal a dynamic object into
// JSON where maps with interface{} key type are not allowed.
//
// Recursion is implemented into values of the following types:
//
//	-map[interface{}]interface{}
//	-map[string]interface{}
//	-[]interface{}
//
// When converting map[interface{}]interface{} to map[string]interface{},
// fmt.Sprint() with default formatting is used to convert the key to a string key.
func ConvertMapI2MapS(v interface{}) interface{} {
	switch x := v.(type) {
	case map[interface{}]interface{}: // 目标转换类型，需要把 key 为 interface{} 类型的转换成 string
		m := map[string]interface{}{}
		for k, v2 := range x {
			switch k2 := k.(type) {
			case string: // 如果 key 已经是 string 类型，则直接使用
				m[k2] = ConvertMapI2MapS(v2)
			default: // 如果 key 是其他类型则需要转换成 string
				m[fmt.Sprint(k)] = ConvertMapI2MapS(v2)
			}
		}
		v = m

	case []interface{}: // 递归处理数组元素
		for i, v2 := range x {
			x[i] = ConvertMapI2MapS(v2)
		}

	case map[string]interface{}: // key 已经是 string，仅递归处理 value
		for k, v2 := range x {
			x[k] = ConvertMapI2MapS(v2)
		}
	}

	return v
}

/*

map[string]interface {}{
	"object": map[interface {}]interface {}{
		"a":1,
		"array":[]interface {}{
			map[string]interface {}{
				"null_value":interface {}(nil)
			},
			map[string]interface {}{
				"boolean":true
			},
			map[string]interface {}{
				"integer":1
			}
		},
		"key":"value",
		1:2
	}
}

*/
