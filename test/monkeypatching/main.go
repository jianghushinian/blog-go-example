package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID   int
	Name string
}

func NewMySQLDB(host, port, user, pass, dbname string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, dbname)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func NewUserHandler(store *gorm.DB) *UserHandler {
	return &UserHandler{store: store}
}

type UserHandler struct {
	store *gorm.DB
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, `{"msg":"%s"}`, err.Error())
		return
	}

	if err := h.store.Create(&u).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"msg":"%s"}`, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps[0].Value
	uid, _ := strconv.Atoi(id)

	w.Header().Set("Content-Type", "application/json")
	var u User
	if err := h.store.First(&u, uid).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"msg":"%s"}`, err.Error())
		return
	}
	_, _ = fmt.Fprintf(w, `{"id":%d,"name":"%s"}`, u.ID, u.Name)
}

func setupRouter(handler *UserHandler) *httprouter.Router {
	router := httprouter.New()
	router.POST("/users", handler.CreateUser)
	router.GET("/users/:id", handler.GetUser)
	return router
}

func main() {
	mysqlDB, _ := NewMySQLDB("localhost", "3306", "user", "password", "test")
	handler := NewUserHandler(mysqlDB)
	router := setupRouter(handler)
	_ = http.ListenAndServe(":8000", router)
}
