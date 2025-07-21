package model

type Comment struct {
	ID        string  `json:"id" db:"id"`
	PostID    string  `json:"postID" db:"post_id"` // ✅ нужно добавить
	ParentID  *string `json:"parentID" db:"parent_id"`
	Text      string  `json:"text" db:"text"`
	Author    string  `json:"author" db:"author"`
	CreatedAt string  `json:"createdAt" db:"created_at"`
}
