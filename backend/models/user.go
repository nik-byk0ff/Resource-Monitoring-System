package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username" validate:"required,min=3"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type AuthRequest struct {
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
}

type AuthResponse struct {
	Token    string `json:"token"`
	Role     string `json:"role"`
	Username string `json:"username"`
}
