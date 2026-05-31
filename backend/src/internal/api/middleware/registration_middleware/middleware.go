package registration_middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ruslanonly/blindtyping/src/internal/api"

	"github.com/ruslanonly/blindtyping/src/internal/api/middleware"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
)

type tokenValidator interface {
	ValidateRegistrationToken(token string) error
}

type cookieManager interface {
	GetRegistrationToken(c *gin.Context) (string, error)
}

type Middleware struct {
	tokenValidator tokenValidator
	cookieManager  cookieManager
}

func (m *Middleware) Handle(c *gin.Context) {
	registrationToken, err := m.cookieManager.GetRegistrationToken(c)
	if err != nil {
		proto.WriteError(c, http.StatusUnauthorized, "registration token cookie not found")
		return
	}

	err = m.tokenValidator.ValidateRegistrationToken(registrationToken)
	if err != nil {
		proto.WriteError(c, http.StatusUnauthorized, "invalid or expired registration token")
		return
	}

	api.SetRegistrationToken(c, registrationToken)

	c.Next()
}

func (m *Middleware) Name() string {
	return middleware.Registration
}

func New(tokenValidator tokenValidator, cookieManager cookieManager) proto.Middleware {
	return &Middleware{
		tokenValidator: tokenValidator,
		cookieManager:  cookieManager,
	}
}
