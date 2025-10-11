package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"notes-service/internal/models"
	"notes-service/internal/storage"
	"notes-service/internal/tools"
)

func CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	
	userID, err := tools.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	
	note := models.Note{
		Title:     req.Title,
		Content:   req.Content,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	
	id, err := storage.CreateNote(r.Context(), note)
	if err != nil {
		log.Println("Error creating note:", err)
		http.Error(w, "Error creating note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Note created successfully",
		"note_id": id,
	})
}

func GetNotesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	
	userID, err := tools.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	
	notes, err := storage.GetUserNotes(r.Context(), userID)
	if err != nil {
		log.Println("Error fetching notes:", err)
		http.Error(w, "Error fetching notes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"notes": notes,
	})
}

func UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}


	userID, err := tools.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	
	noteIDStr := r.URL.Query().Get("id")
	noteID, err := strconv.Atoi(noteIDStr)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	
	err = storage.UpdateNote(r.Context(), int32(noteID), userID, req.Title, req.Content)
	if err != nil {
		log.Println("Error updating note:", err)
		http.Error(w, "Error updating note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Note updated successfully",
	})
}

func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	
	userID, err := tools.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	noteIDStr := r.URL.Query().Get("id")
	noteID, err := strconv.Atoi(noteIDStr)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	err = storage.DeleteNote(r.Context(), int32(noteID), userID)
	if err != nil {
		log.Println("Error deleting note:", err)
		http.Error(w, "Error deleting note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Note deleted successfully",
	})
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}