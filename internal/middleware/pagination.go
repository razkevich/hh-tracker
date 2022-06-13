package middleware

import (
	"github.com/gin-gonic/gin"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/controller/helper"
	mongo "gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/driver"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/internalerrors"
	"strconv"
)

// MiddlewareConfig represents middleware configuration.
type MiddlewareConfig struct {
	EnforceLimits    bool
	DefaultPageLimit int
}

// Middleware is the collection of middlewares used in system
type Middleware struct {
	Config      MiddlewareConfig
	mongoClient *mongo.Client
}

// ProvideMiddleware is a wire function for providing pagination
func ProvideMiddleware(conf MiddlewareConfig, mongoClient *mongo.Client) Middleware {
	return Middleware{Config: conf, mongoClient: mongoClient}
}

// OffsetLimit is the cut off for offset when getting paginated lists
const OffsetLimit = 10000

// PageLengthLimit is the cut off for page limit when getting paginated lists
const PageLengthLimit = 100

// CheckPaginationOffsetLimit checks the offset does not exceed the enforced limitation
func (m Middleware) CheckPaginationOffsetLimit(ctx *gin.Context) {

	offset := ctx.Request.URL.Query().Get("page[offset]")
	intOffset, _ := strconv.ParseInt(offset, 10, 64)

	limit := ctx.Request.URL.Query().Get("page[limit]")
	intLimit, _ := strconv.ParseInt(limit, 10, 64)

	if m.Config.EnforceLimits && intOffset > OffsetLimit {
		res := helper.BuildErrorResponse(ctx, internalerrors.ErrPageOffsetExceed.Error(), internalerrors.ErrPageOffsetExceed)
		ctx.AbortWithStatusJSON(res.Status, res.ErrorResponse)
		return
	}
	if m.Config.EnforceLimits && intLimit > PageLengthLimit {
		res := helper.BuildErrorResponse(ctx, internalerrors.ErrPageLimitExceed.Error(), internalerrors.ErrPageLimitExceed)
		ctx.AbortWithStatusJSON(res.Status, res.ErrorResponse)
		return
	}
	ctx.Next()
}

// SetPaginationSetting gets the page offset and limit from query
// if limit is not set in the query params, it sets from the settings service (which is the X-Moltin-Settings-page_length in the header)
func (m Middleware) SetPaginationSetting(ctx *gin.Context) {
	if ctx.Request.URL.Query().Get("page[limit]") == "" {
		pageLensetting := ctx.Request.Header.Get("X-Moltin-Settings-page_length")
		if pageLensetting == "" {
			// default page limit for local development
			ctx.Set("page[limit]", strconv.Itoa(m.Config.DefaultPageLimit))
		} else {
			if intPageLen, _ := strconv.ParseInt(pageLensetting, 10, 64); intPageLen > PageLengthLimit {
				pageLensetting = strconv.Itoa(PageLengthLimit)
			}
			ctx.Set("page[limit]", pageLensetting)
		}
	} else {
		ctx.Set("page[limit]", ctx.Request.URL.Query().Get("page[limit]"))
	}
	ctx.Set("page[offset]", ctx.Request.URL.Query().Get("page[offset]"))
	ctx.Next()
}
