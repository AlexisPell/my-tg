package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

type kafkaConfig struct {
	Url string
}

type config struct {
	Port  int
	Kafka *kafkaConfig
}

var (
	configSingleton *config
	once            sync.Once
)

func GetConfig() *config {
	once.Do(func() {
		fmt.Println(">>> Initializing config...")
		portStr := os.Getenv("PORT")
		if portStr == "" {
			panic("Env variable PORT is not specified")
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			panic("Non integer PORT value")
		}
		kafkaUrl := os.Getenv("KAFKA_URL")
		if kafkaUrl == "" {
			panic("KAFKA_URL is not specified")
		}

		configSingleton = &config{
			Port: port,
			Kafka: &kafkaConfig{
				Url: kafkaUrl,
			},
		}
	})

	return configSingleton
}
