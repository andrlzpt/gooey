package window

import "github.com/go-gl/glfw/v3.4/glfw"

type Loop func(state *State)

type State struct {
	MouseX              float64
	MouseY              float64
	IsMouseInsideWindow bool
	DeltaTime           float64
	ElapsedTime         float64
	FrameCount          uint64
}

type Config struct {
	Width  int
	Height int
	Title  string
}

func Run(config Config, loop Loop) {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(config.Width, config.Height, config.Title, nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	state := State{}
	window.SetCursorEnterCallback(func(w *glfw.Window, entered bool) {
		state.IsMouseInsideWindow = entered
	})
	lastTime := glfw.GetTime()
	for !window.ShouldClose() {
		currentTime := glfw.GetTime()
		state.DeltaTime = currentTime - lastTime
		state.ElapsedTime = currentTime
		state.FrameCount++
		lastTime = currentTime
		x, y := window.GetCursorPos()
		state.MouseX = x
		state.MouseY = y
		loop(&state)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
