package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/its-asif/go-commerce/db"
	"github.com/its-asif/go-commerce/models"
	"github.com/its-asif/go-commerce/utils"
	"github.com/lib/pq"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name, Email, Password string
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	hashPass, err := utils.HashPass(input.Password)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	query := `INSERT INTO users(name, email, password)
				Values ($1, $2, $3)
				RETURNING id, created_at`

	user := models.User{}
	err = db.DB.QueryRowx(query, input.Name, input.Email, hashPass).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			http.Error(w, "Email already taken", http.StatusBadRequest)
			return
		}
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email, Password string
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	user := &models.User{}

	// Check cache first
	cacheKey := fmt.Sprintf("user_email_%s", strings.TrimSpace(input.Email))
	err = utils.GetCache(cacheKey, user)
	if err != nil {
		// Get from database if not in cache
		err = db.DB.Get(user, "Select * FROM users where email=$1", strings.TrimSpace(input.Email))
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		// Cache the user data
		_ = utils.SetCache(cacheKey, user, time.Minute*15)
	}

	//	match pass
	err = utils.MatchPass(user.Password, input.Password)
	if err != nil {
		http.Error(w, "Wrong Email or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})

}
