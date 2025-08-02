package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/goccy/go-yaml" // YAML 解析库
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

	// 转换为 JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("JSON convert error: %v", err)
	}
	fmt.Printf("Type: %T\nValue: %s\n", jsonData, jsonData)
}
