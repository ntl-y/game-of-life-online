package client

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"
)

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
		currPixel := g.indexOfPixel(mx, my)
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.Pixels[currPixel] == 0 {
			currPixelUp := g.indexOfPixel(mx, my+1)
			currPixelRight := g.indexOfPixel(mx+1, my)

			g.Mu.Lock()
			g.World.Paint(g.Pixels, currPixel, mx, my, g.PlayerColor)
			g.World.Paint(g.Pixels, currPixelUp, mx, my+1, g.PlayerColor)
			g.World.Paint(g.Pixels, currPixelRight, mx+1, my, g.PlayerColor)
			g.Mu.Unlock()

			err := g.Conn.WriteJSON(PixelValidate{
				X:     mx,
				Y:     my,
				Color: g.PlayerColor,
			})

			if err != nil {
				logrus.Error("Error sending pixel data:", err)
			}

		}
	}
}

func (g *Game) PaintEnemy(mx, my int, color []byte) {
	if mx >= 0 && mx < g.World.width && my >= 0 && my < g.World.height {
		currPixel := g.indexOfPixel(mx, my)
		currPixelUp := g.indexOfPixel(mx, my+1)
		currPixelRight := g.indexOfPixel(mx+1, my)

		g.Mu.Lock()
		g.World.Paint(g.Pixels, currPixel, mx, my, color)
		g.World.Paint(g.Pixels, currPixelUp, mx, my+1, color)
		g.World.Paint(g.Pixels, currPixelRight, mx+1, my, color)
		g.Mu.Unlock()
	}
}

func (g *Game) Update() error {
	g.Mu.Lock()
	g.World.UpdatePixels(g.Pixels)
	g.Mu.Unlock()

	go g.PaintPlayer()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Mu.Lock()
	g.World.DrawPixels(g.Pixels)
	g.Mu.Unlock()
	screen.WritePixels(g.Pixels)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.World.width, g.World.height
}

func NewGame(screenWidth, screenHeight int, cellColor []byte, conn *websocket.Conn) *Game {
	return &Game{
		World:       NewWorld(screenWidth, screenHeight),
		Conn:        conn,
		PlayerColor: cellColor,
		Pixels:      make([]byte, screenWidth*screenHeight*4),
	}
}
