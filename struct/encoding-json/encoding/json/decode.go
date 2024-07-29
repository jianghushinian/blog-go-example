package json

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Unmarshal 反序列化
func Unmarshal(data []byte, v any) error {
	parsedData, err := parseJSON(string(data))
	if err != nil {
		return err
	}

	val := reflect.ValueOf(v).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)

		// 获取 JSON 标签
		tag := fieldType.Tag.Get("json")
		if tag == "" {
			tag = fieldType.Name
		}

		// 从解析的数据中获取值
		if value, ok := parsedData[tag]; ok {
			switch fieldVal.Kind() {
			case reflect.String:
				fieldVal.SetString(value)
			case reflect.Int:
				intValue, err := strconv.Atoi(value)
				if err != nil {
					return err
				}
				fieldVal.SetInt(int64(intValue))
			default:
				return fmt.Errorf("unsupported field type: %s", fieldVal.Kind())
			}
		}
	}

	return nil
}

// 简易版 JSON 解析器，仅支持 string/int 且不考虑嵌套
func parseJSON(data string) (map[string]string, error) {
	result := make(map[string]string)

	data = strings.TrimSpace(data)
	if len(data) < 2 || data[0] != '{' || data[len(data)-1] != '}' {
		return nil, errors.New("invalid JSON")
	}

	data = data[1 : len(data)-1]
	parts := strings.Split(data, ",")
	for _, part := range parts {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) != 2 {
			return nil, errors.New("invalid JSON")
		}

		k := strings.Trim(strings.TrimSpace(kv[0]), `"`)
		v := strings.Trim(strings.TrimSpace(kv[1]), `"`)

		result[k] = v
	}

	return result, nil
}
