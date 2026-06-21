package renderer

import (
	"github.com/andrlzpt/gooey/internal/ascii"
	"github.com/go-gl/gl/v3.2-compatibility/gl"
)

type Config struct {
	CellWidth  int
	CellHeight int
}

var initialized = false

func Render(buffer *ascii.Buffer, config Config) {
	if !initialized {
		gl.Init()
		initialized = true
	}
}
