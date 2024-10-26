package v1

import (
	"context"
	"net/http"

	"github.com/adamkirk/panoptes/internal/api/operations"
	"github.com/adamkirk/panoptes/internal/api/v1/responses"
	"github.com/adamkirk/panoptes/internal/domain/ingestion"
	"github.com/danielgtaylor/huma/v2"
)

type GithubIngestor interface {
	Process(e ingestion.GithubEvent) error
}

type GithubWebhookRequest struct {
	GithubEvent string `header:"X-GitHub-Event" required:"true"`
	GithubDelivery string `header:"X-GitHub-Delivery" required:"true"`
	Body map[string]any `doc:"Any webhook structure that github may send"`
}

type IngestionController struct {
	github GithubIngestor
}

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
		Security: []map[string][]string{
			{"scopes": {"ingest.github"}},
		},
	}, ErrorHandler(true, c.IngestGithubWebhook))
}

func NewIngestController(gh GithubIngestor) *IngestionController {
	return &IngestionController{
		github: gh,
	}
}

func (c *IngestionController) IngestGithubWebhook(ctx context.Context, req *GithubWebhookRequest) (*responses.NoContent, error) {
	e := ingestion.GithubEvent{
		Payload: req.Body,
		DeliveryID: req.GithubDelivery,
		Event: req.GithubEvent,
	}

	if err := c.github.Process(e); err != nil {
		return nil, err
	}

	return &responses.NoContent{
		Status: http.StatusNoContent,
	}, nil
}
