package oauth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

type Config struct {
	Google    Provider
	Github    Provider
	CookieKey string
}

type OAuth struct {
}

func New(config Config) *OAuth {
	store := sessions.NewCookieStore([]byte(config.CookieKey))
	gothic.Store = store
	goth.UseProviders(
		google.New(
			config.Google.ClientKey,
			config.Google.Secret,
			config.Google.CallbackURL,
			[]string{"email", "profile"}...,
		),
		github.New(
			config.Github.ClientKey,
			config.Github.Secret,
			config.Github.CallbackURL,
			[]string{"read:user", "user:email"}...,
		),
	)

	return &OAuth{}
}

func (o *OAuth) requestWithProvider(c *gin.Context, provider string) *http.Request {
	return c.Request.WithContext(o.contextWithProvider(c.Request.Context(), provider))
}

func (*OAuth) contextWithProvider(ctx context.Context, provider string) context.Context {
	return context.WithValue(ctx, "provider", provider)
}

func (o *OAuth) Begin(c *gin.Context, provider string) {
	gothic.BeginAuthHandler(
		c.Writer,
		o.requestWithProvider(c, provider),
	)
}

func (o *OAuth) Complete(c *gin.Context, provider string) (*User, error) {
	user, err := gothic.CompleteUserAuth(
		c.Writer,
		o.requestWithProvider(c, provider),
	)
	if err != nil {
		return nil, err
	}

	return newUser(user), nil
}

func (o *OAuth) Logout(c *gin.Context) error {
	return gothic.Logout(c.Writer, c.Request)
}
