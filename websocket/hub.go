package websocket

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Cấu trúc Hub quản lý các kết nối Websocket

type Hub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	mutex     sync.Mutex
	upgrader  websocket.Upgrader
}

var WebSocketHub = &Hub{
	clients:   make(map[*websocket.Conn]bool),
	broadcast: make(chan []byte),
	upgrader: websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	},
}

// Xử lý khi Client kết nối vào Websocket

func (h *Hub) HandleConnections(c *gin.Context) {
	ws, err := h.upgrader.Upgrade(c.Writer, c.Request, nil) // Nâng cấp từ HTTP lên websocket
	if err != nil {
		fmt.Println("Websocket upgrader error :", err)
		return
	}
	defer ws.Close()

	// thêm client vào danh sách
	h.mutex.Lock()
	h.clients[ws] = true
	h.mutex.Unlock()

	fmt.Println("New Websocket connection")

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil { // Nếu lỗi thì xoá dần các client ra khỏi danh sách
			fmt.Println("Websocket Read error: ", err)
			h.mutex.Lock()
			delete(h.clients, ws)
			h.mutex.Unlock()
			break
		}
		fmt.Println("Received:", string(msg))

		// Gửi tin nhắn đến tất cả các client
		h.broadcast <- msg
	}

}

// Gửi tin nhắn đến tất cả các client
func (h *Hub) HandleMessages() {
	for {
		msg := <-h.broadcast
		h.mutex.Lock()
		for client := range h.clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				fmt.Println("write error", err)
				client.Close()
				delete(h.clients, client)
			}
		}
		h.mutex.Unlock()
	}
}
