package consumers

import (
	"context"
	"encoding/json"
	"watcharis/go-poc-kafka/handlers"
	"watcharis/go-poc-kafka/logger"
	"watcharis/go-poc-kafka/models"
	"watcharis/go-poc-kafka/services"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type kafkaConsumerHandlers struct {
	processorKafkaTopicService services.ProcessorKafkaTopic
}

func NewKafkaConsumerHandlers(processorKafkaTopicService services.ProcessorKafkaTopic) handlers.KafkaConsumerHandlers {
	return &kafkaConsumerHandlers{
		processorKafkaTopicService: processorKafkaTopicService,
	}
}

func (h *kafkaConsumerHandlers) ProcessorKafkaPocTopicFirstHanlder(ctx context.Context, message *sarama.ConsumerMessage) error {
	logger.Info("ProcessorKafkaPocTopicFirstHanlder - start handler")

	var messageTopicFirst models.MessageKafkaPocTopicFirst
	if err := json.Unmarshal(message.Value, &messageTopicFirst); err != nil {
		logger.Errorf("ProcessorKafkaPocTopicFirst - cannot json.Unmarshal message kafka topic = %s, meesgae = %+v",
			logger.SugerArgs(message.Topic, string(message.Value)),
			logger.ZapCoreFileds(zap.Error(err)))
		return err
	}

	if err := h.processorKafkaTopicService.ProcessorKafkaPocTopicFirst(ctx, messageTopicFirst); err != nil {
		logger.Error("ProcessorKafkaPocTopicFirst - service error", zap.Error(err))
		return err
	}

	return nil
}
