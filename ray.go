package raytrace

import (
	"image"
	"image/color"
	"math"
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
func BasicColor(r Ray, h Hitter) Vec3 {
	// find any hits
	if hit := h.Hit(r, 0, 10); hit.Valid {
		n := hit.Pos.Sub(Vec3{0, 0, -1}).Unit()
		return n.ScalarAdd(1).ScalarMul(0.5)
	}
	// show background
	t := 0.5 * (r.Dir.Unit().Y() + 1)
	return Vec3{1, 1, 1}.ScalarMul(1 - t).Add(Vec3{0.5, 0.7, 1}.ScalarMul(t))
}

// BasicRay outputs the image described on page 10
func BasicRay() image.Image {
	nx := 200
	ny := 100

	world := World{
		Sphere{
			Center: Vec3{0, 0, -1},
			Radius: 0.5,
		},
	}

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
			c := BasicColor(r, world)
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

// Sphere hitter
type Sphere struct {
	Center Vec3
	Radius float64
}

// Hit implements Hitter
func (s Sphere) Hit(r Ray, tmin, tmax float64) Hit {
	oc := r.Origin.Sub(s.Center)
	a := r.Dir.Dot(r.Dir)
	b := 2 * oc.Dot(r.Dir)
	c := oc.Dot(oc) - s.Radius*s.Radius
	discriminant := b*b - a*c
	if discriminant <= 0 {
		return Hit{}
	}
	if t := (-b - math.Sqrt(b*b-a*c)) / a; tmin <= t && t <= tmax {
		pos := r.Pos(t)
		return Hit{
			Valid: true,
			T:     t,
			Pos:   pos,
			Norm:  pos.Sub(s.Center).ScalarDiv(s.Radius),
		}
	}
	if t := (-b + math.Sqrt(b*b-a*c)) / a; tmin <= t && t <= tmax {
		pos := r.Pos(t)
		return Hit{
			Valid: true,
			T:     t,
			Pos:   pos,
			Norm:  pos.Sub(s.Center).ScalarDiv(s.Radius),
		}
	}
	return Hit{}
}

// Hit is the result of a hit
type Hit struct {
	Valid bool
	T     float64
	Pos   Vec3
	Norm  Vec3
}

// Hitter is an object in the scene which can be hit
type Hitter interface {

	// Hit checks and reports whether r hit ths object
	Hit(r Ray, tmin, tmax float64) Hit
}

// World providers a hitter implementation for multiple underlying hitters
type World []Hitter

// Hit implements Hitter
func (w World) Hit(r Ray, tmin, tmax float64) Hit {
	var closest Hit
	for _, h := range w {
		if hit := h.Hit(r, tmin, tmax); hit.Valid {
			if !closest.Valid || hit.T < closest.T {
				closest = hit
			}
		}
	}
	return closest
}
