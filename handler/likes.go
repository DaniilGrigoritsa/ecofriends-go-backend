package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	repository "github.com/ecofriends/authentication-backend/repository"
	"github.com/ecofriends/authentication-backend/util"
)

type Like struct {
	repo *repository.PostGreSQL
}

func (like *Like) New(repo *repository.PostGreSQL) {
	like.repo = repo
}

// @Summary Like a post
// @Description Like a post as an authenticated user
// @Tags likes
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param like body util.LikePostRequestBody true "Post to like"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 401 {object} util.Response
// @Failure 403 {object} util.Response
// @Router /likes/like [post]
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

// @Summary Unlike a post
// @Description Unlike a post as an authenticated user
// @Tags likes
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param unlike body util.LikePostRequestBody true "Post to unlike"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 401 {object} util.Response
// @Failure 403 {object} util.Response
// @Router /likes/unlike [post]
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

// @Summary Get like count for a post
// @Description Returns number of likes for the given post
// @Tags likes
// @Produce json
// @Param post_id query int true "Post ID"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Router /likes/count [get]
func (like *Like) GetLikeCount(w http.ResponseWriter, r *http.Request) {
	postId := r.URL.Query().Get("post_id")

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

// @Summary Check if user has liked a post
// @Description Returns true if the user has liked the specified post
// @Tags likes
// @Produce json
// @Param post_id query int true "Post ID"
// @Param user_id query string true "User ID"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Router /likes/has_liked [get]
func (like *Like) HasLiked(w http.ResponseWriter, r *http.Request) {
	postId := r.URL.Query().Get("post_id")
	userId := r.URL.Query().Get("user_id")

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

// @Summary Get liked posts by user
// @Description Returns a paginated list of posts liked by the user
// @Tags likes
// @Produce json
// @Param user_id query string true "User ID"
// @Param limit query int true "Limit"
// @Param offset query int true "Offset"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Router /likes/user_likes [get]
func (like *Like) GetLikesByUser(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limit := query.Get("limit")
	offset := query.Get("offset")
	userId := query.Get("user_id")

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
