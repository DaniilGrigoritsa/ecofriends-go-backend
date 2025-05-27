package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	repository "github.com/ecofriends/authentication-backend/repository"
	"github.com/ecofriends/authentication-backend/util"
	"github.com/go-chi/chi/v5"
)

type Like struct {
	repo *repository.PostGreSQL
}

func (like *Like) New(repo *repository.PostGreSQL) {
	like.repo = repo
}

func (like *Like) LikePost(w http.ResponseWriter, r *http.Request) {
	var body = util.LikePostRequestBody{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		util.JsonResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	userID, err := util.ExtractUserIDFromClaims(r.Context())
	if err != nil {
		util.JsonResponse(w, err.Error(), http.StatusUnauthorized, nil)
		return
	}

	if userID != body.UserID.String() {
		msg := "Forbidden: Access to this resource is denied"
		util.JsonResponse(w, msg, http.StatusForbidden, nil)
		return
	}

	if err := like.repo.LikePost(context.Background(), body.UserID.String(), body.PostID); err != nil {
		util.JsonResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	util.JsonResponse(w, "Successfully liked post", http.StatusOK, nil)
}

func (like *Like) UnlikePost(w http.ResponseWriter, r *http.Request) {
	var body = util.LikePostRequestBody{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		util.JsonResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	userID, err := util.ExtractUserIDFromClaims(r.Context())
	if err != nil {
		util.JsonResponse(w, err.Error(), http.StatusUnauthorized, nil)
		return
	}

	if userID != body.UserID.String() {
		msg := "Forbidden: Access to this resource is denied"
		util.JsonResponse(w, msg, http.StatusForbidden, nil)
		return
	}

	if err := like.repo.UnlikePost(context.Background(), body.UserID.String(), body.PostID); err != nil {
		util.JsonResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	util.JsonResponse(w, "Successfully unliked post", http.StatusOK, nil)
}

func (like *Like) GetLikeCount(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "post_id")

	postIdInt, err := strconv.Atoi(postId)
	if err != nil {
		util.JsonResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	likes, err := like.repo.GetLikeCount(context.Background(), postIdInt)
	if err != nil {
		util.JsonResponse(w, err.Error(), http.StatusBadRequest, likes)
		return
	}

	util.JsonResponse(w, "Successfully got like count for the post", http.StatusOK, likes)
}

func (like *Like) HasLiked(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "post_id")
	userId := chi.URLParam(r, "user_id")

	postIdInt, err := strconv.Atoi(postId)
	if err != nil {
		util.JsonResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	hasLiked, err := like.repo.HasLiked(context.Background(), userId, postIdInt)
	if err != nil {
		util.JsonResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	util.JsonResponse(w, "Successfully got whether user has liked the post", http.StatusOK, hasLiked)
}

func (like *Like) GetLikesByUser(w http.ResponseWriter, r *http.Request) {
	limit := chi.URLParam(r, "limit")
	offset := chi.URLParam(r, "offset")
	userId := chi.URLParam(r, "user_id")

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		util.JsonResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		util.JsonResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	postsLike, err := like.repo.GetLikesByUser(context.Background(), userId, limitInt, offsetInt)
	if err != nil {
		util.JsonResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	util.JsonResponse(w, "Successfully got likes by the user", http.StatusOK, postsLike)
}
