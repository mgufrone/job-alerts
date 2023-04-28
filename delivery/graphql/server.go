package main

import (
	"context"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"log"
	"mgufrone.dev/job-alerts/delivery/graphql/graph/resolvers"
	"mgufrone.dev/job-alerts/internal/bootstrap"
	"mgufrone.dev/job-alerts/internal/usecases/job"
	"mgufrone.dev/job-alerts/internal/usecases/notification"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"mgufrone.dev/job-alerts/delivery/graphql/graph/generated"
)

const defaultPort = "8080"

func main() {
	godotenv.Load()
	app := fx.New(
		bootstrap.AppModule,
		job.Module,
		notification.Module,
		fx.Provide(
			resolvers.NewResolver,
			func(rsv *resolvers.Resolver) *handler.Server {
				return handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: rsv}))
			},
		),
		fx.Invoke(func(srv *handler.Server) {
			port := os.Getenv("PORT")
			if port == "" {
				port = defaultPort
			}
			http.Handle("/", playground.Handler("GraphQL playground", "/query"))
			http.Handle("/query", srv)

			log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
			log.Fatal(http.ListenAndServe(":"+port, nil))
		}),
	)

	if err := app.Start(context.TODO()); err != nil {
		os.Exit(1)
	}
	<-app.Done()
}
