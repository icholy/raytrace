package raytrace

import (
	"image"
	"image/color"
	"math"
	"math/rand"
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
func BasicColor(r Ray, h Hitter, depth int) Vec3 {
	// find any hits
	if hit := h.Hit(r, 0.001, math.MaxFloat64); hit.Valid {
		if scattered, attenuation, ok := hit.Mat.Scatter(r, hit); ok && depth <= 50 {
			return attenuation.Mul(BasicColor(scattered, h, depth+1))
		}
		// end of recursion
		return Vec3{}
	}
	// background gradient
	t := 0.5 * (r.Dir.Unit().Y() + 1)
	return Vec3{1, 1, 1}.ScalarMul(1 - t).Add(Vec3{0.5, 0.7, 1}.ScalarMul(t))
}

// BasicRay outputs the image described on page 10
func BasicRay() image.Image {
	nx := 200
	ny := 100
	ns := 100

	world := World{
		Sphere{
			Center: Vec3{0, 0, -1},
			Radius: 0.5,
			Material: Lambertian{
				Albedo: Vec3{0.1, 0.2, 0.5},
			},
		},
		Sphere{
			Center: Vec3{0, -100.5, -1},
			Radius: 100,
			Material: Lambertian{
				Albedo: Vec3{0.8, 0.8, 0},
			},
		},
		Sphere{
			Center: Vec3{1, 0, -1},
			Radius: 0.5,
			Material: Metal{
				Fuzz:   0.3,
				Albedo: Vec3{0.8, 0.6, 0.2},
			},
		},
		Sphere{
			Center: Vec3{-1, 0, -1},
			Radius: 0.5,
			Material: Dielectric{
				RefIndex: 1.5,
			},
		},
	}

	cam := Camera{
		BottomLeft: Vec3{-2, -1, -1},
		Horizontal: Vec3{4, 0, 0},
		Vertical:   Vec3{0, 2, 0},
		Origin:     Vec3{0, 0, 0},
	}

	_ = rand.Float64()

	m := image.NewRGBA(image.Rect(0, 0, nx, ny))
	for y := 0; y < ny; y++ {
		for x := 0; x < nx; x++ {
			var col Vec3
			for i := 0; i < ns; i++ {
				u := (float64(x) + rand.Float64()) / float64(nx)
				v := (float64(ny-y) + rand.Float64()) / float64(ny)
				r := cam.Ray(u, v)
				col = col.Add(BasicColor(r, world, 0))
			}
			col = col.ScalarDiv(float64(ns))
			m.Set(x, y, color.RGBA{
				R: uint8(col.R() * 255),
				G: uint8(col.G() * 255),
				B: uint8(col.B() * 255),
				A: 0xFF,
			})
		}
	}
	return m
}

// Sphere hitter
type Sphere struct {
	Center   Vec3
	Radius   float64
	Material Material
}

// Hit implements Hitter
func (s Sphere) Hit(r Ray, tmin, tmax float64) Hit {
	oc := r.Origin.Sub(s.Center)
	a := r.Dir.Dot(r.Dir)
	b := oc.Dot(r.Dir)
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
			Mat:   s.Material,
		}
	}
	if t := (-b + math.Sqrt(b*b-a*c)) / a; tmin <= t && t <= tmax {
		pos := r.Pos(t)
		return Hit{
			Valid: true,
			T:     t,
			Pos:   pos,
			Norm:  pos.Sub(s.Center).ScalarDiv(s.Radius),
			Mat:   s.Material,
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
	Mat   Material
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

// Camera is the users point of view
type Camera struct {
	Origin     Vec3
	BottomLeft Vec3
	Horizontal Vec3
	Vertical   Vec3
}

// Ray returns a ray for the provided u/v coordinates
func (c Camera) Ray(u, v float64) Ray {
	return Ray{
		Origin: c.Origin,
		Dir:    c.BottomLeft.Add(c.Horizontal.ScalarMul(u)).Add(c.Vertical.ScalarMul(v)).Sub(c.Origin),
	}
}

// RandomInUnitSphere doesn't look very efficient
func RandomInUnitSphere() Vec3 {
	for {
		p := Vec3{rand.Float64(), rand.Float64(), rand.Float64()}.ScalarMul(2).ScalarSub(1)
		if p.SquareLen() < 1 {
			return p
		}
	}
}

// Material scatters rays
type Material interface {
	Scatter(r Ray, h Hit) (scattered Ray, attenuation Vec3, ok bool)
}

// Lambertian is a material which diffuses light
type Lambertian struct {
	Albedo Vec3
}

// Scatter implements Material
func (l Lambertian) Scatter(r Ray, h Hit) (scattered Ray, attenuation Vec3, ok bool) {
	scattered = Ray{
		Origin: h.Pos,
		Dir:    h.Norm.Add(RandomInUnitSphere()),
	}
	return scattered, l.Albedo, true
}

// Metal is a material which reflects light
type Metal struct {
	Albedo Vec3
	Fuzz   float64
}

// reflect the ray
func (Metal) reflect(v, n Vec3) Vec3 {
	return v.Sub(n.ScalarMul(2 * v.Dot(n)))
}

// Scatter implements Material
func (m Metal) Scatter(r Ray, h Hit) (scattered Ray, attenuation Vec3, ok bool) {
	scattered = Ray{
		Origin: h.Pos,
		Dir:    m.reflect(r.Dir.Unit(), h.Norm).Add(RandomInUnitSphere().ScalarMul(m.Fuzz)),
	}
	return scattered, m.Albedo, scattered.Dir.Dot(h.Norm) > 0
}

// Dielectric is a material which both reflects and refracts light
type Dielectric struct {
	RefIndex float64
}

// reflect the ray
func (Dielectric) reflect(v, n Vec3) Vec3 {
	return v.Sub(n.ScalarMul(2 * v.Dot(n)))
}

// refract the ray
func (Dielectric) refract(v, n Vec3, niOverNt float64) (Vec3, bool) {
	uv := v.Unit()
	dt := uv.Dot(n)
	discriminant := 1 - niOverNt*niOverNt*(1-dt*dt)
	if discriminant > 0 {
		return uv.Sub(n.ScalarMul(dt)).ScalarMul(niOverNt).Sub(n.ScalarMul(math.Sqrt(discriminant))), true
	}
	return Vec3{}, false
}

// Scatter implements material
func (d Dielectric) Scatter(r Ray, h Hit) (scattered Ray, attenuation Vec3, ok bool) {
	var outwardNorm Vec3
	var niOverNt float64
	attenuation = Vec3{1, 1, 0}
	if r.Dir.Dot(h.Norm) > 0 {
		outwardNorm = h.Norm.Neg()
		niOverNt = d.RefIndex
	} else {
		outwardNorm = h.Norm
		niOverNt = 1.0 / d.RefIndex
	}
	if refracted, ok := d.refract(r.Dir, outwardNorm, niOverNt); ok {
		scattered = Ray{
			Origin: h.Pos,
			Dir:    refracted,
		}
	} else {
		scattered = Ray{
			Origin: h.Pos,
			Dir:    d.reflect(r.Dir, h.Norm),
		}
		return scattered, attenuation, false
	}
	return scattered, attenuation, true
}
