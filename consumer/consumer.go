package consumer

import (
	"context"
	"flag"
	"log"
	"sync"
	"watcharis/go-poc-kafka/config"
	"watcharis/go-poc-kafka/handlers"
	"watcharis/go-poc-kafka/logger"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

var (
	version  = ""
	assignor = ""
)

func init() {
	flag.StringVar(&assignor, "assignor", "roundrobin", "Consumer group partition assignment strategy (range, roundrobin, sticky)")
	flag.StringVar(&version, "version", sarama.DefaultVersion.String(), "Kafka cluster version")
}

func NewKafkaConsumerGroup(ctx context.Context, processorKafkaTopicHanlers handlers.KafkaConsumerHandlers, cfg config.Config) error {

	version, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}
	logger.Info("parsing kafka version", zap.String("version", version.String()))

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
	logger.Info("consumer-group reblance strategy", zap.String("assignor", assignor))

	// addr := []string{KAFKA_ADDRESS}
	addr := cfg.Secerts.KafkaAddress

	group, err := sarama.NewConsumerGroup(addr, cfg.Kafka.ConsumerGroup, config)
	if err != nil {
		log.Println("[Error] cannot start consumer group err :", err)
		return err
	}

	consumer := ConsumerGroup{
		ready:                      make(chan bool),
		processorKafkaTopicHanlers: processorKafkaTopicHanlers,
		cfg:                        cfg,
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
			// topics = nil
			topics := []string{cfg.Kafka.Topic1, cfg.Kafka.Topic2}

			err := group.Consume(ctx, topics, &consumer)
			if err != nil {
				logger.Error("cannot consume message", zap.Error(err))
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
		logger.Info("Sarama consumer up and running!...")
	}()

	if consumerError := <-errCH; consumerError != nil {
		logger.Info("consumer error", zap.Error(consumerError))
		return consumerError
	}

	wg.Wait()
	return nil
}
