package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type ErrFailedToConnect struct {
	Err error
}

func (err ErrFailedToConnect) Error() {

}

type Config interface {
	DBHost() string 
	DBUser() string
	DBPassword() string 
	DBPort() uint32 
	DBSchema() string 
	DBName() string 
}

type Connector struct {
	db *sql.DB

	cfg Config
}

func (c *Connector) Connection() (*sql.DB, error) {
	if c.db != nil {
		return c.db, nil
	}

	// TODO: support SSL
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?search_path=%s&sslmode=disable",
		c.cfg.DBUser(),
		c.cfg.DBPassword(),
		c.cfg.DBHost(),
		c.cfg.DBPort(),
		c.cfg.DBName(),
		c.cfg.DBSchema(),
	)

	db, err := sql.Open("postgres", connString)

	if err != nil {
		return nil, err
	}
	
		
	c.db = db
	return db, nil
}

func NewConnector(cfg Config) *Connector {
	return &Connector{
		cfg: cfg,
	}
}