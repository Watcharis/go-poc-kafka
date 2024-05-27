package consumer

import (
	"errors"
	"fmt"
	"log"
	"watcharis/go-poc-kafka/handlers"

	"github.com/IBM/sarama"
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

func NewCosumerGroupHabler(ready chan bool) *ConsumerGroup {
	return &ConsumerGroup{
		ready: ready,
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
				log.Printf("message channel was closed")
				return nil
			}

			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s partition = %d", string(message.Value), message.Timestamp, message.Topic, message.Partition)

			switch topic := message.Topic; topic {
			case KAFKA_TOPIC_POC_1:
				if err := cgh.processorKafkaTopicHanlers.ProcessorKafkaPocTopicFirstHanlder(session.Context(), message); err != nil {
					log.Printf("[ERROR] ProcessorKafkaPocTopicFirstHanlder - handler error : %+v\n", err)
					return err
				}

				session.MarkMessage(message, "")

			case KAFKA_TOPIC_POC_2:
				fmt.Println("service topic :", message.Topic)

			default:
				errMsg := errors.New(fmt.Sprintf("consume message not found topic : %s", message.Topic))
				log.Printf("not service process for topic: %s, err : %s", message.Topic, errMsg)
				return errMsg
			}

		case <-session.Context().Done():
			return nil
		}
	}
}
