package kafkaapp

import (
	"fmt"
	"github.com/p1xray/pxr-sso/internal/config"
	"github.com/p1xray/pxr-sso/internal/infrastructure/kafka/data"
	"github.com/p1xray/pxr-sso/pkg/kafka"
	"github.com/p1xray/pxr-sso/pkg/logger/sl"
	"hash/fnv"
	"log/slog"
	"strings"
)

// App is a kafka queue application.
type App struct {
	log            *slog.Logger
	producer       *kafka.Producer
	numberOfTopics int
	input          chan data.KafkaMessage
}

func New(log *slog.Logger, cfg config.KafkaConfig) *App {
	address := strings.Split(cfg.Address, ",")

	producer := kafka.NewAsyncProducer(address)

	return &App{
		log:            log,
		producer:       producer,
		numberOfTopics: cfg.NumberOfTopics,
		input:          make(chan data.KafkaMessage),
	}
}

func (a *App) Start() {
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

				a.producer.ProduceAsync(d.Topic, key, d.Value)
			}
		}()
	}

	a.handleAsyncProducerErrors()
}

func (a *App) Input() chan<- data.KafkaMessage {
	return a.input
}

func (a *App) Stop() {
	const op = "kafkaapp.Stop"

	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("stopping kafka producers")

	defer close(a.input)

	if err := a.producer.Close(); err != nil {
		log.Error("error closing kafka producer", sl.Err(err))
	}
}

func (a *App) fanOut() []<-chan data.KafkaMessage {
	const op = "kafkaapp.fanOut"

	log := a.log.With(slog.String("op", op))

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
				log.Error("error calculation index of topic", sl.Err(err))
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

func (a *App) handleAsyncProducerErrors() {
	const op = "kafkaapp.handleAsyncProducerErrors"

	log := a.log.With(slog.String("op", op))

	go func() {
		err := <-a.producer.Notify()
		if err != nil {
			log.Error("error writing message to kafka", sl.Err(err))
		}
	}()
}
