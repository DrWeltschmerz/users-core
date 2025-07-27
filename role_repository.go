package users

import "context"

type RoleRepository interface {
	Create(ctx context.Context, role Role) (*Role, error)
	Update(ctx context.Context, role Role) (*Role, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Role, error)
	GetByName(ctx context.Context, name string) (*Role, error)
	List(ctx context.Context) ([]Role, error)
}
