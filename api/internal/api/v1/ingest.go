package v1

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/adamkirk/heimdallr/internal/api/operations"
	"github.com/adamkirk/heimdallr/internal/api/v1/responses"
	"github.com/danielgtaylor/huma/v2"
)

type GithubWebhookRequest struct {
	GithubEvent string `header:"X-GitHub-Event" required:"true"`
	GithubDelivery string `header:"X-GitHub-Delivery" required:"true"`
	Body map[string]any `doc:"Any webhook structure that github may send"`
}

type IngestionController struct {}

func (c *IngestionController) RegisterRoutes(api huma.API) {
	huma.Register[GithubWebhookRequest, responses.NoContent](api, huma.Operation{
		OperationID:  "v1.ingest.github",
		Method:       http.MethodPost,
		Path:         "/ingestion/github",
		Summary:      "Ingest a webhook event from github",
		DefaultStatus: http.StatusNoContent,
		Metadata: map[string]any{
			operations.OptDisableNotFound: true,
		},
	}, ErrorHandler(true, c.IngestGithubWebhook))
}

func NewIngestController(
) *IngestionController {
	return &IngestionController{}
}

func (c *IngestionController) IngestGithubWebhook(ctx context.Context, req *GithubWebhookRequest) (*responses.NoContent, error) {

	slog.Debug("github event ingested", "request", req)
	return &responses.NoContent{
		Status: http.StatusNoContent,
	}, nil
}
