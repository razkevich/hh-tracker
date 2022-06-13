package internal

import (
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/controller"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/middleware"
	"io"
	"time"

	// Register Swagger docs
	"net/http"
	"strconv"
)

// ServerConfig is the HTTP server configuration.
type ServerConfig struct {
	Port int
}

// ProvideServer is a Server provider for Gin
func ProvideServer(controller controller.LogEntriesController,
	healthController controller.HealthController, config ServerConfig,
	middleware middleware.Middleware) *http.Server {
	return NewServer(controller, healthController, middleware, config)
}

// NewServer returns an initialized Gin Server
func NewServer(logEntriesController controller.LogEntriesController, healthController controller.HealthController, middleware middleware.Middleware,
	config ServerConfig) *http.Server {

	router := initializeGin()

	v2 := router.Group("/v2", middleware.PopulateRequestID(), middleware.InitLogger())
	{
		logs := v2.Group("/personal-data/logs")
		{
			logs.GET("", middleware.SetPaginationSetting, middleware.CheckPaginationOffsetLimit, logEntriesController.GetLogEntries)
		}

	}

	router.GET("/checks/readiness", healthController.Readiness)
	router.GET("/checks/healthz", healthController.Liveness)
	router.HandleMethodNotAllowed = true
	return &http.Server{
		Addr:    ":" + strconv.Itoa(config.Port),
		Handler: router,
	}
}

func initializeGin() *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())

	r.Use(
		logger.SetLogger(
			logger.WithUTC(true),
			logger.WithLogger(
				func(c *gin.Context, out io.Writer, latency time.Duration) zerolog.Logger {
					obj := struct {
						URL        string `json:"url"`
						StatusCode int    `json:"status_code"`
						Method     string `json:"method"`
						UserAgent  string `json:"useragent"`
					}{
						URL:        c.Request.URL.Path,
						StatusCode: c.Writer.Status(),
						Method:     c.Request.Method,
						UserAgent:  c.Request.UserAgent(),
					}

					return zerolog.New(out).With().
						Str(zerolog.TimestampFieldName, time.Now().UTC().Format(time.RFC3339)).
						Str("x_request_id", c.Request.Header.Get("X-Request-Id")).
						Int64("duration", int64(latency)).
						Interface("http", obj).
						Str("store_uuid", c.Request.Header.Get("X-Moltin-Auth-Store")).
						Logger()
				},
			),
		),
	)

	return r
}
