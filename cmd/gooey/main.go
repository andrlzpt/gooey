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

type Particle struct {
	X  float64
	Y  float64
	VX float64
	VY float64
}

const Gravity = 40.0
const Bounce = 0.8

type Command struct {
	Raw string
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

	particle := Particle{
		X:  float64(buffer.Width / 2),
		Y:  2,
		VX: 20,
		VY: 0,
	}

	commands := make(chan Command, 16)
	go readCommands(commands)

	window.Run(windowConfig, func(state *window.State) {
		loop(state, buffer, &particle, commands)
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

func loop(state *window.State, buffer *ascii.Buffer, particle *Particle, commands <-chan Command) {
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
	// ascii.DrawImage(buffer, img, 0, 0, 60, 40)

	handleCommands(commands, particle)

	particle.VY += Gravity * state.DeltaTime
	particle.X += particle.VX * state.DeltaTime
	particle.Y += particle.VY * state.DeltaTime

	maxX := float64(buffer.Width - 1)
	maxY := float64(buffer.Height - 1)

	if particle.X < 0 {
		particle.X = 0
		particle.VX = -particle.VX * Bounce
	}

	if particle.X > maxX {
		particle.X = maxX
		particle.VX = -particle.VX * Bounce
	}

	if particle.Y < 0 {
		particle.Y = 0
		particle.VY = -particle.VY * Bounce
	}

	if particle.Y > maxY {
		particle.Y = maxY
		particle.VY = -particle.VY * Bounce
	}
	ascii.DrawPoint(buffer, int(particle.X), int(particle.Y), '@')
	renderer.Render(buffer, renderConfig)
}

func handleCommands(commands <-chan Command, particle *Particle) {
	for {
		select {
		case command := <-commands:
			handleCommand(command, particle)
		default:
			return
		}
	}
}

func handleCommand(command Command, particle *Particle) {
	switch command.Raw {
	case "reset":
		particle.X = 160
		particle.Y = 2
		particle.VX = 20
		particle.VY = 0
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
