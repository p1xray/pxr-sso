package kafka

import "github.com/segmentio/kafka-go"

// ProducerOption is how options for the Producer are set up.
type ProducerOption func(*Producer)

// AcksRequireNone sets up a number of acknowledges from partition replicas required before receiving
// a response to a produce request to "none" (fire-and-forget, do not wait for acknowledgements from the).
func AcksRequireNone() ProducerOption {
	return func(p *Producer) {
		p.writer.RequiredAcks = kafka.RequireNone
	}
}

// AcksRequireLeader sets up a number of acknowledges from partition replicas required before receiving
// a response to a produce request to "leader" (wait for the leader to acknowledge the writes).
func AcksRequireLeader() ProducerOption {
	return func(p *Producer) {
		p.writer.RequiredAcks = kafka.RequireOne
	}
}

// AcksRequireAll sets up a number of acknowledges from partition replicas required before receiving
// a response to a produce request to "all" (wait for the full ISR to acknowledge the writes).
func AcksRequireAll() ProducerOption {
	return func(p *Producer) {
		p.writer.RequiredAcks = kafka.RequireAll
	}
}
