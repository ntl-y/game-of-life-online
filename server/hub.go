package server

import (
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	screenWidth  = 64
	screenHeight = 48
)

type Message struct {
	//slice of color indices
	Data []uint8
}

type PixelMessage struct {
	X            int `json:"x"`
	Y            int `json:"y"`
	IndexOfPixel int `json:"index"`
	ColorIndex   int `json:"color"`
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	world      *World
	pixelArray []byte
}

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message),
		world:      NewWorld(screenWidth, screenHeight),
		pixelArray: make([]byte, screenWidth*screenHeight*4),
	}
}

func (h *Hub) Run() {
	ticker := time.NewTicker(time.Second / 20)
	defer ticker.Stop()

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			var pixelMsg PixelMessage
			err := json.Unmarshal(message.Data, &pixelMsg)
			if err != nil {
				logrus.Println(err)
				continue
			}
			h.world.PaintPixel(h.pixelArray, pixelMsg.IndexOfPixel, pixelMsg.X, pixelMsg.Y, pixelMsg.Color)
			h.sendWorld()

		case <-ticker.C:
			h.updateWorld()
		}
	}
}

func (h *Hub) updateWorld() {
	h.world.UpdatePixels(h.pixelArray)
	h.sendWorld()
}

func (h *Hub) sendWorld() {
	message := Message{Data: h.pixelArray}
	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}
