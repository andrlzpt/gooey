package main

import (
	"bufio"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"runtime"
	"strings"

	"github.com/andrlzpt/gooey/internal/ascii"
	"github.com/andrlzpt/gooey/internal/physics"
	"github.com/andrlzpt/gooey/internal/renderer"
	"github.com/andrlzpt/gooey/internal/sprite"
	"github.com/andrlzpt/gooey/internal/window"
)

const WindowWidth = 800
const WindowHeight = 600
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

type Entity struct {
	BodyIndex int
	Sprite    sprite.Animation
}

func main() {

	runtime.LockOSThread()

	bufferWidth := WindowWidth / CellWidth
	bufferHeight := WindowHeight / CellHeight
	buffer := ascii.NewBuffer(bufferWidth, bufferHeight)

	file, err := os.Open("sheet1.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	frames, err := sprite.SplitSheet(img, 64, 128)
	if err != nil {
		panic(err)
	}
	// frames := []image.Image{img}
	sprite, err := sprite.NewAnimation(frames, 0.25)
	if err != nil {
		panic(err)
	}

	w := physics.New(Gravity, Bounce)

	// particle := physics.Body{
	// 	Position: physics.Vector{
	// 		X: float64(buffer.Width / 2),
	// 		Y: 2,
	// 	},
	// 	Velocity: physics.Vector{
	// 		X: 20,
	// 		Y: 0,
	// 	},
	// 	Shape: physics.Shape{
	// 		Kind: physics.ShapePoint,
	// 	},
	// 	Weightless: false,
	// 	Collidable: true,
	// }
	// w.AddBody(particle)

	// rect := physics.Body{
	// 	Position: physics.Vector{
	// 		X: 20,
	// 		Y: 20,
	// 	},
	// 	Velocity: physics.Vector{
	// 		X: 0,
	// 		Y: 0,
	// 	},
	// 	Shape: physics.Shape{
	// 		Kind:   physics.ShapeRect,
	// 		Width:  24,
	// 		Height: 24,
	// 	},
	// 	Weightless: true,
	// 	Collidable: true,
	// }
	// w.AddBody(rect)

	// circle := physics.Body{
	// 	Position: physics.Vector{
	// 		X: 110,
	// 		Y: 12,
	// 	},
	// 	Velocity: physics.Vector{
	// 		X: -15,
	// 		Y: 0,
	// 	},
	// 	Shape: physics.Shape{
	// 		Kind:   physics.ShapeCircle,
	// 		Radius: 8,
	// 	},
	// 	Weightless: false,
	// 	Collidable: true,
	// }
	// w.AddBody(circle)

	// triangle := physics.Body{
	// 	Position: physics.Vector{
	// 		X: 120,
	// 		Y: 15,
	// 	},
	// 	Velocity: physics.Vector{
	// 		X: -15,
	// 		Y: 0,
	// 	},
	// 	Shape: physics.Shape{
	// 		Kind:   physics.ShapeTriangle,
	// 		Width:  21,
	// 		Height: 11,
	// 	},
	// 	Weightless: false,
	// 	Collidable: true,
	// }
	// w.AddBody(triangle)

	spriteBody := physics.Body{
		Position: physics.Vector{
			X: float64(buffer.Width / 2),
			Y: float64(buffer.Height / 2),
		},
		Velocity: physics.Vector{
			X: 0,
			Y: 0,
		},
		Shape: physics.Shape{
			Kind:   physics.ShapeRect,
			Width:  64,
			Height: 128,
		},
		Weightless: true,
		Collidable: true,
	}

	spriteBodyIndex := len(w.Bodies)
	w.AddBody(spriteBody)

	entity := Entity{
		BodyIndex: spriteBodyIndex,
		Sprite:    *sprite,
	}

	entities := []Entity{
		entity,
	}

	commands := make(chan Command, 16)
	go readCommands(commands)

	window.Run(windowConfig, func(state *window.State) {
		loop(state, buffer, w, entities, commands)
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

func loop(state *window.State, buffer *ascii.Buffer, world *physics.World, entities []Entity, commands <-chan Command) {
	buffer.Clear()
	// buffer.FillRandom()
	// ascii.DrawText(buffer, 2, 2, "GOOEY")

	// if state.IsMouseInsideWindow {
	// 	eraseCircle(buffer, state)
	// }

	// ascii.DrawImage(buffer, img, 0, 0, buffer.Width, buffer.Height)
	// ascii.DrawImage(buffer, img, 0, 0, 60, 40)

	handleCommands(commands, entities, world)
	world.Update(state.DeltaTime, buffer.Width, buffer.Height)

	for i := range entities {
		entities[i].Sprite.Update(state.DeltaTime)
	}

	drawWorld(buffer, world)
	drawEntities(buffer, world, entities)

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

func drawEntities(buffer *ascii.Buffer, w *physics.World, entities []Entity) {
	for i := range entities {
		entity := &entities[i]

		if entity.BodyIndex < 0 || entity.BodyIndex >= len(w.Bodies) {
			continue
		}

		body := w.Bodies[entity.BodyIndex]
		frame := entity.Sprite.Frame()
		if frame == nil {
			continue
		}

		ascii.DrawImage(
			buffer,
			frame,
			int(body.Position.X),
			int(body.Position.Y),
			body.Shape.Width,
			body.Shape.Height,
		)
	}
}

func handleCommands(commands <-chan Command, entities []Entity, world *physics.World) {
	for {
		select {
		case command := <-commands:
			handleCommand(command, entities, world)
		default:
			return
		}
	}
}

func handleCommand(command Command, entities []Entity, world *physics.World) {
	switch command.Raw {
	case "pause": // glyph := ascii.Glyphs[len(ascii.Glyphs)-1]
		world.TogglePause()
		for i := range entities {
			entities[i].Sprite.TogglePause()
		}
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
