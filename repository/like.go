package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ecofriends/authentication-backend/model"
	_ "github.com/lib/pq"
)

func (repo *PostGreSQL) LikePost(ctx context.Context, userID string, postID int) error {
	tx, err := repo.Database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}

	defer func() {
		if rErr := tx.Rollback(); rErr != nil && err == nil {
			err = fmt.Errorf("rollback failed: %w", rErr)
		}
	}()

	// First check if the user already liked the post
	hasLiked, err := repo.HasLiked(ctx, userID, postID)
	if err != nil {
		return fmt.Errorf("could not check like status: %w", err)
	}
	if hasLiked {
		return fmt.Errorf("user already liked this post")
	}

	query := `
        INSERT INTO post_likes (user_id, post_id, created_at)
        VALUES ($1, $2, $3)
    `

	_, err = tx.ExecContext(ctx, query, userID, postID, time.Now())
	if err != nil {
		return fmt.Errorf("could not like post: %w", err)
	}

	// Update the like_count in posts table
	updateQuery := `
        UPDATE posts
        SET like_count = like_count + 1
        WHERE id = $1
    `
	_, err = tx.ExecContext(ctx, updateQuery, postID)
	if err != nil {
		return fmt.Errorf("could not update post like count: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func (repo *PostGreSQL) UnlikePost(ctx context.Context, userID string, postID int) error {
	tx, err := repo.Database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}

	defer func() {
		if rErr := tx.Rollback(); rErr != nil && err == nil {
			err = fmt.Errorf("rollback failed: %w", rErr)
		}
	}()

	// First check if the user has liked the post
	hasLiked, err := repo.HasLiked(ctx, userID, postID)
	if err != nil {
		return fmt.Errorf("could not check like status: %w", err)
	}
	if !hasLiked {
		return fmt.Errorf("user hasn't liked this post")
	}

	query := `
        DELETE FROM post_likes
        WHERE user_id = $1 AND post_id = $2
    `

	result, err := tx.ExecContext(ctx, query, userID, postID)
	if err != nil {
		return fmt.Errorf("could not unlike post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("like not found")
	}

	// Update the like_count in posts table
	updateQuery := `
        UPDATE posts
        SET like_count = like_count - 1
        WHERE id = $1
    `
	_, err = tx.ExecContext(ctx, updateQuery, postID)
	if err != nil {
		return fmt.Errorf("could not update post like count: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func (repo *PostGreSQL) GetLikeCount(ctx context.Context, postID int) (int, error) {
	query := `
        SELECT COUNT(*)
        FROM post_likes
        WHERE post_id = $1
    `

	var count int
	err := repo.Database.QueryRowContext(ctx, query, postID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("could not get like count: %w", err)
	}

	return count, nil
}

func (repo *PostGreSQL) HasLiked(ctx context.Context, userID string, postID int) (bool, error) {
	query := `
        SELECT EXISTS (
            SELECT 1 FROM post_likes
            WHERE user_id = $1 AND post_id = $2
        )
    `

	var exists bool
	err := repo.Database.QueryRowContext(ctx, query, userID, postID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("could not check like status: %w", err)
	}

	return exists, nil
}

func (repo *PostGreSQL) GetLikesByUser(ctx context.Context, userID string, limit int, offset int) ([]model.PostLike, error) {
	query := `
        SELECT user_id, post_id, created_at
        FROM post_likes
        WHERE user_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := repo.Database.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not query likes: %w", err)
	}
	defer rows.Close()

	var likes []model.PostLike
	for rows.Next() {
		var like model.PostLike
		err := rows.Scan(
			&like.UserID,
			&like.PostID,
			&like.CreatedAt,
		)
		if err != nil {
			log.Printf("Error scanning like row: %v", err)
			continue
		}
		likes = append(likes, like)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating likes: %w", err)
	}

	return likes, nil
}
