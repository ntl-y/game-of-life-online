package client

type Color []byte

var (
	Black        = Color{0, 0, 0, 1}
	White        = Color{255, 255, 255, 1}
	Red          = Color{255, 0, 0, 1}
	Green        = Color{0, 255, 0, 1}
	Blue         = Color{0, 0, 255, 1}
	Yellow       = Color{255, 255, 0, 1}
	Cyan         = Color{0, 255, 255, 1}
	Magenta      = Color{255, 0, 255, 1}
	Orange       = Color{255, 165, 0, 1}
	Pink         = Color{255, 192, 203, 1}
	Purple       = Color{128, 0, 128, 1}
	Brown        = Color{165, 42, 42, 1}
	Gray         = Color{128, 128, 128, 1}
	LightGray    = Color{211, 211, 211, 1}
	DarkGray     = Color{169, 169, 169, 1}
	Navy         = Color{0, 0, 128, 1}
	Olive        = Color{128, 128, 0, 1}
	Maroon       = Color{128, 0, 0, 1}
	Lime         = Color{0, 255, 0, 1}
	Teal         = Color{0, 128, 128, 1}
	Lavender     = Color{230, 230, 250, 1}
	Gold         = Color{255, 215, 0, 1}
	Silver       = Color{192, 192, 192, 1}
	Coral        = Color{255, 127, 80, 1}
	Tomato       = Color{255, 99, 71, 1}
	Salmon       = Color{250, 128, 114, 1}
	SkyBlue      = Color{135, 206, 235, 1}
	MidnightBlue = Color{25, 25, 112, 1}
	DarkGreen    = Color{0, 100, 0, 1}
	DarkRed      = Color{139, 0, 0, 1}
	ForestGreen  = Color{34, 139, 34, 1}
	DarkOrange   = Color{255, 140, 0, 1}
	Turquoise    = Color{64, 224, 208, 1}
	DeepPink     = Color{255, 20, 147, 1}
	Indigo       = Color{75, 0, 130, 1}
)

func GetAllColors() []Color {
	return []Color{
		Black, Red, Green, Blue, White, Yellow, Cyan, Magenta,
		Orange, Pink, Purple, Brown, Gray, LightGray, DarkGray,
		Navy, Olive, Maroon, Lime, Teal, Lavender, Gold, Silver,
		Coral, Tomato, Salmon, SkyBlue, MidnightBlue, DarkGreen,
		DarkRed, ForestGreen, DarkOrange, Turquoise, DeepPink, Indigo,
	}
}
