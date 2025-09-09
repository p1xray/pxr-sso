package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(address []string) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(address...),
		Balancer: &kafka.Hash{},
	}

	return &Producer{writer: writer}
}

func (p *Producer) Produce(ctx context.Context, topic string, key, message []byte) error {
	kafkaMessage := kafka.Message{
		Topic: topic,
		Key:   key,
		Value: message,
	}

	err := p.writer.WriteMessages(ctx, kafkaMessage)
	if err != nil {
		return err // TODO: wrap to own error
	}

	return nil
}

func (p *Producer) Close() error {
	if err := p.writer.Close(); err != nil {
		return err // TODO: wrap to own error
	}

	return nil
}
