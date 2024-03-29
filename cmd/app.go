package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/adrg/xdg"
	"github.com/alecthomas/kong"
	"github.com/lmittmann/tint"

	"github.com/prskr/git-age-keyring-agent/handlers/cli"
)

type App struct {
	Logging struct {
		Level slog.Level `env:"GIT_AGE_LOG_LEVEL" help:"Log level" default:"warn"`
	} `embed:""`

	Serve cli.ServeCliHandler `cmd:"" name:"serve" help:"serve a keyring agent server"`
}

func (a *App) Execute() error {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	cliCtx := kong.Parse(a,
		kong.Name("git-age-keyring-agent"),
		kong.BindTo(ctx, (*context.Context)(nil)),
		kong.Vars{
			"XDG_RUNTIME_DIR": xdg.RuntimeDir,
		},
	)

	return cliCtx.Run()
}

func (a *App) AfterApply() error {
	opts := &tint.Options{
		Level:      a.Logging.Level,
		TimeFormat: time.RFC3339,
	}
	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, opts)))

	return nil
}
