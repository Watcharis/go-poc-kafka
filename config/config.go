package config

import (
	"sync"
	"watcharis/go-poc-kafka/logger"

	"github.com/spf13/viper"
)

type Config struct {
	Kafka  Kafka  `mapstructure:"kafka"`
	Secert Secert `mapstructure:"secert"`
}

type Kafka struct {
	ConsumerGroup string `mapstructure:"consumer-group"`
	Topic1        string `mapstructure:"topic-poc-kafka-1"`
	Topic2        string `mapstructure:"topic-poc-kafka-2"`
}

type Secert struct {
	KafkaAddress []string `mapstructure:"kafka-address"`
}

func LoadConfig(path string) *Config {
	var config Config
	var configOnce sync.Once
	configOnce.Do(func() {

		logger.InitZapLogger()
		defer logger.Sync()

		viper.AddConfigPath(path)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")

		viper.AutomaticEnv()

		err := viper.ReadInConfig()
		if err != nil {
			logger.Panic(err.Error())
		}
		if err = viper.Unmarshal(&config); err != nil {
			logger.Panic(err.Error())
		}
	})

	return &config
}
