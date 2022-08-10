package handlers

import (
	"encoding/json"
	"io/ioutil"
	"mainApi/models"
	"mainApi/pkg/database"
	"mainApi/utils"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type TodoApi struct {
	DB database.Database
}

func (api *TodoApi) CreateTodo(w http.ResponseWriter, r *http.Request) {
	logger := logrus.New().WithField("func", "todo -> Create()")
	var todo *models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	if err := todo.Verify(); err != nil {
		logger.WithError(err).Warn("Not all field found in todo")
		utils.ErrorResponse(w, err, http.StatusBadRequest)
	}

	err := api.DB.CreateTodo(r.Context(), todo)
	if err != nil {
		utils.ErrorResponse(w, err, http.StatusUnprocessableEntity)
		return
	}
	utils.ToJson(w, todo)
}

func (api *TodoApi) GetTodoListByUserID(w http.ResponseWriter, r *http.Request) {
	logger := logrus.New().WithField("func", "todo -> GetTodoListByUserID()")
	params := mux.Vars(r)
	id, _ := strconv.ParseUint(params["id"], 10, 32)

	todoList, err := api.DB.GetTodoListByUser(r.Context(), uint32(id))
	if err != nil {
		logger.WithError(err).Warn("could not find todo with id")
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	utils.ToJson(w, todoList)
}

func (api *TodoApi) GetTodoByTodoID(w http.ResponseWriter, r *http.Request) {
	logger := logrus.New().WithField("func", "todo -> GetTodoByTodoID()")
	params := mux.Vars(r)
	todoID, err := uuid.FromString(string(params["id"]))
	if err != nil {
		logger.WithError(err).Warn("error of parsing todoId")
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	todo, err := api.DB.GetTodoById(r.Context(), todoID)
	if err != nil {
		logger.WithError(err).Warn("could not find todo with id")
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	utils.ToJson(w, todo)
}

func (api *TodoApi) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	logger := logrus.New().WithField("func", "todo -> UpdateTodoTodoID()")
	params := mux.Vars(r)
	todoID, _ := uuid.FromString(string(params["id"]))
	body, _ := ioutil.ReadAll(r.Body)
	var todo models.Todo
	err := json.Unmarshal(body, &todo)
	if err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.ErrorResponse(w, err, http.StatusUnprocessableEntity)
		return
	}
	todo.ID = todoID
	rows, err := api.DB.UpdateTodoByID(r.Context(), todo)
	if err != nil {
		logger.WithError(err).Warn("could not update todo")
		utils.ErrorResponse(w, err, http.StatusUnprocessableEntity)
		return
	}
	utils.ToJson(w, rows)

}

func (api *TodoApi) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	logger := logrus.New().WithField("func", "todo -> DeleteTodo()")
	params := mux.Vars(r)
	todoID, _ := uuid.FromString(string(params["id"]))
	_, err := api.DB.DeleteTodo(r.Context(), todoID)
	if err != nil {
		logger.WithError(err).Warn("could not delete todo")
		utils.ErrorResponse(w, err, http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
