package kafka

const topicNameSeparator = "-"

func GenerateTopicNameByClientCode(clientCode, event string) string {
	return clientCode + topicNameSeparator + event
}
