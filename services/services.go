package services

import (
	"context"
	"watcharis/go-poc-kafka/models"
)

type ProcessorKafkaTopic interface {
	ProcessorKafkaPocTopicFirst(ctx context.Context, messageTopicFirst models.MessageKafkaPocTopicFirst) error
}
