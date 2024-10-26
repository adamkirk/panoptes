package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	apicmd "github.com/adamkirk/heimdallr/cmd/api"
	"github.com/adamkirk/heimdallr/internal/api"
	v1 "github.com/adamkirk/heimdallr/internal/api/v1"
	"github.com/adamkirk/heimdallr/internal/config"
	"github.com/adamkirk/heimdallr/internal/domain/ingestion"
	"github.com/adamkirk/heimdallr/internal/repository/postgres"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var cfgFile string
var appCfg *config.Config

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Heimdallr Organisations API service",
	Long:  `Generates metrics about development processes by combining info from task management and VCS services.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var apiServeCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the API server",
	Long:  `Blah`,
	Run: func(cmd *cobra.Command, args []string) {
		apicmd.Handler(SharedOpts(appCfg), cmd, args)
	},
}

func newFs() afero.Fs {
	return afero.NewOsFs()
}

func SharedOpts(cfg *config.Config) []fx.Option {
	opts := []fx.Option{
		fx.Provide(buildConfig),
		fx.Provide(
			fx.Annotate(
				buildConfig,
				fx.As(new(api.ApiServerConfig)),
			),
		),
		fx.Provide(api.NewServer),
		fx.Provide(
			fx.Annotate(
				api.NewV1Api,
				fx.ParamTags(`group:"api.v1.controllers"`),
			),
		),
		fx.Provide(
			fx.Annotate(
				v1.NewProbesController,
				fx.As(new(api.Controller)),
				fx.ResultTags(`group:"api.v1.controllers"`),
			),
		),
		fx.Provide(
			fx.Annotate(
				v1.NewIngestController,
				fx.As(new(api.Controller)),
				fx.ResultTags(`group:"api.v1.controllers"`),
			),
		),

		fx.Provide(
			fx.Annotate(
				ingestion.NewGithubIngestor,
				fx.As(new(v1.GithubIngestor)),
			),
		),
	}

	if !cfg.EventStoreDbDriver().IsKnown() {
		slog.Error("Unknown persistent db driver", "driver", string(appCfg.EventStoreDbDriver()))
		os.Exit(1)
	}

	// Register difference implementations based on configured driver
	if cfg.EventStoreDbDriver().IsPostgres() {
		opts = append(opts, []fx.Option{
			fx.Provide(
				func (cfg *config.Config) *postgres.Connector {
					return postgres.NewConnector(cfg.Db.EventStore.Postgres)
				},
			),
			fx.Provide(
				fx.Annotate(
					postgres.NewGithubWebhooksRepository,
					fx.As(new(ingestion.GithubIngestorRepo)),
				),
			),
		}...)
	}
	return opts
}



func buildConfig() *config.Config {
	if appCfg != nil {
		return appCfg
	}

	c := config.NewDefault()
	err := viper.Unmarshal(c)
	cobra.CheckErr(err)

	appCfg = c
	return appCfg
}

func init() {
	cobra.OnInitialize(bootstrap)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is WORKING_DIRECTORY/config.yaml)")
	rootCmd.PersistentFlags().String("log-level", "info", "log level to use")
	rootCmd.PersistentFlags().String("log-format", "json", "log format to use")
	rootCmd.PersistentFlags().Int("port", 8080, "port to serve API on")

	rootCmd.AddCommand(apiServeCmd)

	viper.BindPFlag("logging.level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("logging.format", rootCmd.PersistentFlags().Lookup("log-format"))
	viper.BindPFlag("api.server.port", rootCmd.PersistentFlags().Lookup("port"))
}

func bootstrap() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		currentDir, err := os.Getwd()
		cobra.CheckErr(err)

		viper.AddConfigPath(currentDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// Tell viper to replace . in nested path with underscores
	// e.g. logging.level becomes LOGGING_LEVEL
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetEnvPrefix("heimdallr")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	cobra.CheckErr(err)

	cfg := buildConfig()

	l := slog.Level(slog.LevelInfo)
	err = l.UnmarshalText([]byte(cfg.Logging.Level))

	cobra.CheckErr(err)

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     l,
	}

	var logger *slog.Logger

	if cfg.LogFormat() == "text" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}

	logger = logger.With(slog.String("log_type", "app"))
	slog.SetDefault(logger)
	// fmt.Println("Using config file:", viper.ConfigFileUsed())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
