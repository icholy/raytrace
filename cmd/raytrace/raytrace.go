package main

import (
	"image/png"
	"log"
	"os"

	"github.com/icholy/raytrace"
)

func main() {
	m := raytrace.BasicRay()
	f, err := os.Create("output.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, m); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
