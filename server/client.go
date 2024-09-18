package server

import (
	"fmt"
	"image/color"
	"image/color/palette"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var (
	mu       sync.Mutex
	colorMap = make(map[Color]bool)
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan Message
}

func init() {
	mu.Lock()
	defer mu.Unlock()

	for _, col := range palette.Plan9 {
		colorMap[col] = true
	}
}

func getRandomAvailableColorIndex() int {
	mu.Lock()
	defer mu.Unlock()

	var availableColors []color.Color
	for col, available := range colorMap {
		if available {
			availableColors = append(availableColors, col)
		}
	}

	if len(availableColors) == 0 {
		return -1
	}
	randIndex := rand.Intn(len(availableColors))
	selectedColor := availableColors[randIndex]
	colorMap[selectedColor] = false

	return randIndex
}

func initColor() http.Header {
	index := getRandomAvailableColorIndex()

	responseHeader := http.Header{}
	responseHeader.Set("Color", fmt.Sprintf("%d", index))
	return responseHeader
}

func (c *Client) readPump() {
	defer func() {
		logrus.Printf("Client disconnected: %s", c.conn.RemoteAddr().String())

		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Printf("error: %v", err)
			}
			break
		}
		c.hub.broadcast <- Message{Data: message}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message.Data)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(hub *Hub, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, initColor())

	if err != nil {
		logrus.Println(err)
		return
	}

	logrus.Printf("Client connected: %s", conn.RemoteAddr().String())

	client := &Client{hub: hub, conn: conn, send: make(chan Message)}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
