# users-core

Core domain logic and interfaces for user management.

## Features

- User and Role domain models
- Repository interfaces (`UserRepository`, `RoleRepository`)
- Service layer with business logic (registration, login, password change, etc.)
- Password hashing abstraction

## Usage

This module defines the core types and interfaces for user management.  
It does **not** provide any database implementation—see [users-adapter-gorm](https://github.com/DrWeltschmerz/users-adapter-gorm) for a GORM adapter.

### Example: Using the Service with a Repository

```go
import (
    "context"
    "github.com/DrWeltschmerz/users-core"
    gormadapter "github.com/DrWeltschmerz/users-adapter-gorm/gorm"
    "github.com/DrWeltschmerz/jwt-auth/pkg/authjwt"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func main() {
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    db.AutoMigrate(&gormadapter.GormUser{}, &gormadapter.GormRole{})

    userRepo := gormadapter.NewGormUserRepository(db)
    roleRepo := gormadapter.NewGormRoleRepository(db)
    hasher := authjwt.NewBcryptHasher()

    service := users.NewService(userRepo, roleRepo, hasher)

    user, err := service.Register(context.Background(), users.UserRegisterInput{
        Email:    "test@example.com",
        Username: "testuser",
        Password: "secret",
    })
    // ...
}
```

## Repository Interfaces

The repository interfaces (`UserRepository`, `RoleRepository`) are defined in the main package files and specify the required methods for data access and persistence.  
You can implement these interfaces to connect the service layer to any storage backend.

## Testing

Unit tests use mocks for all dependencies.  
Tests use [Testify](https://github.com/stretchr/testify) for assertions.

Run tests with:

```sh
go test ./...
```

---

## Requirements

- Go 1.24.5 or newer

Dependencies (see [`go.mod`](go.mod)):

- [github.com/stretchr/testify](https://github.com/stretchr/testify) (for testing)

---

## License

This project is licensed under the [GNU General Public License v3.0 (GPL-3.0)](LICENSE).

See the LICENSE file for details.
