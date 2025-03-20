package services

import (
	"github.com/google/uuid"

	"context"
	"errors"
	"time"

	"just-kanban/internal/models"
	"just-kanban/internal/repositories/interfaces"
	"just-kanban/pkg/auth/jwt"
	"just-kanban/pkg/sqlddl"
)

var (
	invalidTokenError = errors.New("invalid token")
	noTokenExistsErr  = errors.New("token does not exist")
)

type (
	TokenService struct {
		interfaces.RefreshTokenRepository
		JWTSecret string
	}
	AccessTokenClaims struct {
		Email     string `json:"email"`
		Avatar    string `json:"avatar"`
		Username  string `json:"username"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		sqlddl.ID `json:"id"`
		jwt.RegisteredClaims
	}
	RefreshTokenClaims struct {
		jwt.RegisteredClaims
	}
)

func NewTokenService(refreshTokenRepo interfaces.RefreshTokenRepository, JWTSecret string) *TokenService {
	return &TokenService{refreshTokenRepo, JWTSecret}
}

func (ts *TokenService) CreateAuthPair(ctx context.Context, user *models.User) (*jwt.AccessTokens, error) {
	accessTokenId := uuid.New().String()
	accessToken, accessTokenErr := jwt.CreateSignedToken(&AccessTokenClaims{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Avatar:    user.Avatar,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   string(user.ID),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        accessTokenId,
		},
	}, ts.JWTSecret)
	if accessTokenErr != nil {
		return nil, accessTokenErr
	}
	refreshToken, refreshTokenErr := jwt.CreateSignedToken(&RefreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   accessTokenId,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}, ts.JWTSecret)
	if refreshTokenErr != nil {
		return nil, accessTokenErr
	}
	if _, searchErr := ts.FindRefreshByUserID(ctx, user.ID); searchErr == nil {
		ts.RemoveRefreshByUserID(ctx, user.ID)
	}
	refreshSaveErr := ts.RefreshTokenRepository.Create(ctx, &models.RefreshToken{
		Model:  models.Model{ID: sqlddl.ID(uuid.NewString())},
		UserID: user.ID,
		Token:  refreshToken,
	})
	if refreshSaveErr != nil {
		return nil, refreshSaveErr
	}
	return &jwt.AccessTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (ts *TokenService) FindUserIDByRefresh(ctx context.Context, token string) (sqlddl.ID, error) {
	userId, searchErr := ts.RefreshTokenRepository.FindUserIDByToken(ctx, token)
	return userId, searchErr
}

func (ts *TokenService) RemoveRefreshByUserID(ctx context.Context, userID sqlddl.ID) error {
	// todo: include redis storage for tokens invalidation
	_, searchErr := ts.RefreshTokenRepository.FindByUserID(ctx, userID)
	if searchErr != nil {
		return noTokenExistsErr
	}
	deleteErr := ts.RefreshTokenRepository.DeleteByUserID(ctx, userID)
	return deleteErr
}

func (ts *TokenService) FindRefreshByUserID(ctx context.Context, userID sqlddl.ID) (*models.RefreshToken, error) {
	token, searchErr := ts.RefreshTokenRepository.FindByUserID(ctx, userID)
	if searchErr != nil {
		return nil, noTokenExistsErr
	}
	return token, nil
}

func (ts *TokenService) ParseAccessToken(accessToken string) (*AccessTokenClaims, error) {
	var claims AccessTokenClaims
	_, parseErr := jwt.ParseWithClaims(&claims, accessToken, ts.JWTSecret)
	if parseErr != nil {
		return nil, invalidTokenError
	}
	return &claims, nil
}

func (ts *TokenService) CheckRefreshToken(ctx context.Context, refreshToken string) error {
	var claims RefreshTokenClaims
	_, parseErr := jwt.ParseWithClaims(&claims, refreshToken, ts.JWTSecret)
	if parseErr != nil {
		return invalidTokenError
	}
	findToken, searchErr := ts.RefreshTokenRepository.FindByToken(ctx, refreshToken)
	if searchErr != nil {
		return invalidTokenError
	}
	if sqlddl.ID(claims.Subject) != findToken.UserID {
		return invalidTokenError
	}
	return nil
}
