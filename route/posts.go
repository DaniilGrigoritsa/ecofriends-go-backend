package route

import (
	"database/sql"

	"github.com/ecofriends/authentication-backend/handler"
	"github.com/ecofriends/authentication-backend/middleware"
	repository "github.com/ecofriends/authentication-backend/repository"
	"github.com/go-chi/chi/v5"
)

func LoadPostRoutes(router chi.Router, db *sql.DB) {
	post := &handler.Post{}
	post.New(&repository.PostGreSQL{Database: db})

	router.Get("/{id}", post.GetPostByID)
	router.Get("/all", post.GetAllPosts)
	router.Get("/user", post.GetPostsByUser)

	router.With(middleware.AuthenticateMiddleware).Post("/create", post.CreatePost)
	router.With(middleware.AuthenticateMiddleware).Delete("/delete", post.DeletePost)
}
