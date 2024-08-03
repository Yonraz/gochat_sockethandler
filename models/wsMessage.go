package models

import (
	"github.com/yonraz/gochat_sockethandler/constants"
	"gorm.io/gorm"
)

type WsMessage struct {
	gorm.Model
	ID      	string 					`json:"id" gorm:"primary key"`
	Content 	string					`json:"content"`	
	Sender 		string 					`json:"sender"`	
	Receiver 	string					`json:"receiver"`
	Status  	constants.RoutingKey	`json:"status"`	
	Type 		constants.MessageType	`json:"type"`
	Read 		bool					`json:"read"`
	Sent 		bool					`json:"sent"`
}
