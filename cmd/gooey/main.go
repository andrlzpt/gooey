package main

import (
	"bufio"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"runtime"
	"strings"

	"github.com/andrlzpt/gooey/internal/ascii"
	"github.com/andrlzpt/gooey/internal/physics"
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

const Gravity = 40.0
const Bounce = 0.8

type Command struct {
	Raw string
}

type AppState struct {
	Paused bool
}

func main() {

	runtime.LockOSThread()

	bufferWidth := WindowWidth / CellWidth
	bufferHeight := WindowHeight / CellHeight
	buffer := ascii.NewBuffer(bufferWidth, bufferHeight)

	// file, err := os.Open("test.png")
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()

	// img, _, err := image.Decode(file)
	// if err != nil {
	// 	panic(err)
	// }
	w := physics.New(Gravity, Bounce)
	particle := physics.Body{
		Position: physics.Vector{
			X: float64(buffer.Width / 2),
			Y: 2,
		},
		Velocity: physics.Vector{
			X: 20,
			Y: 0,
		},
		Shape: physics.Shape{
			Kind: physics.ShapePoint,
		},
		Weightless: false,
		Collidable: true,
	}
	w.AddBody(particle)

	rect := physics.Body{
		Position: physics.Vector{
			X: 20,
			Y: 20,
		},
		Velocity: physics.Vector{
			X: 0,
			Y: 0,
		},
		Shape: physics.Shape{
			Kind:   physics.ShapeRect,
			Width:  24,
			Height: 24,
		},
		Weightless: true,
		Collidable: true,
	}
	w.AddBody(rect)

	circle := physics.Body{
		Position: physics.Vector{
			X: 110,
			Y: 12,
		},
		Velocity: physics.Vector{
			X: -15,
			Y: 0,
		},
		Shape: physics.Shape{
			Kind:   physics.ShapeCircle,
			Radius: 8,
		},
		Weightless: false,
		Collidable: true,
	}
	w.AddBody(circle)

	triangle := physics.Body{
		Position: physics.Vector{
			X: 120,
			Y: 15,
		},
		Velocity: physics.Vector{
			X: -15,
			Y: 0,
		},
		Shape: physics.Shape{
			Kind:   physics.ShapeTriangle,
			Width:  21,
			Height: 11,
		},
		Weightless: false,
		Collidable: true,
	}
	w.AddBody(triangle)

	commands := make(chan Command, 16)
	go readCommands(commands)

	window.Run(windowConfig, func(state *window.State) {
		loop(state, buffer, w, commands)
	})
}

func readCommands(commands chan<- Command) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		commands <- Command{Raw: line}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("reading command err = %v", err)
	}

}

func loop(state *window.State, buffer *ascii.Buffer, world *physics.World, commands <-chan Command) {
	buffer.Clear()
	// buffer.FillRandom()
	// ascii.DrawText(buffer, 2, 2, "GOOEY")

	// if state.IsMouseInsideWindow {
	// 	eraseCircle(buffer, state)
	// }

	// ascii.DrawImage(buffer, img, 0, 0, buffer.Width, buffer.Height)
	// ascii.DrawImage(buffer, img, 0, 0, 60, 40)

	handleCommands(commands, world)
	world.Update(state.DeltaTime, buffer.Width, buffer.Height)
	drawWorld(buffer, world)
	renderer.Render(buffer, renderConfig)
}

func drawWorld(buffer *ascii.Buffer, w *physics.World) {
	glyph := ascii.Glyphs[len(ascii.Glyphs)-1]

	for _, body := range w.Bodies {
		x := int(body.Position.X)
		y := int(body.Position.Y)

		switch body.Shape.Kind {
		case physics.ShapePoint:
			ascii.DrawPoint(buffer, x, y, glyph)
		case physics.ShapeRect:
			ascii.DrawRect(buffer, x, y, body.Shape.Width, body.Shape.Height, glyph)
		case physics.ShapeCircle:
			ascii.DrawCircle(buffer, x, y, body.Shape.Radius, glyph)
		case physics.ShapeTriangle:
			ascii.DrawTriangle(buffer, x, y, body.Shape.Width, body.Shape.Height, glyph)
		}
	}
}

func handleCommands(commands <-chan Command, world *physics.World) {
	for {
		select {
		case command := <-commands:
			handleCommand(command, world)
		default:
			return
		}
	}
}

func handleCommand(command Command, world *physics.World) {
	switch command.Raw {
	case "pause": // glyph := ascii.Glyphs[len(ascii.Glyphs)-1]
		world.TogglePause()
	default:
		fmt.Println("unknown command:", command.Raw)
	}
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
