package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/coursework/monitor-backend/internal/db"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func getJWTKey() []byte {
    key := os.Getenv("JWT_SECRET")
    if key == "" {
        key = "default_secret" // fallback for tests
    }
    return []byte(key)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var user User
	err := db.Conn.QueryRow(context.Background(), "SELECT id, username, password_hash, role FROM users WHERE username=$1", creds.Username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.RegisteredClaims{
		Subject:   user.Username,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		ID:        user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getJWTKey())
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

// return json
    w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TokenResponse{Token: tokenString, Role: user.Role})
}

// Middleware to protect routes
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		tokenStr := parts[1]
		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return getJWTKey(), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "role", claims.ID)
		ctx = context.WithValue(ctx, "username", claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// Role Middleware
func RoleMiddleware(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value("role").(string)
		if role != "admin" && role != requiredRole {
			http.Error(w, "Unauthorized role", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
