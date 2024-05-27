package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"watcharis/go-poc-kafka/consumer"
	"watcharis/go-poc-kafka/handlers/kafkaconsumerhandlers"
	"watcharis/go-poc-kafka/services"
)

func main() {
	ctx := context.Background()

	// init kafka services
	processorKafkaTopicService := services.NewProcessorKafkaTopic()

	//init kafka handler
	processorKafkaTopicHandlers := kafkaconsumerhandlers.NewKafkaConsumerHandlers(processorKafkaTopicService)

	//init kafka consumer group
	wg := new(sync.WaitGroup)
	go func(ctx context.Context) {
		log.Println("start kafka consumer")
		if err := consumer.NewKafkaConsumerGroup(ctx, processorKafkaTopicHandlers); err != nil {
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
		log.Println("exit :", exit)
	}()
	wg.Wait()
}
