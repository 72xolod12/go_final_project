package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


var todoPassword string


func InitAuthConfig() {
	todoPassword = os.Getenv("TODO_PASSWORD")
}

type AuthRequest struct {
	Password string `json:"password"`
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendError(w, "Ошибка формата запроса", http.StatusBadRequest)
		return
	}

	
	if req.Password != todoPassword {
		SendError(w, "Неверный пароль", http.StatusUnauthorized)
		return
	}

	hash := sha256.Sum256([]byte(todoPassword))
	passHash := hex.EncodeToString(hash[:])

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": passHash,
		"exp":  time.Now().Add(time.Hour * 8).Unix(),
	})

	tokenString, err := token.SignedString([]byte(todoPassword))
	if err != nil {
		SendError(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		if len(todoPassword) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			tokenString := cookie.Value
			claims := jwt.MapClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(todoPassword), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			hash := sha256.Sum256([]byte(todoPassword))
			expectedHash := hex.EncodeToString(hash[:])
			if claims["hash"] != expectedHash {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}