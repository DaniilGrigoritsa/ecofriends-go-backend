package model

import "time"

/*
Comment model struct

Fields:
  - ID:         int           - Unique identifier for the comment
  - UserID:     string        - ID of the user who created the comment
  - PostID:     int           - ID of the post being commented on
  - Text:       string        - Content of the comment
  - CreatedAt:  time.Time     - When the comment was created
  - UpdatedAt:  *time.Time    - When the comment was last updated (nullable)
*/
type Comment struct {
	ID        int        `json:"id"`
	UserID    string     `json:"user_id"`
	PostID    int        `json:"post_id"`
	Text      string     `json:"text"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// CommentWithUser includes basic user information with the comment
type CommentWithUser struct {
	Comment
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
}
