package util

import "math"

type Vec4D struct {
	X, Y, Z, T int
}

func MaxVec4D() Vec4D {
	return Vec4D{math.MaxInt32, math.MaxInt32, math.MaxInt32, math.MaxInt32}
}

func MinVec4D() Vec4D {
	return Vec4D{math.MinInt32, math.MinInt32, math.MinInt32, math.MinInt32}
}

func (v Vec4D) Add(o Vec4D) Vec4D {
	result := v
	result.AddInPlace(o)
	return result
}

func (v Vec4D) Sub(o Vec4D) Vec4D {
	result := v
	result.SubInPlace(o)
	return result
}

func (v Vec4D) Scale(s int) Vec4D {
	result := v
	result.X *= s
	result.Y *= s
	result.Z *= s
	result.T *= s
	return result
}

func (v *Vec4D) AddInPlace(o Vec4D) {
	v.X += o.X
	v.Y += o.Y
	v.Z += o.Z
	v.T += o.T
}

func (v *Vec4D) SubInPlace(o Vec4D) {
	v.X -= o.X
	v.Y -= o.Y
	v.Z -= o.Z
	v.T -= o.T
}

func (v *Vec4D) MinInPlace(o Vec4D) {
	if o.X < v.X { v.X = o.X }
	if o.Y < v.Y { v.Y = o.Y }
	if o.Z < v.Z { v.Z = o.Z }
	if o.T < v.T { v.T = o.T }
}

func (v *Vec4D) MaxInPlace(o Vec4D) {
	if o.X > v.X { v.X = o.X }
	if o.Y > v.Y { v.Y = o.Y }
	if o.Z > v.Z { v.Z = o.Z }
	if o.T > v.T { v.T = o.T }
}

func (v Vec4D) Length() float64 {
	return math.Sqrt(math.Pow(float64(v.X), 2) + math.Pow(float64(v.Y), 2) + math.Pow(float64(v.Z), 2) + math.Pow(float64(v.T), 2))
}

func (v Vec4D) Manhattan() int {
	return AbsInt(v.X) + AbsInt(v.Y) + AbsInt(v.Z) + AbsInt(v.T)
}
