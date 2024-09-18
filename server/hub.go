package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	screenWidth  = 64
	screenHeight = 48
)

type IndicesMessage struct {
	//slice of color indices
	Data []int `json:"Data"`
}

type PixelMessage struct {
	X            int `json:"x"`
	Y            int `json:"y"`
	IndexOfPixel int `json:"index"`
	ColorIndex   int `json:"color"`
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	world      *World
	pixelArray []byte
}

func NewHub() *Hub {
	h := &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		world:      NewWorld(screenWidth, screenHeight),
		pixelArray: make([]byte, screenWidth*screenHeight*4),
	}

	for i := 0; i < len(h.pixelArray); i += 4 {
		h.pixelArray[i] = 0   // R
		h.pixelArray[i+1] = 0 // G
		h.pixelArray[i+2] = 0 // B
		h.pixelArray[i+3] = 1 // A
	}

	return h
}

func (h *Hub) Run() {
	updateTicker := time.NewTicker(time.Second / 20)
	defer updateTicker.Stop()

	go h.sendWorld()
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.sendWorld()
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				// if len(h.clients) == 0 {
				// 	h.resetWorld()
				// }
			}
		case message := <-h.broadcast:
			var pixelMsg PixelMessage
			err := json.Unmarshal(message, &pixelMsg)
			if err != nil {
				logrus.Println(err)
				continue
			}
			h.world.PaintPixel(h.pixelArray, pixelMsg.IndexOfPixel, pixelMsg.X, pixelMsg.Y, allColors[pixelMsg.ColorIndex])
			h.sendWorld()

		case <-updateTicker.C:
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
	message, err := json.Marshal(IndicesMessage{Data: indices})
	if err != nil {
		logrus.Println("Error:", err)
		return
	}
	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

func pixelsToColorIndices(pixels []byte) ([]int, error) {
	indices := make([]int, 0, len(pixels)/4)

	for i := 0; i < len(pixels); i += 4 {
		pixel := pixels[i : i+4] // RGBA
		index, found := findColorIndex(pixel)
		if !found {
			return nil, errors.New("color not found in palette")
		}
		indices = append(indices, index)
	}

	return indices, nil
}

func findColorIndex(pixel []byte) (int, bool) {
	for i, color := range allColors {
		if bytes.Equal(pixel, color) {
			return i, true
		}
	}
	return 0, false
}

func (h *Hub) resetWorld() {
	h.world = NewWorld(screenWidth, screenHeight)
	for i := 0; i < len(h.pixelArray); i += 4 {
		h.pixelArray[i] = 0   // R
		h.pixelArray[i+1] = 0 // G
		h.pixelArray[i+2] = 0 // B
		h.pixelArray[i+3] = 1 // A
	}
}
