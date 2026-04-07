package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"monitoring-backend/db"
	"monitoring-backend/models"
	"monitoring-backend/utils"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	err := db.Pool.QueryRow(context.Background(),
		"SELECT id, username, password_hash, role FROM users WHERE username=$1", req.Username).
		Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)

	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(user.Username, user.Role)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	res := models.AuthResponse{
		Token:    token,
		Role:     user.Role,
		Username: user.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	_, err = db.Pool.Exec(context.Background(),
		"INSERT INTO users (username, password_hash, role) VALUES ($1, $2, 'user')",
		req.Username, hash)

	if err != nil {
		http.Error(w, "failed to create user, username might be taken", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
