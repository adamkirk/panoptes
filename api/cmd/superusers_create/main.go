package superuserscreate

import (
	"context"

	"github.com/adamkirk/panoptes/internal/domain/users"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

type UsersService interface {
	Create(def users.CreateDTO) (*users.User, error)
}

type Action struct {
	sh       fx.Shutdowner
	cmd      *cobra.Command
	svc UsersService
	args     []string
}

type actionInput struct {
	cmd  *cobra.Command
	args []string
}

func newAction(
	lc fx.Lifecycle,
	sh fx.Shutdowner,
	svc UsersService,
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
	password, err := act.cmd.Flags().GetString("password")

	if err != nil {
		color.Red("Failed to get password option: %s", err.Error())
		act.sh.Shutdown(fx.ExitCode(1))
		return
	}

	firstName, err := act.cmd.Flags().GetString("first-name")

	if err != nil {
		color.Red("Failed to get first name option: %s", err.Error())
		act.sh.Shutdown(fx.ExitCode(1))
		return
	}

	lastName, err := act.cmd.Flags().GetString("last-name")

	if err != nil {
		color.Red("Failed to get last name option: %s", err.Error())
		act.sh.Shutdown(fx.ExitCode(1))
		return
	}

	email, err := act.cmd.Flags().GetString("email")

	if err != nil {
		color.Red("Failed to get email option: %s", err.Error())
		act.sh.Shutdown(fx.ExitCode(1))
		return
	}
	
	user, err := act.svc.Create(users.CreateDTO{
		FirstName: firstName,
		LastName: lastName,
		Email: email,
		Password: password,
		Roles: []string{"superuser"},
	})

	if err != nil {
		color.Red("Failed to create token: %s", err.Error())
		act.sh.Shutdown(fx.ExitCode(1))
		return
	}

	color.Cyan("ID: '%s'", user.ID)
	color.Cyan("Email: '%s'", user.Email)
	color.Cyan("FirstName: '%s'", user.FirstName)
	color.Cyan("LastName: '%s'", user.LastName)

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
