package physics

type Vector struct {
	X float64
	Y float64
}

type ShapeKind int

const (
	ShapePoint ShapeKind = iota
	ShapeRect
	ShapeCircle
	ShapeTriangle
)

type Shape struct {
	Kind   ShapeKind
	Width  int
	Height int
	Radius int
	X2     int
	Y2     int
}

type Body struct {
	Position   Vector
	Velocity   Vector
	Shape      Shape
	Weightless bool
}

type Bounds struct {
	Left   float64
	Right  float64
	Top    float64
	Bottom float64
}

func (b Body) Bounds() Bounds {
	bounds := Bounds{}
	switch b.Shape.Kind {
	case ShapeRect:
		bounds.Left = b.Position.X
		bounds.Right = b.Position.X + float64(b.Shape.Width) - 1
		bounds.Top = b.Position.Y
		bounds.Bottom = b.Position.Y + float64(b.Shape.Height) - 1
	case ShapeCircle:
		bounds.Left = b.Position.X - float64(b.Shape.Radius)
		bounds.Right = b.Position.X + float64(b.Shape.Radius)
		bounds.Top = b.Position.Y - float64(b.Shape.Radius)
		bounds.Bottom = b.Position.Y + float64(b.Shape.Radius)
	case ShapeTriangle:
		bounds.Left = b.Position.X - float64(b.Shape.Width)/2
		bounds.Right = b.Position.X + float64(b.Shape.Width)/2
		bounds.Top = b.Position.Y
		bounds.Bottom = b.Position.Y + float64(b.Shape.Height)
	default:
		bounds.Left = b.Position.X
		bounds.Right = b.Position.X
		bounds.Top = b.Position.Y
		bounds.Bottom = b.Position.Y
	}

	return bounds

}
