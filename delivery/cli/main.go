package main

import (
	"context"
	"github.com/alecthomas/kong"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"mgufrone.dev/job-alerts/delivery/cli/app"
	"mgufrone.dev/job-alerts/delivery/cli/entrypoints"
	"mgufrone.dev/job-alerts/delivery/cli/handlers"
	"mgufrone.dev/job-alerts/internal/bootstrap"
	"mgufrone.dev/job-alerts/internal/usecases/notification"
	"os"
)

func main() {
	godotenv.Load()
	app := fx.New(
		bootstrap.AppModule,
		notification.Module,
		handlers.Module,
		fx.Provide(
			app.NewKernel,
			func() *kong.Context {
				return kong.Parse(&entrypoints.CLI{})
			},
		),
		fx.Invoke(func(ctx *kong.Context, kernel *app.Kernel) error {
			err := ctx.Run(kernel)
			return err
		}),
	)
	ctx := context.TODO()
	if err := app.Start(ctx); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	app.Stop(ctx)
}
