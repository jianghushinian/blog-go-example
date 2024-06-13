package main

import (
	"fmt"
	"log"

	"go.uber.org/dig"
)

// NOTE: dig 示例
// ref: https://www.jetbrains.com/guide/go/tutorials/dependency_injection_part_one/di_with_dig/

type User struct {
	name string
}

// NewUser - Creates a new instance of User
func NewUser(name string) User {
	return User{name: name}
}

// Get - A method with user as dependency
func (u *User) Get(message string) string {
	return fmt.Sprintf("Hello %s - %s", u.name, message)
}

// Run - Depends on user and calls the Get method on User
func Run(user User) {
	result := user.Get("It's nice to meet you!")
	fmt.Println(result)
}

func main() {
	// Initialize a new dig container
	container := dig.New()
	// Provide a name parameter to the container
	container.Provide(func() string { return "jianghushinian" })
	// Provide a new User instance to the container using the name injected above
	if err := container.Provide(NewUser); err != nil {
		log.Fatal(err)
	}
	// Invoke the Run function; Dig automatically injects the User instance provided above
	if err := container.Invoke(Run); err != nil {
		log.Fatal(err)
	}
}
