package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"auth-service/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func PasswordToHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ValidatePassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func MakeCookieAfterLogin(w http.ResponseWriter, id int32, username, email string) {
	if err := godotenv.Load(); err != nil {
		http.Error(w, "Error loading environment variables", http.StatusInternalServerError)
		return
	}

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       id,
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Error signing token", http.StatusInternalServerError)
		return
	}


	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   24 * 60 * 60,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Login successful",
		"user_id":  id,
		"username": username,
	})
}

func ExtractTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("token")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func ValidateToken(inputToken string) (jwt.Claims, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(inputToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	
	if err != nil {
		return nil, err
	}

	return token.Claims, nil
}

func GetUserFromToken(r *http.Request) (*models.User, error) {
	token := ExtractTokenFromCookie(r)
	claims, err := ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claimsMap, ok := claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	id, ok := claimsMap["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	username, ok := claimsMap["username"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid username in token")
	}

	email, ok := claimsMap["email"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid email in token")
	}

	return &models.User{
		ID:       int32(id),
		Username: username,
		Email:    email,
	}, nil
}