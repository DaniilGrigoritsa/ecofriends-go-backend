package route

import (
	"database/sql"

	handler "github.com/ecofriends/authentication-backend/handler/auth"
	repository "github.com/ecofriends/authentication-backend/repository"
	"github.com/ecofriends/authentication-backend/service"
	"github.com/go-chi/chi/v5"
)

func LoadAuthRoutes(router chi.Router, db *sql.DB) {
	authDBService := &service.DatabaseProvider{}
	authDBService.New(&repository.PostGreSQL{Database: db})

	authHandler := &handler.AuthHandler{}
	authHandler.WithService(authDBService)

	router.Get("/", authHandler.Home)
	router.Post("/sign-up", authHandler.SignUp)
	router.Post("/sign-in", authHandler.SignIn)
	router.Post("/sign-out", authHandler.SignOut)
	router.Get("/oauth/google", authHandler.GoogleSignIn)
	router.Get("/oauth/google/callback", authHandler.GoogleSignInCallback)
	router.Get("/oauth/{x}/failure", authHandler.OAuthFailure)
}
