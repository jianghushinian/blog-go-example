package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var users = []User{}

func createUser(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u.ID = len(users) + 1
	users = append(users, u)
	c.JSON(http.StatusOK, u)
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	for _, u := range users {
		if id == string(rune(u.ID+'0')) {
			c.JSON(http.StatusOK, u)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/users", createUser)
	r.GET("/users/:id", getUser)
	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
