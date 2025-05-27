package handler

import (
	"fmt"
	"net/http"

	oauth "github.com/ecofriends/authentication-backend/handler/auth/oauth"
	password "github.com/ecofriends/authentication-backend/handler/auth/password"
	shared "github.com/ecofriends/authentication-backend/handler/auth/shared"
	"github.com/ecofriends/authentication-backend/service"
	"github.com/ecofriends/authentication-backend/util"
	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	dbService *service.DatabaseProvider
}

func (authHandler *AuthHandler) WithService(service *service.DatabaseProvider) {
	authHandler.dbService = service
}

func (auth *AuthHandler) Home(w http.ResponseWriter, r *http.Request) {
	msg := "Auth route home"
	util.JsonResponse(w, msg, http.StatusOK, nil)
}

func (auth *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	password.SignUp(auth.dbService, w, r)
}

func (auth *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	password.SignIn(auth.dbService, w, r)
}

func (auth *AuthHandler) GoogleSignIn(w http.ResponseWriter, r *http.Request) {
	oauth.GoogleSignIn(w, r)
}

func (auth *AuthHandler) GoogleSignInCallback(w http.ResponseWriter, r *http.Request) {
	oauth.GoogleSignInCallback(auth.dbService, w, r)
}

func (auth *AuthHandler) OAuthFailure(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "x")
	msg := fmt.Sprintf("Failed to sign-in using %s OAuth", util.CapitalizeFirstLetter(provider))
	util.JsonResponse(w, msg, http.StatusUnauthorized, nil)
}

func (auth *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	shared.SignOut(w, r)
}
