package raytrace

import (
	"fmt"
	"math"
)

// Vec3 is a 3 dimentional vector
type Vec3 [3]float64

// String returns a string representation of v
func (v Vec3) String() string { return fmt.Sprintf("(%f, %f, %f)", v[0], v[1], v[2]) }

// X returns the x component
func (v Vec3) X() float64 { return v[0] }

// Y returns the y component
func (v Vec3) Y() float64 { return v[1] }

// Z returns the z component
func (v Vec3) Z() float64 { return v[2] }

// R returns the red component
func (v Vec3) R() float64 { return v[0] }

// G returns the green component
func (v Vec3) G() float64 { return v[1] }

// B returns the blue component
func (v Vec3) B() float64 { return v[2] }

// Neg returns negated v
func (v Vec3) Neg() Vec3 { return Vec3{-v[0], -v[1], -v[2]} }

// Add returns v added to p
func (v Vec3) Add(p Vec3) Vec3 { return Vec3{v[0] + p[0], v[1] + p[1], v[2] + p[2]} }

// Sub returns a p substracted from v
func (v Vec3) Sub(p Vec3) Vec3 { return Vec3{v[0] - p[0], v[1] - p[1], v[2] - p[2]} }

// Mul returns v multiplied by p
func (v Vec3) Mul(p Vec3) Vec3 { return Vec3{v[0] * p[0], v[1] * p[1], v[2] * p[2]} }

// Div returns v divided by p
func (v Vec3) Div(p Vec3) Vec3 { return Vec3{v[0] / p[0], v[1] / p[1], v[2] / p[2]} }

// ScalarMul returns v multiplied by x
func (v Vec3) ScalarMul(x float64) Vec3 { return Vec3{v[0] * x, v[1] * x, v[2] * x} }

// ScalarDiv returns v divided by x
func (v Vec3) ScalarDiv(x float64) Vec3 { return Vec3{v[0] / x, v[1] / x, v[2] / x} }

// ScalarAdd returns x added to v
func (v Vec3) ScalarAdd(x float64) Vec3 { return Vec3{v[0] + x, v[1] + x, v[2] + x} }

// ScalarSub returns x subtracted from v
func (v Vec3) ScalarSub(x float64) Vec3 { return Vec3{v[0] - x, v[1] - x, v[2] - x} }

// Len returns the vector length
func (v Vec3) Len() float64 { return math.Sqrt(v.SquareLen()) }

// SquareLen returns the length squared
func (v Vec3) SquareLen() float64 { return v[0]*v[0] + v[1]*v[1] + v[2]*v[2] }

// Unit returns the unit vector of v
func (v Vec3) Unit() Vec3 { return v.ScalarDiv(v.Len()) }

// Dot returns the dot product of v and p
func (v Vec3) Dot(p Vec3) float64 { return v[0]*p[0] + v[1]*p[1] + v[2]*p[2] }

// Cross returns the cross product of v and p
func (v Vec3) Cross(p Vec3) Vec3 {
	return Vec3{
		v[1]*p[2] - v[2]*p[1],
		v[2]*p[0] - v[0]*p[2],
		v[0]*p[1] - v[1]*p[0],
	}
}
