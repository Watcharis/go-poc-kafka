package services

import (
	"context"
	"watcharis/go-poc-kafka/logger"
	"watcharis/go-poc-kafka/models"
	"watcharis/go-poc-kafka/util"

	"go.uber.org/zap"
)

type processorKafkaTopic struct{}

func NewProcessorKafkaTopic() ProcessorKafkaTopic {
	return &processorKafkaTopic{}
}

func (s *processorKafkaTopic) ProcessorKafkaPocTopicFirst(ctx context.Context, messageTopicFirst models.MessageKafkaPocTopicFirst) error {

	messageJson, err := util.JSONString(messageTopicFirst)
	if err != nil {
		logger.Error("ProcessorKafkaPocTopicFirst - cannot convert message to json string", zap.Error(err))
		return err
	}

	logger.Info("ProcessorKafkaPocTopicFirst - start service", zap.String("message", messageJson))
	// want to insert data to elastic search

	return nil
}
