package renderer

import (
	"github.com/andrlzpt/gooey/internal/ascii"
	"github.com/go-gl/gl/v3.2-compatibility/gl"
)

type Config struct {
	WindowWidth  int
	WindowHeight int
	CellWidth    int
	CellHeight   int
}

var initialized = false

func Render(buffer *ascii.Buffer, config Config) {
	if !initialized {
		if err := gl.Init(); err != nil {
			panic(err)
		}
		initialized = true
		gl.ClearColor(0, 0, 0, 0)
	}

	gl.Viewport(0, 0, int32(config.WindowWidth), int32(config.WindowHeight))
	gl.Clear(gl.COLOR_BUFFER_BIT)

	for y := 0; y < buffer.Height; y++ {
		for x := 0; x < buffer.Width; x++ {
			r := buffer.At(x, y)
			brightness := brightnessFor(r)
			drawCell(x, y, brightness, config)
		}

	}
	// drawCell(0, 0, 1, config)
}

func brightnessFor(r rune) float32 {
	brightness, ok := brightnessByRune[r]
	if !ok {
		return 0
	}
	return brightness
}

var brightnessByRune = map[rune]float32{
	' ': 0.0,
	'.': 0.1,
	':': 0.2,
	'-': 0.3,
	'=': 0.4,
	'+': 0.5,
	'*': 0.6,
	'#': 0.75,
	'%': 0.9,
	'@': 1.0,
}

func drawCell(x, y int, brightness float32, config Config) {
	leftPx := float32(x * config.CellWidth)
	rightPx := float32(x+1) * float32(config.CellWidth)
	topPx := float32(y * config.CellHeight)
	bottomPx := float32(y+1) * float32(config.CellHeight)

	left := pixelToClipX(leftPx, config)
	right := pixelToClipX(rightPx, config)
	top := pixelToClipY(topPx, config)
	bottom := pixelToClipY(bottomPx, config)

	drawRect(left, top, right, bottom, brightness)
}

func drawRect(left, top, right, bottom, brightness float32) {
	gl.Color3f(brightness, brightness, brightness)
	gl.Begin(gl.QUADS)
	gl.Vertex2f(left, top)
	gl.Vertex2f(right, top)
	gl.Vertex2f(right, bottom)
	gl.Vertex2f(left, bottom)
	gl.End()
}

func pixelToClipX(x float32, config Config) float32 {
	return (x/float32(config.WindowWidth))*2 - 1
}

func pixelToClipY(y float32, config Config) float32 {
	return 1 - (y/float32(config.WindowHeight))*2
}
