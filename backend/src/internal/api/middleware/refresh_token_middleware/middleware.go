package refresh_token_middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ruslanonly/blindtyping/src/internal/api"
	"github.com/ruslanonly/blindtyping/src/internal/api/middleware"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
)

type cookieManager interface {
	GetRefreshToken(ctx *gin.Context) (string, error)
}

type Middleware struct {
	cookieManager cookieManager
}

func (m *Middleware) Handle(c *gin.Context) {
	refreshToken, err := m.cookieManager.GetRefreshToken(c)
	if err != nil {
		proto.WriteError(c, http.StatusUnauthorized, "refresh token cookie not found")
		return
	}

	api.SetRefreshToken(c, refreshToken)

	c.Next()
}

func (m *Middleware) Name() string {
	return middleware.RefreshToken
}

func New(cookieManager cookieManager) proto.Middleware {
	return &Middleware{
		cookieManager: cookieManager,
	}
}
