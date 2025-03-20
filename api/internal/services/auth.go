package services

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"just-kanban/internal/models"
	"just-kanban/pkg/auth/jwt"
	"just-kanban/pkg/sqlddl"
)

var (
	wrongCredentialsErr = errors.New("wrong credentials")
)

type (
	AuthService struct {
		*TokenService
		UserService
	}
	LoginData struct {
		Identifier string `json:"identifier" validate:"required"`
		Password   string `json:"password" validate:"required"`
	}
)

func NewAuthService(ts *TokenService, us UserService) *AuthService {
	return &AuthService{ts, us}
}

func (as *AuthService) RegisterUser(ctx context.Context, registrationData *CreateUserData) (*jwt.AccessTokens, error) {
	hashedPassword, hashingErr := bcrypt.GenerateFromPassword([]byte(registrationData.Password), 10)
	if hashingErr != nil {
		return nil, hashingErr
	}
	createdUser, creationErr := as.UserService.CreateUser(ctx, &CreateUserData{
		Email:     registrationData.Email,
		Username:  registrationData.Username,
		FirstName: registrationData.FirstName,
		LastName:  registrationData.LastName,
		Password:  string(hashedPassword),
	})
	if creationErr != nil {
		return nil, creationErr
	}
	tokens, tokensErr := as.TokenService.CreateAuthPair(ctx, createdUser)
	if tokensErr != nil {
		return nil, tokensErr
	}
	return tokens, nil
}

func (as *AuthService) Login(ctx context.Context, loginData *LoginData) (*jwt.AccessTokens, error) {
	var searchUser *models.User
	emailUser, searchEmailUserErr := as.UserService.FindByEmail(ctx, loginData.Identifier)
	if searchEmailUserErr == nil {
		searchUser = emailUser
	} else {
		usernameUser, searchUsernameErr := as.UserService.FindByUsername(ctx, loginData.Identifier)
		if searchUsernameErr == nil {
			searchUser = usernameUser
		}
	}
	if searchUser == nil {
		return nil, wrongCredentialsErr
	}
	if compareErr := bcrypt.CompareHashAndPassword([]byte(searchUser.Password), []byte(loginData.Password)); compareErr != nil {
		return nil, wrongCredentialsErr
	}
	tokens, tokensErr := as.CreateAuthPair(ctx, searchUser)
	if tokensErr != nil {
		return nil, tokensErr
	}
	return tokens, nil
}

func (as *AuthService) Logout(ctx context.Context, userID sqlddl.ID) error {
	removeErr := as.TokenService.RemoveRefreshByUserID(ctx, userID)
	return removeErr
}
