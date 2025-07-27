package users

import (
	"context"
	"time"
)

type Service struct {
	userRepo UserRepository
	roleRepo RoleRepository
	hasher   PasswordHasher
}

func NewService(userRepo UserRepository, roleRepo RoleRepository, hasher PasswordHasher) *Service {
	return &Service{
		userRepo: userRepo,
		roleRepo: roleRepo,
		hasher:   hasher,
	}
}

func (s *Service) Register(ctx context.Context, input UserRegisterInput) (*User, error) {
	hashedPassword, err := s.hasher.Hash(input.Password)
	if err != nil {
		return nil, err
	}

	role, err := s.roleRepo.GetByName(ctx, RoleUser)
	if err != nil {
		role, err = s.roleRepo.Create(ctx, Role{Name: RoleUser})
		if err != nil {
			return nil, ErrFailedToCreateRole
		}
	}

	user := User{
		Email:          input.Email,
		Username:       input.Username,
		HashedPassword: hashedPassword,
		LastSeen:       time.Now(),
		RoleID:         role.ID,
	}

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *Service) Login(ctx context.Context, input UserLoginInput) (*User, error) {
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, ErrUserNotFound
	}
	if !s.hasher.Verify(user.HashedPassword, input.Password) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *Service) GetUserByID(ctx context.Context, id string) (*User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *Service) UpdateUser(ctx context.Context, user User) (*User, error) {
	updatedUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, ErrFailedToUpdateUser
	}
	return updatedUser, nil
}

func (s *Service) ListUsers(ctx context.Context) ([]User, error) {
	users, err := s.userRepo.List(ctx)
	if err != nil {
		return nil, ErrFailedToListUsers
	}
	return users, nil
}

func (s *Service) DeleteUser(ctx context.Context, id string) error {
	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		return ErrFailedToDeleteUser
	}
	return nil
}

func (s *Service) GetRoleByID(ctx context.Context, id string) (*Role, error) {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

func (s *Service) CreateRole(ctx context.Context, role Role) (*Role, error) {
	createdRole, err := s.roleRepo.Create(ctx, role)
	if err != nil {
		return nil, ErrFailedToCreateRole
	}
	return createdRole, nil
}

func (s *Service) AssignRoleToUser(ctx context.Context, userID, roleID string) (*User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return nil, ErrRoleNotFound
	}

	user.RoleID = role.ID
	updatedUser, err := s.userRepo.Update(ctx, *user)
	if err != nil {
		return nil, ErrFailedToUpdateUser
	}

	return updatedUser, nil
}

func (s *Service) ListRoles(ctx context.Context) ([]Role, error) {
	roles, err := s.roleRepo.List(ctx)
	if err != nil {
		return nil, ErrFailedToListUsers
	}
	return roles, nil
}

func (s *Service) IsAdmin(user *User) bool {
	if user.RoleID == "" {
		return false
	}
	// IsAdmin nie ma contextu, więc nie zmieniamy sygnatury
	role, err := s.roleRepo.GetByID(context.Background(), user.RoleID)
	if err != nil {
		return false
	}
	return role.Name == RoleAdmin
}

func (s *Service) UpdateLastSeen(ctx context.Context, userID string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	user.LastSeen = time.Now()
	_, err = s.userRepo.Update(ctx, *user)
	if err != nil {
		return ErrFailedToUpdateUser
	}

	return nil
}

func (s *Service) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) (*User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if oldPassword == newPassword {
		return nil, ErrCannotUseSamePassword
	}

	if !s.hasher.Verify(user.HashedPassword, oldPassword) {
		return nil, ErrInvalidCredentials
	}

	hashedNewPassword, err := s.hasher.Hash(newPassword)
	if err != nil {
		return nil, ErrFailedToHashPassword
	}

	user.HashedPassword = hashedNewPassword
	updatedUser, err := s.userRepo.Update(ctx, *user)
	if err != nil {
		return nil, ErrFailedToUpdateUser
	}

	return updatedUser, nil
}

func (s *Service) ResetPassword(ctx context.Context, userID, newPassword string) (*User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	hashedPassword, err := s.hasher.Hash(newPassword)
	if err != nil {
		return nil, ErrFailedToHashPassword
	}

	user.HashedPassword = hashedPassword
	updatedUser, err := s.userRepo.Update(ctx, *user)
	if err != nil {
		return nil, ErrFailedToUpdateUser
	}

	return updatedUser, nil
}
