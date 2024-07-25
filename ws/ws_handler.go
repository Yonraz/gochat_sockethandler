package ws

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	hub        *redis.Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
	Clients    map[string]*Client
}

func NewHandler(h *redis.Client) *Handler {
	return &Handler{
		hub:        h,
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
		Clients:    make(map[string]*Client),
	}
}

type CreateRoomRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) CreateRoom(ctx *gin.Context) {
	var req CreateRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c := context.Background()
	roomKey := "room:" + req.ID

	h.hub.HSet(c, roomKey, "name", req.Name)
	ctx.JSON(http.StatusOK, req)
}

func (h *Handler) JoinRoom(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomID := ctx.Param("roomId")
	username := ctx.Param("username")
	client := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		Username: username,
		RoomID:   roomID,
	}

	c := context.Background()
	roomKey := "room:" + roomID

	h.Clients[username] = client
	h.hub.SAdd(c, roomKey+":clients", username)

	msg := &Message{
		Content: "User " + username + " has joined the room!",
		RoomID:  roomID,
		Sender:  username,
	}

	h.Register <- client
	h.Broadcast <- msg

	go client.readPump(h)
	go client.writePump()
}

func (h *Handler) handleDisconnection(client *Client) {
	c := context.Background()

	h.hub.SRem(c, "room:"+client.RoomID+":clients", client.Username)
	h.hub.Del(c, client.Username)

	delete(h.Clients, client.Username)
}

func (h *Handler) Run() {
	for {
		select {
		case client := <-h.Unregister:
			h.handleDisconnection(client)
			client.Conn.Close()
		case message := <-h.Broadcast:
			c := context.Background()
			// Get list of connected clients from Redis
			clients, err := h.hub.SMembers(c, "room:"+message.RoomID+":clients").Result()
			if err != nil {
				log.Println("Error fetching clients:", err)
				continue
			}

			// Send the message to each client
			for _, username := range clients {
				if client, ok := h.Clients[username]; ok {
					select {
					case client.Message <- message:
						// Successfully sent message to client
					default:
						// Handle case where message channel might be full
						log.Println("Message channel full for client:", username)
					}
				}
			}
		}
	}
}
