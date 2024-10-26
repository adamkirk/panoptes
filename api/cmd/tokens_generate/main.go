package tokensgenerate

import (
	"context"

	"github.com/adamkirk/panoptes/internal/domain/users"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

type TokensService interface {
	Create(def users.TokenDefinition) (*users.AccessToken, error)
}

type Action struct {
	sh       fx.Shutdowner
	cmd      *cobra.Command
	svc TokensService
	args     []string
}

type actionInput struct {
	cmd  *cobra.Command
	args []string
}

func newAction(
	lc fx.Lifecycle,
	sh fx.Shutdowner,
	svc TokensService,
	input *actionInput,
) *Action {
	act := &Action{
		sh:       sh,
		cmd:      input.cmd,
		svc: svc,
		args:     input.args,
	}

	lc.Append(fx.Hook{
		OnStart: act.start,
		OnStop:  act.stop,
	})

	return act
}

func (act *Action) start(ctx context.Context) error {
	go act.run()
	return nil
}

func (act *Action) stop(ctx context.Context) error {
	return nil
}

func (act *Action) run() {
	user, err := act.cmd.Flags().GetString("user")

	if err != nil {
		color.Red("Failed to get user option: %s", err.Error())
		act.sh.Shutdown(fx.ExitCode(1))
		return
	}

	expire, err := act.cmd.Flags().GetInt("expire")

	if err != nil {
		color.Red("Failed to get expiry option: %s", err.Error())
		act.sh.Shutdown(fx.ExitCode(1))
		return
	}

	if expire < -1 || expire == 0 {
		color.Red("Invalid expire option, must '-1' or greater than 0")
		act.sh.Shutdown(fx.ExitCode(1))
		return
	}

	if expire == -1 {
		color.Cyan("This token will never expire!")
	}

	id, err := uuid.Parse(user)

	if err != nil {
		color.Red("The user option must be a valid uuid")
		act.sh.Shutdown(fx.ExitCode(1))
		return
	}

	token, err := act.svc.Create(users.TokenDefinition{
		ExpiryDays: expire,
		UserID: id,
	})

	if err != nil {
		color.Red("Failed to create token: %s", err.Error())
		act.sh.Shutdown(fx.ExitCode(1))
		return
	}

	color.Cyan("ID: '%s'", token.ID)
	color.Cyan("Raw token: '%s'", *token.Secret)
	color.Cyan("Hash: '%s'", token.SecretHash)

	act.sh.Shutdown()
}

func Handler(opts []fx.Option, cmd *cobra.Command, args []string) {
	opts = append(opts, []fx.Option{
		// Prevents all the logging noise when building the service container
		fx.NopLogger,
		fx.Provide(func() *actionInput {
			return &actionInput{
				cmd:  cmd,
				args: args,
			}
		}),
		fx.Provide(newAction),
		fx.Invoke(func(*Action) {}),
	}...)

	fx.New(
		opts...,
	).Run()
}
