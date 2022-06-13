package controller

import (
	"github.com/gin-gonic/gin"
	VV "github.com/go-playground/validator/v10"
	"github.com/thoas/go-funk"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/controller/helper"
	Links "gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/controller/links"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/internalerrors"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/service"
	"net/http"
	"strconv"
)

// ControllerConfig is the HTTP server configuration.
type ControllerConfig struct {
	SvcBaseURL string
}

// LogEntriesController contains controller context including business service.
type LogEntriesController struct {
	service service.ILogEntriesService
	config  ControllerConfig
}

var validate *VV.Validate

// ProvideLogEntriesController is the wire function for the controller.
func ProvideLogEntriesController(service service.LogEntriesService, config ControllerConfig) LogEntriesController {
	validate = VV.New()
	_ = validate.RegisterValidation("positive-pagination-limit", helper.ValidatePaginationLimit)
	_ = validate.RegisterValidation("non-negative-pagination-offset", helper.ValidatePaginationOffset)
	return LogEntriesController{
		service: service,
		config:  config,
	}
}

// GetLogEntries get log entries
func (a LogEntriesController) GetLogEntries(ctx *gin.Context) {
	var validateAuthStoreHeader helper.ValidateAuthStoreHeader
	if err := ctx.ShouldBindHeader(&validateAuthStoreHeader); err != nil {
		res := helper.BuildErrorResponse(ctx, GetValidationErrorMessage(err), internalerrors.ErrBadRequest)
		ctx.JSON(res.Status, res.ErrorResponse)
		return
	}
	var validatePagination helper.ValidatePagination
	if err := ctx.ShouldBindQuery(&validatePagination); err != nil {
		res := helper.BuildErrorResponse(ctx, err.Error(), internalerrors.ErrBadRequest)
		ctx.JSON(res.Status, res.ErrorResponse)
		return
	}
	if err := validate.Struct(validatePagination); err != nil {
		res := helper.BuildErrorResponse(ctx, GetValidationErrorMessage(err), internalerrors.ErrBadRequest)
		ctx.JSON(res.Status, res.ErrorResponse)
		return
	}
	limit, _ := ctx.Get("page[limit]")
	offset, _ := ctx.Get("page[offset]")
	storeID := ctx.Request.Header.Get("X-Moltin-Auth-Store")

	search := ctx.Request.Header.Get("ep-internal-search-json")
	searchJSON := ""
	if !funk.IsEmpty(search) {
		data, err := helper.ConvertRQLToJSON(search)
		if err != nil {
			res := helper.BuildErrorResponse(ctx, err.Error(), internalerrors.ErrBadRequest)
			ctx.JSON(res.Status, res.ErrorResponse)
			return
		}
		searchJSON = data
	}

	intLimit, _ := strconv.ParseInt(limit.(string), 10, 64)
	intOffset, _ := strconv.ParseInt(offset.(string), 10, 64)

	list, err := a.service.GetLogEntries(ctx, storeID, int(intOffset), int(intLimit), searchJSON)
	if err != nil {
		res := helper.BuildErrorResponse(ctx, GetValidationErrorMessage(err), err)
		ctx.JSON(res.Status, res.ErrorResponse)
		return
	}

	links := Links.CreatePaginationLinks(ctx.Request, a.config.SvcBaseURL, intOffset, intLimit, list.TotalCount)
	page, results := helper.GetPaginationPageAndResults(intLimit, intOffset, list.TotalCount)
	meta := helper.Meta{
		Page:    &page,
		Results: &results,
	}
	fullLogEntries := make([]helper.FullLogEntry, 0, len(list.LogEntries))
	for _, relations := range list.LogEntries {
		selfLinks := helper.Links{
			Self: relations.GetSelfLink(a.config.SvcBaseURL),
		}
		fullLogEntry := helper.BuildLogEntryWithLinks(*relations, selfLinks)
		fullLogEntries = append(fullLogEntries, fullLogEntry)
	}

	res := helper.BuildFullResponseWithMeta(fullLogEntries, links, meta)
	ctx.JSON(http.StatusOK, res)
}
