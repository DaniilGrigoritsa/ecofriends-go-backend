package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ecofriends/authentication-backend/authentication"
	"github.com/ecofriends/authentication-backend/service"
	"github.com/ecofriends/authentication-backend/util"
)

func SignIn(dbService *service.DatabaseProvider, w http.ResponseWriter, r *http.Request) {
	// Store the auth request body
	var body = util.SignInRequestBody{}

	// Read response body into body struct
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		msg := "Bad request, username, email or password not present"
		util.JsonResponse(w, msg, http.StatusBadRequest, nil)
		return
	}

	// Expire the token cookie
	util.ExpireCookie(w, "token")

	// Sanitize user input
	util.SanitizeUserInput(&body)

	// Check if the user exists
	userExists, err := dbService.Repo.UserExists(r.Context(), body.Email, "")
	if err != nil {
		msg := "Internal server error, could not check if user already exists"
		util.JsonResponse(w, msg, http.StatusInternalServerError, nil)
		return
	}

	// If the user does not exist, respond with an error
	if !userExists {
		msg := "A user with those credentials does not exist"
		util.JsonResponse(w, msg, http.StatusBadRequest, nil)
		return
	}

	// Query the database and obtain the user the the provided email
	user, err := dbService.Repo.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		msg := "Internal server error, could not check if a user with that email exists"
		util.JsonResponse(w, msg, http.StatusInternalServerError, nil)
		return
	}

	// Check that the password matches the hash
	if !util.CompareWithHash([]byte(user.Password), body.Password) {
		msg := "Provided passwords mismatch"
		util.JsonResponse(w, msg, http.StatusUnauthorized, nil)
		return
	}

	// Generate a new token and send response
	token, err := authentication.CreateJWToken(user.ID)
	if err != nil {
		log.Println(err)
		msg := "Failed to create token"
		util.JsonResponse(w, msg, http.StatusInternalServerError, nil)
		return
	}

	// Create the user payload
	var userPayload = util.UserPayload{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	// Set the token cookie and send the response
	cookie := util.CreateTokenCookie(token)
	http.SetCookie(w, &cookie)
	util.JsonResponse(w, "Successfully signed-in", http.StatusOK, userPayload)
}
