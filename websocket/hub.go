package websocket_manager

import (
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Hub duy trì các máy khách đang hoạt động
type Hub struct {
	// Client đã được đăng kí
	clients map[*Client]bool

	// Tin nhắn (sự kiện) được gửi từ client
	broadcast chan []byte

	// Cấu hình nâng cấp từ HTTP sang Websocket
	upgrader *websocket.Upgrader

	sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:   make(map[*Client]bool),
		broadcast: make(chan []byte),
		upgrader: &websocket.Upgrader{
			WriteBufferSize: 1024,
			ReadBufferSize:  1024,
		},
	}
}

func (h *Hub) registerClient(client *Client) {
	h.Lock()
	defer h.Unlock()

	h.clients[client] = true
}

func (h *Hub) unregisterClient(client *Client) {
	h.Lock()

	defer h.Unlock()

	// Nếu có client trong map thì đóng websocket của client đó lại -> xoá client đó ra khỏi map
	if _, ok := h.clients[client]; ok {
		client.connection.Close()
		delete(h.clients, client)
	}
}

func (h *Hub) ServeWS(c *gin.Context) {
	log.Println("new connection")
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Println(err)
		return
	}

	// Tạo 1 client mới
	client := NewClient(conn, h)

	h.registerClient(client)

	// Bắt đầu các tiến trình của client
	go client.readMessages()
}
