package handler

import (
	"fmt"
	"net/http"
	"strings"

	repository "github.com/ecofriends/authentication-backend/repository"
	"github.com/ecofriends/authentication-backend/util"
	"github.com/go-chi/chi/v5"
)

type User struct {
	repo *repository.PostGreSQL
}

func (user *User) New(repo *repository.PostGreSQL) {
	user.repo = repo
}

func (user *User) Home(w http.ResponseWriter, r *http.Request) {
	msg := "User route home"
	util.JsonResponse(w, msg, http.StatusOK, nil)
}

func (user *User) GetUserByID(w http.ResponseWriter, r *http.Request) {
	requestedID := chi.URLParam(r, "id")
	var msg = ""

	userID, err := util.ExtractUserIDFromClaims(r.Context())
	if err != nil {
		util.JsonResponse(w, err.Error(), http.StatusUnauthorized, nil)
		return
	}

	fmt.Println("USER ID:", userID)

	if userID != requestedID {
		msg = "Forbidden: Access to this resource is denied"
		util.JsonResponse(w, msg, http.StatusForbidden, nil)
		return
	}

	theUser, err := user.repo.GetUserByID(r.Context(), requestedID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			msg = "A user with that id doesn't exist"
			util.JsonResponse(w, msg, http.StatusBadRequest, nil)
			return
		}
		msg = "Internal server error, failed to get user with that id"
		util.JsonResponse(w, msg, http.StatusInternalServerError, nil)
		return
	}

	msg = fmt.Sprintf("Successfully fetched user with the id: %s", requestedID)
	util.JsonResponse(w, msg, http.StatusOK, theUser)
}
