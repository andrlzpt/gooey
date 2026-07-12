package sprite

import (
	"errors"
	"image"
)

type Animation struct {
	Frames       []image.Image
	FrameSeconds float64
	Elapsed      float64
	Current      int
	Paused       bool
}

var (
	ErrNilAnimationFrame   = errors.New("animation contains nil frames")
	ErrNoAnimationFrames   = errors.New("animation has no frames")
	ErrInvalidFrameSeconds = errors.New("invalid frame seconds")
)

func NewAnimation(frames []image.Image, frameSeconds float64) (*Animation, error) {
	if frames == nil {
		return nil, ErrNilAnimationFrame
	}

	if len(frames) == 0 {
		return nil, ErrNoAnimationFrames
	}

	for _, frame := range frames {
		if frame == nil {
			return nil, ErrNilAnimationFrame
		}
	}

	if frameSeconds < 0 {
		return nil, ErrInvalidFrameSeconds
	}
	return &Animation{
		Frames:       frames,
		FrameSeconds: frameSeconds,
		Elapsed:      0,
		Current:      0,
	}, nil
}

func (a *Animation) Update(dt float64) {
	if a.Paused {
		return
	}
	if len(a.Frames) == 0 || a.FrameSeconds <= 0 {
		return
	}

	a.Elapsed += dt

	for a.Elapsed >= a.FrameSeconds {
		a.Elapsed -= a.FrameSeconds
		a.Current = (a.Current + 1) % len(a.Frames)
	}
}

func (a *Animation) Frame() image.Image {
	if len(a.Frames) == 0 {
		return nil
	}

	if a.Current < 0 || a.Current >= len(a.Frames) {
		a.Current = 0
	}

	return a.Frames[a.Current]
}

func (a *Animation) TogglePause() {
	a.Paused = !a.Paused
}
