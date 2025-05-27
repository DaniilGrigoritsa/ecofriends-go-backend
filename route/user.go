package route

import (
	"database/sql"

	"github.com/ecofriends/authentication-backend/handler"
	"github.com/ecofriends/authentication-backend/middleware"
	repository "github.com/ecofriends/authentication-backend/repository"
	"github.com/go-chi/chi/v5"
)

func LoadUserRoutes(router chi.Router, db *sql.DB) {
	user := &handler.User{}
	user.New(&repository.PostGreSQL{Database: db})

	router.Get("/", user.Home)
	router.With(middleware.AuthenticateMiddleware).Get("/{id}", user.GetUserByID)
}
