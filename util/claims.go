package util

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const TokenClaimsKey = ContextKey("tokenClaims")

func ExtractUserIDFromClaims(ctx context.Context) (string, error) {
	// Extract token claims from context
	claims, ok := ctx.Value(TokenClaimsKey).(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("Unauthorized")
	}

	fmt.Println("CLAIMS: ", claims)

	// Extract user ID from token
	userID, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("%v", "Forbidden: Error extracting the user id")
	}

	return userID, nil
}
