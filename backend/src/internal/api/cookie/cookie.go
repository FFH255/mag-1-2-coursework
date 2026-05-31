package cookie

import (
	"github.com/gin-gonic/gin"
)

var (
	accessTokenCookieName = "access_token"
	accessTokenCookieTTL  = 60 * 60 * 2 // 2 hours

	refreshTokenCookieName = "refresh_token"
	refreshTokenCookieTTL  = 60 * 60 * 24 * 14 // 2 weeks

	registrationTokenCookieName = "registration_token"
	registrationTokenCookieTTL  = 60 * 5 // 5 min
)

type Manager struct {
	path   string
	domain string
	secure bool
}

func NewManager(secure bool) *Manager {
	return &Manager{
		path:   "/",
		domain: "",
		secure: secure,
	}
}

func (m *Manager) SetAccessToken(ctx *gin.Context, accessToken string) {
	ctx.SetCookie(
		accessTokenCookieName,
		accessToken,
		accessTokenCookieTTL,
		m.path,
		m.domain,
		m.secure,
		false,
	)
}

func (m *Manager) GetAccessToken(ctx *gin.Context) (string, error) {
	return ctx.Cookie(accessTokenCookieName)
}

func (m *Manager) DeleteAccessToken(ctx *gin.Context) {
	ctx.SetCookie(
		accessTokenCookieName,
		"",
		-1,
		m.path,
		m.domain,
		m.secure,
		false,
	)
}

func (m *Manager) SetRefreshToken(ctx *gin.Context, refreshToken string) {
	ctx.SetCookie(
		refreshTokenCookieName,
		refreshToken,
		refreshTokenCookieTTL,
		m.path,
		m.domain,
		m.secure,
		true,
	)
}

func (m *Manager) GetRefreshToken(ctx *gin.Context) (string, error) {
	return ctx.Cookie(refreshTokenCookieName)
}

func (m *Manager) DeleteRefreshToken(ctx *gin.Context) {
	ctx.SetCookie(
		refreshTokenCookieName,
		"",
		-1,
		m.path,
		m.domain,
		m.secure,
		true,
	)
}

func (m *Manager) SetRegistrationToken(ctx *gin.Context, refreshToken string) {
	ctx.SetCookie(
		registrationTokenCookieName,
		refreshToken,
		registrationTokenCookieTTL,
		m.path,
		m.domain,
		m.secure,
		true,
	)
}

func (m *Manager) GetRegistrationToken(ctx *gin.Context) (string, error) {
	return ctx.Cookie(registrationTokenCookieName)
}

func (m *Manager) DeleteRegistrationToken(ctx *gin.Context) {
	ctx.SetCookie(
		registrationTokenCookieName,
		"",
		-1,
		m.path,
		m.domain,
		m.secure,
		true,
	)
}
