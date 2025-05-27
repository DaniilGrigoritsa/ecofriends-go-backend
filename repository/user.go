package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/ecofriends/authentication-backend/model"
	"github.com/ecofriends/authentication-backend/util"
	_ "github.com/lib/pq"
)

func (repo *PostGreSQL) InsertUser(ctx context.Context, user model.User) error {
	// Begin a new database transaction
	tx, err := repo.Database.BeginTx(ctx, nil)
	if err != nil {
		_, err = util.Fail(err, "[FAIL]: could not begin database transaction")
		return err
	}

	// Rollback transaction incase of failure (deferred)
	defer func() {
		if rErr := tx.Rollback(); rErr != nil && err == nil {
			err = fmt.Errorf("[FAIL]: rollback failed: %w", rErr)
		}
	}()

	// Hash the user password
	user.Password, err = util.GenerateHash(user.Password, util.DefaultHashCost)
	if err != nil {
		return err
	}

	log.Println("Hashed user password:", user.Password)

	// Construct a query to insert the user from the model data
	var insertQuery = `
		INSERT INTO users (id, username, email, password)
		VALUES ($1, $2, $3, $4)
	`

	// Execute the insertion query
	_, err = tx.ExecContext(ctx, insertQuery, user.ID, user.Username, user.Email, user.Password)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("[FAIL]: could not execute insert query")
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("[FAIL]: could not commit transaction")
	}

	return nil
}

func (repo *PostGreSQL) UserExists(ctx context.Context, email string, username string) (bool, error) {
	// Construct a query to check if a column with the email exists
	var checkUserExistsQuery = `
		SELECT EXISTS (SELECT 1 FROM users WHERE email = $1 OR username = $2)
	`

	// Exists is set to false by default
	var exists = false

	// Check if the user is already stored in the database
	err := repo.Database.QueryRowContext(ctx, checkUserExistsQuery, email, username).Scan(&exists)
	if err != nil {
		log.Println(err)

		// Return false if the table doesn't exist
		if strings.Contains(err.Error(), "does not exist") {
			return false, nil
		}

		return false, fmt.Errorf("[FAIL]: could not check if user already exists")
	}

	return exists, nil
}

func (repo *PostGreSQL) GetUserByID(ctx context.Context, id string) (model.User, error) {
	// Begin the database transaction
	tx, err := repo.Database.BeginTx(ctx, nil)
	if err != nil {
		_, err = util.Fail(err, "[FAIL]: could not begin database transaction")
		return model.User{}, err
	}

	// Construct a query to return the user details from the provided id
	var getUserByIDQuery = `
		SELECT id, username, email, password FROM users WHERE id = $1
	`

	// Allocate memory for the user model data
	var user model.User

	// Execute the query, returns the row with the details
	err = tx.QueryRowContext(ctx, getUserByIDQuery, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return model.User{}, fmt.Errorf("[FAIL]: user with ID %s not found", id)
		}
		return model.User{}, fmt.Errorf("[FAIL]: could not execute query: %w", err)
	}

	return user, nil
}

func (repo *PostGreSQL) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	// Begin the transaction
	tx, err := repo.Database.BeginTx(ctx, nil)
	if err != nil {
		_, err = util.Fail(err, "[FAIL]: could not begin database transaction")
		return model.User{}, err
	}

	// construct a query to return the data model using the email provided
	var getUserByIDQuery = `
		SELECT id, username, email, password FROM users WHERE email = $1
	`

	// Allocate memory for the user data
	var user model.User

	// Execute the query
	err = tx.QueryRowContext(ctx, getUserByIDQuery, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return model.User{}, fmt.Errorf("[FAIL]: user with email %s not found", email)
		}
		return model.User{}, fmt.Errorf("[FAIL]: could not execute query: %w", err)
	}

	return user, nil
}
