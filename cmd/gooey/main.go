package main

import (
	"runtime"

	"github.com/andrlzpt/gooey/internal/ascii"
	"github.com/andrlzpt/gooey/internal/renderer"
	"github.com/andrlzpt/gooey/internal/window"
)

const WindowWidth = 640
const WindowHeight = 480
const CellWidth = 8
const CellHeight = 12

var windowConfig = window.Config{
	Width:  WindowWidth,
	Height: WindowHeight,
	Title:  "gooey",
}

var renderConfig = renderer.Config{
	CellWidth:  CellWidth,
	CellHeight: CellHeight,
}

func main() {

	runtime.LockOSThread()
	bufferWidth := WindowWidth / CellWidth
	bufferHeight := WindowHeight / CellHeight
	buffer := ascii.NewBuffer(bufferWidth, bufferHeight)
	window.Run(windowConfig, func(state *window.State) {
		loop(state, buffer)
	})
}

func loop(state *window.State, buffer *ascii.Buffer) {
	buffer.Clear()
	buffer.FillRandom()
	renderer.Render(buffer, renderConfig)
}
