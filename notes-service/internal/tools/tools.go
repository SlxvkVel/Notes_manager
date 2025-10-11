package tools

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func ExtractUserIDFromToken(r *http.Request) (int32, error) {
	tokenString := extractTokenFromHeader(r)
	if tokenString == "" {
		return 0, fmt.Errorf("no token provided")
	}

	userID, err := validateTokenAndExtractUserID(tokenString)
	if err != nil {
		return 0, fmt.Errorf("invalid token: %w", err)
	}

	return userID, nil
}

func extractTokenFromHeader(r *http.Request) string {
	// Проверяем Authorization header
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}

	// Проверяем cookie
	cookie, err := r.Cookie("token")
	if err == nil {
		return cookie.Value
	}

	return ""
}

func validateTokenAndExtractUserID(tokenString string) (int32, error) {
	if err := godotenv.Load(); err != nil {
		return 0, err
	}

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["id"]
	if !ok {
		return 0, fmt.Errorf("user ID not found in token")
	}

	// Конвертируем userID в int32
	switch v := userID.(type) {
	case float64:
		return int32(v), nil
	case int:
		return int32(v), nil
	case string:
		id, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("invalid user ID format")
		}
		return int32(id), nil
	default:
		return 0, fmt.Errorf("unexpected user ID type")
	}
}