package websockets

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var WsClients = make(map[*websocket.Conn]bool)
var mutex = sync.Mutex{}

func HandleConnections(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(">>> Read msg from WebSocket error :", err)
		return
	}
	defer ws.Close()

	mutex.Lock()
	WsClients[ws] = true
	mutex.Unlock()
	log.Default().Println(">>> New Websocket Connection. addr: ", ws.NetConn().RemoteAddr())

	for {
		_, msg, err := ws.ReadMessage()
		fmt.Println(">>> Message from ws client: ", msg)
		if err != nil {
			log.Println(">>> Client disconnected")
			mutex.Lock()
			delete(WsClients, ws)
			mutex.Unlock()
			break
		}
	}
}
