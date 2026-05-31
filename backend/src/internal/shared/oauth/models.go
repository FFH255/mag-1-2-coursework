package oauth

import "github.com/markbates/goth"

type User struct {
	ID       string
	Email    string
	Provider string
}

type Provider struct {
	ClientKey   string
	Secret      string
	CallbackURL string
}

func newUser(user goth.User) *User {
	return &User{
		ID:       user.UserID,
		Email:    user.Email,
		Provider: user.Provider,
	}
}
