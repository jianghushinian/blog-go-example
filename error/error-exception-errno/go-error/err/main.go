package main

import "log"

type MyError struct {
	msg string
}

func (e *MyError) Error() string {
	return e.msg
}

type SmsCode struct{}

func (c SmsCode) Verify2() (bool, error) {
	// ...
	return false, nil
}

func (c SmsCode) Verify1() error {
	// ...
	return nil
}

func a() error {
	// ...
	return nil
}

func b() error {
	// ...
	return nil
}

func c() error {
	// ...
	return nil
}

func main() {
	err := a()
	if err != nil {
		log.Fatal(err)
	}
	err = b()
	if err != nil {
		log.Fatal(err)
	}
	err = c()
	if err != nil {
		log.Fatal(err)
	}
}
