package models

import (
	"time"

	"github.com/yonraz/gochat_sockethandler/constants"
)

type WsMessage struct {
	ID      	string 					`json:"id" gorm:"primary key"`
	ConversationID string 				`json:"conversationId" gorm:"foreign key"`
	Content 	string					`json:"content"`	
	Sender 		string 					`json:"sender"`	
	Receiver 	string					`json:"receiver"`
	Status  	constants.RoutingKey	`json:"status"`	
	Type 		constants.MessageType	`json:"type"`
	Read 		bool					`json:"read"`
	Sent 		bool					`json:"sent"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

