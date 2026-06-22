package main

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"runtime"

	"github.com/andrlzpt/gooey/internal/ascii"
	"github.com/andrlzpt/gooey/internal/renderer"
	"github.com/andrlzpt/gooey/internal/window"
)

const WindowWidth = 640
const WindowHeight = 480
const CellWidth = 4
const CellHeight = 4

const EraseCircleRadius = 48

var windowConfig = window.Config{
	Width:  WindowWidth,
	Height: WindowHeight,
	Title:  "gooey",
}

var renderConfig = renderer.Config{
	WindowWidth:  WindowWidth,
	WindowHeight: WindowHeight,
	CellWidth:    CellWidth,
	CellHeight:   CellHeight,
}

func main() {

	runtime.LockOSThread()
	bufferWidth := WindowWidth / CellWidth
	bufferHeight := WindowHeight / CellHeight
	buffer := ascii.NewBuffer(bufferWidth, bufferHeight)

	file, err := os.Open("test.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	window.Run(windowConfig, func(state *window.State) {
		loop(state, buffer, img)
	})
}

func loop(state *window.State, buffer *ascii.Buffer, img image.Image) {
	buffer.Clear()
	// buffer.FillRandom()
	// ascii.DrawText(buffer, 2, 2, "GOOEY")

	// if state.IsMouseInsideWindow {
	// 	eraseCircle(buffer, state)
	// }
	// glyph := ascii.Glyphs[len(ascii.Glyphs)-1]
	// ascii.DrawRect(buffer, 0, 0, 12, 12, glyph)

	// ascii.FillRect(buffer, 24, 0, 12, 12, glyph)

	// ascii.DrawCircle(buffer, 54, 8, 8, glyph)

	// ascii.FillCircle(buffer, 72, 8, 8, glyph)

	// ascii.DrawImage(buffer, img, 0, 0, buffer.Width, buffer.Height)
	ascii.DrawImage(buffer, img, 0, 0, 60, 40)

	renderer.Render(buffer, renderConfig)
}

func eraseCircle(buffer *ascii.Buffer, state *window.State) {
	mouseX := int(state.MouseX)
	mouseY := int(state.MouseY)

	mouseCellX := mouseX / CellWidth
	mouseCellY := mouseY / CellHeight

	radius := EraseCircleRadius
	radiusSquared := radius * radius

	radiusCellsX := radius/CellWidth + 1
	radiusCellsY := radius/CellHeight + 1

	for y := mouseCellY - radiusCellsY; y <= mouseCellY+radiusCellsY; y++ {
		for x := mouseCellX - radiusCellsX; x <= mouseCellX+radiusCellsX; x++ {
			cellCenterX := x*CellWidth + CellWidth/2
			cellCenterY := y*CellHeight + CellHeight/2

			dx := cellCenterX - mouseX
			dy := cellCenterY - mouseY

			if dx*dx+dy*dy <= radiusSquared {
				buffer.ClearAt(x, y)
			}
		}
	}
}
