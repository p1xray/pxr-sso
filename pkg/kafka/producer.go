package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
	notify chan error
}

func NewAsyncProducer(address []string) *Producer {
	notify := make(chan error)

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(address...),
		Balancer:               &kafka.Hash{},
		AllowAutoTopicCreation: false,
		Async:                  true,
		Completion: func(messages []kafka.Message, err error) {
			if err != nil {
				notify <- err
			}
		},
	}

	return &Producer{
		writer: writer,
		notify: notify,
	}
}

func (p *Producer) ProduceAsync(topic string, key, message []byte) {
	kafkaMessage := kafka.Message{
		Topic: topic,
		Key:   key,
		Value: message,
	}

	_ = p.writer.WriteMessages(context.Background(), kafkaMessage)
}

func (p *Producer) Notify() <-chan error {
	return p.notify
}

func (p *Producer) Close() error {
	defer close(p.notify)

	if err := p.writer.Close(); err != nil {
		return fmt.Errorf("%w: %w", ErrKafkaWriterClose, err)
	}

	return nil
}
