package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

// Producer provides access to the kafka writer.
type Producer struct {
	writer *kafka.Writer
	notify chan error
}

// NewAsyncProducer returns new instance of Producer which configured to be used asynchronously.
func NewAsyncProducer(address []string, setters ...ProducerOption) *Producer {
	notify := make(chan error)

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(address...),
		Balancer:               &kafka.Hash{},
		AllowAutoTopicCreation: false,
		RequiredAcks:           kafka.RequireAll,
		Async:                  true,
		Completion: func(messages []kafka.Message, err error) {
			if err != nil {
				notify <- err
			}
		},
	}

	producer := &Producer{
		writer: writer,
		notify: notify,
	}

	for _, setter := range setters {
		setter(producer)
	}

	return producer
}

// ProduceAsync writes message to kafka asynchronously.
func (p *Producer) ProduceAsync(topic string, key, message []byte) {
	kafkaMessage := kafka.Message{
		Topic: topic,
		Key:   key,
		Value: message,
	}

	_ = p.writer.WriteMessages(context.Background(), kafkaMessage)
}

// Notify - notifies about kafka writer errors.
func (p *Producer) Notify() <-chan error {
	return p.notify
}

// Close - closes connection to kafka writer.
func (p *Producer) Close() error {
	defer close(p.notify)

	if err := p.writer.Close(); err != nil {
		return fmt.Errorf("%w: %w", ErrKafkaWriterClose, err)
	}

	return nil
}
