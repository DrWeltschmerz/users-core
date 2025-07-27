package users

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// --- Mock Implementations ---

type mockUserRepo struct {
	users                                            map[string]*User
	createErr, getErr, updateErr, deleteErr, listErr error
}

func (m *mockUserRepo) Create(ctx context.Context, u User) (*User, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	m.users[u.ID] = &u
	return &u, nil
}
func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, ErrUserNotFound
}
func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*User, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	u, ok := m.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return u, nil
}
func (m *mockUserRepo) Update(ctx context.Context, u User) (*User, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	m.users[u.ID] = &u
	return &u, nil
}
func (m *mockUserRepo) List(ctx context.Context) ([]User, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	var us []User
	for _, u := range m.users {
		us = append(us, *u)
	}
	return us, nil
}
func (m *mockUserRepo) Delete(ctx context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.users, id)
	return nil
}
func (m *mockUserRepo) GetByUsername(ctx context.Context, username string) (*User, error) {
	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, ErrUserNotFound
}

type mockRoleRepo struct {
	roles                                    map[string]*Role
	getByNameErr, createErr, getErr, listErr error
}

func (m *mockRoleRepo) Update(ctx context.Context, role Role) (*Role, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	m.roles[role.ID] = &role
	return &role, nil
}
func (m *mockRoleRepo) GetByName(ctx context.Context, name string) (*Role, error) {
	if m.getByNameErr != nil {
		return nil, m.getByNameErr
	}
	for _, r := range m.roles {
		if r.Name == name {
			return r, nil
		}
	}
	return nil, ErrRoleNotFound
}
func (m *mockRoleRepo) Create(ctx context.Context, r Role) (*Role, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	m.roles[r.ID] = &r
	return &r, nil
}
func (m *mockRoleRepo) GetByID(ctx context.Context, id string) (*Role, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	r, ok := m.roles[id]
	if !ok {
		return nil, ErrRoleNotFound
	}
	return r, nil
}
func (m *mockRoleRepo) List(ctx context.Context) ([]Role, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	var rs []Role
	for _, r := range m.roles {
		rs = append(rs, *r)
	}
	return rs, nil
}
func (m *mockRoleRepo) Delete(ctx context.Context, id string) error {
	delete(m.roles, id)
	return nil
}

type mockHasher struct {
	hashErr, verifyErr error
}

func (m *mockHasher) Hash(pw string) (string, error) {
	if m.hashErr != nil {
		return "", m.hashErr
	}
	return "hashed:" + pw, nil
}
func (m *mockHasher) Verify(hashed, pw string) bool {
	if m.verifyErr != nil {
		return false
	}
	return hashed == "hashed:"+pw
}

type mockTokenizer struct {
	generateToken string
	generateErr   error
	verifyUserID  string
	verifyErr     error
}

func (m *mockTokenizer) GenerateToken(email, userID string) (string, error) {
	if m.generateErr != nil {
		return "", m.generateErr
	}
	return "token", nil
}

func (m *mockTokenizer) ValidateToken(token string) (string, error) {
	if m.verifyErr != nil {
		return "", m.verifyErr
	}
	return m.verifyUserID, nil
}

// --- Test Data ---

var (
	testUser = &User{
		ID:             "u1",
		Email:          "test@example.com",
		Username:       "testuser",
		HashedPassword: "hashed:password",
		RoleID:         "r1",
		LastSeen:       time.Now(),
	}
	testRole = &Role{
		ID:   "r2",
		Name: "admin",
	}
)

// --- Tests ---

func TestAssignRoleToUser(t *testing.T) {
	ctx := context.Background()
	userRepo := &mockUserRepo{users: map[string]*User{"u1": testUser}}
	roleRepo := &mockRoleRepo{roles: map[string]*Role{"r2": testRole}}
	tokenizer := &mockTokenizer{}
	svc := NewService(userRepo, roleRepo, &mockHasher{}, tokenizer)

	t.Run("success", func(t *testing.T) {
		u, err := svc.AssignRoleToUser(ctx, "u1", "r2")
		require.NoError(t, err)
		require.Equal(t, "r2", u.RoleID)
	})

	t.Run("user not found", func(t *testing.T) {
		_, err := svc.AssignRoleToUser(ctx, "notfound", "r2")
		require.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("role not found", func(t *testing.T) {
		_, err := svc.AssignRoleToUser(ctx, "u1", "notfound")
		require.ErrorIs(t, err, ErrRoleNotFound)
	})

	t.Run("update fails", func(t *testing.T) {
		userRepo.updateErr = errors.New("fail")
		_, err := svc.AssignRoleToUser(ctx, "u1", "r2")
		require.ErrorIs(t, err, ErrFailedToUpdateUser)
		userRepo.updateErr = nil
	})
}

func TestRegister(t *testing.T) {
	ctx := context.Background()
	userRepo := &mockUserRepo{users: map[string]*User{}}
	roleRepo := &mockRoleRepo{roles: map[string]*Role{"user": {ID: "r1", Name: RoleUser}}}
	tokenizer := &mockTokenizer{}
	svc := NewService(userRepo, roleRepo, &mockHasher{}, tokenizer)

	t.Run("success", func(t *testing.T) {
		input := UserRegisterInput{Email: "a@b.com", Username: "a", Password: "pw"}
		u, err := svc.Register(ctx, input)
		require.NoError(t, err)
		require.Equal(t, "a@b.com", u.Email)
	})

	t.Run("hash error", func(t *testing.T) {
		svc := NewService(userRepo, roleRepo, &mockHasher{hashErr: errors.New("fail")}, tokenizer)
		_, err := svc.Register(ctx, UserRegisterInput{Email: "x", Username: "x", Password: "x"})
		require.Error(t, err)
	})

	t.Run("role create fallback", func(t *testing.T) {
		roleRepo := &mockRoleRepo{
			roles:        map[string]*Role{},
			getByNameErr: errors.New("not found"),
			createErr:    errors.New("fail"),
		}
		svc := NewService(userRepo, roleRepo, &mockHasher{}, tokenizer)
		input := UserRegisterInput{Email: "b@b.com", Username: "b", Password: "pw"}
		_, err := svc.Register(ctx, input)
		require.ErrorIs(t, err, ErrFailedToCreateRole)
	})
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	userRepo := &mockUserRepo{users: map[string]*User{"test@example.com": testUser}}
	tokenizer := &mockTokenizer{}
	svc := NewService(userRepo, &mockRoleRepo{}, &mockHasher{}, tokenizer)

	t.Run("success", func(t *testing.T) {
		token, err := svc.Login(ctx, UserLoginInput{Email: "test@example.com", Password: "password"})
		require.NoError(t, err)
		require.NotEqual(t, "", token)
	})

	t.Run("user not found", func(t *testing.T) {
		_, err := svc.Login(ctx, UserLoginInput{Email: "notfound", Password: "pw"})
		require.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		_, err := svc.Login(ctx, UserLoginInput{Email: "test@example.com", Password: "wrong"})
		require.ErrorIs(t, err, ErrInvalidCredentials)
	})
}

func TestGetUserByID(t *testing.T) {
	ctx := context.Background()
	userRepo := &mockUserRepo{users: map[string]*User{"u1": testUser}}
	tokenizer := &mockTokenizer{}
	svc := NewService(userRepo, &mockRoleRepo{}, &mockHasher{}, tokenizer)

	t.Run("success", func(t *testing.T) {
		u, err := svc.GetUserByID(ctx, "u1")
		require.NoError(t, err)
		require.Equal(t, "u1", u.ID)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := svc.GetUserByID(ctx, "notfound")
		require.ErrorIs(t, err, ErrUserNotFound)
	})
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()
	userRepo := &mockUserRepo{users: map[string]*User{"u1": testUser}}
	tokenizer := &mockTokenizer{}
	svc := NewService(userRepo, &mockRoleRepo{}, &mockHasher{}, tokenizer)

	t.Run("success", func(t *testing.T) {
		u, err := svc.UpdateUser(ctx, *testUser)
		require.NoError(t, err)
		require.Equal(t, testUser.ID, u.ID)
	})

	t.Run("fail", func(t *testing.T) {
		userRepo.updateErr = errors.New("fail")
		_, err := svc.UpdateUser(ctx, *testUser)
		require.ErrorIs(t, err, ErrFailedToUpdateUser)
		userRepo.updateErr = nil
	})
}

func TestListUsers(t *testing.T) {
	ctx := context.Background()
	userRepo := &mockUserRepo{users: map[string]*User{"u1": testUser}}
	tokenizer := &mockTokenizer{}
	svc := NewService(userRepo, &mockRoleRepo{}, &mockHasher{}, tokenizer)

	t.Run("success", func(t *testing.T) {
		us, err := svc.ListUsers(ctx)
		require.NoError(t, err)
		require.Len(t, us, 1)
	})

	t.Run("fail", func(t *testing.T) {
		userRepo.listErr = errors.New("fail")
		_, err := svc.ListUsers(ctx)
		require.ErrorIs(t, err, ErrFailedToListUsers)
		userRepo.listErr = nil
	})
}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	userRepo := &mockUserRepo{users: map[string]*User{"u1": testUser}}
	tokenizer := &mockTokenizer{}
	svc := NewService(userRepo, &mockRoleRepo{}, &mockHasher{}, tokenizer)

	t.Run("success", func(t *testing.T) {
		err := svc.DeleteUser(ctx, "u1")
		require.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		userRepo.deleteErr = errors.New("fail")
		err := svc.DeleteUser(ctx, "u1")
		require.ErrorIs(t, err, ErrFailedToDeleteUser)
		userRepo.deleteErr = nil
	})
}

func TestGetRoleByID(t *testing.T) {
	ctx := context.Background()
	roleRepo := &mockRoleRepo{roles: map[string]*Role{"r2": testRole}}
	tokenizer := &mockTokenizer{}
	svc := NewService(&mockUserRepo{}, roleRepo, &mockHasher{}, tokenizer)

	t.Run("success", func(t *testing.T) {
		r, err := svc.GetRoleByID(ctx, "r2")
		require.NoError(t, err)
		require.Equal(t, "r2", r.ID)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := svc.GetRoleByID(ctx, "notfound")
		require.ErrorIs(t, err, ErrRoleNotFound)
	})
}

func TestCreateRole(t *testing.T) {
	ctx := context.Background()
	roleRepo := &mockRoleRepo{roles: map[string]*Role{}}
	tokenizer := &mockTokenizer{}
	svc := NewService(&mockUserRepo{}, roleRepo, &mockHasher{}, tokenizer)

	t.Run("success", func(t *testing.T) {
		r, err := svc.CreateRole(ctx, Role{ID: "r3", Name: "user"})
		require.NoError(t, err)
		require.Equal(t, "user", r.Name)
	})

	t.Run("fail", func(t *testing.T) {
		roleRepo.createErr = errors.New("fail")
		_, err := svc.CreateRole(ctx, Role{ID: "r4", Name: "fail"})
		require.ErrorIs(t, err, ErrFailedToCreateRole)
		roleRepo.createErr = nil
	})
}

func TestListRoles(t *testing.T) {
	ctx := context.Background()
	roleRepo := &mockRoleRepo{roles: map[string]*Role{"r2": testRole}}
	tokenizer := &mockTokenizer{}
	svc := NewService(&mockUserRepo{}, roleRepo, &mockHasher{}, tokenizer)

	t.Run("success", func(t *testing.T) {
		rs, err := svc.ListRoles(ctx)
		require.NoError(t, err)
		require.Len(t, rs, 1)
	})

	t.Run("fail", func(t *testing.T) {
		roleRepo.listErr = errors.New("fail")
		_, err := svc.ListRoles(ctx)
		require.ErrorIs(t, err, ErrFailedToListUsers)
		roleRepo.listErr = nil
	})
}

func TestIsAdmin(t *testing.T) {
	roleRepo := &mockRoleRepo{roles: map[string]*Role{"r2": {ID: "r2", Name: RoleAdmin}}}
	tokenizer := &mockTokenizer{}
	svc := NewService(&mockUserRepo{}, roleRepo, &mockHasher{}, tokenizer)

	t.Run("is admin", func(t *testing.T) {
		u := &User{RoleID: "r2"}
		require.True(t, svc.IsAdmin(u))
	})

	t.Run("not admin", func(t *testing.T) {
		u := &User{RoleID: "notadmin"}
		require.False(t, svc.IsAdmin(u))
	})

	t.Run("no role id", func(t *testing.T) {
		u := &User{}
		require.False(t, svc.IsAdmin(u))
	})
}

func TestUpdateLastSeen(t *testing.T) {
	ctx := context.Background()
	userRepo := &mockUserRepo{users: map[string]*User{"u1": testUser}}
	tokenizer := &mockTokenizer{}
	svc := NewService(userRepo, &mockRoleRepo{}, &mockHasher{}, tokenizer)

	t.Run("success", func(t *testing.T) {
		err := svc.UpdateLastSeen(ctx, "u1")
		require.NoError(t, err)
	})

	t.Run("user not found", func(t *testing.T) {
		err := svc.UpdateLastSeen(ctx, "notfound")
		require.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("update fails", func(t *testing.T) {
		userRepo.updateErr = errors.New("fail")
		err := svc.UpdateLastSeen(ctx, "u1")
		require.ErrorIs(t, err, ErrFailedToUpdateUser)
		userRepo.updateErr = nil
	})
}

func TestChangePassword(t *testing.T) {
	ctx := context.Background()
	userRepo := &mockUserRepo{users: map[string]*User{"u1": testUser}}
	tokenizer := &mockTokenizer{}
	svc := NewService(userRepo, &mockRoleRepo{}, &mockHasher{}, tokenizer)

	t.Run("success", func(t *testing.T) {
		u, err := svc.ChangePassword(ctx, "u1", "password", "newpw")
		require.NoError(t, err)
		require.Equal(t, "hashed:newpw", u.HashedPassword)
	})

	t.Run("user not found", func(t *testing.T) {
		_, err := svc.ChangePassword(ctx, "notfound", "x", "y")
		require.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("same password", func(t *testing.T) {
		_, err := svc.ChangePassword(ctx, "u1", "password", "password")
		require.ErrorIs(t, err, ErrCannotUseSamePassword)
	})

	t.Run("incorrect old password", func(t *testing.T) {
		_, err := svc.ChangePassword(ctx, "u1", "wrong", "newpw")
		require.ErrorIs(t, err, ErrInvalidCredentials)
	})

	t.Run("hash error", func(t *testing.T) {
		svc := NewService(userRepo, &mockRoleRepo{}, &mockHasher{hashErr: errors.New("fail")}, tokenizer)
		_, err := svc.ChangePassword(ctx, "u1", "password", "newpw")
		require.Error(t, err)
	})

	t.Run("update fails", func(t *testing.T) {
		userRepo.updateErr = errors.New("fail")
		_, err := svc.ChangePassword(ctx, "u1", "newpw", "newpw2")
		require.Contains(t, userRepo.users, "u1")
		t.Logf("error returned: %v", err)
		require.ErrorIs(t, err, ErrFailedToUpdateUser)
		userRepo.updateErr = nil
	})

}
