// To override the name of a config field in the yaml file,  you need to use the
// mapstructure tag instead of the yaml tag as you may expect. Viper starts by
// unmarshalling the yaml to map[string]interface. access_log is one example
package config

type EventStoreDbDriver string
type ProjectionDbDriver string

const EventStoreDbDriverPostgres EventStoreDbDriver = "postgres"
const ProjectionDbDriverOpensearch ProjectionDbDriver = "opensearch"
const ProjectionDbDriverPostgres ProjectionDbDriver = "postgres"

var availableEventStoreDbDrivers = []EventStoreDbDriver{
	EventStoreDbDriverPostgres,
}

func (val EventStoreDbDriver) IsKnown() bool {
	for _, chosen := range availableEventStoreDbDrivers {
		if chosen == val {
			return true
		}
	}

	return false
}

func (val EventStoreDbDriver) IsPostgres() bool {
	return val == EventStoreDbDriverPostgres
}

var availableProjectionDbDrivers = []ProjectionDbDriver{
	ProjectionDbDriverOpensearch,
	ProjectionDbDriverPostgres,
}

func (val ProjectionDbDriver) IsKnown() bool {
	for _, chosen := range availableProjectionDbDrivers {
		if chosen == val {
			return true
		}
	}

	return false
}

func (val ProjectionDbDriver) IsPostgres() bool {
	return val == ProjectionDbDriverPostgres
}

func (val ProjectionDbDriver) IsOpensearch() bool {
	return val == ProjectionDbDriverOpensearch
}

type ConfigLogging struct {
	Level  string
	Format string
}

type ConfigApiServerAccessLog struct {
	Enabled bool
	Format  string
}

type ConfigApiServer struct {
	DebugErrorsEnabled bool `yaml:"debug_errors_enabled" mapstructure:"debug_errors_enabled"`
	Port               int
	AccessLog          ConfigApiServerAccessLog `yaml:"access_log" mapstructure:"access_log"`
}

type ConfigApi struct {
	Server ConfigApiServer
}

type ConfigDbPostgres struct {
	User string
	Host string
	Password string
	Port uint32
	Database           string
	Schema string
	ConnectionRetries  int    `mapstructure:"connection_retries"`
}

func (cfg ConfigDbPostgres) DBHost() string {
	return cfg.Host
}

func (cfg ConfigDbPostgres) DBUser() string {
	return cfg.User
}

func (cfg ConfigDbPostgres) DBPassword() string {
	return cfg.Password
}

func (cfg ConfigDbPostgres) DBPort() uint32 {
	return cfg.Port
}

func (cfg ConfigDbPostgres) DBSchema() string {
	return cfg.Schema
}

func (cfg ConfigDbPostgres) DBName() string {
	return cfg.Database
}

func (cfg ConfigDbPostgres) DBConnectionRetries() int {
	return cfg.ConnectionRetries
}

type ConfigDbEventStore struct {
	Driver  EventStoreDbDriver
	Postgres ConfigDbPostgres
}

type ConfigDbProjection struct {
	Driver  ProjectionDbDriver
	Postgres ConfigDbPostgres
}

type ConfigDb struct {
	EventStore ConfigDbEventStore `mapstructure:"event_store"`
	Projection ConfigDbProjection
}

type ConfigAuth struct {
	MasterToken string `mapstructure:"master_token"`
	Bcrypt ConfigAuthBcrypt
}

type ConfigAuthBcrypt struct {
	Cost int
}

type Config struct {
	Auth ConfigAuth
	Logging        ConfigLogging
	Api            ConfigApi
	Db             ConfigDb
}

func (c *Config) LogLevel() string {
	return c.Logging.Level
}

func (c *Config) LogFormat() string {
	return c.Logging.Format
}

func (c *Config) ApiServerPort() int {
	return c.Api.Server.Port
}

func (c *Config) ApiServerAccessLogEnabled() bool {
	return c.Api.Server.AccessLog.Enabled
}

func (c *Config) ApiServerAccessLogFormat() string {
	return c.Api.Server.AccessLog.Format
}

func (c *Config) ApiServerDebugErrorsEnabled() bool {
	return c.Api.Server.DebugErrorsEnabled
}

func (c *Config) EventStoreDbDriver() EventStoreDbDriver {
	return c.Db.EventStore.Driver
}

func (c *Config) ProjectionDbDriver() ProjectionDbDriver {
	return c.Db.Projection.Driver
}

func (c *Config) AuthMasterToken() string {
	return c.Auth.MasterToken
}

func NewDefault() *Config {
	return &Config{
		Logging: ConfigLogging{
			Level:  "info",
			Format: "json",
		},
		Api: ConfigApi{
			Server: ConfigApiServer{
				DebugErrorsEnabled: false,
				Port:               8080,
				AccessLog: ConfigApiServerAccessLog{
					Enabled: true,
					Format:  "json",
				},
			},
		},
		Auth: ConfigAuth{
			Bcrypt: ConfigAuthBcrypt{
				Cost: 12,
			},
		},
		Db: ConfigDb{
			EventStore: ConfigDbEventStore{
				Driver: EventStoreDbDriverPostgres,
				Postgres: ConfigDbPostgres{
					User: "",
					Host: "",
					Password: "",
					Port: 5432,
					Database: "panoptes",
					Schema: "public",
					ConnectionRetries: 3,
				},
			},
			Projection: ConfigDbProjection{
				Driver: ProjectionDbDriverPostgres,
				Postgres: ConfigDbPostgres{
					User: "",
					Host: "",
					Password: "",
					Port: 5432,
					Database: "panoptes",
					Schema: "public",
					ConnectionRetries: 3,
				},
			},
		},
	}
}
