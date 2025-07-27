package users

import (
	"time"
)

type User struct {
	ID             string
	Email          string
	HashedPassword string
	Username       string
	LastSeen       time.Time
	RoleID         string
}
