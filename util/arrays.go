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

func (g *GenericGrid) Width() int {
	return len(*g)
}

func (g *GenericGrid) Height() int {
	return len((*g)[0])
}

func (g *GenericGrid) Valid(p Vec2D) bool {
	return p.X >= 0 && p.X < g.Width() && p.Y >= 0 && p.Y < g.Height()
}

func (g *GenericGrid) At(p Vec2D) *Generic {
	return &(*g)[p.X][p.Y]
}

/*
Traverse the grid column-by-column, row-by-row, calling visit with each set of
coordinates and a pointer to the element at those coordinates.
 */
func (g *GenericGrid) Traverse(visit func(p Vec2D, data *Generic)) {
	p := Vec2D{0, 0}
	for p.X = 0; p.X < g.Width(); p.X++ {
		for p.Y = 0; p.Y < g.Height(); p.Y++ {
			visit(p, g.At(p))
		}
	}
}

func (g *GenericGrid) Initialize(value Generic) {
	g.Traverse(func(p Vec2D, data *Generic) {
		*data = value
	})
}
