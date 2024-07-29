package main

import (
	"encoding/json"
	"fmt"

	simplejson "github.com/jianghushinian/blog-go-example/struct/encoding-json/encoding/json"
)

type User struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string
}

func main() {
	// NOTE: Go 内置 encoding/json
	{
		user := User{
			Name:  "江湖十年",
			Age:   20,
			Email: "jianghushinian007@outlook.com",
		}

		jsonData, err := json.Marshal(user)
		if err != nil {
			fmt.Println("Error marshal to JSON:", err)
			return
		}

		fmt.Printf("JSON data: %s\n", jsonData)
	}

	{
		jsonData := `{"name": "江湖十年", "age": 20, "Email": "jianghushinian007@outlook.com"}`

		var user User
		err := json.Unmarshal([]byte(jsonData), &user)
		if err != nil {
			fmt.Println("Error unmarshal from JSON:", err)
			return
		}

		fmt.Printf("User struct: %+v\n", user)
	}

	// NOTE: 自己实现 encoding/json
	{
		user := User{
			Name:  "江湖十年",
			Age:   20,
			Email: "jianghushinian007@outlook.com",
		}

		jsonData, err := simplejson.Marshal(user)
		if err != nil {
			fmt.Println("Error marshal to JSON:", err)
			return
		}

		fmt.Printf("JSON data: %s\n", jsonData)
	}

	{
		jsonData := `{"name": "江湖十年", "age": 20, "Email": "jianghushinian007@outlook.com"}`

		var user User
		err := simplejson.Unmarshal([]byte(jsonData), &user)
		if err != nil {
			fmt.Println("Error unmarshal from JSON:", err)
			return
		}

		fmt.Printf("User struct: %+v\n", user)
	}
}
