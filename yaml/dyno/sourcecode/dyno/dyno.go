/*
Package dyno is a utility to work with dynamic objects at ease.

Primary goal is to easily handle dynamic objects and arrays (and a mixture of these)
that are the result of unmarshaling a JSON or YAML text into an interface{}
for example. When unmarshaling into interface{}, libraries usually choose
either map[string]interface{} or map[interface{}]interface{} to represent objects,
and []interface{} to represent arrays. Package dyno supports a mixture of
these in any depth and combination.

When operating on a dynamic object, you designate a value you're interested
in by specifying a path. A path is a navigation; it is a series of map keys
and int slice indices that tells how to get to the value.

Should you need to marshal a dynamic object to JSON which contains maps with
interface{} key type (which is not supported by encoding/json), you may use
the ConvertMapI2MapS converter function.

The implementation does not use reflection at all, so performance is rather good.

Let's see a simple example editing a JSON text to mask out a password. This is
a simplified version of the Example_jsonEdit example function:

	src := `{"login":{"password":"secret","user":"bob"},"name":"cmpA"}`
	var v interface{}
	if err := json.Unmarshal([]byte(src), &v); err != nil {
		panic(err)
	}
	// Edit (mask out) password:
	if err = dyno.Set(v, "xxx", "login", "password"); err != nil {
		fmt.Printf("Failed to set password: %v\n", err)
	}
	edited, err := json.Marshal(v)
	fmt.Printf("Edited JSON: %s, error: %v\n", edited, err)

Output will be:

	Edited JSON: {"login":{"password":"xxx","user":"bob"},"name":"cmpA"}, error: <nil>
*/
package dyno

import (
	"fmt"
)

// Get 从动态结构（如 map[string]interface{} 或 []interface{}）中按路径获取值
// 如果 path 参数为空，则直接返回 v
func Get(v interface{}, path ...interface{}) (interface{}, error) {
	// 遍历路径参数 path
	for i, el := range path {
		switch node := v.(type) { // 断言 v 的类型，注意每一轮循环中这个 v 都是新的值
		// 情况 1：处理键为 string 的 map
		case map[string]interface{}:
			key, ok := el.(string) // 路径元素必须是 string
			if !ok {
				return nil, fmt.Errorf("expected string path element, got: %T (path element idx: %d)", el, i)
			}
			v, ok = node[key] // 从 map 中获取值
			if !ok {
				return nil, fmt.Errorf("missing key: %s (path element idx: %d)", key, i)
			}

		// 情况 2：处理键为任意类型的 map
		case map[interface{}]interface{}:
			var ok bool
			v, ok = node[el] // 直接使用路径元素作为键，来获取值
			if !ok {
				return nil, fmt.Errorf("missing key: %v (path element idx: %d)", el, i)
			}

		// 情况 3：处理切片
		case []interface{}:
			idx, ok := el.(int) // 路径元素必须是 int
			if !ok {
				return nil, fmt.Errorf("expected int path element, got: %T (path element idx: %d)", el, i)
			}
			if idx < 0 || idx >= len(node) { // 索引越界检查
				return nil, fmt.Errorf("index out of range: %d (path element idx: %d)", idx, i)
			}
			v = node[idx] // 从切片中获取值

		// 情况 4：不支持的类型
		default:
			return nil, fmt.Errorf("expected map or slice node, got: %T (path element idx: %d)", node, i)
		}
	}

	// 返回最终 v 的值
	return v, nil
}

// GetInt returns an int value denoted by the path.
//
// If path is empty or nil, v is returned as an int.
func GetInt(v interface{}, path ...interface{}) (int, error) {
	v, err := Get(v, path...)
	if err != nil {
		return 0, err
	}
	i, ok := v.(int)
	if !ok {
		return 0, fmt.Errorf("expected int value, got: %T", v)
	}
	return i, nil
}

// GetSlice returns a slice denoted by the path.
//
// If path is empty or nil, v is returned as a slice.
func GetSlice(v interface{}, path ...interface{}) ([]interface{}, error) {
	v, err := Get(v, path...)
	if err != nil {
		return nil, err
	}
	s, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected slice node, got: %T", v)
	}
	return s, nil
}

// GetMapI returns a map with interface{} keys denoted by the path.
//
// If path is empty or nil, v is returned as a slice.
func GetMapI(v interface{}, path ...interface{}) (map[interface{}]interface{}, error) {
	v, err := Get(v, path...)
	if err != nil {
		return nil, err
	}
	m, ok := v.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("expected map with interface keys node, got: %T", v)
	}
	return m, nil
}

// GetMapS returns a map with string keys denoted by the path.
//
// If path is empty or nil, v is returned as a slice.
func GetMapS(v interface{}, path ...interface{}) (map[string]interface{}, error) {
	v, err := Get(v, path...)
	if err != nil {
		return nil, err
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected map with string keys node, got: %T", v)
	}
	return m, nil
}

// GetInteger returns an int64 value denoted by the path.
//
// This function accepts many different types and converts them to int64, namely:
//
//	-integer types (int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64)
//	 (which implies the aliases byte and rune too)
//	-floating point types (float64, float32)
//	-string (fmt.Sscan() will be used for parsing)
//	-any type with an Int64() (int64, error) method (e.g. json.Number)
//
// If path is empty or nil, v is returned as an int64.
func GetInteger(v interface{}, path ...interface{}) (int64, error) {
	v, err := Get(v, path...)
	if err != nil {
		return 0, err
	}

	switch i := v.(type) {
	case int64:
		return i, nil
	case int:
		return int64(i), nil
	case int32:
		return int64(i), nil
	case int16:
		return int64(i), nil
	case int8:
		return int64(i), nil
	case uint:
		return int64(i), nil
	case uint64:
		return int64(i), nil
	case uint32:
		return int64(i), nil
	case uint16:
		return int64(i), nil
	case uint8:
		return int64(i), nil
	case float64:
		return int64(i), nil
	case float32:
		return int64(i), nil
	case string:
		var n int64
		_, err := fmt.Sscan(i, &n)
		return n, err
	case interface {
		Int64() (int64, error)
	}:
		return i.Int64()
	default:
		return 0, fmt.Errorf("expected some form of integer number, got: %T", v)
	}
}

// GetFloat64 returns a float64 value denoted by the path.
//
// If path is empty or nil, v is returned as a float64.
func GetFloat64(v interface{}, path ...interface{}) (float64, error) {
	v, err := Get(v, path...)
	if err != nil {
		return 0, err
	}
	f, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("expected float64 value, got: %T", v)
	}
	return f, nil
}

// GetFloating returns a float64 value denoted by the path.
//
// This function accepts many different types and converts them to float64, namely:
//
//	-floating point types (float64, float32)
//	-integer types (int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64)
//	 (which implies the aliases byte and rune too)
//	-string (fmt.Sscan() will be used for parsing)
//	-any type with a Float64() (float64, error) method (e.g. json.Number)
//
// If path is empty or nil, v is returned as an int64.
func GetFloating(v interface{}, path ...interface{}) (float64, error) {
	v, err := Get(v, path...)
	if err != nil {
		return 0, err
	}

	switch f := v.(type) {
	case float64:
		return f, nil
	case float32:
		return float64(f), nil
	case int64:
		return float64(f), nil
	case int:
		return float64(f), nil
	case int32:
		return float64(f), nil
	case int16:
		return float64(f), nil
	case int8:
		return float64(f), nil
	case uint:
		return float64(f), nil
	case uint64:
		return float64(f), nil
	case uint32:
		return float64(f), nil
	case uint16:
		return float64(f), nil
	case uint8:
		return float64(f), nil
	case string:
		var n float64
		_, err := fmt.Sscan(f, &n)
		return n, err
	case interface {
		Float64() (float64, error)
	}:
		return f.Float64()
	default:
		return 0, fmt.Errorf("expected some form of floating point number, got: %T", v)
	}
}

// GetString returns a string value denoted by the path.
//
// If path is empty or nil, v is returned as a string.
func GetString(v interface{}, path ...interface{}) (string, error) {
	v, err := Get(v, path...)
	if err != nil {
		return "", err
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("expected string value, got: %T", v)
	}
	return s, nil
}

// GetBoolean returns a bool value denoted by the path.
//
// This function accepts many different types and converts them to bool, namely:
//
//	-boolean type
//	-integer and floating point types (false for zero values, true otherwise)
//	-string (fmt.Sscan() will be used for parsing)
//
// If path is empty or nil, v is returned as a bool.
func GetBoolean(v interface{}, path ...interface{}) (bool, error) {
	v, err := Get(v, path...)
	if err != nil {
		return false, err
	}

	switch f := v.(type) {
	case bool:
		return f, nil
	case int:
		return f != 0, nil
	case int64:
		return f != 0, nil
	case int32:
		return f != 0, nil
	case int16:
		return f != 0, nil
	case int8:
		return f != 0, nil
	case uint:
		return f != 0, nil
	case uint64:
		return f != 0, nil
	case uint32:
		return f != 0, nil
	case uint16:
		return f != 0, nil
	case uint8:
		return f != 0, nil
	case float64:
		return f != 0, nil
	case float32:
		return f != 0, nil
	case string:
		var n bool
		_, err := fmt.Sscan(f, &n)
		return n, err
	case interface {
		Float64() (float64, error)
	}:
		val, err := f.Float64()
		if err != nil {
			return false, err
		}
		return val != 0, err
	default:
		return false, fmt.Errorf("expected bool, got: %T", v)
	}
}

// SGet 从嵌套的 map[string]interface{} 结构中，通过纯字符串路径（不支持切片索引）获取值
//
// SGet is an optimized and specialized version of the general Get.
// The path may only contain string map keys (no slice indices), and
// each value associated with the keys (being the path elements) must also
// be maps with string keys, except the value asssociated with the last
// path element.
//
// If path is empty or nil, m is returned.
func SGet(m map[string]interface{}, path ...string) (interface{}, error) {
	if len(path) == 0 {
		return m, nil
	}

	lastIdx := len(path) - 1
	var value interface{}
	var ok bool

	// 遍历 path
	for i, key := range path {
		// 1. 检查当前键是否存在
		if value, ok = m[key]; !ok {
			return nil, fmt.Errorf("missing key: %s (path element idx: %d)", key, i)
		}
		// 2. 若是最后一个路径元素，直接推出循环
		if i == lastIdx {
			break
		}
		// 3. 检查中间节点是否为 map[string]interface{}
		m2, ok := value.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expected map with string keys node, got: %T (path element idx: %d)", value, i)
		}
		// 4. 更新当前节点，继续深入下一层
		m = m2
	}

	return value, nil
}

// Set 通过路径在动态结构（嵌套的 map 或 slice）中修改指定位置元素的值
//
// The last element of the path must be a map key or a slice index, and the
// preceding path must denote a map or a slice respectively which must already exist.
//
// Path cannot be empty or nil, else an error is returned.
func Set(v interface{}, value interface{}, path ...interface{}) error {
	// 1. 路径非空校验
	if len(path) == 0 {
		return fmt.Errorf("path cannot be empty")
	}

	// 2. 分离末位元素
	i := len(path) - 1 // The last index
	if len(path) > 1 {
		var err error
		v, err = Get(v, path[:i]...) // 末位元素的值
		if err != nil {
			return err
		}
	}

	el := path[i] // 末位元素的键或索引

	// 3. 根据末位元素值的类型进行赋值
	switch node := v.(type) {
	case map[string]interface{}: // 情况 1：处理键为 string 的 map
		key, ok := el.(string)
		if !ok {
			return fmt.Errorf("expected string path element, got: %T (path element idx: %d)", el, i)
		}
		node[key] = value

	case map[interface{}]interface{}: // 情况 2：处理键为任意类型的 map
		node[el] = value

	case []interface{}: // 情况 3：处理切片
		idx, ok := el.(int)
		if !ok {
			return fmt.Errorf("expected int path element, got: %T (path element idx: %d)", el, i)
		}
		if idx < 0 || idx >= len(node) {
			return fmt.Errorf("index out of range: %d (path element idx: %d)", idx, i)
		}
		node[idx] = value

	default: // 情况 4：不支持的类型
		return fmt.Errorf("expected map or slice node, got: %T (path element idx: %d)", node, i)
	}

	return nil
}

// SSet sets a map element with string key type, denoted by the path
// consisting of only string keys.
//
// SSet is an optimized and specialized version of the general Set. The
// path may only contain string map keys (no slice indices), and each
// value associated with the keys (being the path elements) must also be
// a maps with string keys, except the value associated with the last path
// element.
//
// The map denoted by the preceding path before the last path element
// must already exist.
//
// Path cannot be empty or nil, else an error is returned.
func SSet(m map[string]interface{}, value interface{}, path ...string) error {
	if len(path) == 0 {
		return fmt.Errorf("path cannot be empty")
	}

	i := len(path) - 1 // The last index
	if len(path) > 1 {
		v, err := SGet(m, path[:i]...)
		if err != nil {
			return err
		}

		var ok bool
		m, ok = v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected map with string keys node, got: %T (path element idx: %d)", value, i)
		}
	}

	m[path[i]] = value
	return nil
}

// Append 在动态结构（如嵌套的 map 或 slice）中，向路径指向的切片末尾追加元素
//
// The slice denoted by path must already exist.
//
// Path cannot be empty or nil, else an error is returned.
func Append(v interface{}, value interface{}, path ...interface{}) error {
	// 1. 路径非空校验
	if len(path) == 0 {
		return fmt.Errorf("path cannot be empty")
	}

	// 2. 获取路径指向的节点（必须是切片类型）
	node, err := Get(v, path...)
	if err != nil {
		return err
	}

	// 3. 类型断言：验证是否为切片
	s, ok := node.([]interface{})
	if !ok {
		return fmt.Errorf("expected slice node, got: %T (path element idx: %d)", node, len(path)-1)
	}

	// 4. 追加元素并更新原切片
	return Set(v, append(s, value), path...)
}

// AppendMore appends values to a slice denoted by the path.
//
// The slice denoted by path must already exist.
//
// Path cannot be empty or nil, else an error is returned.
func AppendMore(v interface{}, values []interface{}, path ...interface{}) error {
	if len(path) == 0 {
		return fmt.Errorf("path cannot be empty")
	}

	node, err := Get(v, path...)
	if err != nil {
		return err
	}

	s, ok := node.([]interface{})
	if !ok {
		return fmt.Errorf("expected slice node, got: %T (path element idx: %d)", node, len(path))
	}

	// Must set the new slice value:
	return Set(v, append(s, values...), path...)
}

// Delete 根据给定的 path，从 map 中删除指定的键值对，或从 slice 中删除指定元素
//
// Deleting a non-existing map key is a no-op. Attempting to delete a slice
// element from a slice with invalid index is an error.
//
// Path cannot be empty or nil if v itself is a slice, else an error is returned.
func Delete(v interface{}, key interface{}, path ...interface{}) error {
	// 1. 路径校验：若 v 是切片，路径不能为空
	if len(path) == 0 {
		if _, ok := v.([]interface{}); ok {
			return fmt.Errorf("path cannot be empty if v is a slice")
		}
	}

	// 2. 获取路径末位节点
	node, err := Get(v, path...)
	if err != nil {
		return err
	}

	// 3. 根据末位节点类型删除
	switch node2 := node.(type) {
	case map[string]interface{}: // 情况 1：处理键为 string 的 map
		skey, ok := key.(string)
		if !ok {
			return fmt.Errorf("expected string key, got: %T", key)
		}
		delete(node2, skey) // 删除 string 类型的键

	case map[interface{}]interface{}: // 情况 2：处理键为任意类型的 map
		delete(node2, key) // 直接删除键

	case []interface{}: // 情况 3：处理切片
		idx, ok := key.(int)
		if !ok {
			return fmt.Errorf("expected int key, got: %T", key)
		}
		if idx < 0 || idx >= len(node2) {
			return fmt.Errorf("index out of range: %d", idx)
		}
		// 删除元素：移位 + 截断
		copy(node2[idx:], node2[idx+1:]) // 后续元素前移
		// Clear the emptied element:
		node2[len(node2)-1] = nil // 清空尾部引用（防内存泄漏）
		// Must set the new slice value:
		return Set(v, node2[:len(node2)-1], path...) // 将新的节点赋值给 v

	default: // 情况 4：不支持的类型
		return fmt.Errorf("expected map or slice node, got: %T (path element idx: %d)", node, len(path)-1)
	}

	return nil
}

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
