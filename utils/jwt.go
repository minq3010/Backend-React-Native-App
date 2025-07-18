package utils

import "github.com/golang-jwt/jwt"

func GenerateJWT(claims jwt.Claims, method jwt.SigningMethod, jwtSecret string) (string, error) {
	return jwt.NewWithClaims(method, claims).SignedString([]byte(jwtSecret))
}