package main

import (
	"log"
	"net/http"

	"notes-service/internal/cache"
	"notes-service/internal/handlers"
)

func main() {
    cache.InitRedis()
    http.HandleFunc("/api/notes", handlers.CreateNoteHandler)          
    http.HandleFunc("/api/notes/list", handlers.GetNotesHandler)       
    http.HandleFunc("/api/notes/", handlers.NoteDetailHandler)         
    http.HandleFunc("/health", handlers.HealthHandler)
    http.HandleFunc("/api/notes/update", handlers.UpdateNoteHandler)
    http.HandleFunc("/api/notes/delete", handlers.DeleteNoteHandler)
    
    log.Println("📝 Notes service starting on port 8081...")
    log.Fatal(http.ListenAndServe(":8081", nil))
}