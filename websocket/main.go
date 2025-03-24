package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// func wsHandler(c *gin.Context) {
// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		log.Println("WebSocket Upgrade error:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	for {
// 		_, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("Read error:", err)
// 			break
// 		}
// 		log.Println("Received:", string(msg))

// 		err = conn.WriteMessage(websocket.TextMessage, []byte("Echo: "+string(msg)))
// 		if err != nil {
// 			log.Println("Write error:", err)
// 			break
// 		}
// 	}
// }

var Clients = make(map[*websocket.Conn]bool)
var mutex = sync.Mutex{}

func HandleConnections(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Read msg from WebSocket error :", err)
		return
	}
	defer ws.Close()

	mutex.Lock()
	Clients[ws] = true
	mutex.Unlock()
	log.Default().Println("New Websocket Connection. addr: ", ws.NetConn().RemoteAddr())

	for {
		_, msg, err := ws.ReadMessage()
		fmt.Println("Message from ws client: ", msg)
		if err != nil {
			log.Println("Client disconnected")
			mutex.Lock()
			delete(Clients, ws)
			mutex.Unlock()
			break
		}
	}
}

func main() {
	// Config
	port := os.Getenv("PORT")
	if port == "" {
		panic("Env variable PORT is not specified")
	}
	kafkaUrl := os.Getenv("KAFKA_URL")
	if kafkaUrl == "" {
		panic("KAFKA_URL is not specified")
	}

	// Services
	go ConsumeKafkaMessages([]string{"kafka:9092"})

	// Routing
	r := gin.Default()

	r.GET("/api/v1/ws", HandleConnections)

	r.GET("/api/v1/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Execution
	log.Printf("Server running on port %s \n", port)
	r.Run(fmt.Sprintf(":%s", port))
}
