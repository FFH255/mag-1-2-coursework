package request_id_middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/ruslanonly/blindtyping/src/internal"
	"github.com/ruslanonly/blindtyping/src/internal/shared/uuid_generator"
)

func New(logger internal.Logger, generator *uuid_generator.Generator) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := generator.Generate()
		withRequestID(c, logger, requestID)
	}
}

func withRequestID(c *gin.Context, logger internal.Logger, requestID string) {
	c.Request = c.Request.WithContext(logger.WithRequestID(c.Request.Context(), requestID))
}
