package main

import (
	"fmt"
	"log"

	"github.com/facebookgo/inject"
)

// NOTE: inject 示例

type User struct {
	Name string `inject:"name"`
}

// Get - A method with user as dependency
func (u *User) Get(message string) string {
	return fmt.Sprintf("Hello %s - %s", u.Name, message)
}

// Run - Depends on user and calls the Get method on User
func Run(user *User) {
	result := user.Get("It's nice to meet you!")
	fmt.Println(result)
}

func main() {
	// new an inject Graph
	var g inject.Graph

	// inject name
	name := "jianghushinian"

	// provide string value
	err := g.Provide(&inject.Object{Value: name, Name: "name"})
	if err != nil {
		log.Fatal(err)
	}

	// create a User instance and supply it to the dependency graph
	user := &User{}
	err = g.Provide(&inject.Object{Value: user})
	if err != nil {
		log.Fatal(err)
	}

	// resolve all dependencies
	err = g.Populate()
	if err != nil {
		log.Fatal(err)
	}

	Run(user)
}
