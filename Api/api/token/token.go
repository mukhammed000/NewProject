package token

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/cast"
)

type JWTHandler struct {
	Sub        string
	Exp        string
	Iat        string
	Role       string
	SigningKey string
	Token      string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

var tokenKey = "my_secret_key"

func GenerateJWTToken(userEmail, userRole string) *Tokens {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessClaims := accessToken.Claims.(jwt.MapClaims)
	accessClaims["email"] = userEmail
	accessClaims["role"] = userRole
	accessClaims["iat"] = time.Now().Unix()
	accessClaims["exp"] = time.Now().Add(48 * time.Hour).Unix()
	access, err := accessToken.SignedString([]byte(tokenKey))
	if err != nil {
		log.Fatal("Error while generating access token: ", err)
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["email"] = userEmail
	refreshClaims["role"] = userRole
	refreshClaims["iat"] = time.Now().Unix()
	refreshClaims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	refresh, err := refreshToken.SignedString([]byte(tokenKey))
	if err != nil {
		log.Fatal("Error while generating refresh token: ", err)
	}

	return &Tokens{
		AccessToken:  access,
		RefreshToken: refresh,
	}
}

func ExtractClaim(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, fmt.Errorf("invalid token or claims")
	}

	return claims, nil
}

func (jwtHandler *JWTHandler) ExtractClaims() (jwt.MapClaims, error) {
	token, err := jwt.Parse(jwtHandler.Token, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtHandler.SigningKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, fmt.Errorf("invalid token or claims")
	}

	return claims, nil
}

func GetIdFromToken(tokenStr string) (string, error) {
	if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
		return "", fmt.Errorf("missing or malformed JWT")
	}

	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	claims, err := ExtractClaim(tokenStr)
	if err != nil {
		return "", err
	}

	return cast.ToString(claims["email"]), nil
}
