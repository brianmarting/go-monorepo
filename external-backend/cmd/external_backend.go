package external_backend_cmd

import (
	"context"
	"external-backend/grpc"
	pb "external-backend/grpc/generated"
	"external-backend/rest/handlers"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	googleGrpc "google.golang.org/grpc"
)

func ExternalBackendCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "external-backend",
		Short: "Runs the REST and GRPC server",
		Run: func(cmd *cobra.Command, args []string) {
			go startRestServer()
			go startGrpcServer()

			<-ctx.Done()
		},
	}
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
