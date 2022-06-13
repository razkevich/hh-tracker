package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/controller/helper"
	mongo "gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/driver"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/health"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/internalerrors"
	"net/http"
)

// HealthController contains HealthController context including health service.
type HealthController struct {
	health      health.IHealth
	mongoClient *mongo.Client
}

// ProvideHealthController is the wire function for checks HealthController.
func ProvideHealthController(health health.IHealth, mongoClient *mongo.Client) HealthController {
	return HealthController{
		health:      health,
		mongoClient: mongoClient,
	}
}

// Liveness the K8S Liveness handler.
func (h *HealthController) Liveness(ctx *gin.Context) {
	res := helper.WrapData(helper.EmptyObj{})

	ctx.JSON(http.StatusNoContent, res)
}

// Readiness the K8S Readiness handler.
func (h *HealthController) Readiness(ctx *gin.Context) {
	res := helper.WrapData(helper.EmptyObj{})

	if !h.health.Readiness() {
		log.Err(internalerrors.ErrSvcNotAvailable)
		ctx.JSON(http.StatusServiceUnavailable, res)
		return
	}

	ctx.JSON(http.StatusNoContent, res)
}
