package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"watcharis/go-poc-kafka/config"
	"watcharis/go-poc-kafka/consumer"
	"watcharis/go-poc-kafka/handlers/consumers"
	"watcharis/go-poc-kafka/logger"
	"watcharis/go-poc-kafka/services"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/labstack/echo/v4"
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
	_, err := initElasticsearch(*cfg)
	if err != nil {
		log.Panicf("[ERROR] cannot connect elastic err : %+v\n", err)
	}

	e := echo.New()
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Service is running!!!")
	})

	//init kafka services
	processorKafkaTopicService := services.NewProcessorKafkaTopic()

	//init kafka handler
	processorKafkaTopicHandlers := consumers.NewKafkaConsumerHandlers(processorKafkaTopicService)

	//init kafka consumer group
	wg := new(sync.WaitGroup)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	go func(ctx context.Context) {
		logger.Info("start kafka consumer ...")
		if err := consumer.NewKafkaConsumerGroup(ctx, processorKafkaTopicHandlers, *cfg); err != nil {
			log.Panicf("[ERROR] kafka consumer group error : %+v\n", err)
		}
	}(ctx)

	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)
	// Block until a signal is received
	go func() {
		defer func() {
			<-ctx.Done()
			stop()
			wg.Done()
		}()

		exit := <-quit
		logger.Info(fmt.Sprintf("exit : %v", exit))
	}()
	wg.Wait()

	if err := e.Shutdown(ctx); err != nil {
		logger.Fatal("graceful shutdown failed", zap.Error(err))
	} else {
		logger.Info("graceful shutdown success !!")
	}
}

func initElasticsearch(cfg config.Config) (*elasticsearch.Client, error) {

	address := strings.Split(cfg.Secerts.ElasticAddress, ",")
	fmt.Println("address :", address)

	elasticConfig := elasticsearch.Config{
		Addresses: address,
		Username:  "",
		Password:  "",
	}

	es, err := elasticsearch.NewClient(elasticConfig)
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
		zap.Any("address", address),
		zap.String("username", ""),
		zap.String("password", ""))

	return es, nil
}
