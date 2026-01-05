package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("super-secret-key")

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// Extracts the Bearer token from the Authorizatoin header
// Expected format: "Authorization: Bearer <token>"
func GetTokenFromHeader(r *http.Request) string {
	// Read the Authorization header from the request
	// If the Authorization header is missing or does not match the expected format
	// return an empty string
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Split the header into two parts: "Bearer" and the token value
	// Validate the split format before returning the token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}

// Create a signed JWT token for the given user ID
func GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // the token expires after 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Parse and validate a JWT string
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			// Provide the secret key used to verify the token signature
			return jwtSecret, nil
		},
	)

	if err != nil || !token.Valid {
		return nil, err
	}

	// Extract claims and ensure the token is valid
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("Invalid token")
	}

	return claims, nil
}
