package websocket_manager

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	connection *websocket.Conn

	hub *Hub
}

func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		connection: conn,
		hub:        hub,
	}
}

// Nhận msg từ client -> server
func (c *Client) readMessages() {

	defer func() {
		// Cleanup connection
		c.hub.unregisterClient(c)
	}()
	for {
		// readMessages được sử dụng để đọc tin nhắn tiếp theo trong hàng đợi
		_, payload, err := c.connection.ReadMessage()
		if err != nil {
			log.Println(err)
			break // Break the Loop để đóng Conn & Dọn dẹp
		}

		// Đọc data từ client gửi lên và in ra tại đây
		c.hub.broadcast <- payload
		log.Println("Payload from client to server: ", string(payload))
	}
}

// Gửi msg từ server -> client
func (c *Client) writeMessages() {
	// Trả lại msg mà client gửi lên server
	for {
		msg := <-c.hub.broadcast
		if err := c.connection.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println(err)
			break
		}
		log.Println("sent message from server to client:", string(msg))
	}
}
