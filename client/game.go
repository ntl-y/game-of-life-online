package client

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"
)

var (
	allColors = GetAllColors()
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

type Game struct {
	width       int
	height      int
	Pixels      []byte
	PlayerColor []byte
	ColorIndex  int
	Mu          sync.Mutex
	Conn        *websocket.Conn
}

func (g *Game) indexOfPixel(x, y int) int {
	return (y*g.width + x) * 4
}

func colorIndicesToPixels(indices []uint8, colorPalette []Color) []byte {
	pixels := make([]byte, 0, len(indices)*4)

	for _, index := range indices {
		color := colorPalette[index]
		pixels = append(pixels, color...)
	}

	return pixels
}

func (g *Game) ReadFromSocket() {
	var message Message
	if err := g.Conn.ReadJSON(&message); err != nil {
		logrus.Error("Error sending pixel data:", err)
	} else {
		pixels := colorIndicesToPixels(message.Data, allColors)
		g.Pixels = pixels
	}
}

func (g *Game) SendToSocket(pixelMsg PixelMessage) {
	err := g.Conn.WriteJSON(pixelMsg)
	if err != nil {
		logrus.Error("Error sending pixel data:", err)
	}
}

func (g *Game) PaintOnScreen(pix []byte, pixel PixelMessage) {
	if len(pix) > 0 && pixel.IndexOfPixel < len(pix) {
		g.Pixels[pixel.IndexOfPixel] = g.PlayerColor[0]
		g.Pixels[pixel.IndexOfPixel+1] = g.PlayerColor[1]
		g.Pixels[pixel.IndexOfPixel+2] = g.PlayerColor[2]
		g.Pixels[pixel.IndexOfPixel+3] = g.PlayerColor[3]

		g.SendToSocket(pixel)
	}
}

func (g *Game) PaintPlayer() {
	mx, my := ebiten.CursorPosition()
	if mx >= 0 && mx < g.width && my >= 0 && my < g.height {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			currPixel := PixelMessage{
				X:            mx,
				Y:            my,
				IndexOfPixel: g.indexOfPixel(mx, my),
				ColorIndex:   g.ColorIndex,
			}

			currPixelUp := PixelMessage{
				X:            mx,
				Y:            my + 1,
				IndexOfPixel: g.indexOfPixel(mx, my+1),
				ColorIndex:   g.ColorIndex,
			}

			currPixelRight := PixelMessage{
				X:            mx + 1,
				Y:            my,
				IndexOfPixel: g.indexOfPixel(mx+1, my),
				ColorIndex:   g.ColorIndex,
			}
			g.PaintOnScreen(g.Pixels, currPixel)
			g.PaintOnScreen(g.Pixels, currPixelUp)
			g.PaintOnScreen(g.Pixels, currPixelRight)

		}

	}
}

func (g *Game) Update() error {
	go g.PaintPlayer()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	go g.ReadFromSocket()
	screen.WritePixels(g.Pixels)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.width, g.height
}

func NewGame(screenWidth, screenHeight, colorIndex int, cellColor []byte, conn *websocket.Conn) *Game {
	g := &Game{
		Conn:        conn,
		PlayerColor: cellColor,
		ColorIndex:  colorIndex,
		Pixels:      nil,
	}
	g.ReadFromSocket()

	return g
}
