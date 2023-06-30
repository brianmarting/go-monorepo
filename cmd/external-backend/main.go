package main

import (
	"context"
	"fmt"
	"go-monorepo/external-backend/grpc"
	pb "go-monorepo/external-backend/grpc/generated"
	"go-monorepo/external-backend/rest/handlers"
	"go-monorepo/internal/logging"
	"go-monorepo/internal/observability/tracing"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/rs/zerolog/log"
	googleGrpc "google.golang.org/grpc"
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

	go startRestServer()
	go startGrpcServer()

	<-ctx.Done()
}

func startRestServer() {
	h := handlers.NewHandler()
	h.CreateAllRoutes()

	log.Info().Msg("starting REST server")

	if err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("REST_PORT")), h); err != nil {
		log.Fatal().Err(err).Msg("failed to start http listener ")
	}
}

func startGrpcServer() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("GRPC_PORT")))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start grpc server listener")
	}

	server := googleGrpc.NewServer()
	pb.RegisterMineralServiceServer(server, grpc.NewMineralGrpcServer())

	log.Info().Msg("starting GRPC server")
	if err := server.Serve(listener); err != nil {
		log.Fatal().Err(err).Msg("failed to start GRPC server")
	}
}
