# users-core

Core domain logic and interfaces for user management.

## Features

- User and Role domain models
- Repository interfaces (`UserRepository`, `RoleRepository`)
- Service layer with business logic (registration, login, password change, etc.)
- Password hashing abstraction

## How to Use With Adapters


This module defines the core types, interfaces, and business logic for user management. It does **not** provide any database, HTTP, or JWT implementation.

To use this module in a real application, combine it with one or more adapters and a JWT implementation:

- [users-adapter-gorm](https://github.com/DrWeltschmerz/users-adapter-gorm): GORM-based repository implementations
- [users-adapter-gin](https://github.com/DrWeltschmerz/users-adapter-gin): Gin HTTP REST API adapter
- [jwt-auth](https://github.com/DrWeltschmerz/jwt-auth): JWT tokenizer and password hasher implementations

### Example: Wiring Everything Together

```go
import (
    "github.com/DrWeltschmerz/users-core"
    gormadapter "github.com/DrWeltschmerz/users-adapter-gorm/gorm"
    ginadapter "github.com/DrWeltschmerz/users-adapter-gin/ginadapter"
    "github.com/DrWeltschmerz/jwt-auth/pkg/authjwt"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "github.com/gin-gonic/gin"
)

func main() {
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    db.AutoMigrate(&gormadapter.GormUser{}, &gormadapter.GormRole{})

    userRepo := gormadapter.NewGormUserRepository(db)
    roleRepo := gormadapter.NewGormRoleRepository(db)

    // Use jwt-auth or your own implementation for hasher and tokenizer
    hasher := authjwt.NewBcryptHasher()
    tokenizer := authjwt.NewJWTTokenizer()

    service := users.NewService(userRepo, roleRepo, hasher, tokenizer)

    r := gin.Default()
    ginadapter.RegisterRoutes(r, service, tokenizer)
    r.Run()
}
```

See the [users-adapter-gorm](https://github.com/DrWeltschmerz/users-adapter-gorm) and [users-adapter-gin](https://github.com/DrWeltschmerz/users-adapter-gin) READMEs for details on their own APIs and extension points.

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
