package server

type Color []byte

var (
	Black        = Color{0, 0, 0, 255}
	White        = Color{255, 255, 255, 255}
	Red          = Color{255, 0, 0, 255}
	Green        = Color{0, 255, 0, 255}
	Blue         = Color{0, 0, 255, 255}
	Yellow       = Color{255, 255, 0, 255}
	Cyan         = Color{0, 255, 255, 255}
	Magenta      = Color{255, 0, 255, 255}
	Orange       = Color{255, 165, 0, 255}
	Pink         = Color{255, 192, 203, 255}
	Purple       = Color{128, 0, 128, 255}
	Brown        = Color{165, 42, 42, 255}
	Gray         = Color{128, 128, 128, 255}
	LightGray    = Color{211, 211, 211, 255}
	DarkGray     = Color{169, 169, 169, 255}
	Navy         = Color{0, 0, 128, 255}
	Olive        = Color{128, 128, 0, 255}
	Maroon       = Color{128, 0, 0, 255}
	Lime         = Color{0, 255, 0, 255}
	Teal         = Color{0, 128, 128, 255}
	Lavender     = Color{230, 230, 250, 255}
	Gold         = Color{255, 215, 0, 255}
	Silver       = Color{192, 192, 192, 255}
	Coral        = Color{255, 127, 80, 255}
	Tomato       = Color{255, 99, 71, 255}
	Salmon       = Color{250, 128, 114, 255}
	SkyBlue      = Color{135, 206, 235, 255}
	MidnightBlue = Color{25, 25, 112, 255}
	DarkGreen    = Color{0, 100, 0, 255}
	DarkRed      = Color{139, 0, 0, 255}
	ForestGreen  = Color{34, 139, 34, 255}
	DarkOrange   = Color{255, 140, 0, 255}
	Turquoise    = Color{64, 224, 208, 255}
	DeepPink     = Color{255, 20, 147, 255}
	Indigo       = Color{75, 0, 130, 255}
)

func GetAllColors() []Color {
	return []Color{Red, Green, Blue, White, Black, Yellow, Cyan, Magenta}
}

func GetRGBA(color Color) (red, green, blue, alpha byte) {
	if len(color) != 4 {
		return 0, 0, 0, 0
	}

	red = color[0]
	green = color[1]
	blue = color[2]
	alpha = color[3]
	return
}
