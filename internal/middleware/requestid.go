package middleware

import (
	"github.com/gin-gonic/gin"
)

// PopulateRequestID generates the request id in the context
func (m Middleware) PopulateRequestID() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Set("x_request_id", ctx.GetHeader("X-Request-Id"))
		ctx.Next()
	}
}
