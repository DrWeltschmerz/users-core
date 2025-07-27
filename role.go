package users

type Role struct {
	ID   string
	Name string
}

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)
