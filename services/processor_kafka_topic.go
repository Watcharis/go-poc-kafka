package services

import (
	"context"
	"watcharis/go-poc-kafka/logger"
	"watcharis/go-poc-kafka/models"

	"go.uber.org/zap"
)

type processorKafkaTopic struct{}

func NewProcessorKafkaTopic() ProcessorKafkaTopic {
	return &processorKafkaTopic{}
}

func (s *processorKafkaTopic) ProcessorKafkaPocTopicFirst(ctx context.Context, messageTopicFirst models.MessageKafkaPocTopicFirst) error {

	logger.Info("ProcessorKafkaPocTopicFirst - start service", zap.Any("message", messageTopicFirst))
	// want to insert data to elastic search

	return nil
}
