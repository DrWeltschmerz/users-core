package users

type UserRegisterInput struct {
	Email    string
	Username string
	Password string
}

type UserLoginInput struct {
	Email    string
	Password string
}
