package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strconv"

	"github.com/adamkirk/heimdallr/internal/api/operations"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type RouteHandler func(e echo.Context) error

type Route struct {
	Handler RouteHandler
	Path    string
	Name    string
	Method  string
}

type Routes []Route

type Controller interface {
	RegisterRoutes(g huma.API)
}

type ApiServerConfig interface {
	ApiServerPort() int
	ApiServerAccessLogEnabled() bool
	ApiServerAccessLogFormat() string
	ApiServerDebugErrorsEnabled() bool
}

type Server struct {
	cfg ApiServerConfig
	e   *echo.Echo
}

func (s *Server) Start() error {
	slog.Info(fmt.Sprintf("listening on port %d", s.cfg.ApiServerPort()))
	return s.e.Start(fmt.Sprintf(":%d", s.cfg.ApiServerPort()))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}

var opsWithoutBodies = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodOptions,
	http.MethodDelete,

	// Probably don't need this one, but leaving for good measure
	http.MethodTrace,
}

func ConfigureDefaultResponses(api *huma.OpenAPI, op *huma.Operation) {

	if _, ok := op.Responses["default"]; ok {
		// Remove the default as it's an error, but has no status code
		// Maybe there's another way to turn it off
		op.Responses["default"] = nil
	}

	if v, ok := op.Metadata[operations.OptDisableAllDefaults]; ok {
		if optAsBool, ok := v.(bool); ok && optAsBool {
			return
		}
	}

	validationStatus := strconv.Itoa(http.StatusUnprocessableEntity)

	if slices.Contains(opsWithoutBodies, op.Method) {
		validationStatus = strconv.Itoa(http.StatusBadRequest)
	}

	if _, ok := op.Responses["default"]; ok {
		// Remove the default as it's an error, but has no status code
		// Maybe there's another way to turn it off
		op.Responses["default"] = nil
	}

	if _, ok := op.Responses[validationStatus]; !ok {
		op.Responses[validationStatus] = &huma.Response{
			Description: "validation error",
			Content: map[string]*huma.MediaType{
				"application/problem+json": {
					Schema: &huma.Schema{
						Ref: "#/components/schemas/ErrorModel",
					},
				},
			},
		}
	}

	internalError := strconv.Itoa(http.StatusInternalServerError)

	if _, ok := op.Responses[internalError]; !ok {
		op.Responses[internalError] = &huma.Response{
			Description: "validation error",
			Content: map[string]*huma.MediaType{
				"application/problem+json": {
					Schema: &huma.Schema{
						Ref: "#/components/schemas/ErrorModel",
					},
				},
			},
		}
	}

	var notFoundEnabled = true

	if v, ok := op.Metadata[operations.OptDisableNotFound]; ok {
		if optAsBool, ok := v.(bool); ok {
			notFoundEnabled = ! optAsBool
		}
	}

	notFound := strconv.Itoa(http.StatusNotFound)

	if _, ok := op.Responses[notFound]; !ok && notFoundEnabled {
		op.Responses[notFound] = &huma.Response{
			Description: "Resource Not Found",
			Content: map[string]*huma.MediaType{
				"application/problem+json": {
					Schema: &huma.Schema{
						Ref: "#/components/schemas/ErrorModel",
					},
				},
			},
		}
	}
}

func NewServer(v1Api *V1Api, cfg ApiServerConfig) *Server {
	e := echo.New()
	
	e.HideBanner = true
	e.HidePort = true
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RequestID())

	if cfg.ApiServerAccessLogEnabled() {
		e.Use(buildLoggingMiddleware(cfg.ApiServerAccessLogFormat()))
	}
	
	for _, c := range v1Api.Controllers() {
		apiBase := fmt.Sprintf("/api/%s", v1Api.Version())
		api := e.Group(apiBase)
		cfg := huma.DefaultConfig("Organisations", v1Api.Version())

		// Needed to get the docs displaying properly.
		cfg.OpenAPI.Servers = []*huma.Server{
			{
				URL: apiBase,
			},
		}

		hg := humaecho.NewWithGroup(e, api, cfg)
		hg.OpenAPI().OnAddOperation = append(hg.OpenAPI().OnAddOperation, ConfigureDefaultResponses)
		c.RegisterRoutes(hg)
	}

	srv := &Server{
		cfg: cfg,
		e:   e,
	}

	return srv
}
