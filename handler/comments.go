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

type Comment struct {
	repo *repository.PostGreSQL
}

func (comment *Comment) New(repo *repository.PostGreSQL) {
	comment.repo = repo
}

func (comment *Comment) CreateComment(w http.ResponseWriter, r *http.Request) {
	var body = util.CreateCommentRequestBody{}

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

	theComment, err := comment.repo.CreateComment(context.Background(), body.UserID.String(), body.PostID, body.Text)
	if err != nil {
		util.JsonResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	util.JsonResponse(w, "Successfully created comment", http.StatusOK, theComment)
}

func (comment *Comment) DeleteComment(w http.ResponseWriter, r *http.Request) {
	var body = util.DeleteCommentRequestBody{}

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

	if err := comment.repo.DeleteComment(context.Background(), body.CommentID, body.UserID.String()); err != nil {
		util.JsonResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	util.JsonResponse(w, "Successfully deleted comment", http.StatusOK, nil)
}

func (comment *Comment) UpdateComment(w http.ResponseWriter, r *http.Request) {
	var body = util.UpdateCommentRequestBody{}

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

	if err := comment.repo.UpdateComment(context.Background(), body.CommentID, body.UserID.String(), body.Text); err != nil {
		util.JsonResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	util.JsonResponse(w, "Successfully updated comment", http.StatusOK, nil)
}

func (comment *Comment) GetCommentByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	commentIdInt, err := strconv.Atoi(id)
	if err != nil {
		util.JsonResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	theComment, err := comment.repo.GetCommentByID(context.Background(), commentIdInt)
	if err != nil {
		util.JsonResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	util.JsonResponse(w, "Successfully got comments by id", http.StatusOK, theComment)
}

func (comment *Comment) GetCommentsByPost(w http.ResponseWriter, r *http.Request) {
	limit := chi.URLParam(r, "limit")
	offset := chi.URLParam(r, "offset")
	postId := chi.URLParam(r, "post_id")

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

	postIdInt, err := strconv.Atoi(postId)
	if err != nil {
		util.JsonResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	comments, err := comment.repo.GetCommentsByPost(context.Background(), postIdInt, limitInt, offsetInt)
	if err != nil {
		util.JsonResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	util.JsonResponse(w, "Successfully got comments by post", http.StatusOK, comments)
}
