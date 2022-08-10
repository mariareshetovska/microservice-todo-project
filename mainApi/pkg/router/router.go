package router

import (
	"mainApi/pkg/database"
	"net/http"

	"github.com/gorilla/mux"

	"mainApi/pkg/handlers"
)

func NewRouter(db database.Database) (http.Handler, error) {
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	userApi := &handlers.UserAPI{
		DB: db,
	}

	todoApi := &handlers.TodoApi{
		DB: db,
	}
	apiRouter.HandleFunc("/register", userApi.SignUP).Methods("POST")
	apiRouter.HandleFunc("/login", userApi.Login).Methods("POST")

	apiRouter.HandleFunc("/users/todos", todoApi.CreateTodo).Methods("POST")
	apiRouter.HandleFunc("/users/{id}/todos", todoApi.GetTodoListByUserID).Methods("GET")
	apiRouter.HandleFunc("/users/todos/{id}", todoApi.GetTodoByTodoID).Methods("GET")
	apiRouter.HandleFunc("/users/todos/{id}", todoApi.UpdateTodo).Methods("PUT")
	apiRouter.HandleFunc("/users/todos/{id}", todoApi.DeleteTodo).Methods("DELETE")

	return router, nil
}
