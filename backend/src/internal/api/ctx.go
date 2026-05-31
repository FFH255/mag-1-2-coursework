package api

import "github.com/gin-gonic/gin"

const (
	accessTokenCtxKey       = "access_token"
	refreshTokenCtxKey      = "refresh_token"
	registrationTokenCtxKey = "registration_token"
	userIDCtxKey            = "user_id"
)

func getString(c *gin.Context, key string) (value string, exists bool) {
	v, ok := c.Get(key)
	if !ok {
		return "", false
	}

	value, ok = v.(string)
	if !ok {
		return "", false
	}

	if value == "" {
		return "", false
	}

	return value, true
}

func SetAccessToken(c *gin.Context, token string) {
	c.Set(accessTokenCtxKey, token)
}

func GetAccessToken(c *gin.Context) string {
	token, exists := getString(c, accessTokenCtxKey)
	if !exists {
		panic("access token not found")
	}

	return token
}

func SetRefreshToken(c *gin.Context, token string) {
	c.Set(refreshTokenCtxKey, token)
}

func GetRefreshToken(c *gin.Context) string {
	token, exists := getString(c, refreshTokenCtxKey)
	if !exists {
		panic("refresh token not found")
	}

	return token
}

func SetRegistrationToken(c *gin.Context, token string) {
	c.Set(registrationTokenCtxKey, token)
}

func GetRegistrationToken(c *gin.Context) string {
	token, exists := getString(c, registrationTokenCtxKey)
	if !exists {
		panic("registration token not found")
	}

	return token
}

func SetUserID(c *gin.Context, userID uint64) {
	c.Set(userIDCtxKey, userID)
}

func GetUserID(c *gin.Context) uint64 {
	v, ok := c.Get(userIDCtxKey)
	if !ok {
		panic("user id not found")
	}

	id, ok := v.(uint64)
	if !ok {
		panic("user id not found")
	}

	return id
}
