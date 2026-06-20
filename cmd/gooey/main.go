package main

import (
	"fmt"
	"runtime"

	"github.com/andrlzpt/gooey/internal/window"
)

func main() {
	runtime.LockOSThread()
	window.Run(loop)
}

func loop(state *window.State) {
	fmt.Printf("State: %#v\n", state)
	if state.IsMouseInsideWindow {
		fmt.Println("Mouse inside the window.")
	}
}
