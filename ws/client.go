package ws

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/yonraz/gochat_sockethandler/constants"
	"github.com/yonraz/gochat_sockethandler/events/publishers"
	"github.com/yonraz/gochat_sockethandler/initializers"
	"github.com/yonraz/gochat_sockethandler/models"
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
	Receiver string `json:"receiver"`
	Content string `json:"content"`
	RoomID string `json:"roomId"`
	Type constants.MessageType `json:"type"`
	Status constants.RoutingKey `json:"status"`
}

func (client *Client) readPump(handler *Handler) {
	defer func() {
			handler.Unregister <- client
			client.Conn.Close()
	}()
	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v\n", err)
			}
			break
		}
		var msg *models.WsMessage
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Panicf("error unmarshaling while reading message: %v", err)
		}

		
		handler.Broadcast <- &Message{
			Content: msg.Content,
			RoomID: client.RoomID,
			Sender: client.Username,
			Receiver: msg.Receiver,
			Type: msg.Type,
			Status: msg.Status,
		}
		p := publishers.NewPublisher(initializers.RmqChannel)

		p.MessageEvent(msg.Status, message)
	}
}

func (client *Client) writePump() {
	defer func() {
			client.Conn.Close()
	}()
	for message := range client.Message {
		client.Conn.WriteJSON(message)
	}
}