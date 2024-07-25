package ws

import (
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	Conn *websocket.Conn
	Message chan *Message
	Username string
	RoomID string
}

type Message struct {
	gorm.Model
	Sender string `json:"sender"`
	Content string `json:"content"`
	RoomID string `json:"roomId"`
}

func (client *Client) readPump(handler *Handler) {
	defer func() {
			handler.handleDisconnection(client)
			client.Conn.Close()
		}()
		for {
			_, message, err := client.Conn.ReadMessage()
			if err != nil {
				break
			}
			handler.Broadcast <- &Message{
				Content: string(message),
				RoomID: client.RoomID,
				Sender: client.Username,
			}
		}
}

func (client *Client) writePump() {
	for message := range client.Message {
		client.Conn.WriteJSON(message)
	}
}