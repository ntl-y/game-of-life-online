package server

import (
	"bytes"
	"encoding/json"
	"fmt"
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
			h.world.PaintPixel(h.pixelArray, pixelMsg.IndexOfPixel, pixelMsg.X, pixelMsg.Y, allColors[pixelMsg.ColorIndex])
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
	indices, err := pixelsToColorIndices(h.pixelArray)
	if err != nil {
		logrus.Println("Error:", err)
		return
	}
	message := Message{Data: indices}
	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

func pixelsToColorIndices(pixels []byte) ([]uint8, error) {
	indices := make([]uint8, 0, len(pixels)/4)

	for i := 0; i < len(pixels); i += 4 {
		pixel := pixels[i : i+4] // RGBA
		index, found := findColorIndex(pixel)
		if !found {
			return nil, fmt.Errorf("color not found in palette")
		}
		indices = append(indices, index)
	}

	return indices, nil
}

func findColorIndex(pixel []byte) (uint8, bool) {
	for i, color := range allColors {
		if bytes.Equal(pixel, color) {
			return uint8(i), true
		}
	}
	return 0, false
}
