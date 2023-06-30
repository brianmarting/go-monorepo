package grpc

import (
	"context"
	pb "go-monorepo/external-backend/grpc/generated"
	"go-monorepo/external-backend/service"
	"go-monorepo/internal/model"
	"go-monorepo/internal/observability/tracing"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
)

type MineralServer struct {
	pb.UnimplementedMineralServiceServer

	tracer  trace.Tracer
	service service.MineralService
}

func NewMineralGrpcServer() *MineralServer {
	return &MineralServer{
		tracer:  tracing.GetTracer(),
		service: service.GetMineralService(),
	}
}

func (m MineralServer) SendStreaming(stream pb.MineralService_SendStreamingServer) error {
	for {
		handleStreamMsg(stream, m.service)
	}
}

func handleStreamMsg(stream pb.MineralService_SendStreamingServer, service service.MineralService) {
	recv, err := stream.Recv()

	tracer := tracing.GetTracer()
	spanCtx, span := tracer.Start(context.Background(), "receive-mineral-msg-grpc-stream")
	defer span.End()

	if err != nil {
		log.Error().Err(err).Msg("failed to receive data from grpc endpoint")
		return
	}

	if err = service.AddMineral(spanCtx, convertToMineral(recv)); err != nil {
		log.Error().Err(err).Msg("failed to receive data from grpc endpoint")
		return
	}
}

func convertToMineral(dto *pb.MineralDto) model.Mineral {
	return model.Mineral{
		Name: dto.Name,
	}
}
