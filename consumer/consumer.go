package consumer

import (
	"context"
	"flag"
	"log"
	"sync"
	"watcharis/go-poc-kafka/handlers"

	"github.com/IBM/sarama"
)

const (
	KAFKA_ADDRESS     = "localhost:9092"
	GROUP_KAFKA_POC_1 = "kafka-poc-group-1"
	KAFKA_TOPIC_POC_1 = "kafka-poc-topic-1"
	KAFKA_TOPIC_POC_2 = "kafka-poc-topic-2"
)

var (
	version  = ""
	assignor = ""
)

func init() {
	flag.StringVar(&assignor, "assignor", "roundrobin", "Consumer group partition assignment strategy (range, roundrobin, sticky)")
	flag.StringVar(&version, "version", sarama.DefaultVersion.String(), "Kafka cluster version")
}

func NewKafkaConsumerGroup(ctx context.Context, processorKafkaTopicHanlers handlers.KafkaConsumerHandlers) error {

	version, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}
	log.Println("version :", version)

	config := sarama.NewConfig()
	config.Version = version

	switch assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", assignor)
	}
	log.Println("assignor :", assignor)

	addr := []string{KAFKA_ADDRESS}

	group, err := sarama.NewConsumerGroup(addr, GROUP_KAFKA_POC_1, config)
	if err != nil {
		log.Println("[Error] cannot start consumer group err :", err)
		return err
	}

	consumer := ConsumerGroup{
		ready:                      make(chan bool),
		processorKafkaTopicHanlers: processorKafkaTopicHanlers,
	}
	// consumer := NewCosumerGroupHabler(make(chan bool))

	// fmt.Printf("client kafka : %+v\n", client)
	errCH := make(chan error)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func(ctx context.Context, errCH chan error) {
		defer func() {
			wg.Done()
		}()
		for {
			topics := []string{KAFKA_TOPIC_POC_1, KAFKA_TOPIC_POC_2}
			// topics = nil

			err := group.Consume(ctx, topics, &consumer)
			if err != nil {
				log.Println("[ERROR] cannot consuner message err :", err)
				errCH <- err
				return
			}

			if ctx.Err() != nil {
				errCH <- ctx.Err()
				return
			}

			errCH <- nil
			consumer.ready = make(chan bool)
		}
	}(ctx, errCH)

	go func() {
		<-consumer.ready // Await till the consumer has been set up
		log.Println("Sarama consumer up and running!...")
	}()

	if consumerError := <-errCH; consumerError != nil {
		log.Println("consumerError :", consumerError)
		return consumerError
	}

	wg.Wait()
	return nil
}
