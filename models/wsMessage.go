package models

import (
	"github.com/yonraz/gochat_sockethandler/constants"
	"gorm.io/gorm"
)

type WsMessage struct {
	ID      	string 					`json:"id" gorm:"primary key"`
	Content 	string					`json:"content"`		
	Receiver 	string					`json:"receiver"`
	Status  	constants.RoutingKey	`json:"status"`	
	Type 		constants.MessageType	`json:"type"`
	gorm.Model
}