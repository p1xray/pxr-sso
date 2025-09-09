package kafkaapp

import (
	"context"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/infrastructure/kafka/data"
	"github.com/p1xray/pxr-sso/pkg/kafka"
	"hash/fnv"
	"log/slog"
)

// App is a kafka queue application.
type App struct {
	log            *slog.Logger
	producer       *kafka.Producer
	numberOfTopics int
	input          chan data.KafkaMessage
	notify         chan error
}

func New(log *slog.Logger, numberOfTopics int) *App {
	address := []string{"localhost:9092"} // TODO: get this from config

	producer := kafka.NewProducer(address)

	return &App{
		log:            log,
		producer:       producer,
		numberOfTopics: numberOfTopics,
		input:          make(chan data.KafkaMessage),
	}
}

func (a *App) Start(ctx context.Context) {
	const op = "kafkaapp.Start"

	log := a.log.With(
		slog.String("op", op),
	)
	log.Info(fmt.Sprintf("running kafka %d producers", a.numberOfTopics))

	dataByTopics := a.fanOut()
	for i := range a.numberOfTopics {
		go func() {
			for d := range dataByTopics[i] {
				var key []byte
				if len(d.Key) > 0 {
					key = []byte(d.Key)
				}

				if err := a.producer.Produce(ctx, d.Topic, key, d.Value); err != nil {
					a.notifyError(err)
				}
			}
		}()
	}
}

func (a *App) Input() chan<- data.KafkaMessage {
	return a.input
}

func (a *App) Notify() <-chan error {
	return a.notify
}

func (a *App) Stop() {
	const op = "kafkaapp.Stop"

	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("stopping kafka producers")

	close(a.input)

	if err := a.producer.Close(); err != nil {
		a.notifyError(err)
	}
}

func (a *App) notifyError(err error) {
	a.notify <- err
}

func (a *App) fanOut() []<-chan data.KafkaMessage {
	out := make([]chan data.KafkaMessage, a.numberOfTopics)
	for i := range a.numberOfTopics {
		out[i] = make(chan data.KafkaMessage)
	}

	go func() {
		defer func() {
			for _, ch := range out {
				close(ch)
			}
		}()

		for receiveData := range a.input {
			idx, err := a.topicIndex(receiveData.Topic)
			if err != nil {
				a.notifyError(err)
				continue
			}
			out[idx] <- receiveData
		}

	}()

	res := make([]<-chan data.KafkaMessage, a.numberOfTopics)
	for i := range a.numberOfTopics {
		res[i] = out[i]
	}

	return res
}

func (a *App) topicIndex(value string) (uint32, error) {
	h := fnv.New32a()
	if _, err := h.Write([]byte(value)); err != nil {
		return 0, err
	}

	idx := h.Sum32() % uint32(a.numberOfTopics)

	return idx, nil
}
