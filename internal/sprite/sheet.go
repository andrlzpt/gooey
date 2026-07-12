package sprite

import (
	"errors"
	"image"
)

var (
	ErrNilSheet         = errors.New("nil sprite sheet")
	ErrInvalidFrameSize = errors.New("invalid frame size")
	ErrUnsupportedSheet = errors.New("spritesheet does not support sub-image")
	ErrNoFrames         = errors.New("spritesheet contains no complete frames")
)

func SplitSheet(sheet image.Image, frameWidth, frameHeight int) ([]image.Image, error) {
	if sheet == nil {
		return nil, ErrNilSheet
	}

	if frameWidth <= 0 || frameHeight <= 0 {
		return nil, ErrInvalidFrameSize
	}

	subImager, ok := sheet.(interface {
		SubImage(r image.Rectangle) image.Image
	})
	if !ok {
		return nil, ErrUnsupportedSheet
	}

	bounds := sheet.Bounds()
	sheetWidth := bounds.Dx()
	sheetHeight := bounds.Dy()

	columns := sheetWidth / frameWidth
	rows := sheetHeight / frameHeight

	if columns == 0 || rows == 0 {
		return nil, ErrNoFrames
	}

	frames := []image.Image{}

	for row := range rows {
		for col := range columns {
			x0 := bounds.Min.X + col*frameWidth
			y0 := bounds.Min.Y + row*frameHeight

			rect := image.Rect(x0, y0, x0+frameWidth, y0+frameHeight)
			frames = append(frames, subImager.SubImage(rect))
		}
	}

	return frames, nil
}
