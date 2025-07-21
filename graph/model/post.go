package model

type Post struct {
	ID              string     `json:"id" db:"id"`
	Title           string     `json:"title" db:"title"`
	Content         string     `json:"content" db:"content"`
	Author          string     `json:"author" db:"author"`
	CreatedAt       string     `json:"createdAt" db:"created_at"` // 👈 важно
	CommentsEnabled bool       `json:"commentsEnabled" db:"comments_enabled"`
	Comments        []*Comment `json:"comments"`
}
