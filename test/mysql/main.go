package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/jianghushinian/blog-go-example/test/mysql/store"
)

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		store: store.NewUserStore(db),
	}
}

type UserHandler struct {
	store store.UserStore
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

	u := store.User{}
	if err := json.Unmarshal(body, &u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"msg":"%s"}`, err.Error())
		return
	}

	if err := h.store.Create(&u); err != nil {
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
	u, err := h.store.Get(uid)
	if err != nil {
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
	mysqlDB, _ := store.NewMySQLDB("localhost", "3306", "user", "password", "test")
	handler := NewUserHandler(mysqlDB)
	router := setupRouter(handler)

	_ = http.ListenAndServe(":8000", router)
}
