package main

import (
	"common/logging"
	"common/observability/tracing"
	"context"
	externalbackendcmd "external-backend/cmd"
	internalbackendcmd "internal-backend/cmd"
	"os"
	"os/signal"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
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

	cmd := &cobra.Command{
		Short: "Go monorepo cmd",
	}
	cmd.AddCommand(externalbackendcmd.ExternalBackendCmd(ctx))
	cmd.AddCommand(internalbackendcmd.InternalBackendCmd(ctx))

	if err := cmd.Execute(); err != nil {
		log.Info().Err(err).Msg("failed to start cmd")
	}
}
