package kafka

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/alexispell/my-tg/messages/internal/websockets"
)

type KafkaConsumer struct {
	Consumer sarama.Consumer
}

func NewKafkaConsumer(brokers []string) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}
	// defer Consumer.Close()

	return &KafkaConsumer{Consumer: consumer}, nil
}

func (kc *KafkaConsumer) ConsumeKafkaMessages(topic string) {
	partitionConsumer, err := kc.Consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf(">>> Error subscribing on kafka. Topic: %s", err)
	}
	defer partitionConsumer.Close()

	for message := range partitionConsumer.Messages() {
		log.Println(">>> Reading new message from kafka:", string(message.Value))
		broadcastMessage(string(message.Value))
	}
}

func broadcastMessage(message string) {
	for client := range websockets.WsClients {
		err := client.WriteMessage(1, []byte(message))
		if err != nil {
			log.Println(">>> Error sending WebSocket-message:", err)
			client.Close()
			delete(websockets.WsClients, client)
		}
	}
}
