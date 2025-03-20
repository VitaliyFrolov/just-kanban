package services

import (
	"context"
	"errors"

	"just-kanban/internal/contextkeys"
	"just-kanban/internal/models"
	"just-kanban/internal/repositories/interfaces"
	"just-kanban/pkg/identifier"
	"just-kanban/pkg/sqlddl"
)

type (
	UserService interface {
		interfaces.UserRepository
		CreateUser(ctx context.Context, d *CreateUserData) (*models.User, error)
		UpdateUser(ctx context.Context, id sqlddl.ID, d *UpdateUserData) (*models.User, error)
		IsUpdateAllowed(ctx context.Context, userId, targetId sqlddl.ID) bool
		ListUsers(ctx context.Context) ([]models.User, error)
		DeleteUser(ctx context.Context, id sqlddl.ID) error
	}

	userService struct {
		interfaces.UserRepository
	}

	CreateUserData struct {
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required,min=6,max=70,trimmed"`
		Username  string `json:"username" validate:"required,min=4,max=30,trimmed"`
		FirstName string `json:"first_name" validate:"required,min=4,max=50,trimmed"`
		LastName  string `json:"last_name" validate:"required,min=4,max=50,trimmed"`
	}

	UpdateUserData struct {
		FirstName string `json:"first_name" validate:"omitempty,min=4,max=50,trimmed"`
		LastName  string `json:"last_name" validate:"omitempty,min=5,max=50,trimmed"`
		Avatar    string `json:"avatar" validate:"omitempty,url"`
	}
)

var (
	ErrorUserEmailTaken = errors.New("this e-mail address is taken")
	ErrorUsernameTaken  = errors.New("username is unavailable, try another")
	userNotExistsErr    = errors.New("user not exists")
	updateNotAllowedErr = errors.New("not allowed")
)

func NewUserService(userRepository interfaces.UserRepository) UserService {
	return &userService{userRepository}
}

func (us *userService) CreateUser(ctx context.Context, d *CreateUserData) (*models.User, error) {
	id := sqlddl.ID(identifier.GenerateUUID())
	_, searchEmailErr := us.UserRepository.FindByEmail(ctx, d.Email)
	if searchEmailErr == nil {
		return nil, ErrorUserEmailTaken
	}
	_, searchUsernameErr := us.UserRepository.FindByUsername(ctx, d.Username)
	if searchUsernameErr == nil {
		return nil, ErrorUsernameTaken
	}
	creationErr := us.UserRepository.Create(ctx, &models.User{
		Model:     models.Model{ID: id},
		Email:     d.Email,
		Password:  d.Password,
		Username:  d.Username,
		FirstName: d.FirstName,
		LastName:  d.LastName,
	})
	if creationErr != nil {
		return nil, creationErr
	}
	newUser, searchEmailErr := us.UserRepository.FindByID(ctx, id)
	return newUser, nil
}

func (us *userService) UpdateUser(ctx context.Context, id sqlddl.ID, d *UpdateUserData) (*models.User, error) {
	_, searchErr := us.UserRepository.FindByID(ctx, id)
	if searchErr != nil {
		return nil, userNotExistsErr
	}
	userId, requesterOk := ctx.Value(contextkeys.KeyUserId).(sqlddl.ID)
	if !requesterOk || !us.IsUpdateAllowed(ctx, userId, id) {
		return nil, updateNotAllowedErr
	}
	updateErr := us.UserRepository.Update(ctx, id, &models.UpdateUser{
		FirstName: &d.FirstName,
		LastName:  &d.LastName,
		Avatar:    &d.Avatar,
	})
	if updateErr != nil {
		return nil, updateErr
	}
	updatedUser, searchErr := us.UserRepository.FindByID(ctx, id)
	return updatedUser, searchErr
}

func (us *userService) IsUpdateAllowed(ctx context.Context, userId, targetId sqlddl.ID) bool {
	return userId == targetId
}

func (us *userService) ListUsers(ctx context.Context) ([]models.User, error) {
	users, searchErr := us.UserRepository.FindAll(ctx)
	if searchErr != nil {
		return nil, searchErr
	}
	return users, nil
}

func (us *userService) DeleteUser(ctx context.Context, id sqlddl.ID) error {
	_, searchErr := us.UserRepository.FindByID(ctx, id)
	if searchErr != nil {
		return userNotExistsErr
	}
	delErr := us.UserRepository.Delete(ctx, id)
	if delErr != nil {
		return delErr
	}
	return nil
}
