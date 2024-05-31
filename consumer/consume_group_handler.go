package consumer

import (
	"fmt"
	"watcharis/go-poc-kafka/handlers"
	"watcharis/go-poc-kafka/logger"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type ConsumerGroupHandler interface {
	// Setup is run at the beginning of a new session, before ConsumeClaim.
	Setup(sarama.ConsumerGroupSession) error

	// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
	// but before the offsets are committed for the very last time.
	Cleanup(sarama.ConsumerGroupSession) error

	// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
	// Once the Messages() channel is closed, the Handler must finish its processing
	// loop and exit.
	ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error
}

type ConsumerGroup struct {
	ready                      chan bool
	processorKafkaTopicHanlers handlers.KafkaConsumerHandlers
}

func NewCosumerGroupHabler(ready chan bool, processorKafkaTopicHanlers handlers.KafkaConsumerHandlers) *ConsumerGroup {
	return &ConsumerGroup{
		ready:                      ready,
		processorKafkaTopicHanlers: processorKafkaTopicHanlers,
	}
}

func (cgh *ConsumerGroup) Setup(sarama.ConsumerGroupSession) error {
	close(cgh.ready)
	return nil
}

func (cgh *ConsumerGroup) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (cgh *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():

			if !ok {
				logger.Info("message channel was closed")
				return nil
			}

			logger.Info("claimed message", zap.String("value", string(message.Value)), zap.Time("timestamp", message.Timestamp),
				zap.String("topic", message.Topic), zap.Int("partition", int(message.Partition)))

			switch topic := message.Topic; topic {
			case KAFKA_TOPIC_POC_1:
				if err := cgh.processorKafkaTopicHanlers.ProcessorKafkaPocTopicFirstHanlder(session.Context(), message); err != nil {
					logger.Error("ProcessorKafkaPocTopicFirstHanlder - handler error", zap.Error(err))
					return err
				}

				session.MarkMessage(message, "")
			case KAFKA_TOPIC_POC_2:
				fmt.Println("service topic :", message.Topic)

				session.MarkMessage(message, "")

			default:
				err := fmt.Errorf("consume message not found topic : %s", message.Topic)
				logger.Error(fmt.Sprintf("not service process for topic: %s", message.Topic), zap.Error(err))
				return err
			}

		case <-session.Context().Done():
			return nil
		}
	}
}
