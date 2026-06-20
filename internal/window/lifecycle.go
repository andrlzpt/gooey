package window

import "github.com/go-gl/glfw/v3.4/glfw"

type Loop func(state *State)

type State struct {
	MouseX              float64
	MouseY              float64
	IsMouseInsideWindow bool
}

func Run(loop Loop) {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(640, 480, "gooey", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	state := State{}
	window.SetCursorEnterCallback(func(w *glfw.Window, entered bool) {
		state.IsMouseInsideWindow = entered
	})
	for !window.ShouldClose() {
		x, y := window.GetCursorPos()
		state.MouseX = x
		state.MouseY = y
		loop(&state)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
