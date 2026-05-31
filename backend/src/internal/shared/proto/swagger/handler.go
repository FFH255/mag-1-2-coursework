package swagger

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
)

type Handler struct {
	handler gin.HandlerFunc
}

// Handle
// @title Blindtyping
// @version 0.1
// @BasePath /
func (h Handler) Handle(c *gin.Context) {
	h.handler(c)
}

func (h Handler) Path() string {
	return "/swagger/*any"
}

func (h Handler) Method() string {
	return http.MethodGet
}

func (h Handler) Middleware() []string {
	return nil
}

func New(login, password string) proto.Handler {
	basicAuth := gin.BasicAuth(gin.Accounts{
		login: password,
	})

	handler := func(c *gin.Context) {
		basicAuth(c)

		if !c.IsAborted() {
			ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
		}
	}

	return &Handler{
		handler: handler,
	}
}
