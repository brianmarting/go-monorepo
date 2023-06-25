package main

import (
	"common/logging"
	"common/observability/tracing"
	"context"
	"external_backend/rest/handlers"
	"net/http"

	"github.com/rs/zerolog/log"
)

func main() {
	logging.EnableLogging()

	tracer := tracing.InitTracerProvider()
	defer func() {
		if err := tracer.Shutdown(context.Background()); err != nil {
			log.Info().Err(err).Msg("failed to shut down tracer provider")
		}
	}()

	h := handlers.NewHandler()
	h.CreateAllRoutes()

	if err := http.ListenAndServe(":8888", h); err != nil {
		log.Info().Err(err).Msg("failed to start http listener ")
	}
}
