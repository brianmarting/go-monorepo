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

	doneChan <-chan struct{}
	tracer   trace.Tracer
	service  service.MineralService
}

func NewMineralGrpcServer(done <-chan struct{}) *MineralServer {
	return &MineralServer{
		doneChan: done,
		tracer:   tracing.GetTracer(),
		service:  service.GetMineralService(),
	}
}

func (m *MineralServer) SendStreaming(stream pb.MineralService_SendStreamingServer) error {
	msgsChan := make(chan *pb.MineralDto)

	go func() {
		for {
			recv, err := stream.Recv()
			if err != nil {
				log.Error().Msg("failed to receive msg from grpc stream")
				return
			}

			msgsChan <- recv
		}
	}()

	for {
		select {
		case <-m.doneChan:
			return nil
		case msg := <-msgsChan:
			if err := handleStreamMsg(msg, m.service); err != nil {
				log.Error().Err(err)
			}
		}
	}
}

func handleStreamMsg(recv *pb.MineralDto, service service.MineralService) error {
	tracer := tracing.GetTracer()
	spanCtx, span := tracer.Start(context.Background(), "receive-mineral-msg-grpc-stream")
	defer span.End()

	if err := service.AddMineral(spanCtx, convertToMineral(recv)); err != nil {
		log.Error().Err(err).Msg("failed to receive data from grpc endpoint")
	}

	return nil
}

func convertToMineral(dto *pb.MineralDto) model.Mineral {
	return model.Mineral{
		Name:   dto.Name,
		Amount: int(dto.Amount),
	}
}
