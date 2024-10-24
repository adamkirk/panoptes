package api

import (
	"github.com/labstack/echo/v4"
)

type V1Api struct {
	controllers []Controller
}

func (c *V1Api) Version() string {
	return "v1"
}

func (c *V1Api) Middleware(cfg ApiServerConfig) []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		// v1.NewErrorHandler(cfg.ApiServerDebugErrorsEnabled()),
	}
}

func (c *V1Api) Controllers() []Controller {
	return c.controllers
}

func NewV1Api(controllers []Controller) *V1Api {
	return &V1Api{
		controllers: controllers,
	}
}
