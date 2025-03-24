package main

import (
	"log"

	"github.com/IBM/sarama"
)

func ConsumeKafkaMessages(brokers []string) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatal("Error creating kafka consumer:", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("chat-messages", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error subscribing on kafka. Topic: %s", err)
	}
	defer partitionConsumer.Close()

	for message := range partitionConsumer.Messages() {
		log.Println("Reading new message from kafka:", string(message.Value))
		broadcastMessage(string(message.Value))
	}
}

func broadcastMessage(message string) {
	for client := range Clients {
		err := client.WriteMessage(1, []byte(message))
		if err != nil {
			log.Println("Error sending WebSocket-message:", err)
			client.Close()
			delete(Clients, client)
		}
	}
}
