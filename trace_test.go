package raytrace

import (
	"image/png"
	"os"
	"testing"
)

func TestGradient(t *testing.T) {
	m := Gradient()
	f, err := os.Create("gradient.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, m); err != nil {
		t.Fatal(err)
	}
}
