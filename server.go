package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/ecofriends/authentication-backend/application"
	"github.com/ecofriends/authentication-backend/database"
)

/*
Connects to the PostGreSQL database using the credentials

Params:
  - No parameters

Returns:
  - A pointer to the SQL database
  - An error if the connection failed
*/
func initDatabaseConnection() (*sql.DB, error) {
	db, err := database.ConnectDatabase()
	if err != nil {
		msg := "[FAIL]: unable to connect database"
		return nil, fmt.Errorf("%s: %w", msg, err)
	}

	return db, nil
}

/*
Main is the server entry point

Objectives:
  - Initializes a database connection
  - Spins up a context channel to handle OS interrupts
  - Starts the server

Params:
  - No parameters

Returns:
  - No return value
*/

// @title Ecofriends Go Backend
// @version 1.0

// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name token
// @description Enter your auth cookie (e.g., "token=abc123")

// @host giving-vision-production.up.railway.app
// @BasePath /
func main() {
	appDatabase, err := initDatabaseConnection()
	if err != nil {
		log.Fatal(err)
	}

	// Setup a context channel to handle OS interrupts
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Initialize the application with the database connection
	app := application.New(appDatabase)

	// Start the app
	app.Start(ctx)
}
