package main

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

func main() {
	type A struct {
		noCopy noCopy
		a      string
	}

	type B struct {
		b string
	}

	a := A{a: "a"}
	b := B{b: "b"}

	_ = a
	_ = b
}

// $ go vet main.go
//
// # command-line-arguments
// # [command-line-arguments]
// ./main.go:21:6: assignment copies lock value to _: command-line-arguments.A contains command-line-arguments.noCopy
