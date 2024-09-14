package client

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"
)

var deadPixel = []byte{byte(0), byte(0), byte(0), byte(1)}

type PixelValidate struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Color []byte `json:"color"`
}

type Game struct {
	World       *World
	Pixels      []byte
	PlayerColor []byte
	Mu          sync.Mutex
	Conn        *websocket.Conn
}

func (g *Game) indexOfPixel(x, y int) int {
	return (y*g.World.width + x) * 4
}

func (g *Game) PaintPlayer() {

	mx, my := ebiten.CursorPosition()
	if mx >= 0 && mx < g.World.width && my >= 0 && my < g.World.height {

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			currPixel := g.indexOfPixel(mx, my)
			currPixelUp := g.indexOfPixel(mx, my+1)
			currPixelRight := g.indexOfPixel(mx+1, my)

			g.Mu.Lock()
			g.World.PaintPixel(g.Pixels, currPixel, mx, my, g.PlayerColor)
			g.World.PaintPixel(g.Pixels, currPixelUp, mx, my+1, g.PlayerColor)
			g.World.PaintPixel(g.Pixels, currPixelRight, mx+1, my, g.PlayerColor)
			g.Mu.Unlock()

			err := g.Conn.WriteJSON(PixelValidate{
				X:     mx,
				Y:     my,
				Color: g.PlayerColor,
			})

			if err != nil {
				logrus.Error("Error sending pixel data:", err)
			}
			logrus.Println("Sending pixel data from: ", g.Conn.LocalAddr())

		}

	}
}

func (g *Game) PaintEnemy() {
	for {
		var newPixel PixelValidate
		if err := g.Conn.ReadJSON(&newPixel); err != nil {
			logrus.Fatal(err)
		}
		logrus.Println("Recived pixel data, currentAdress: ", g.Conn.LocalAddr())

		if newPixel.X >= 0 && newPixel.X < g.World.width && newPixel.Y >= 0 && newPixel.Y < g.World.height {
			currPixel := g.indexOfPixel(newPixel.X, newPixel.Y)
			currPixelUp := g.indexOfPixel(newPixel.X, newPixel.Y+1)
			currPixelRight := g.indexOfPixel(newPixel.X+1, newPixel.Y)

			g.Mu.Lock()
			g.World.PaintPixel(g.Pixels, currPixel, newPixel.X, newPixel.Y, newPixel.Color)
			g.World.PaintPixel(g.Pixels, currPixelUp, newPixel.X, newPixel.Y+1, newPixel.Color)
			g.World.PaintPixel(g.Pixels, currPixelRight, newPixel.X+1, newPixel.Y, newPixel.Color)
			g.Mu.Unlock()

		}
	}
}

func (g *Game) Update() error {
	go g.PaintEnemy()
	go g.PaintPlayer()

	g.Mu.Lock()
	g.World.UpdatePixels(g.Pixels)
	g.Mu.Unlock()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Mu.Lock()
	g.World.CellsToPixels(g.Pixels)
	g.Mu.Unlock()

	screen.WritePixels(g.Pixels)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.World.width, g.World.height
}

func NewGame(screenWidth, screenHeight int, cellColor []byte, conn *websocket.Conn) *Game {
	game := &Game{
		World:       NewWorld(screenWidth, screenHeight),
		Conn:        conn,
		PlayerColor: cellColor,
		Pixels:      make([]byte, screenWidth*screenHeight*4),
	}
	return game
}
