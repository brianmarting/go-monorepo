package main

import (
	"context"
	"go-monorepo/internal-backend/service"
	"go-monorepo/internal/logging"
	"go-monorepo/internal/observability/tracing"
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

	mineralService := service.NewMineralService()
	if err := mineralService.StartConsumingMessages(); err != nil {
		log.Fatal().Err(err).Msg("failed to start consuming messages from broker")
	}

	<-ctx.Done()
}
