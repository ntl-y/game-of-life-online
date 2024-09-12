package main

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"

	client "github.com/ntl-y/gameoflife/client"
	"github.com/sirupsen/logrus"
)

const (
	screenFrameWidth  = 640
	screenFrameHeight = 480
	screenWidth       = 64
	screenHeight      = 48
)

func getHeaderColor(r *http.Response) ([]byte, error) {
	cellColorEncoded := r.Header.Get("Color")
	if cellColorEncoded == "" {
		return nil, errors.New("Color header not found")
	}
	cellColor, err := base64.StdEncoding.DecodeString(cellColorEncoded)
	if err != nil {
		return nil, err
	}
	if len(cellColor) != 4 {
		return nil, errors.New("Invalid cell color length")
	}
	return cellColor, nil
}

func main() {
	ebiten.SetTPS(20)
	ebiten.SetWindowSize(screenFrameWidth, screenFrameHeight)
	ebiten.SetWindowTitle("Game of Life")

	conn, respColor, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:3000/", nil)
	if err != nil {
		logrus.Fatal(err)
	}
	defer conn.Close()

	cellColor, err := getHeaderColor(respColor)
	if err != nil {
		logrus.Fatal(err)
	}

	game := client.NewGame(screenWidth, screenHeight, cellColor, conn)

	go func() {
		defer game.Conn.Close()
		for {
			var newPixel client.PixelValidate
			if err := game.Conn.ReadJSON(&newPixel); err != nil {
				logrus.Fatal(err)
			}
			game.PaintEnemy(newPixel.X, newPixel.Y, newPixel.Color)
		}
	}()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
