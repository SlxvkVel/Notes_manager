package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"sync"
	"notes-service/internal/cache"
	"notes-service/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var (
	dbPool *pgxpool.Pool
	once   sync.Once
)

func initDB() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env: %s", err)
		return
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), 
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), 
		os.Getenv("DB_NAME"), os.Getenv("DB_SSLMODE"))

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Error parsing database config: %s", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Error creating database pool: %s", err)
	}

	dbPool = pool
	log.Println("Notes service: Database connection established")
}

func CreateNote(ctx context.Context, note models.Note) (int32, error) {
	once.Do(initDB)
	
	var id int32
	err := dbPool.QueryRow(ctx, 
		"INSERT INTO notes (title, content, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		note.Title, note.Content, note.UserID, note.CreatedAt, note.UpdatedAt).Scan(&id)
		
	if err != nil {
		return 0, fmt.Errorf("error inserting note: %w", err)
	}

	// Очищаем кэш пользователя
	cache.InvalidateUserCache(note.UserID)

	return id, nil
}

func GetUserNotes(ctx context.Context, userID int32) ([]models.Note, error) {
	once.Do(initDB)
	
	// Сначала пробуем получить из кэша
	cachedNotes, err := cache.GetCachedUserNotes(userID)
	if err == nil && len(cachedNotes) > 0 {
		log.Printf("✅ Got notes from cache for user %d", userID)
		return cachedNotes, nil
	}
	
	// Если нет в кэше - идем в БД
	rows, err := dbPool.Query(ctx, 
		"SELECT id, title, content, user_id, created_at, updated_at FROM notes WHERE user_id = $1 ORDER BY created_at DESC",
		userID)
		
	if err != nil {
		return nil, fmt.Errorf("error fetching notes: %w", err)
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.UserID, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning note: %w", err)
		}
		notes = append(notes, note)
	}

	// Сохраняем в кэш на 2 минуты
	if len(notes) > 0 {
		cache.CacheUserNotes(userID, notes, 2*time.Minute)
	}

	return notes, nil
}

func UpdateNote(ctx context.Context, noteID int32, userID int32, title, content string) error {
	once.Do(initDB)
	
	log.Printf("🔄 Storage UpdateNote - noteID: %d, userID: %d, title: %s", noteID, userID, title)
	
	result, err := dbPool.Exec(ctx,
		"UPDATE notes SET title = $1, content = $2, updated_at = NOW() WHERE id = $3 AND user_id = $4",
		title, content, noteID, userID)
		
	if err != nil {
		log.Printf("❌ Storage UpdateNote - DB error: %v", err)
		return fmt.Errorf("error updating note: %w", err)
	}

	rowsAffected := result.RowsAffected()
	log.Printf("📊 Storage UpdateNote - Rows affected: %d", rowsAffected)
	
	if rowsAffected == 0 {
		log.Printf("❌ Storage UpdateNote - No rows affected: noteID=%d, userID=%d", noteID, userID)
		return fmt.Errorf("note not found or access denied")
	}

	// Очищаем кэш пользователя после обновления
	cache.InvalidateUserCache(userID)
	log.Printf("✅ Storage UpdateNote - Successfully updated note %d for user %d", noteID, userID)
	
	return nil
}

func DeleteNote(ctx context.Context, noteID int32, userID int32) error {
	once.Do(initDB)
	
	log.Printf("🔄 Storage DeleteNote - noteID: %d, userID: %d", noteID, userID)
	
	result, err := dbPool.Exec(ctx,
		"DELETE FROM notes WHERE id = $1 AND user_id = $2",
		noteID, userID)
		
	if err != nil {
		log.Printf("❌ Storage DeleteNote - DB error: %v", err)
		return fmt.Errorf("error deleting note: %w", err)
	}

	rowsAffected := result.RowsAffected()
	log.Printf("📊 Storage DeleteNote - Rows affected: %d", rowsAffected)
	
	if rowsAffected == 0 {
		log.Printf("❌ Storage DeleteNote - No rows affected: noteID=%d, userID=%d", noteID, userID)
		return fmt.Errorf("note not found or access denied")
	}

	// Очищаем кэш пользователя после удаления
	cache.InvalidateUserCache(userID)
	log.Printf("✅ Storage DeleteNote - Successfully deleted note %d for user %d", noteID, userID)
	
	return nil
}

func GetNoteByID(ctx context.Context, noteID int32, userID int32) (*models.Note, error) {
	once.Do(initDB)
	
	log.Printf("🔍 Storage GetNoteByID - noteID: %d, userID: %d", noteID, userID)
	
	var note models.Note
	err := dbPool.QueryRow(ctx,
		"SELECT id, title, content, user_id, created_at, updated_at FROM notes WHERE id = $1 AND user_id = $2",
		noteID, userID).Scan(&note.ID, &note.Title, &note.Content, &note.UserID, &note.CreatedAt, &note.UpdatedAt)
		
	if err != nil {
		log.Printf("❌ Storage GetNoteByID - DB error: %v", err)
		return nil, fmt.Errorf("error fetching note: %w", err)
	}

	log.Printf("✅ Storage GetNoteByID - Found note: ID=%d, Title=%s", note.ID, note.Title)
	return &note, nil
}