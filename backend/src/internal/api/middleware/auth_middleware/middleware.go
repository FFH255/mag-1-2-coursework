package auth_middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ruslanonly/blindtyping/src/internal"
	"github.com/ruslanonly/blindtyping/src/internal/api"
	"github.com/ruslanonly/blindtyping/src/internal/api/middleware"
	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/services/auth_service"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
)

type tokenValidator interface {
	ValidateToken(token string) (models.ID, error)
}

type cookieManager interface {
	GetAccessToken(ctx *gin.Context) (string, error)
	GetRefreshToken(ctx *gin.Context) (string, error)
}

type Middleware struct {
	tokenValidator tokenValidator
	cookieManager  cookieManager
	logger         internal.Logger
}

func (m *Middleware) withUserID(c *gin.Context, userID models.ID) {
	ctx := c.Request.Context()
	ctx = m.logger.WithUserID(ctx, int64(userID))
	c.Request = c.Request.WithContext(ctx)
}

func (m *Middleware) Handle(c *gin.Context) {
	accessToken, err := m.cookieManager.GetAccessToken(c)
	if err != nil {
		proto.WriteError(c, http.StatusUnauthorized, "access token cookie not found")
		return
	}

	refreshToken, err := m.cookieManager.GetRefreshToken(c)
	if err != nil {
		proto.WriteError(c, http.StatusUnauthorized, "refresh token cookie not found")
		return
	}

	userID, err := m.tokenValidator.ValidateToken(accessToken)
	if err != nil {
		handleError(c, err)
		return
	}

	api.SetAccessToken(c, accessToken)
	api.SetRefreshToken(c, refreshToken)
	api.SetUserID(c, uint64(userID))

	m.withUserID(c, userID)

	c.Next()
}

func (m *Middleware) Name() string {
	return middleware.Auth
}

func New(tokenValidator tokenValidator, cookieManager cookieManager, logger internal.Logger) proto.Middleware {
	return &Middleware{
		tokenValidator: tokenValidator,
		cookieManager:  cookieManager,
		logger:         logger,
	}
}

func handleError(c *gin.Context, err error) {
	var (
		status  = http.StatusUnauthorized
		message = "unauthorized"
	)

	switch {
	case auth_service.IsBlockedAccessTokenError(err):
		status = http.StatusUnauthorized
		message = "access token is blocked"
	}

	proto.WriteError(c, status, message)
}
