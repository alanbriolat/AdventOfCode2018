package util

import "math"

type Vec3D struct {
	X, Y, Z int
}

func MaxVec3D() Vec3D {
	return Vec3D{math.MaxInt32, math.MaxInt32, math.MaxInt32}
}

func MinVec3D() Vec3D {
	return Vec3D{math.MinInt32, math.MinInt32, math.MinInt32}
}

func (v Vec3D) Add(o Vec3D) Vec3D {
	result := v
	result.AddInPlace(o)
	return result
}

func (v Vec3D) Sub(o Vec3D) Vec3D {
	result := v
	result.SubInPlace(o)
	return result
}

func (v Vec3D) Scale(s int) Vec3D {
	result := v
	result.X *= s
	result.Y *= s
	result.Z *= s
	return result
}

func (v *Vec3D) AddInPlace(o Vec3D) {
	v.X += o.X
	v.Y += o.Y
	v.Z += o.Z
}

func (v *Vec3D) SubInPlace(o Vec3D) {
	v.X -= o.X
	v.Y -= o.Y
	v.Z -= o.Z
}

func (v *Vec3D) MinInPlace(o Vec3D) {
	if o.X < v.X { v.X = o.X }
	if o.Y < v.Y { v.Y = o.Y }
	if o.Z < v.Z { v.Z = o.Z }
}

func (v *Vec3D) MaxInPlace(o Vec3D) {
	if o.X > v.X { v.X = o.X }
	if o.Y > v.Y { v.Y = o.Y }
	if o.Z > v.Z { v.Z = o.Z }
}

func (v Vec3D) Volume() int {
	return AbsInt(v.X) * AbsInt(v.Y) * AbsInt(v.Z)
}

func (v Vec3D) Length() float64 {
	return math.Sqrt(math.Pow(float64(v.X), 2) + math.Pow(float64(v.Y), 2) + math.Pow(float64(v.Z), 2))
}

func (v Vec3D) Manhattan() int {
	return AbsInt(v.X) + AbsInt(v.Y) + AbsInt(v.Z)
}
