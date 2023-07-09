package animal

import "fmt"

func ExampleAnimal_shout() {
	bird := Animal{Name: "bird"}
	fmt.Print(bird.shout())
	// Output:
	// 吼～
}
