package json

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Marshal 序列化
func Marshal(v any) (string, error) {
	// 拿到对象 v 的 reflect.Value 和 reflect.Type
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return "", fmt.Errorf("only structs are supported")
	}
	typ := val.Type()

	// 用来保存 JSON 字符串
	jsonBuilder := strings.Builder{}

	// NOTE: 三步走拼接 JSON 字符串

	// 1. JSON 左花括号
	jsonBuilder.WriteString("{")

	// 2. key/value
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)

		// 获取 JSON 标签
		tag := fieldType.Tag.Get("json")
		if tag == "" {
			tag = fieldType.Name
		}

		jsonBuilder.WriteString(`"` + tag + `":`)

		// 根据字段类型转换，仅支持 string/int
		switch fieldVal.Kind() {
		case reflect.String:
			jsonBuilder.WriteString(`"` + fieldVal.String() + `"`)
		case reflect.Int:
			jsonBuilder.WriteString(strconv.FormatInt(fieldVal.Int(), 10))
		default:
			return "", fmt.Errorf("unsupported field type: %s", fieldVal.Kind())
		}

		if i < val.NumField()-1 {
			jsonBuilder.WriteString(",")
		}
	}

	// 3. JSON 右花括号
	jsonBuilder.WriteString("}")

	return jsonBuilder.String(), nil
}
