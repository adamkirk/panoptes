package api

import (
	"context"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// setupLoggingMiddleare adds a curtom logger to the given echo server
// We're purposely using slgo here so that it will default to whatever type of
// logger we initially used .e.g JSON or TEXT
func buildLoggingMiddleware(format string) echo.MiddlewareFunc {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	// logger := slog.Default().With(slog.String("source", "access_log"))
	var logger *slog.Logger

	// TODO: error if the format is invalid
	if format == "json" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	}

	// See https://echo.labstack.com/docs/middleware/logger#examples
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:    true,
		LogURI:       true,
		LogError:     true,
		LogLatency:   true,
		LogMethod:    true,
		LogRequestID: true,
		HandleError:  true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {

			level := slog.LevelInfo
			errorMsg := "nil"

			if v.Error != nil {
				errorMsg = v.Error.Error()
				level = slog.LevelError
			}

			logger.LogAttrs(context.Background(), level, "REQUEST",
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.String("method", v.Method),
				slog.String("request-id", v.RequestID),
				slog.String("log_type", "access"),

				// Gives us a consistent id that we can use for filtering,
				// aggregation rather than regexing our way through life
				slog.String("route_id", c.Path()),
				// Convert to milliseconds
				slog.Float64("duration", float64(v.Latency.Microseconds())/1000),
				slog.String("err", errorMsg),
			)
			return nil
		},
	})
}
