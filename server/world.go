package server

var background = allColors[0]

type World struct {
	width    int
	height   int
	area     []bool
	colorMap map[int][]byte
}

func NewWorld(width int, height int) *World {
	w := &World{
		width:    width,
		height:   height,
		area:     make([]bool, width*height),
		colorMap: make(map[int][]byte),
	}

	//fill with background
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			index := w.indexInArea(x, y)
			w.colorMap[index] = background
		}
	}
	return w
}

func (w *World) indexInArea(x, y int) int {
	return y*w.width + x
}

func (w *World) colorCell(index int, color []byte) {
	w.area[index] = true
	w.colorMap[index] = color
}

func (w *World) updateCells() {
	newArea := make([]bool, len(w.area))
	newColorMap := make(map[int][]byte)

	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			index := w.indexInArea(x, y)

			if w.area[index] {
				cellColor := w.colorMap[index]
				neighbours := w.countNeighboursForAliveCells(x, y)
				if neighbours == 2 || neighbours == 3 {
					newArea[index] = true
					newColorMap[index] = cellColor
				} else {
					newArea[index] = false
					newColorMap[index] = background
				}
			} else {
				neighbours, color := w.countNeighboursForDeadCells(x, y)
				if neighbours == 3 && color != nil {
					newArea[index] = true
					newColorMap[index] = color
				} else {
					newArea[index] = false
					newColorMap[index] = background
				}
			}
		}
	}
	w.area = newArea
	w.colorMap = newColorMap
}

func (w *World) countNeighboursForAliveCells(x, y int) int {
	neighbours := 0
	directions := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1}, {1, -1},
		{1, 0}, {1, 1},
	}

	for _, dir := range directions {
		nx, ny := x+dir[0], y+dir[1]
		if nx >= 0 && nx < w.width && ny >= 0 && ny < w.height {
			neighbourIndex := w.indexInArea(nx, ny)
			if w.area[neighbourIndex] {
				neighbours++
			}
		}
	}
	return neighbours
}

func (w *World) countNeighboursForDeadCells(x, y int) (int, []byte) {
	neighbours := 0
	colorCount := make(map[string]int)
	var dominantColor []byte

	directions := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1}, {1, -1},
		{1, 0}, {1, 1},
	}

	for _, dir := range directions {
		nx, ny := x+dir[0], y+dir[1]
		if nx >= 0 && nx < w.width && ny >= 0 && ny < w.height {
			neighbourIndex := w.indexInArea(nx, ny)
			if w.area[neighbourIndex] {
				neighbours++
				color := string(w.colorMap[neighbourIndex])
				colorCount[color]++
				if dominantColor == nil || colorCount[color] > colorCount[string(dominantColor)] {
					dominantColor = w.colorMap[neighbourIndex]
				}
			}
		}
	}

	if neighbours == 3 && dominantColor != nil {
		return neighbours, dominantColor
	}
	return neighbours, nil
}

func (w *World) UpdatePixels(pixelArray []byte) {
	w.updateCells()
	w.cellsToPixels(pixelArray)

}

func (w *World) cellsToPixels(pixels []byte) {
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			pixelIndex := (y*w.width + x) * 4
			index := w.indexInArea(x, y)
			color := w.colorMap[index]
			w.colorPixel(pixels, pixelIndex, color)

		}
	}
}

func (w *World) PaintPixel(pix []byte, pixelIndex int, x, y int, color []byte) {
	if len(pix) > 0 && pixelIndex < len(pix) {
		i := w.indexInArea(x, y)
		w.colorCell(i, color)
		w.colorPixel(pix, pixelIndex, color)

	}
}

func (w *World) colorPixel(pixels []byte, pixelIndex int, color []byte) {
	pixels[pixelIndex] = color[0]
	pixels[pixelIndex+1] = color[1]
	pixels[pixelIndex+2] = color[2]
	pixels[pixelIndex+3] = color[3]
}
