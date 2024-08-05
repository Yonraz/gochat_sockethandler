package ws

import (
	"encoding/json"
	"log"
	"time"

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
	ID string `json:"id"  gorm:"primary key"`
	Sender string `json:"sender"`
	Receiver string `json:"receiver"`
	Content string `json:"content"`
	RoomID string `json:"roomId"`
	Type constants.MessageType `json:"type"`
	Status constants.RoutingKey `json:"status"`
	Read bool `json:"read"`
	Sent bool `json:"sent"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
		log.Printf("reading message: %v", string(message))
		var msg *models.WsMessage
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Panicf("error unmarshaling while reading message %v: %v", string(message), err)
			break
		}
		msg.Sender = client.Username

		if (msg.Type == constants.MessageCreate) {
			msg.CreatedAt = time.Now()
		}
		messageToSend := &Message{
			ID: msg.ID,
			Content: msg.Content,
			RoomID: client.RoomID,
			Sender: client.Username,
			Receiver: msg.Receiver,
			Type: msg.Type,
			Status: msg.Status,
			Read: msg.Status == constants.MessageReadKey,
			Sent: msg.Sent,
			CreatedAt: msg.CreatedAt,
			UpdatedAt: msg.UpdatedAt,
		}

		handler.Broadcast <- messageToSend
		p := publishers.NewPublisher(initializers.RmqChannel)

		p.MessageEvent(msg.Status, msg)
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