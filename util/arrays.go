package util

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "Generic=int,byte,bool"

// Already defined in stack.go
//type Generic generic.Type

type GenericGrid [][]Generic

func NewGenericGrid(w, h int) GenericGrid {
	raw := make([]Generic, w*h)
	grid := make([][]Generic, w)
	for x := 0; x < w; x++ {
		grid[x], raw = raw[:h], raw[h:]
	}
	return grid
}

func (g *GenericGrid) Valid(p Vec2D) bool {
	return p.X >= 0 && p.X < g.Width() && p.Y >= 0 && p.Y < g.Height()
}

func (g *GenericGrid) At(p Vec2D) *Generic {
	return &(*g)[p.X][p.Y]
}
