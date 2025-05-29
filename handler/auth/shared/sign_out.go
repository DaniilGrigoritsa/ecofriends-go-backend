package handler

import (
	"log"
	"net/http"

	"github.com/ecofriends/authentication-backend/util"
)

// SignOut handles user log out
// @Summary Log out the user
// @Description Log out an existing user
// @Tags authentication
// @Accept json
// @Produce json
// @Success 200 {object} util.Response
// @Failure 500 {object} util.Response
// @Router /auth/sign-out [post]
func SignOut(w http.ResponseWriter, r *http.Request) {
	// Expire the token cookie
	util.ExpireCookie(w, "token")

	util.JsonResponse(w, "Successfully signed-out", http.StatusOK, nil)
	log.Println("[LOG]: Successfully signed user out")
}
