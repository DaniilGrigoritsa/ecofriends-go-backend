package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	repository "github.com/ecofriends/authentication-backend/repository"
	"github.com/ecofriends/authentication-backend/util"
	"github.com/go-chi/chi/v5"
)

type Post struct {
	repo *repository.PostGreSQL
}

func (post *Post) New(repo *repository.PostGreSQL) {
	post.repo = repo
}

// CreatePost creates a new post
// @Summary Create a new post
// @Description Allows an authenticated user to create a new post
// @Tags posts
// @Accept json
// @Produce json
// @Param request body util.CreatePostRequestBody true "Post creation payload"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 401 {object} util.Response
// @Failure 403 {object} util.Response
// @Security CookieAuth
// @Router /posts/create [post]
func (post *Post) CreatePost(w http.ResponseWriter, r *http.Request) {
	var body = util.CreatePostRequestBody{}

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

	thePost, err := post.repo.CreatePost(context.Background(), body.UserID.String(), body.Text)
	if err != nil {
		util.JsonResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	util.JsonResponse(w, "Successfully created post", http.StatusOK, thePost)
}

// DeletePost deletes a post
// @Summary Delete a post
// @Description Allows an authenticated user to delete their post
// @Tags posts
// @Accept json
// @Produce json
// @Param request body util.DeletePostRequestBody true "Post deletion payload"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 401 {object} util.Response
// @Failure 403 {object} util.Response
// @Security CookieAuth
// @Router /posts/delete [delete]
func (post *Post) DeletePost(w http.ResponseWriter, r *http.Request) {
	var body = util.DeletePostRequestBody{}

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

	if err := post.repo.DeletePost(context.Background(), body.PostID, body.UserID.String()); err != nil {
		util.JsonResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	util.JsonResponse(w, "Successfully deleted post", http.StatusOK, nil)
}

// GetPostByID fetches a single post by its ID
// @Summary Get post by ID
// @Description Returns a post based on its ID
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 500 {object} util.Response
// @Router /posts/{id} [get]
func (post *Post) GetPostByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var msg = ""

	idInt, err := strconv.Atoi(id)
	if err != nil {
		msg = err.Error()
		util.JsonResponse(w, msg, http.StatusBadRequest, nil)
		return
	}

	thePost, err := post.repo.GetPostByID(r.Context(), idInt)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			msg = "A post with that id doesn't exist"
			util.JsonResponse(w, msg, http.StatusBadRequest, nil)
			return
		}
		msg = "Internal server error, failed to get post with that id"
		util.JsonResponse(w, msg, http.StatusInternalServerError, nil)
		return
	}

	msg = fmt.Sprintf("Successfully fetched user with the id: %s", id)
	util.JsonResponse(w, msg, http.StatusOK, thePost)
}

// @Summary Get all posts
// @Description Returns all posts with pagination
// @Tags posts
// @Produce json
// @Param limit query int true "Limit number of posts"
// @Param offset query int true "Offset for pagination"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Router /posts/all [get]
func (post *Post) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

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

	posts, err := post.repo.GetAllPosts(context.Background(), limitInt, offsetInt)
	if err != nil {
		util.JsonResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	util.JsonResponse(w, "Successfully got all posts", http.StatusOK, posts)
}

// @Summary Get posts by user
// @Description Returns posts made by a specific user with pagination
// @Tags posts
// @Produce json
// @Param user_id query string true "User ID"
// @Param limit query int true "Limit number of posts"
// @Param offset query int true "Offset for pagination"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Router /posts/user [get]
func (post *Post) GetPostsByUser(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	userId := r.URL.Query().Get("user_id")

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

	posts, err := post.repo.GetPostsByUser(context.Background(), userId, limitInt, offsetInt)
	if err != nil {
		util.JsonResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	util.JsonResponse(w, "Successfully got posts by user", http.StatusOK, posts)
}
