package kafkaconsumerhandlers

import (
	"context"
	"encoding/json"
	"log"
	"watcharis/go-poc-kafka/handlers"
	"watcharis/go-poc-kafka/models"
	"watcharis/go-poc-kafka/services"

	"github.com/IBM/sarama"
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

	var messageTopicFirst models.MessageKafkaPocTopicFirst
	if err := json.Unmarshal(message.Value, &messageTopicFirst); err != nil {
		log.Printf("[ERROR] cannot json.Unmarshal message kafka topic = %s, meesgae = %+v, error : %+v\n", message.Topic, string(message.Value), err)
		return err
	}

	if err := h.processorKafkaTopicService.ProcessorKafkaPocTopicFirst(ctx, messageTopicFirst); err != nil {
		log.Printf("[ERROR] ProcessorKafkaPocTopicFirst - service error : %+v\n", err)
		return err
	}

	return nil
}
