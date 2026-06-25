package world

type World struct {
	Bodies  []Body
	Gravity float64
	Bounce  float64
	Paused  bool
}

func New(gravity, bounce float64) *World {
	return &World{
		Gravity: gravity,
		Bounce:  bounce,
		Bodies:  []Body{},
	}
}

func (w *World) AddBody(body Body) {
	w.Bodies = append(w.Bodies, body)
}

func (w *World) TogglePause() {
	w.Paused = !w.Paused
}

func (w *World) Update(dt float64, width, height int) {
	if w.Paused {
		return
	}

	for i := range w.Bodies {
		body := &w.Bodies[i]

		if !body.Weightless {
			body.Velocity.Y += w.Gravity * dt
		}

		body.Position.X += body.Velocity.X * dt
		body.Position.Y += body.Velocity.Y * dt

		w.bounceBody(body, width, height)

	}
}

func (w *World) bounceBody(body *Body, width, height int) {
	bounds := body.Bounds()

	maxX := float64(width - 1)
	maxY := float64(height - 1)

	if bounds.Left < 0 {
		body.Position.X += -bounds.Left
		body.Velocity.X = -body.Velocity.X * w.Bounce
	}

	if bounds.Right > maxX {
		body.Position.X -= bounds.Right - maxX
		body.Velocity.X = -body.Velocity.X * w.Bounce
	}

	if bounds.Top < 0 {
		body.Position.Y += -bounds.Top
		body.Velocity.Y = -body.Velocity.Y * w.Bounce
	}

	if bounds.Bottom > maxY {
		body.Position.Y -= bounds.Bottom - maxY
		body.Velocity.Y = -body.Velocity.Y * w.Bounce
	}
}
