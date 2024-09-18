package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

var (
	allColors = client.GetAllColors()
)

func getHeaderColor(r *http.Response) ([]byte, int, error) {
	colorIndexStr := r.Header.Get("Color")
	if colorIndexStr == "" {
		return nil, 0, errors.New("Color header not found")
	}
	colorIndex, err := strconv.Atoi(colorIndexStr)
	if err != nil {
		return nil, 0, fmt.Errorf("Invalid color index: %v", err)
	}
	if colorIndex < 0 || colorIndex >= len(allColors) {
		return nil, 0, errors.New("Color index out of range")
	}
	return allColors[colorIndex], colorIndex, nil
}

func main() {
	ebiten.SetTPS(20)
	ebiten.SetWindowSize(screenFrameWidth, screenFrameHeight)
	ebiten.SetWindowTitle("Game of Life")

	conn, respColor, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:3000/", nil)
	if err != nil {
		logrus.Fatal(err)
	}

	cellColor, colorIndex, err := getHeaderColor(respColor)
	if err != nil {
		logrus.Fatal(err)
	}

	game := client.NewGame(screenWidth, screenHeight, colorIndex, cellColor, conn)
	defer game.Conn.Close()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
