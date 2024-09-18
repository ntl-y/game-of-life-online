package server

import (
	"fmt"
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
	pongWait   = 120 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var (
	mu        sync.Mutex
	colorMap  map[string]bool
	allColors = GetAllColors()
	upgrader  = websocket.Upgrader{
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
	colorMap := make(map[string]bool)
	for i, color := range allColors {
		if i == 0 {
			colorStr := string(color)
			colorMap[colorStr] = true //skip black as background
		}
		colorStr := string(color)
		colorMap[colorStr] = false
	}
}

func getRandomAvailableColorIndex() int {
	availableIndices := make([]int, 0)
	for i, color := range allColors {
		colorStr := string(color)
		if !colorMap[colorStr] {
			availableIndices = append(availableIndices, i)
		}
	}

	if len(availableIndices) == 0 {
		return -1
	}

	randomIndex := availableIndices[rand.Intn(len(availableIndices))]
	selectedColorStr := string(allColors[randomIndex])
	colorMap[selectedColorStr] = true

	return randomIndex
}

func initColor() http.Header {
	index := getRandomAvailableColorIndex()
	if index == -1 {
		return nil
	}

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
	color := initColor()
	if color == nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, map[string]interface{}{
			"message": "no available room for new players",
		})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, color)

	if err != nil {
		logrus.Println(err)
		c.AbortWithStatusJSON(http.StatusBadGateway, map[string]interface{}{
			"message": err,
		})
		return
	}

	logrus.Printf("Client connected: %s", conn.RemoteAddr().String())

	client := &Client{hub: hub, conn: conn, send: make(chan Message)}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
