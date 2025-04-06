package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/alexispell/my-tg/messages/internal/config"
	"github.com/alexispell/my-tg/messages/internal/kafka"
	"github.com/alexispell/my-tg/messages/internal/websockets"
)

func main() {
	// Config
	cfg := config.GetConfig()

	// Services
	// Listen to kafka
	kafkaConsumer, err := kafka.NewKafkaConsumer([]string{cfg.Kafka.Url})
	if err != nil {
		panic(fmt.Sprintf(">>> Error trying to create kafka consumer: %s", err))
	}
	fmt.Println(">>> Kafka consumer created")
	defer kafkaConsumer.Consumer.Close()

	go kafkaConsumer.ConsumeKafkaMessages("chat-messages")

	// Routing
	r := gin.Default()

	r.GET("/api/v1/ws", websockets.HandleConnections)

	r.GET("/api/v1/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Listener
	log.Printf(">>> Server running on port %d \n", cfg.Port)
	r.Run(fmt.Sprintf(":%d", cfg.Port))
}
