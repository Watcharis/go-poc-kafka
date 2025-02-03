package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"watcharis/go-poc-kafka/config"
	"watcharis/go-poc-kafka/logger"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"github.com/xdg/scram"
	otrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	SHA256 scram.HashGeneratorFcn = sha256.New
	SHA512 scram.HashGeneratorFcn = sha512.New
)

type XDGSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}

func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}

var kafkaProducer sarama.SyncProducer

func KafkaInitialize(ctx context.Context, tracerProvider otrace.TracerProvider, cfg *config.Config) {
	logger.Info("KafkaInitialize - start initialize")
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Net.SASL.Enable = true
	config.Net.SASL.User = ""
	config.Net.SASL.Password = ""
	config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
		return &XDGSCRAMClient{
			HashGeneratorFcn: SHA512,
		}
	}
	config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = createTLSConfigurationInternal(ctx)

	producer, err := sarama.NewSyncProducer(cfg.Secerts.KafkaAddress, config)
	if err != nil {
		logger.Panic(fmt.Sprintf("Could not create producer: %v", err), zap.Error(err))
	}

	// producer = otelsarama.WrapSyncProducer(config, producer, otelsarama.WithTracerProvider(tracerProvider))
	kafkaProducer = producer

	logger.Info("KafkaInitialize - init kafka complete")
}

func ProduceKafkaTransactionInternal(ctx context.Context, request map[string]interface{}) error {

	logger.Info("ProduceKafkaTransactionInternal - send data to kafka", zap.Any("request", request))

	topic := viper.GetString("kafka.topic-internal")
	producer := kafkaProducer

	in, err := json.Marshal(request)
	if err != nil {
		logger.Error(fmt.Sprintf("ProduceKafkaTransactionInternal - Could not marshall: %v", err), zap.Any("request", request), zap.Error(err))
	}
	bPrefix := []byte("CBSTRANS")

	body := bytes.Join([][]byte{bPrefix, in}, []byte("||"))

	message := sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(body)}

	// for telemetry trace format
	// propagator := propagation.NewCompositeTextMapPropagator(
	// 	// gcppropagator.CloudTraceOneWayPropagator{},
	// 	gcppropagator.CloudTraceFormatPropagator{},
	// 	propagation.TraceContext{},
	// 	propagation.Baggage{},
	// )

	// propagator.Inject(otelsarama.NewProducerMessageCarrier(&message))

	partition, offset, err := producer.SendMessage(&message)
	if err != nil {
		logger.Error(fmt.Sprintf("ProduceKafkaTransactionInternal - Error sending message: %v", err), zap.Error(err))
	} else {
		logger.Info(fmt.Sprintf("ProduceKafkaTransactionInternal - Sent message value at partition = %d, offset = %d", partition, offset))
	}

	return err
}

func createTLSConfigurationInternal(ctx context.Context) (t *tls.Config) {

	cert := []byte(viper.GetString("kafka-internal.tls"))

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(cert)
	if !ok {
		logger.Error("failed to parse root certificate")
	}

	t = &tls.Config{
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}

	return t
}
