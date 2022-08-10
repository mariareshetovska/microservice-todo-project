package handlers

import (
	"encoding/json"
	"mainApi/models"
	"mainApi/pkg/database"
	"mainApi/utils"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type UserAPI struct {
	DB database.Database // will represent all database inrefaces
}

func (api *UserAPI) SignUP(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "user.go-> Create")
	var user *models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.WithError(err).Warn("could not decode parameters")
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	if err := user.Verify(); err != nil {
		logger.WithError(err).Warn("Not all field found")
		utils.ErrorResponse(w, err, http.StatusBadRequest)
	}

	if err := api.DB.CreateUser(r.Context(), user); err == database.ErrUserExists {
		logger.WithError(err).Warn("User already exists")
		utils.ErrorResponse(w, err, http.StatusConflict)
	} else if err != nil {
		logger.WithError(err).Warn("Error creating user")
		utils.ErrorResponse(w, err, http.StatusConflict)
	}
	utils.ToJson(w, user)
}

func (api *UserAPI) Login(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("func", "user.go-> Login")
	params := mux.Vars(r)
	firstname := string(params["id"])
	var user *models.User
	reqPassword := user.Password
	user, err := api.DB.GetUserByCredentials(r.Context(), firstname)
	if err != nil {
		logger.WithError(err).Warn("Error logging in")
		utils.ErrorResponse(w, err, http.StatusConflict)
	}
	err = utils.VerifyPassword(user.Password, reqPassword)
	if err != nil {
		logger.WithError(err).Warn("Error logging in")
		utils.ErrorResponse(w, err, http.StatusBadRequest)
	}

	utils.ToJson(w, user)
}
