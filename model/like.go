package model

import "time"

/*
PostLike model struct

Fields:
  - UserID:     string    - ID of the user who liked the post
  - PostID:     int       - ID of the post being liked
  - CreatedAt:  time.Time - When the like was created
*/
type PostLike struct {
	UserID    string    `json:"user_id"`
	PostID    int       `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}
