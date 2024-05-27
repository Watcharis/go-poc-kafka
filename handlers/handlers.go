package handlers

import (
	"context"

	"github.com/IBM/sarama"
)

type KafkaConsumerHandlers interface {
	ProcessorKafkaPocTopicFirstHanlder(ctx context.Context, message *sarama.ConsumerMessage) error
}
