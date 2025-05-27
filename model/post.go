package model

import "time"

/*
Post model struct

Fields:
  - ID:         int           - Unique identifier for the post
  - UserID:     string        - ID of the user who created the post (references users.id)
  - Text:       string        - Content of the post
  - LikeCount:  int           - Number of likes the post has received
  - CreatedAt:  time.Time     - When the post was created
  - UpdatedAt:  *time.Time    - When the post was last updated (nullable)
*/
type Post struct {
	ID        int        `json:"id"`
	UserID    string     `json:"user_id"`
	Text      string     `json:"text"`
	LikeCount int        `json:"like_count"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"` // Pointer to allow null
}
