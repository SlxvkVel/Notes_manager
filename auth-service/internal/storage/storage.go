package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"auth-service/internal/models"
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
	log.Println("Database connection established")
}


func AddUser(ctx context.Context, user models.User) (int32, error) {
	once.Do(initDB)
	
	var id int32
	err := dbPool.QueryRow(ctx, 
		"INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id",
		user.Username, user.Email, user.Password).Scan(&id)
		
	if err != nil {
		return 0, fmt.Errorf("error inserting user: %w", err)
	}
	return id, nil
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	once.Do(initDB)
	
	var user models.User
	err := dbPool.QueryRow(ctx, 
		"SELECT id, username, email, password FROM users WHERE email = $1", 
		email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
		
	if err != nil {
		return nil, fmt.Errorf("error fetching user by email: %w", err)
	}
	return &user, nil
}

func GetUserByID(ctx context.Context, id int32) (*models.User, error) {
	once.Do(initDB)
	
	var user models.User
	err := dbPool.QueryRow(ctx, 
		"SELECT id, username, email, password FROM users WHERE id = $1", 
		id).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
		
	if err != nil {
		return nil, fmt.Errorf("error fetching user by id: %w", err)
	}
	return &user, nil
}