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

func (repo *PostGreSQL) CreateComment(ctx context.Context, userID string, postID int, text string) (model.Comment, error) {
	tx, err := repo.Database.BeginTx(ctx, nil)
	if err != nil {
		return model.Comment{}, fmt.Errorf("could not begin transaction: %w", err)
	}

	defer func() {
		if rErr := tx.Rollback(); rErr != nil && err == nil {
			err = fmt.Errorf("rollback failed: %w", rErr)
		}
	}()

	query := `
		INSERT INTO comments (user_id, post_id, text, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	var createdComment model.Comment
	err = tx.QueryRowContext(ctx, query,
		userID,
		postID,
		text,
		time.Now(),
	).Scan(
		&createdComment.ID,
		&createdComment.CreatedAt,
	)

	if err != nil {
		return model.Comment{}, fmt.Errorf("could not create comment: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return model.Comment{}, fmt.Errorf("could not commit transaction: %w", err)
	}

	createdComment.UserID = userID
	createdComment.PostID = postID
	createdComment.Text = text

	return createdComment, nil
}

func (repo *PostGreSQL) DeleteComment(ctx context.Context, commentID int, userID string) error {
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
		DELETE FROM comments
		WHERE id = $1 AND user_id = $2
	`

	result, err := tx.ExecContext(ctx, query, commentID, userID)
	if err != nil {
		return fmt.Errorf("could not delete comment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comment not found or not owned by user")
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func (repo *PostGreSQL) GetCommentsByPost(ctx context.Context, postID int, limit int, offset int) ([]model.CommentWithUser, error) {
	query := `
		SELECT c.id, c.user_id, c.post_id, c.text, c.created_at, c.updated_at, u.username
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := repo.Database.QueryContext(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not query comments: %w", err)
	}
	defer rows.Close()

	var comments []model.CommentWithUser
	for rows.Next() {
		var comment model.CommentWithUser
		var updatedAt sql.NullTime

		err := rows.Scan(
			&comment.ID,
			&comment.UserID,
			&comment.PostID,
			&comment.Text,
			&comment.CreatedAt,
			&updatedAt,
			&comment.Username,
		)
		if err != nil {
			log.Printf("Error scanning comment row: %v", err)
			continue
		}

		if updatedAt.Valid {
			comment.UpdatedAt = &updatedAt.Time
		}

		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating comments: %w", err)
	}

	return comments, nil
}

func (repo *PostGreSQL) GetCommentByID(ctx context.Context, commentID int) (model.Comment, error) {
	query := `
		SELECT id, user_id, post_id, text, created_at, updated_at
		FROM comments
		WHERE id = $1
	`

	var comment model.Comment
	var updatedAt sql.NullTime

	err := repo.Database.QueryRowContext(ctx, query, commentID).Scan(
		&comment.ID,
		&comment.UserID,
		&comment.PostID,
		&comment.Text,
		&comment.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.Comment{}, fmt.Errorf("comment not found")
		}
		return model.Comment{}, fmt.Errorf("could not get comment: %w", err)
	}

	if updatedAt.Valid {
		comment.UpdatedAt = &updatedAt.Time
	}

	return comment, nil
}

func (repo *PostGreSQL) UpdateComment(ctx context.Context, commentID int, userID string, text string) error {
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
		UPDATE comments
		SET text = $1, updated_at = $2
		WHERE id = $3 AND user_id = $4
	`

	result, err := tx.ExecContext(ctx, query,
		text,
		time.Now(),
		commentID,
		userID,
	)

	if err != nil {
		return fmt.Errorf("could not update comment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comment not found or not owned by user")
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}
