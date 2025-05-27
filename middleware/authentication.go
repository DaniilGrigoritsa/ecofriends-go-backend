package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/ecofriends/authentication-backend/authentication"
	"github.com/ecofriends/authentication-backend/util"
	"github.com/golang-jwt/jwt/v5"
)

func AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("[LOG]: authentication requested on:", r.URL)

		tokenCookie, err := r.Cookie("token")
		if err != nil {
			log.Println("[FAIL]: token not present in cookie")
			msg := "Unauthorized request to a protected endpoint"
			util.JsonResponse(w, msg, http.StatusUnauthorized, nil)
			return
		}

		token, err := authentication.VerifyToken(tokenCookie.Value)
		if err != nil {
			log.Printf("[FAIL]: token verification failed: %v", err)
			msg := "Failed to verify token"
			util.JsonResponse(w, msg, http.StatusUnauthorized, nil)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("[FAIL]: could not parse token claims")
			msg := "Invalid token claims"
			util.JsonResponse(w, msg, http.StatusUnauthorized, nil)
			return
		}

		log.Printf("[SUCCESS]: token successfully verified: %v", claims)

		ctx := context.WithValue(r.Context(), util.TokenClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
