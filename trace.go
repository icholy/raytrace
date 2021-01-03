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
			r := float64(x) / float64(nx)
			g := float64(y) / float64(ny)
			b := 0.2
			m.Set(x, y, color.RGBA{
				R: uint8(r * 255),
				G: uint8(g * 255),
				B: uint8(b * 255),
				A: 0xFF,
			})
		}
	}
	return m
}
