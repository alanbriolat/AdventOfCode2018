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

