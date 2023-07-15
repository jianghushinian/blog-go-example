package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var users = []User{
	{ID: 1, Name: "user1"},
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, `{"msg":"%s"}`, err.Error())
		return
	}
	defer func() { _ = r.Body.Close() }()

	u := User{}
	if err := json.Unmarshal(body, &u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"msg":"%s"}`, err.Error())
		return
	}
	u.ID = users[len(users)-1].ID + 1
	users = append(users, u)

	w.WriteHeader(http.StatusCreated)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, _ := strconv.Atoi(ps[0].Value)
	w.Header().Set("Content-Type", "application/json")

	for _, u := range users {
		if u.ID == userID {
			user, _ := json.Marshal(u)
			_, _ = w.Write(user)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte(`{"msg":"notfound"}`))
}

func setupRouter() *httprouter.Router {
	router := httprouter.New()
	router.POST("/users", CreateUserHandler)
	router.GET("/users/:id", GetUserHandler)
	return router
}

func main() {
	router := setupRouter()
	_ = http.ListenAndServe(":8000", router)
}
