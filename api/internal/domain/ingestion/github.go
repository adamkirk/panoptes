package ingestion

import (
	"time"

	"github.com/adamkirk/heimdallr/internal/util/dt"
	"github.com/google/uuid"
)


type GithubIngestorRepo interface {
	Create(id uuid.UUID, ts time.Time, payload map[string]any) error
}

type GithubEvent struct {
	Payload map[string]any
	Event string
	DeliveryID string
}

type GithubIngestorOpt func(*GithubIngestor)

// WithCustomNowProvider allows you to override the way we generate a timestamp
// for now. By default it will use the db.NowUTC function. Changing this can be 
// useful for testing and if we wanna change the timezone and such.
func WithCustomNowProvider(f func () time.Time) GithubIngestorOpt {
	return func (gi *GithubIngestor) {
		gi.getNow = f
	}
}

type GithubIngestor struct {
	repo GithubIngestorRepo
	getNow func() time.Time
}

func (gi *GithubIngestor) Process(e GithubEvent) error {
	return gi.repo.Create(uuid.New(), gi.getNow(), e.Payload)
}

func NewGithubIngestor(repo GithubIngestorRepo, opts... GithubIngestorOpt) *GithubIngestor {
	gi := &GithubIngestor{
		repo: repo,
		getNow: dt.NowUTC,
	}

	for _, opt := range opts {
		opt(gi)
	}

	return gi
}