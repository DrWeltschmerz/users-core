package users

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(hashedPassword, password string) bool
}

type Tokenizer interface {
	GenerateToken(email, userID string) (string, error)
	ValidateToken(token string) (string, error)
}
