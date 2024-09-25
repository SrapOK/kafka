package client

import (
	"strings"

	"github.com/segmentio/kafka-go"
)

func GetKafkaReader(kafkaURL, topic string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
	})
}
