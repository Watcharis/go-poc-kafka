package config

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"watcharis/go-poc-kafka/logger"

	"github.com/spf13/viper"
)

type Config struct {
	Kafka   Kafka   `mapstructure:"kafka"`
	Secerts Secerts `mapstructure:"secerts"`
}

type Kafka struct {
	ConsumerGroup string `mapstructure:"consumer-group"`
	Topic1        string `mapstructure:"topic-poc-kafka-1"`
	Topic2        string `mapstructure:"topic-poc-kafka-2"`
}

type Secerts struct {
	KafkaAddress   []string `mapstructure:"kafka-address"`
	ElasticAddress string   `mapstructure:"elastic-address"`
	RedisAddress   string   `mapstructure:"redis-address"`
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

		if err := viper.ReadInConfig(); err != nil {
			logger.Panic(err.Error())
		}

		GetSecretValue()

		if err := viper.MergeConfig(strings.NewReader(viper.GetString("config"))); err != nil {
			logger.Panic(err.Error())
		}

		if err := viper.Unmarshal(&config); err != nil {
			logger.Panic(err.Error())
		}
	})

	return &config
}

func GetSecretValue() {
	for _, value := range os.Environ() {
		fmt.Println("env value :", value)
		pair := strings.SplitN(value, "=", 2)
		if strings.Contains(pair[0], "SECRET_") {
			keys := strings.Replace(pair[0], "SECRET_", "secrets.", -1)
			fmt.Println("keys :", keys)
			keys = strings.Replace(keys, "_", "-", -1)
			fmt.Println("keys :", keys)
			newKey := strings.Trim(keys, " ")
			newValue := strings.Trim(pair[1], " ")
			viper.Set(newKey, newValue)
		}
	}
}
