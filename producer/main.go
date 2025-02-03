package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"watcharis/go-poc-kafka/config"
	"watcharis/go-poc-kafka/logger"
	"watcharis/go-poc-kafka/models"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// ctx := context.Background()

	// init zap logger
	logger.InitZapLogger()
	defer logger.Sync()

	logger.Info("initial logger success", zap.String("log-status", "ok"))

	cfg := config.LoadConfig("./config")
	logger.Info("app load config success", zap.Any("cfg", cfg))

	saramaProducer := InitProducer(cfg)
	defer func() {
		if err := saramaProducer.Close(); err != nil {
			logger.Fatal(err.Error())
		}
	}()

	filename := "./produce-message.json"
	f, err := os.Open(filename)
	if err != nil {
		logger.Error("cannot open file", zap.Error(err))
		return
	}

	responseFile, err := io.ReadAll(f)
	if err != nil {
		logger.Error("cannot io.ReadAll file", zap.Error(err))
		return
	}

	var result []models.MessageKafkaPocTopicFirst
	if err := json.Unmarshal(responseFile, &result); err != nil {
		logger.Error("cannot json.Unmarshal", zap.Error(err))
		return
	}

	for _, v := range result {
		var topic string
		if v.Title == cfg.Kafka.Topic1 {
			topic = cfg.Kafka.Topic1
		} else {
			topic = cfg.Kafka.Topic2
		}

		messageByte, err := json.Marshal(v)
		if err != nil {
			logger.Error("cannot json.Marshal", zap.Error(err))
			return
		}

		if err := Producer(saramaProducer, topic, string(messageByte)); err != nil {
			logger.Error("produces failed", zap.Error(err))
			return
		}
	}

}

func InitProducer(cfg *config.Config) sarama.SyncProducer {
	producer, err := sarama.NewSyncProducer(cfg.Secerts.KafkaAddress, nil)
	if err != nil {
		logger.Fatal(err.Error())
	}
	return producer
}

func Producer(saramaProducer sarama.SyncProducer, topic string, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := saramaProducer.SendMessage(msg)
	if err != nil {
		logger.Error("FAILED to send message", zap.String("topic", topic), zap.Error(err))
		return err
	}

	logger.Infof(fmt.Sprintf("> message sent to partition %d at offset %d\n", partition, offset), nil, []zapcore.Field{
		zap.String("topic", topic),
	})
	return nil
}
