package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

var Producer sarama.SyncProducer

func InitKafkaProducer(brokers []string) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	var err error
	Producer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatal(">>> Error instantiating Kafka Producer:", err)
	}
	log.Println(">>> Kafka Producer instantiated successfully :)")
}

func SendMessageToKafka(msg *sarama.ProducerMessage) error {
	_, _, err := Producer.SendMessage(msg)
	if err != nil {
		log.Println(">>> Error on pulling msg to Kafka:", err)
		return err
	}
	return nil
}

func CloseKafkaProducer() {
	Producer.Close()
}
