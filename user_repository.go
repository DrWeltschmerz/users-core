package users

import "context"

type UserRepository interface {
	Create(ctx context.Context, user User) (*User, error)
	Update(ctx context.Context, user User) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	List(ctx context.Context) ([]User, error)
	Delete(ctx context.Context, id string) error
}
