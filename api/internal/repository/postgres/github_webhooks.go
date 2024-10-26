package postgres

import (
	// dot import so go code would resemble as much as native SQL
	// dot import is not mandatory
	"encoding/json"
	"time"

	"github.com/adamkirk/heimdallr/internal/repository/postgres/schema/heimdallr/public/table"
	"github.com/google/uuid"
)

type GithubWebhooksRepository struct {
	conn *Connector
}

func (r *GithubWebhooksRepository) Create(id uuid.UUID, ts time.Time, payload map[string]any) error {
	conn, err := r.conn.Connection()

	if err != nil {
		return err
	}

	var payloadJSON []byte
	if payloadJSON, err = json.Marshal(payload); err != nil {
		return err
	}

	stmt := table.GithubWebhooks.INSERT(table.GithubWebhooks.ID, table.GithubWebhooks.OccurredAt, table.GithubWebhooks.Payload).
		VALUES(uuid.New(), ts, string(payloadJSON))

	if _, err := stmt.Exec(conn); err != nil {
		return err
	}

	return nil
}

func NewGithubWebhooksRepository(conn *Connector) *GithubWebhooksRepository {
	return &GithubWebhooksRepository{
		conn: conn,
	}
}