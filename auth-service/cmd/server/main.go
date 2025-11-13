package main

import (
	"log"
	"net/http"
	"auth-service/internal/handlers"
)

func main() {
	http.HandleFunc("/api/auth/register", handlers.RegisterHandler)
	http.HandleFunc("/api/auth/login", handlers.LoginHandler)
	http.HandleFunc("/api/auth/logout", handlers.LogoutHandler)
	http.HandleFunc("/api/auth/me", handlers.MeHandler)
	http.HandleFunc("/health", handlers.HealthHandler)

	log.Println(" Auth service starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}