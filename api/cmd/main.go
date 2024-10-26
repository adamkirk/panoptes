package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	apicmd "github.com/adamkirk/panoptes/cmd/api"
	superuserscreate "github.com/adamkirk/panoptes/cmd/superusers_create"
	tokensgenerate "github.com/adamkirk/panoptes/cmd/tokens_generate"
	"github.com/adamkirk/panoptes/internal/api"
	v1 "github.com/adamkirk/panoptes/internal/api/v1"
	"github.com/adamkirk/panoptes/internal/config"
	"github.com/adamkirk/panoptes/internal/domain/ingestion"
	"github.com/adamkirk/panoptes/internal/domain/users"
	"github.com/adamkirk/panoptes/internal/domain/validation"
	"github.com/adamkirk/panoptes/internal/repository/postgres"
	"github.com/adamkirk/panoptes/internal/util/encryption"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var cfgFile string
var appCfg *config.Config

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Panoptes",
	Long:  `Generates metrics about development processes by combining info from task management and VCS services.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var apiServeCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the API server",
	Run: func(cmd *cobra.Command, args []string) {
		apicmd.Handler(SharedOpts(appCfg), cmd, args)
	},
}

var superusersCmd = &cobra.Command{
	Use:   "superusers",
	Short: "Commands for managing super users.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var superusersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Commands for creating super users.",
	Run: func(cmd *cobra.Command, args []string) {
		superuserscreate.Handler(SharedOpts(appCfg), cmd, args)
	},
}

var tokensCmd = &cobra.Command{
	Use:   "tokens",
	Short: "Commands for managing access tokens.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var tokensGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates an access token",
	Run: func(cmd *cobra.Command, args []string) {
		tokensgenerate.Handler(SharedOpts(appCfg), cmd, args)
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
				v1.NewUsersController,
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

		fx.Provide(
			fx.Annotate(
				func (* config.Config) *encryption.Bcrypter {
					return encryption.NewBcrypter(encryption.WithCost(cfg.Auth.Bcrypt.Cost))
				},
				fx.As(new(users.Encrypter)),
				fx.As(new(api.TokenVerifier)),
			),
		),

		fx.Provide(
			fx.Annotate(
				users.NewAccessTokensService,
				fx.As(new(tokensgenerate.TokensService)),
			),
		),

		fx.Provide(
			fx.Annotate(
				users.NewUsersService,
				fx.As(new(superuserscreate.UsersService)),
				fx.As(new(v1.UsersService)),
			),
		),

		fx.Provide(validation.NewValidator),
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
			fx.Provide(
				fx.Annotate(
					postgres.NewUsersRepository,
					fx.As(new(users.UsersRepo)),
				),
			),
			fx.Provide(
				fx.Annotate(
					postgres.NewRolesRepository,
					fx.As(new(users.RolesRepo)),
				),
			),
			fx.Provide(
				fx.Annotate(
					postgres.NewUserAccessTokensRepository,
					fx.As(new(users.AccessTokensRepo)),
					fx.As(new(api.AuthRepo)),
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

	tokensGenerateCmd.Flags().StringP("user", "u", "", "The ID of the user the token belongs to. The token will have the same permissions as the given user.")
	tokensGenerateCmd.Flags().Int("expire", 6*30, "Days that the token to be valid for. -1 makes it valid forever.")

	superusersCreateCmd.Flags().StringP("email", "e", "", "Users email address")
	superusersCreateCmd.Flags().StringP("first-name", "f", "", "Users first name")
	superusersCreateCmd.Flags().StringP("last-name", "l", "", "Users last name")
	superusersCreateCmd.Flags().StringP("password", "p", "", "Users password")
	superusersCreateCmd.MarkFlagRequired("email")
	superusersCreateCmd.MarkFlagRequired("first-name")
	superusersCreateCmd.MarkFlagRequired("last-name")
	superusersCreateCmd.MarkFlagRequired("password")

	rootCmd.AddCommand(apiServeCmd)
	rootCmd.AddCommand(tokensCmd)
	tokensCmd.AddCommand(tokensGenerateCmd)

	rootCmd.AddCommand(superusersCmd)
	superusersCmd.AddCommand(superusersCreateCmd)

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

	viper.SetEnvPrefix("panoptes")
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
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
