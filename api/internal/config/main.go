// To override the name of a config field in the yaml file,  you need to use the
// mapstructure tag instead of the yaml tag as you may expect. Viper starts by
// unmarshalling the yaml to map[string]interface. access_log is one example
package config

type DbDriver string

const DbDriverPostgres DbDriver = "postgres"

var availableDbDrivers = []DbDriver{
	DbDriverPostgres,
}

func (val DbDriver) IsKnown() bool {
	for _, chosen := range availableDbDrivers {
		if chosen == val {
			return true
		}
	}

	return false
}

func (val DbDriver) IsPostgres() bool {
	return val == DbDriverPostgres
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
	Host string
	Password string
	Port uint32
	Database           string
	Schema string
	ConnectionRetries  int    `mapstructure:"connection_retries"`
}

type ConfigDb struct {
	Driver  DbDriver
	Postgres ConfigDbPostgres
}


type Config struct {
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

func (c *Config) DbDriver() DbDriver {
	return c.Db.Driver
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
		Db: ConfigDb{
			Driver: DbDriverPostgres,
			Postgres: ConfigDbPostgres{
				Host: "",
				Password: "",
				Port: 5432,
				Database: "heimdallr",
				Schema: "public",
				ConnectionRetries: 3,
			},
		},
	}
}
