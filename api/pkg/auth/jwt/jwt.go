package jwt

import (
	gjwt "github.com/golang-jwt/jwt/v5"
)

type (
	AccessTokens struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	RegisteredClaims = gjwt.RegisteredClaims
	Token            = gjwt.Token
)

const (
	AuthTypeBearer  = "Bearer"
	RefreshTokenKey = "refresh_token"
)

var (
	NewNumericDate = gjwt.NewNumericDate
)

func CreateSignedToken(claims gjwt.Claims, secret string) (string, error) {
	token := gjwt.NewWithClaims(gjwt.SigningMethodHS256, claims)
	signedToken, signingErr := token.SignedString([]byte(secret))
	if signingErr != nil {
		return "", signingErr
	}
	return signedToken, nil
}

func ParseWithClaims[TClaims gjwt.Claims](claims TClaims, tokenString string, secret string) (*Token, error) {
	token, err := gjwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *gjwt.Token) (interface{}, error) { return []byte(secret), nil },
	)
	if err != nil {
		return nil, err
	}
	return token, nil
}
