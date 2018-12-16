package util

import "math"

type Vec2D struct {
	X, Y int
}

func MaxVec2D() Vec2D {
	return Vec2D{math.MaxInt32, math.MaxInt32}
}

func MinVec2D() Vec2D {
	return Vec2D{math.MinInt32, math.MinInt32}
}

func (v Vec2D) Add(o Vec2D) Vec2D {
	result := v
	result.AddInPlace(o)
	return result
}

func (v Vec2D) Sub(o Vec2D) Vec2D {
	result := v
	result.SubInPlace(o)
	return result
}

func (v Vec2D) Scale(s int) Vec2D {
	result := v
	result.X *= s
	result.Y *= s
	return result
}

func (v *Vec2D) AddInPlace(o Vec2D) {
	v.X += o.X
	v.Y += o.Y
}

func (v *Vec2D) SubInPlace(o Vec2D) {
	v.X -= o.X
	v.Y -= o.Y
}

func (v *Vec2D) MinInPlace(o Vec2D) {
	if o.X < v.X { v.X = o.X }
	if o.Y < v.Y { v.Y = o.Y }
}

func (v *Vec2D) MaxInPlace(o Vec2D) {
	if o.X > v.X { v.X = o.X }
	if o.Y > v.Y { v.Y = o.Y }
}

func (v Vec2D) Area() int {
	return AbsInt(v.X) * AbsInt(v.Y)
}

func (v Vec2D) Length() float64 {
	return math.Sqrt(math.Pow(float64(v.X), 2) + math.Pow(float64(v.Y), 2))
}

func (v Vec2D) Manhattan() int {
	return AbsInt(v.X) + AbsInt(v.Y)
}
