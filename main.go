package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yonraz/gochat_sockethandler/initializers"
)

func init() {
	time.Sleep(15 * time.Second)
	fmt.Println("Application starting...")
	initializers.LoadEnvVariables()
	initializers.ConnectToRabbitmq()
}

func main() {
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
}