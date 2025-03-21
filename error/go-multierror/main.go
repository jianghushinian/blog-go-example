package main

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

func main() {
	var errs *multierror.Error

	// 模拟多个操作可能失败
	if err := step1(); err != nil {
		errs = multierror.Append(errs, err)
	}

	if err := step2(); err != nil {
		errs = multierror.Append(errs, err)
	}

	// 如果有任何错误，将返回聚合的错误
	if errs != nil {
		// 2 errors occurred:
		//        * step1 failed
		//        * step2 failed
		//
		fmt.Println(errs)
		fmt.Println(errors.Unwrap(errs))                               // step1 failed
		fmt.Println(errors.Unwrap(errors.Unwrap(errs)))                // step2 failed
		fmt.Println(errors.Unwrap(errors.Unwrap(errors.Unwrap(errs)))) // <nil>

		var e *multierror.Error
		fmt.Println(errors.As(errs, &e)) // true
		fmt.Println(errors.Is(errs, &multierror.Error{
			Errors: []error{
				errors.New("step1 failed"),
				errors.New("step2 failed"),
			},
		})) // false
		fmt.Println(errors.Is(Sentinel(), ErrSentinel)) // true

		errs.ErrorFormat = func(e []error) string {
			return e[0].Error()
		}
		fmt.Println("Errors occurred:")
		fmt.Println(errs.Error())                       // step1 failed
		fmt.Println(errors.Unwrap(errs))                // step1 failed
		fmt.Println(errors.Unwrap(errors.Unwrap(errs))) // step2 failed
	} else {
		fmt.Println("All steps succeeded!")
	}
}

func step1() error {
	return errors.New("step1 failed")
}

func step2() error {
	return errors.New("step2 failed")
}

var ErrSentinel = &multierror.Error{
	Errors: []error{errors.New("step1 failed")},
}

func Sentinel() error {
	return ErrSentinel
}
