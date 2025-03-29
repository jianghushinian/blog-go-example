/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package errors

import (
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"
)

// MessageCountMap contains occurrence for each error message.
type MessageCountMap map[string]int

// Aggregate 表示一个包含多个错误的对象，但这些错误并不一定具有单一的语义含义。
// 该聚合错误可通过 `errors.Is()` 检查是否包含特定类型的错误。
// 不支持 Errors.As()，因为调用方可能需要关注多个匹配给定类型的错误中的特定错误
type Aggregate interface {
	error
	Errors() []error // 暴露内部错误列表
	Is(error) bool   // 兼容 Go 1.13+ 错误链判断
}

// NewAggregate 错误聚合入口，将一个 error 切片转换成 Aggregate 接口，如果传入的切片为空，返回 nil
func NewAggregate(errlist []error) Aggregate {
	if len(errlist) == 0 {
		return nil
	}
	// 过滤 nil 错误（防止空指针）
	var errs []error
	for _, e := range errlist {
		if e != nil {
			errs = append(errs, e)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return aggregate(errs) // 返回 Aggregate 接口的实现 aggregate 对象
}

// Aggregate 接口的具体实现
// 该辅助类型实现了 error 接口和 Errors 接口。
// 保持其私有性可防止外部创建包含 0 个错误的聚合对象，
// 空错误集合虽然满足 error 接口，但语义上并非真正的错误。
type aggregate []error

// Error 实现 error 接口
func (agg aggregate) Error() string {
	// 错误为空直接返回 ""
	if len(agg) == 0 {
		// This should never happen, really.
		return ""
	}

	// 单错误直接返回
	if len(agg) == 1 {
		return agg[0].Error()
	}

	// 使用集合去重
	seenerrs := sets.NewString() // 使用 map 实现 set：map[string]struct{}
	result := ""
	agg.visit(func(err error) bool { // 这里使用 visit 递归判断 agg 中每一个错误对象
		msg := err.Error()
		if seenerrs.Has(msg) { // 在闭包函数中实现去重
			return false
		}
		seenerrs.Insert(msg)
		if len(seenerrs) > 1 {
			result += ", " // 多错误时中间使用逗号分割
		}
		result += msg // 拼接去重后的错误信息
		return false
	})

	// 单错误直接返回
	if len(seenerrs) == 1 {
		return result
	}

	return "[" + result + "]" // 多错误用方括号包裹
}

// Is 兼容 Go 1.13+
func (agg aggregate) Is(target error) bool {
	return agg.visit(func(err error) bool { // 递归判断每一个错误对象，是否等于 target
		return errors.Is(err, target)
	})
}

// 递归遍历聚合错误树，对每个错误执行判断函数 f
// 返回值 true 表示存在满足条件的错误，false 表示未找到
func (agg aggregate) visit(f func(err error) bool) bool {
	// 遍历错误列表
	for _, err := range agg {
		switch err := err.(type) {
		case aggregate: // 嵌套的私有聚合类型
			// 递归遍历子聚合错误
			if match := err.visit(f); match {
				return match
			}
		case Aggregate: // 实现了 Aggregate 接口的其他类型
			// 遍历接口公开的错误列表
			for _, nestedErr := range err.Errors() {
				if match := f(nestedErr); match { // 将嵌套的错误传给 f 函数进行检查
					return match // 嵌套的错误匹配则终止
				}
			}
		default: // 其他错误类型
			if match := f(err); match { // 直接应用判断函数
				return match
			}
		}
	}

	return false
}

// Errors 返回错误列表
func (agg aggregate) Errors() []error {
	return []error(agg)
}

// Matcher 用于匹配错误，如果匹配则返回 true
type Matcher func(error) bool

// FilterOut 从输入错误中移除所有匹配任意 Matcher 的错误。
// 如果输入是单一错误，仅对该错误进行测试。
// 如果输入实现了 Aggregate 接口，将递归处理错误列表。
//
// 例如：可用于从错误列表中移除已知无害的错误（如 io.EOF 或 os.PathNotFound）
// 使用示例：TestFilterOut
// e.g., https://github.com/kubernetes/kubernetes/blob/v1.32.0/staging/src/k8s.io/component-helpers/auth/rbac/reconciliation/namespace.go#L37
func FilterOut(err error, fns ...Matcher) error {
	if err == nil { // err 为 nil 直接返回
		return nil
	}
	if agg, ok := err.(Aggregate); ok { // 如果是 Aggregate 类型
		return NewAggregate(filterErrors(agg.Errors(), fns...)) // 递归处理错误列表
	}
	if !matchesError(err, fns...) { // 如果全部不匹配，返回原 err
		return err
	}
	return nil
}

// matchesError 如果任意 Matcher 返回 true，则返回 true
func matchesError(err error, fns ...Matcher) bool {
	for _, fn := range fns {
		if fn(err) {
			return true
		}
	}
	return false
}

// filterErrors 返回所有未被 Matcher（所有 fns 返回 false）过滤掉的错误，
// 包含嵌套错误（如果列表中存在嵌套的 Errors 对象）。
// 如果没有剩余错误，则返回 nil 列表。
// 副作用：返回的结果切片会将所有嵌套的错误切片扁平化。
func filterErrors(list []error, fns ...Matcher) []error {
	result := []error{}
	for _, err := range list {
		r := FilterOut(err, fns...)
		if r != nil {
			result = append(result, r)
		}
	}
	return result
}

// Flatten 接收一个嵌套任意层的 Aggregate，并递归的将其扁平化
// 使用示例：TestFlatten
// e.g., https://github.com/kubernetes/kubernetes/blob/v1.32.0/pkg/scheduler/apis/config/validation/validation.go#L81
func Flatten(agg Aggregate) Aggregate {
	result := []error{} // 保存扁平化的单层错误列表
	if agg == nil {     // 如果为 nil 直接返回
		return nil
	}

	// 遍历当前层错误列表
	for _, err := range agg.Errors() {
		if a, ok := err.(Aggregate); ok { // 如果嵌套了 Aggregate 类型
			r := Flatten(a) // 递归展开嵌套结构
			if r != nil {
				result = append(result, r.Errors()...)
			}
		} else {
			if err != nil {
				result = append(result, err)
			}
		}
	}
	return NewAggregate(result)
}

// CreateAggregateFromMessageCountMap 将给定的 MessageCountMap 转换为 Aggregate
// 使用示例：TestCreateAggregateFromMessageCountMap
func CreateAggregateFromMessageCountMap(m MessageCountMap) Aggregate {
	if m == nil { // 如果 map 为 nil 直接返回
		return nil
	}
	result := make([]error, 0, len(m))
	for errStr, count := range m {
		var countStr string
		if count > 1 {
			countStr = fmt.Sprintf(" (repeated %v times)", count)
		}
		result = append(result, fmt.Errorf("%v%v", errStr, countStr))
	}
	return NewAggregate(result)
}

// Reduce 如果错误是一个 Aggregate 类型且只有一项，将会返回错误或 nil，即返回 aggregate 中的第一项。
// 使用示例：Reduce(Flatten(NewAggregate(errs)))
// e.g., https://github.com/kubernetes/kubernetes/blob/v1.32.0/staging/src/k8s.io/kubectl/pkg/cmd/get/get.go#L729
func Reduce(err error) error {
	if agg, ok := err.(Aggregate); ok && err != nil {
		switch len(agg.Errors()) {
		case 1: // 单错误提取
			return agg.Errors()[0]
		case 0: // Aggregate 为空
			return nil
		}
	}
	return err // 非 Aggregate 类型直接返回
}

// AggregateGoroutines 并行运行提供的函数，将所有非 nil 错误收集到返回的 Aggregate 中。
// 如果所有函数都成功完成，则返回 nil。
// e.g., https://github.com/kubernetes/kubernetes/blob/v1.32.0/staging/src/k8s.io/apiserver/pkg/audit/union.go#L56
func AggregateGoroutines(funcs ...func() error) Aggregate {
	errChan := make(chan error, len(funcs)) // 创建容量等于函数数量的缓冲 channel
	for _, f := range funcs {
		go func(f func() error) { errChan <- f() }(f) // 并行执行
	}
	errs := make([]error, 0)
	for i := 0; i < cap(errChan); i++ { // 按容量遍历，确保处理所有已启动的任务
		if err := <-errChan; err != nil { // 同步等待每个任务完成
			errs = append(errs, err)
		}
	}
	return NewAggregate(errs) // 将错误列表封装为 Aggregate 接口
}

// ErrPreconditionViolated is returned when the precondition is violated
var ErrPreconditionViolated = errors.New("precondition is violated")
