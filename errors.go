package users

import "errors"

var (
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrUserNotFound          = errors.New("user not found")
	ErrEmailTaken            = errors.New("email already taken")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrFailedToCreateRole    = errors.New("failed to create role")
	ErrFailedToUpdateUser    = errors.New("failed to update user")
	ErrFailedToDeleteUser    = errors.New("failed to delete user")
	ErrFailedToListUsers     = errors.New("failed to list users")
	ErrRoleNotFound          = errors.New("role not found")
	ErrFailedToHashPassword  = errors.New("failed to hash password")
	ErrCannotUseSamePassword = errors.New("cannot use the same password")
)
