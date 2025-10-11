package models

import "time"

type Note struct {
	ID        int32     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int32     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}