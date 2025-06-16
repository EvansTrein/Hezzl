package event

import (
	"encoding/json"
	"hezzl/internal/model"
	"hezzl/pkg/broker/nats"
	"log/slog"
)

type ILogsRepo interface {
	Create(data *model.Product)
}

type logging struct {
	log    *slog.Logger
	Broker *nats.NatsBroker
	repo   ILogsRepo
}

type LoggingDeps struct {
	*slog.Logger
	*nats.NatsBroker
	ILogsRepo
}

func NewLogging(deps *LoggingDeps) *logging {
	return &logging{
		log:    deps.Logger,
		Broker: deps.NatsBroker,
		repo:   deps.ILogsRepo,
	}
}

func (e *logging) SendToBroker(data *model.Product) {
	op := "event logging: send to broker"
	log := e.log.With(slog.String("operation", op))
	log.Debug("Call func SendToBroker", "data", data)

	mess, err := json.Marshal(data)
	if err != nil {
		log.Error("failed to marshal json", "error", err)
		return
	}

	if err := e.Broker.Conn.Publish(e.Broker.NameMess, mess); err != nil {
		log.Error("failed to publish message to nats", "error", err)
		return
	}

	log.Info("successfully sent to broker")
}

func (e *logging) SendLogToDB(data []byte) {
	op := "event logging: send log to repository"
	log := e.log.With(slog.String("operation", op))
	log.Debug("Call func SendLogToDB", "data", data)

	var product model.Product
	if err := json.Unmarshal(data, &product); err != nil {
		e.log.Error("failed to unmarshal message", "error", err)
		return
	}

	e.repo.Create(&product)

	log.Info("successfully sent to repo")
}
