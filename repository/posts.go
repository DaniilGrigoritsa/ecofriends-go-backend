package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ecofriends/authentication-backend/model"
	_ "github.com/lib/pq"
)

func (repo *PostGreSQL) CreatePost(ctx context.Context, userID string, text string) (model.Post, error) {
	tx, err := repo.Database.BeginTx(ctx, nil)
	if err != nil {
		return model.Post{}, fmt.Errorf("could not begin transaction: %w", err)
	}

	defer func() {
		if rErr := tx.Rollback(); rErr != nil && err == nil {
			err = fmt.Errorf("rollback failed: %w", rErr)
		}
	}()

	query := `
		INSERT INTO posts (user_id, text, created_at)
		VALUES ($1, $2, $3)
		RETURNING id, like_count, created_at
	`

	var createdPost model.Post
	err = tx.QueryRowContext(ctx, query,
		userID,
		text,
		time.Now(),
	).Scan(
		&createdPost.ID,
		&createdPost.LikeCount,
		&createdPost.CreatedAt,
	)

	if err != nil {
		return model.Post{}, fmt.Errorf("could not create post: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return model.Post{}, fmt.Errorf("could not commit transaction: %w", err)
	}

	createdPost.UserID = userID
	createdPost.Text = text

	return createdPost, nil
}

func (repo *PostGreSQL) DeletePost(ctx context.Context, postID int, userID string) error {
	tx, err := repo.Database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}

	defer func() {
		if rErr := tx.Rollback(); rErr != nil && err == nil {
			err = fmt.Errorf("rollback failed: %w", rErr)
		}
	}()

	query := `
		DELETE FROM posts
		WHERE id = $1 AND user_id = $2
	`

	result, err := tx.ExecContext(ctx, query, postID, userID)
	if err != nil {
		return fmt.Errorf("could not delete post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("post not found or not owned by user")
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func (repo *PostGreSQL) GetAllPosts(ctx context.Context, limit int, offset int) ([]model.Post, error) {
	query := `
		SELECT id, user_id, text, like_count, created_at, updated_at
		FROM posts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := repo.Database.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not query posts: %w", err)
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		var updatedAt sql.NullTime

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Text,
			&post.LikeCount,
			&post.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			log.Printf("Error scanning post row: %v", err)
			continue
		}

		if updatedAt.Valid {
			post.UpdatedAt = &updatedAt.Time
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating posts: %w", err)
	}

	return posts, nil
}

func (repo *PostGreSQL) GetPostsByUser(ctx context.Context, userID string, limit int, offset int) ([]model.Post, error) {
	query := `
		SELECT id, user_id, text, like_count, created_at, updated_at
		FROM posts
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := repo.Database.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not query user posts: %w", err)
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		var updatedAt sql.NullTime

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Text,
			&post.LikeCount,
			&post.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			log.Printf("Error scanning user post row: %v", err)
			continue
		}

		if updatedAt.Valid {
			post.UpdatedAt = &updatedAt.Time
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user posts: %w", err)
	}

	return posts, nil
}

func (repo *PostGreSQL) GetPostByID(ctx context.Context, postID int) (model.Post, error) {
	query := `
		SELECT id, user_id, text, like_count, created_at, updated_at
		FROM posts
		WHERE id = $1
	`

	var post model.Post
	var updatedAt sql.NullTime

	err := repo.Database.QueryRowContext(ctx, query, postID).Scan(
		&post.ID,
		&post.UserID,
		&post.Text,
		&post.LikeCount,
		&post.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.Post{}, fmt.Errorf("post not found")
		}
		return model.Post{}, fmt.Errorf("could not get post: %w", err)
	}

	if updatedAt.Valid {
		post.UpdatedAt = &updatedAt.Time
	}

	return post, nil
}
