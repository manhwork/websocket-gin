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

func (c *Client) readMessages() {

	defer func() {
		// Cleanup connection
		c.hub.unregisterClient(c)
	}()
	for {
		messageType, payload, err := c.connection.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(messageType)
		log.Println(string(payload))
	}
}
