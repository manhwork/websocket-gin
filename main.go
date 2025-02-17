package main

import (
	"websocket_gin/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/ws", websocket.WebSocketHub.HandleConnections)

	go websocket.WebSocketHub.HandleMessages()

	r.Run(":8080")
}
