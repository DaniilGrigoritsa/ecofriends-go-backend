package route

import (
	"database/sql"

	"github.com/ecofriends/authentication-backend/handler"
	"github.com/ecofriends/authentication-backend/middleware"
	repository "github.com/ecofriends/authentication-backend/repository"
	"github.com/go-chi/chi/v5"
)

func LoadCommentRoutes(router chi.Router, db *sql.DB) {
	comment := &handler.Comment{}
	comment.New(&repository.PostGreSQL{Database: db})

	router.Get("/{id}", comment.GetCommentByID)
	router.Get("/post", comment.GetCommentsByPost)

	router.With(middleware.AuthenticateMiddleware).Post("/create", comment.CreateComment)
	router.With(middleware.AuthenticateMiddleware).Put("/update", comment.UpdateComment)
	router.With(middleware.AuthenticateMiddleware).Delete("/delete", comment.DeleteComment)
}
