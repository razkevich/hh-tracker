package middleware

import (
	"github.com/gin-gonic/gin"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/logging"
)

// InitLogger initializes a logger in context
func (m Middleware) InitLogger() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Set("logger", logging.BuildLogger(ctx.GetHeader("X-Moltin-Auth-Store"), ctx.GetHeader("X-Request-Id")))
		ctx.Next()
	}
}
