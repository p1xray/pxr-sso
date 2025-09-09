package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/entity"
	"github.com/p1xray/pxr-sso/internal/infrastructure"
	"github.com/p1xray/pxr-sso/internal/infrastructure/converter"
	"github.com/p1xray/pxr-sso/internal/infrastructure/kafka"
	"github.com/p1xray/pxr-sso/internal/infrastructure/kafka/data"
	"github.com/p1xray/pxr-sso/pkg/logger/sl"
	"log/slog"
)

const event = "register"

type UserHasRegistered struct {
	log                *slog.Logger
	receiveDataChannel chan<- data.KafkaMessage
}

func NewUserHasRegistered(log *slog.Logger, receiveDataChannel chan<- data.KafkaMessage) *UserHasRegistered {
	return &UserHasRegistered{
		log:                log,
		receiveDataChannel: receiveDataChannel,
	}
}

func (uhr *UserHasRegistered) SendToKafka(clientCode string, user entity.User) error {
	const op = "handlers.register.SendToKafka"

	log := uhr.log.With(
		slog.String("op", op),
		slog.Int64("user ID", user.ID),
	)

	registeredUserKafkaData := converter.ToRegisteredUserKafkaData(user)

	jsonKafkaDataAsByte, err := json.Marshal(registeredUserKafkaData)
	if err != nil {
		log.Error("error marshaling data for sending to kafka", sl.Err(err))

		return fmt.Errorf("%w: %w", infrastructure.ErrMarshalData, err)
	}

	kafkaMessage := data.KafkaMessage{
		Topic: kafka.GenerateTopicNameByClientCode(clientCode, event),
		Value: jsonKafkaDataAsByte,
	}
	uhr.receiveDataChannel <- kafkaMessage

	return nil
}
