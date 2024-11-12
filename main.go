package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"watcharis/go-poc-kafka/config"
	"watcharis/go-poc-kafka/consumer"
	"watcharis/go-poc-kafka/handlers/consumers"
	"watcharis/go-poc-kafka/logger"
	"watcharis/go-poc-kafka/services"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	// init zap logger
	logger.InitZapLogger()
	defer logger.Sync()

	logger.Info("initial logger success", zap.String("log-status", "ok"))

	cfg := config.LoadConfig("./config")
	logger.Info("app load config success", zap.Any("cfg", cfg))

	// init elasticsearch
	// _, err := initElasticsearch()
	// if err != nil {
	// 	log.Panicf("[ERROR] cannot connect elastic err : %+v\n", err)
	// }

	//init kafka services
	processorKafkaTopicService := services.NewProcessorKafkaTopic()

	//init kafka handler
	processorKafkaTopicHandlers := consumers.NewKafkaConsumerHandlers(processorKafkaTopicService)

	//init kafka consumer group
	wg := new(sync.WaitGroup)
	go func(ctx context.Context) {
		logger.Info("start kafka consumer ...")
		if err := consumer.NewKafkaConsumerGroup(ctx, processorKafkaTopicHandlers, *cfg); err != nil {
			log.Panicf("[ERROR] kafka consumer group error : %+v\n", err)
		}
	}(ctx)

	// Channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)
	// Block until a signal is received
	go func() {
		defer wg.Done()
		exit := <-quit
		logger.Info(fmt.Sprintf("exit : %v", exit))
	}()
	wg.Wait()
}

func initElasticsearch() (*elasticsearch.Client, error) {

	address := viper.GetStringSlice("secert.kafka-address")
	fmt.Println("address :", address)

	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
		Username:  "",
		Password:  "",
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err

	}

	res, err := es.Info()
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		log.Panic("invalid status code connect els")
	}

	logger.Info("connect elasticsearch success",
		zap.String("address", "http://localhost:9200"),
		zap.String("username", ""),
		zap.String("password", ""))

	return es, nil
}
