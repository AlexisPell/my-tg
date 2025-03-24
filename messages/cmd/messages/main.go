package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/alexispell/my-tg/messages/internal"
	"github.com/alexispell/my-tg/messages/internal/kafka"
	"github.com/alexispell/my-tg/messages/pb"

	"google.golang.org/grpc"
)

func main() {
	// config
	port := os.Getenv("PORT")
	if port == "" {
		panic("Env variable PORT is not specified")
	}

	// r := gin.Default()

	// db.InitScylla()
	// defer db.CloseScylla()

	kafka.InitKafkaProducer([]string{"kafka:9092"})
	defer kafka.CloseKafkaProducer()

	messageService := internal.NewMessageService()

	// tcp listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal("Error starting tcp server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMessageServiceServer(grpcServer, messageService)

	log.Println("MessageService starting on port: ", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Global error on running messages grpc service: ", err)
	}

	// r.GET("/api/v1/messages/ping", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{"message": "pong"})
	// })

	// r.POST("/api/v1/messages", func(c *gin.Context) {
	// 	payload := &pb.SendMessageRequest{
	// 		ChatId:   "chat-id-1",
	// 		SenderId: "sender-id-1",
	// 		Message:  "Hello world!!!",
	// 	}
	// 	fmt.Println("Payload: ", payload)
	// 	service.SendMessage(context.Background(), payload)
	// })

	// log.Printf("Server running on port %s \n", port)
	// r.Run(fmt.Sprintf(":%s", port))
}
