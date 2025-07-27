package users

type UserRepository interface {
	Create(user User) (*User, error)
	Update(user User) (*User, error)
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	List() ([]User, error)
	Delete(id string) error
}

type RoleRepository interface {
	Create(role Role) (*Role, error)
	Update(role Role) (*Role, error)
	Delete(id string) error
	GetByID(id string) (*Role, error)
	GetByName(name string) (*Role, error)
	List() ([]Role, error)
}
