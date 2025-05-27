package handler

import (
	"log"
	"net/http"

	"github.com/ecofriends/authentication-backend/util"
)

func SignOut(w http.ResponseWriter, r *http.Request) {
	// Expire the token cookie
	util.ExpireCookie(w, "token")

	util.JsonResponse(w, "Successfully signed-out", http.StatusOK, nil)
	log.Println("[LOG]: Successfully signed user out")
}
