package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yonraz/gochat_sockethandler/initializers"
	"github.com/yonraz/gochat_sockethandler/middlewares"
	"github.com/yonraz/gochat_sockethandler/ws"
)

func init() {
	fmt.Println("Application starting...")
	time.Sleep(1 * time.Minute)
	initializers.LoadEnvVariables()
	initializers.ConnectToRabbitmq()
	initializers.ConnectToRedis()
}

func main() {
	router := gin.Default()
	defer func() {
		if err := initializers.RmqChannel.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ channel: %v", err)
		}
	}()
	defer func() {
		if err := initializers.RmqConn.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ connection: %v", err)
		}
	}()
	wsHandler := ws.NewHandler(initializers.RedisClient)
	router.Use(middlewares.CurrentUser)
	router.Use(middlewares.RequireAuth)
	router.POST("/ws/chat/createRoom", wsHandler.CreateRoom)
	router.GET("/ws/chat/joinRoom/:roomId", middlewares.CurrentUser, middlewares.RequireAuth, wsHandler.JoinRoom)

	go wsHandler.Run()
	router.Run()
}