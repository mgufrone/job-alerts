package main

import (
	"context"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"go.uber.org/fx"
	"mgufrone.dev/job-alerts/delivery/cli/entrypoints"
	"mgufrone.dev/job-alerts/delivery/cli/handlers"
	"mgufrone.dev/job-alerts/internal/bootstrap"
	"mgufrone.dev/job-alerts/internal/usecases/notification"
	_ "mgufrone.dev/job-alerts/migrations"

	"os"
)

func client(commands []cli.Command, lc fx.Lifecycle) *cli.App {
	app := &cli.App{
		Commands: commands,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "dry-run",
				Usage: "set operation to dry run (no-op)",
			},
			cli.BoolFlag{
				Name:  "debug",
				Usage: "enable debug mode",
			},
		},
	}
	lc.Append(fx.StartHook(func() error {
		return app.Run(os.Args)
	}))
	return app
}
func main() {
	godotenv.Load()
	app := fx.New(
		fx.NopLogger,
		bootstrap.AppModule,
		notification.Module,
		handlers.Module,
		entrypoints.Module,
		fx.Provide(
			fx.Annotate(
				client,
				fx.ParamTags(`group:"commands"`),
			),
		),
		fx.Invoke(func(*cli.App) {}),
	)
	ctx := context.TODO()
	if err := app.Start(ctx); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	app.Stop(ctx)
}
