package route

import (
	"database/sql"

	"github.com/ecofriends/authentication-backend/handler"
	"github.com/ecofriends/authentication-backend/middleware"
	repository "github.com/ecofriends/authentication-backend/repository"
	"github.com/go-chi/chi/v5"
)

func LoadLikeRoutes(router chi.Router, db *sql.DB) {
	like := &handler.Like{}
	like.New(&repository.PostGreSQL{Database: db})

	router.Get("/count/post_id={post_id}", like.GetLikeCount)
	router.Get("/has_liked/post_id={post_id}/user_id={user_id}", like.HasLiked)
	router.Get("/user_id={user_id}/limit={limit}/offset={offset}", like.GetLikesByUser)

	router.With(middleware.AuthenticateMiddleware).Post("/like", like.LikePost)
	router.With(middleware.AuthenticateMiddleware).Post("/unlike", like.UnlikePost)
}
