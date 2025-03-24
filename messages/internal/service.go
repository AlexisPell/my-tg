package internal

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/alexispell/my-tg/messages/internal/kafka"
	"github.com/alexispell/my-tg/messages/pb"
	"github.com/gocql/gocql"
)

type MessageService struct {
	DB       *gocql.Session
	Producer sarama.SyncProducer
	pb.UnimplementedMessageServiceServer
}

func NewMessageService(
// db *gocql.Session, kafkaProducer sarama.SyncProducer
) *MessageService {
	return &MessageService{
		// DB: db, Producer: kafkaProducer
	}
}

func (s *MessageService) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	// Генерация UUID для сообщения
	messageID := gocql.TimeUUID().String()
	timestamp := time.Now()

	// Сохранение в ScyllaDB
	// err := db.Session.Query(`INSERT INTO messages (chat_id, message_id, sender_id, content, timestamp)
	// VALUES (?, ?, ?, ?, ?)`, req.ChatId, messageID, req.SenderId, req.Message, timestamp).Exec()
	// if err != nil {
	// 	log.Println("Error writing in ScyllaDB:", err)
	// 	return nil, err
	// }
	// err := s.DB.Query(`INSERT INTO messages (chat_id, message_id, sender_id, content, timestamp)
	// 	VALUES (?, ?, ?, ?, ?)`, req.ChatId, messageID, req.SenderId, req.Message, timestamp).Exec()

	// Отправка сообщения в Kafka
	msg := &sarama.ProducerMessage{
		Topic: "chat-messages",
		Key:   sarama.StringEncoder(req.ChatId),
		Value: sarama.StringEncoder(messageID + "|" + req.SenderId + "|" + req.Message),
	}
	err := kafka.SendMessageToKafka(msg)
	if err != nil {
		return nil, err
	}
	// _, _, err = s.Producer.SendMessage(msg)
	// if err != nil {
	// 	log.Println("Error sending msg to kafka:", err)
	// 	return nil, err
	// }

	return &pb.SendMessageResponse{MessageId: messageID, Timestamp: timestamp.String()}, nil
}
