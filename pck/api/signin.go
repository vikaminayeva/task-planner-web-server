package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type SigninRequest struct {
	Password string `json:"password"`
}

type SigninResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func getPasswordHash() string {
	pass := os.Getenv("TODO_PASSWORD")
	hash := sha256.Sum256([]byte(pass))
	return hex.EncodeToString(hash[:])
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	var req SigninRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJson(w, http.StatusBadRequest, SigninResponse{Error: "JSON deserialization error: " + err.Error()})
		return
	}

	expectedPassword := os.Getenv("TODO_PASSWORD")
	if req.Password != expectedPassword {
		writeJson(w, http.StatusBadRequest, SigninResponse{Error: "Wrong password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": getPasswordHash(),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(getPasswordHash()))
	if err != nil {
		writeJson(w, http.StatusInternalServerError, SigninResponse{Error: "Token generation error: " + err.Error()})
		return
	}

	writeJson(w, http.StatusOK, SigninResponse{Token: tokenString})
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			tokenStr := cookie.Value
			claims := jwt.MapClaims{}

			token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return []byte(getPasswordHash()), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			tokenHash, ok := claims["hash"].(string)
			if !ok || tokenHash != getPasswordHash() {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
