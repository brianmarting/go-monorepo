package service

import (
	"context"
	"encoding/json"
	"go-monorepo/internal-backend/persistence/db"
	"go-monorepo/internal-backend/persistence/db/psql"
	"go-monorepo/internal-backend/queue"
	"go-monorepo/internal/model"
	"go-monorepo/internal/observability/tracing"
	commonQueue "go-monorepo/internal/queue"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
)

type MineralService interface {
	StartConsumingMessages() error
}

type mineralService struct {
	tracer       trace.Tracer
	mineralStore db.MineralStore
}

func NewMineralService() MineralService {
	return &mineralService{
		tracer:       tracing.GetTracer(),
		mineralStore: psql.NewMineralStore(),
	}
}

func (m mineralService) StartConsumingMessages() error {
	consumer := queue.NewMineralConsumer()
	msgs, err := consumer.StartConsuming("mineral-deposit-queue", "mineral.deposit")
	if err != nil {
		return err
	}

	for msg := range msgs {
		func(msg commonQueue.Message) {
			_, span := m.tracer.Start(context.Background(), "mineral-service-consume-messages")
			//_, span := m.tracer.Start(createContextFromMsgHeaders(msg), "mineral-service-consume-messages") // TODO check why not working
			defer span.End()

			mineral := &model.Mineral{}

			err := json.Unmarshal(msg.GetBytes(), mineral)
			if err != nil {
				log.Error().Err(err).Msg("failed to parse message body to mineral")
				nackMsg(msg)
				return
			}

			dbMineral, err := m.mineralStore.GetByName(mineral.Name)
			if err != nil {
				log.Error().Err(err).Msg("failed to fetch mineral from DB")
				nackMsg(msg)
				return
			}

			amount := dbMineral.Amount + mineral.Amount
			if err := m.mineralStore.UpdateAmount(dbMineral.Name, amount); err != nil {
				log.Error().Err(err).Msg("failed update amount of mineral")
				nackMsg(msg)
				return
			}

			if err := msg.Ack(); err != nil {
				log.Error().Err(err).Msg("failed to ack msg")
			}
		}(msg)
	}

	return nil
}

func nackMsg(msg commonQueue.Message) {
	if err := msg.Nack(); err != nil {
		log.Error().Err(err).Msg("failed to nack msg")
	}
}
