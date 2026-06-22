package ascii

func FillCircle(buffer *Buffer, centerX, centerY, radius int, r rune) {
	radiusSquared := radius * radius

	for y := centerY - radius; y <= centerY+radius; y++ {
		for x := centerX - radius; x <= centerX+radius; x++ {
			dx := x - centerX
			dy := y - centerY

			if dx*dx+dy*dy <= radiusSquared {
				buffer.Set(x, y, r)
			}
		}
	}
}

func DrawCircle(buffer *Buffer, centerX, centerY, radius int, r rune) {
	if radius <= 0 {
		buffer.Set(centerX, centerY, r)
		return
	}

	x := radius
	y := 0
	err := 0

	for x >= y {
		plotCirclePoints(buffer, centerX, centerY, x, y, r)

		y++
		if err <= 0 {
			err += 2*y + 1
		}
		if err > 0 {
			x--
			err -= 2*x + 1
		}
	}
}

func plotCirclePoints(buffer *Buffer, centerX, centerY, x, y int, r rune) {
	buffer.Set(centerX+x, centerY+y, r)
	buffer.Set(centerX+y, centerY+x, r)
	buffer.Set(centerX-y, centerY+x, r)
	buffer.Set(centerX-x, centerY+y, r)
	buffer.Set(centerX-x, centerY-y, r)
	buffer.Set(centerX-y, centerY-x, r)
	buffer.Set(centerX+y, centerY-x, r)
	buffer.Set(centerX+x, centerY-y, r)
}

func FillRect(buffer *Buffer, x, y, width, height int, r rune) {
	for row := y; row < y+height; row++ {
		for col := x; col < x+width; col++ {
			buffer.Set(col, row, r)
		}
	}
}

func DrawRect(buffer *Buffer, x, y, width, height int, r rune) {
	DrawHorizontalLine(buffer, x, y, x+width, r)
	DrawVerticalLine(buffer, x+width, y, y+height, r)
	DrawHorizontalLine(buffer, x, y+height, y+height, r)
	DrawVerticalLine(buffer, x, y, y+height, r)

}

func DrawHorizontalLine(buffer *Buffer, x, y, length int, r rune) {
	for i := 0; i <= length; i++ {
		x := x + i
		buffer.Set(x, y, r)
	}
}

func DrawVerticalLine(buffer *Buffer, x, y, length int, r rune) {
	for i := 0; i <= length; i++ {
		y := y + i
		buffer.Set(x, y, r)
	}
}

func DrawPoint(buffer *Buffer, x, y int, r rune) {
	buffer.Set(x, y, r)
}

func DrawText(buffer *Buffer, x, y int, text string) {
	for offset, r := range text {
		buffer.Set(x+offset, y, r)
	}
}
