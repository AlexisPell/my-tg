package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/alexispell/my-tg/messages/internal"
	"github.com/alexispell/my-tg/messages/internal/config"
	"github.com/alexispell/my-tg/messages/internal/kafka"
	"github.com/alexispell/my-tg/messages/pb"
)

func main() {
	// Config
	cfg := config.GetConfig()

	// // Init Prometheus metrics
	// requests := prometheus.NewCounterVec(
	// 	prometheus.CounterOpts{
	// 		Name: "grpc_requests_total",
	// 		Help: "Total number of gRPC requests",
	// 	},
	// 	[]string{"method"},
	// )
	// prometheus.MustRegister(requests)
	// go func() {
	// 	http.Handle("/metrics", promhttp.Handler())
	// 	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.MetricsPort), nil))
	// }()

	// r := gin.Default()

	// Init ScyllaDB
	// db.InitScylla()
	// defer db.CloseScylla()

	kafka.InitKafkaProducer([]string{cfg.Kafka.Url})
	defer kafka.CloseKafkaProducer()

	messageService := internal.NewMessageService()

	// tcp listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port)) // here was s
	if err != nil {
		log.Fatal(">>> Error starting tcp server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMessageServiceServer(grpcServer, messageService)

	log.Println(">>> MessageService starting on port: ", cfg.Port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(">>> Global error on running messages grpc service: ", err)
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
	// 	fmt.Println(">>> Payload: ", payload)
	// 	service.SendMessage(context.Background(), payload)
	// })

	// log.Printf(">>> Server running on port %s \n", port)
	// r.Run(fmt.Sprintf(":%s", port))
}
