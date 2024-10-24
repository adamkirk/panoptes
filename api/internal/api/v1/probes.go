package v1

import (
	"context"
	"net/http"

	"github.com/adamkirk/heimdallr/internal/api/operations"
	"github.com/adamkirk/heimdallr/internal/api/v1/responses"
	"github.com/danielgtaylor/huma/v2"
)

type StartupRequest struct {}

type ProbesController struct {}

func (c *ProbesController) RegisterRoutes(api huma.API) {
	huma.Register[StartupRequest, responses.NoContent](api, huma.Operation{
		OperationID:  "v1.probes.startup",
		Method:       http.MethodGet,
		Path:         "/_probes/startup",
		Summary:      "Check if the app is started up",
		DefaultStatus: http.StatusNoContent,
		Metadata: map[string]any{
			operations.OptDisableAllDefaults: true,
		},
	}, ErrorHandler(true, c.Startup))
}

func NewProbesController(
) *ProbesController {
	return &ProbesController{}
}

func (c *ProbesController) Startup(ctx context.Context, req *StartupRequest) (*responses.NoContent, error) {
	return &responses.NoContent{
		Status: http.StatusNoContent,
	}, nil
}
