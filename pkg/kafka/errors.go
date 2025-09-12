package kafka

import "errors"

var (
	ErrKafkaWriterClose = errors.New("error closing kafka writer")
)
