package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"auth-service/internal/models"
	"auth-service/internal/storage"
	"auth-service/internal/tools"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}


	hashedPassword, err := tools.PasswordToHash(data.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	
	user := models.User{
		Username: data.Username,
		Email:    data.Email,
		Password: hashedPassword,
	}


	id, err := storage.AddUser(r.Context(), user)
	if err != nil {
		log.Println("Error adding user:", err)
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	
	tools.MakeCookieAfterLogin(w, id, data.Username, data.Email)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}


	user, err := storage.GetUserByEmail(r.Context(), data.Email)
	if err != nil {
		log.Println("Error retrieving user:", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}


	if !tools.ValidatePassword(data.Password, user.Password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	
	tools.MakeCookieAfterLogin(w, user.ID, user.Username, user.Email)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully",
	})
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := tools.GetUserFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}