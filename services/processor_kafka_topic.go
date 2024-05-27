package services

import (
	"context"
	"watcharis/go-poc-kafka/models"
)

type processorKafkaTopic struct{}

func NewProcessorKafkaTopic() ProcessorKafkaTopic {
	return &processorKafkaTopic{}
}

func (s *processorKafkaTopic) ProcessorKafkaPocTopicFirst(ctx context.Context, messageTopicFirst models.MessageKafkaPocTopicFirst) error {

	// want to insert data to elastic search

	return nil
}
