package raytrace

import (
	"image"
	"image/color"
)

// Gradient outputs the gradient image described in the pdf
func Gradient() image.Image {
	nx := 200
	ny := 100
	m := image.NewRGBA(image.Rect(0, 0, nx, ny))
	for y := 0; y < ny; y++ {
		for x := 0; x < nx; x++ {
			v := Vec3{
				float64(x) / float64(nx),
				float64(y) / float64(ny),
				0.2,
			}
			m.Set(x, y, color.RGBA{
				R: uint8(v.R() * 255),
				G: uint8(v.G() * 255),
				B: uint8(v.B() * 255),
				A: 0xFF,
			})
		}
	}
	return m
}
