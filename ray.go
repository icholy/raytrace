package raytrace

import (
	"image"
	"image/color"
)

// Ray is a ray of light
type Ray struct {
	Origin Vec3
	Dir    Vec3
}

// Pos returns the ray position at t
func (r Ray) Pos(t float64) Vec3 {
	return r.Origin.Add(r.Dir.ScalarMul(t))
}

// BasicColor is the color function described on page 10
func BasicColor(r Ray) Vec3 {
	unitdir := r.Dir.Unit()
	t := 0.5 * (unitdir.Y() + 1)
	return Vec3{1, 1, 1}.ScalarMul(1 - t).Add(Vec3{0.5, 0.7, 1}.ScalarMul(t))
}

// BasicRay outputs the image described on page 10
func BasicRay() image.Image {
	nx := 200
	ny := 100

	bottomleft := Vec3{-2, -1, -1}
	horizontal := Vec3{4, 0, 0}
	vertical := Vec3{0, 2, 0}
	origin := Vec3{0, 0, 0}

	m := image.NewRGBA(image.Rect(0, 0, nx, ny))
	for y := 0; y < ny; y++ {
		for x := 0; x < nx; x++ {
			u := float64(x) / float64(nx)
			v := float64(y) / float64(ny)
			r := Ray{
				Origin: origin,
				Dir:    bottomleft.Add(horizontal.ScalarMul(u)).Add(vertical.ScalarMul(v)),
			}
			c := BasicColor(r)
			m.Set(x, y, color.RGBA{
				R: uint8(c.R() * 255),
				G: uint8(c.G() * 255),
				B: uint8(c.B() * 255),
				A: 0xFF,
			})
		}
	}
	return m
}
