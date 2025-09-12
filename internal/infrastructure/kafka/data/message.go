package data

type KafkaMessage struct {
	Topic string
	Key   string
	Value []byte
}
