package main

import (
	"errors"
	"fmt"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

func main() {
	var errs []error

	// 模拟多个操作可能失败
	err1 := step1()
	if err1 != nil {
		errs = append(errs, err1)
	}

	err2 := step2()
	if err2 != nil {
		errs = append(errs, err2)
	}

	agg := utilerrors.NewAggregate(errs)

	fmt.Printf("errs: %s\n", errs)                                                // errs: [step1 failed step2 failed]
	fmt.Printf("aggregate: %s\n", agg)                                            // aggregate: [step1 failed, step2 failed]
	fmt.Printf("errs len: %d, aggregate len: %d\n", len(errs), len(agg.Errors())) // errs len: 2, aggregate len: 2
	fmt.Printf("errors: %s\n", agg.Errors())                                      // errors: [step1 failed step2 failed]
	fmt.Printf("err1: %s, err2: %s\n", agg.Errors()[0], agg.Errors()[1])          // err1: step1 failed, err2: step2 failed

	fmt.Println(
		errors.Is(agg, err1),               // true
		errors.Is(agg, err2),               // true
		errors.Is(agg, errors.New("err3")), // false
	)
}

func step1() error {
	return errors.New("step1 failed")
}

func step2() error {
	return errors.New("step2 failed")
}
