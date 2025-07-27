# users-core

This repository contains the **core business logic** for user and role management in Go. It is designed to be framework-agnostic and easily integrated with other modules (such as HTTP handlers, database adapters, etc.) that may live in separate repositories.

---

## Features

- User registration and authentication
- Password hashing and verification (pluggable)
- Role management (create, assign, list)
- User CRUD operations
- Last seen tracking
- Password change and reset
- Error handling with domain-specific errors

---

## Structure

- `service.go` — Core service logic for users and roles
- `repository.go` — Repository interfaces for users and roles (to be implemented elsewhere)
- `hasher.go` — Password hasher interface (to be implemented elsewhere)
- `user.go` — User domain model
- `role.go` — Role domain model and constants
- `DTO.go` — Data transfer objects for user input
- `service_test.go` — Unit tests with mocks for all service logic

---

## Usage

This module is intended to be imported by other modules (such as HTTP handlers or database adapters) that provide concrete implementations for the repository and hasher interfaces.

Example usage:

```go
import "github.com/DrWeltschmerz/users-core"

userRepo := NewYourUserRepo()     // your implementation
roleRepo := NewYourRoleRepo()     // your implementation
hasher := NewYourPasswordHasher() // your implementation

svc := users.NewService(userRepo, roleRepo, hasher)
```

---

## Interfaces

You must provide implementations for:

- `UserRepository` (see `repository.go`)
- `RoleRepository` (see `repository.go`)
- `PasswordHasher` (see `hasher.go`)

---

## Extending

- Add your own adapters for persistence (e.g., Postgres, MongoDB, etc.) in separate repositories.
- Add HTTP or gRPC handlers in separate repositories.

---

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
