package api

import (
	"context"

	"github.com/adamkirk/heimdallr/internal/api"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func Handler(opts []fx.Option, cmd *cobra.Command, args []string) {
	opts = append(opts, []fx.Option{
		fx.Invoke(startServer),
	}...)

	fx.New(
		opts...,
	).Run()
}

func startServer(lc fx.Lifecycle, srv *api.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go srv.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}
