package ascii

import "math/rand"

type Cell struct {
	Rune rune
}

type Buffer struct {
	Width  int
	Height int
	Cells  []Cell
}

func NewBuffer(width, height int) *Buffer {
	return &Buffer{
		Width:  width,
		Height: height,
		Cells:  make([]Cell, width*height),
	}
}

func (b *Buffer) Set(x, y int, r rune) {
	if !b.inBounds(x, y) {
		return
	}

	b.Cells[b.index(x, y)].Rune = r
}

func (b *Buffer) At(x, y int) rune {
	if !b.inBounds(x, y) {
		return ' '
	}

	r := b.Cells[b.index(x, y)].Rune
	if r == 0 {
		return ' '
	}

	return r
}

func (b *Buffer) ClearAt(x, y int) {
	if !b.inBounds(x, y) {
		return
	}

	b.Cells[b.index(x, y)].Rune = ' '
}

func (b *Buffer) Clear() {
	for i := range b.Cells {
		b.Cells[i] = Cell{}
	}
}

var Glyphs = []rune(" .:-=+*#%@")

func (b *Buffer) FillRandom() {
	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			b.Set(x, y, Glyphs[rand.Intn(len(Glyphs))])
		}
	}
}

func (b *Buffer) inBounds(x, y int) bool {
	return x >= 0 && x < b.Width &&
		y >= 0 && y < b.Height
}

func (b *Buffer) index(x, y int) int {
	return y*b.Width + x
}
