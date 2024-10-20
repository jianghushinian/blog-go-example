package main

import "fmt"

func main() {
	fmt.Println(Hello("World"))
}

func Hello(name string) string {
	return "Hello " + name
}
