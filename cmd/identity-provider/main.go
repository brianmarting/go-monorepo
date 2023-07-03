package main

import (
	"context"
	"fmt"
	"go-monorepo/identity-provider/rest/handlers"
	"go-monorepo/internal/logging"
	"go-monorepo/internal/observability/tracing"
	"net/http"
	"os"
	"os/signal"

	"github.com/rs/zerolog/log"
)

func main() {
	logging.EnableLogging()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	tracer := tracing.InitTracerProvider()
	defer func() {
		cancel()
		if err := tracer.Shutdown(context.Background()); err != nil {
			log.Info().Err(err).Msg("failed to shut down tracer provider")
		}
	}()

	go func() {
		handler := handlers.NewHandler()
		handler.CreateAllRoutes()

		port := os.Getenv("IPS_PORT")
		log.Info().Msg(fmt.Sprintf("starting idp service on %s", port))

		if err := http.ListenAndServe(fmt.Sprintf(":%s", port), handler); err != nil {
			log.Fatal().Err(err).Msg("failed to start idp service")
		}
	}()

	<-ctx.Done()
}
