package gokafka

import "github.com/confluentinc/confluent-kafka-go/kafka"

func FindHeaders(key string, headers []kafka.Header) string {
	for _, v := range headers {
		if v.Key == key {
			return string(v.Value)
		}
	}

	return ""
}
