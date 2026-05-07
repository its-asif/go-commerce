package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/its-asif/go-commerce/db"
	"github.com/its-asif/go-commerce/models"
	"net/http"
	"strconv"
)

// @Summary		Get All Users
// @Description	Get all the users from db
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			Authorization	header		string	true	"Bearer + JWT_Token"
// @Success		200				{object}	[]models.User
// @Failure		400				{string}	string	"Bad Request"
// @Router			/api/users [GET]
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := db.GetAllUsers()
	if err != nil {
		http.Error(w, "error getting all users", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(users)
}

// @Summary		Get Single Users by ID
// @Description	Get a single user from db
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			Authorization	header		string	true	"Bearer + JWT_Token"
// @Param			id	path		int		true	"User ID"
// @Success		200				{object}	models.User
// @Failure		400				{string}	string	"Bad Request"
// @Router			/api/users/id/{id} [GET]
func GetSingleUserByID(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	var user models.User
	user, err = db.GetSingleUser("id", id)
	if err != nil {
		http.Error(w, "error getting single users", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}

// @Summary		Get Single Users By Email
// @Description	Get a single user from db
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			Authorization	header		string	true	"Bearer + JWT_Token"
// @Param			email	path		string		true	"User Email"
// @Success		200				{object}	models.User
// @Failure		400				{string}	string	"Bad Request"
// @Router			/api/users/email/{email} [GET]
func GetSingleUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["email"]

	user, err := db.GetSingleUser("email", email)
	if err != nil {
		http.Error(w, "error getting single users", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}
