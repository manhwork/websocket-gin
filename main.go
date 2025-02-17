package main

import (
	websocket_manager "websocket_gin/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/ws", websocket_manager.NewHub().ServeWS)

	r.Run(":8080")
}
