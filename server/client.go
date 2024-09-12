package server

import (
	"encoding/base64"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan Message
	id   string
}

func initColor() (color http.Header) {
	c := []byte{
		byte(rand.Intn(256)),
		byte(rand.Intn(256)),
		byte(rand.Intn(256)),
		byte(rand.Float64() * 255),
	}

	cellColorEncoded := base64.StdEncoding.EncodeToString(c)

	responseHeader := http.Header{}
	responseHeader.Set("Color", cellColorEncoded)
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
		c.hub.broadcast <- Message{ClientID: c.id, Data: message}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
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

			if message.ClientID == c.id {
				continue
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

func generateClientID() string {
	return base64.StdEncoding.EncodeToString([]byte(time.Now().String() + string(rand.Int())))
}

func ServeWs(hub *Hub, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, initColor())

	if err != nil {
		log.Println(err)
		return
	}

	logrus.Printf("Client connected: %s", conn.RemoteAddr().String())

	clientID := generateClientID()
	client := &Client{hub: hub, conn: conn, send: make(chan Message), id: clientID}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
