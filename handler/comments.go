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

// @Summary Create comment
// @Description Creates a new comment on a post
// @Tags comments
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param comment body util.CreateCommentRequestBody true "Create comment body"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 401 {object} util.Response
// @Failure 403 {object} util.Response
// @Router /comments/create [post]
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

// @Summary Delete comment
// @Description Deletes an existing comment
// @Tags comments
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param comment body util.DeleteCommentRequestBody true "Delete comment body"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 401 {object} util.Response
// @Failure 403 {object} util.Response
// @Router /comments/delete [delete]
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

// @Summary Update comment
// @Description Updates an existing comment
// @Tags comments
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param comment body util.UpdateCommentRequestBody true "Update comment body"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 401 {object} util.Response
// @Failure 403 {object} util.Response
// @Router /comments/update [put]
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

// @Summary Get comment by ID
// @Description Returns a single comment by its ID
// @Tags comments
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Router /comments/{id} [get]
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

// @Summary Get comments by post
// @Description Returns comments on a specific post with pagination
// @Tags comments
// @Produce json
// @Param post_id query int true "Post ID"
// @Param limit query int true "Limit number of comments"
// @Param offset query int true "Offset for pagination"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Router /comments/post [get]
func (comment *Comment) GetCommentsByPost(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limit := query.Get("limit")
	offset := query.Get("offset")
	postId := query.Get("post_id")

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
