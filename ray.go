package raytrace

import (
	"fmt"
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
	nx := 600
	ny := 300
	ns := 100

	world := RandomScene()

	cam := NewCamera(
		Vec3{-2, 2, 1},          // from
		Vec3{0, 0, -1},          // to
		Vec3{0, 1, 0},           // vup
		90,                      // fov
		float64(nx)/float64(ny), // aspec
	)

	m := image.NewRGBA(image.Rect(0, 0, nx, ny))
	total := nx * ny
	count := 0
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
			count++
			if count%100 == 0 {
				fmt.Printf("%d/%d: %f %%\n", count, total, (float64(count)/float64(total))*100)
			}
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

// NewCamera constructs a camera with the provided field of view and aspect ratio
func NewCamera(from, to, vup Vec3, vfov, aspect float64) Camera {
	theta := vfov * math.Pi / 180
	halfheight := math.Tan(theta / 2)
	halfwidth := aspect * halfheight
	w := from.Sub(to).Unit()
	u := vup.Cross(w).Unit()
	v := w.Cross(u)
	return Camera{
		Origin:     from,
		BottomLeft: from.Sub(u.ScalarMul(halfwidth)).Sub(v.ScalarMul(halfheight)).Sub(w),
		Horizontal: u.ScalarMul(2 * halfwidth),
		Vertical:   v.ScalarMul(2 * halfheight),
	}
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

// RandomMaterial generates a random material
func RandomMaterial() Material {
	x := rand.Float64()
	switch {
	case x < 0.8: // diffuse
		return Lambertian{
			Albedo: Vec3{
				rand.Float64() * rand.Float64(),
				rand.Float64() * rand.Float64(),
				rand.Float64() * rand.Float64()},
		}
	case x < 0.95: // metal
		return Metal{
			Albedo: Vec3{
				0.5 * (1 + rand.Float64()),
				0.5 * (1 + rand.Float64()),
				0.5 * rand.Float64(),
			},
		}
	default: // glass
		return Dielectric{RefIndex: 1.5}
	}
}

// RandomScene generates a random Hittable world
func RandomScene() World {
	w := World{
		Sphere{
			Center: Vec3{0, -1000, 0},
			Radius: 1000,
			Material: Lambertian{
				Albedo: Vec3{0.5, 0.5, 0.5},
			},
		},
		Sphere{
			Center:   Vec3{0, 1, 0},
			Radius:   1,
			Material: Dielectric{RefIndex: 1.5},
		},
		Sphere{
			Center: Vec3{-4, 1, 0},
			Radius: 1,
			Material: Lambertian{
				Albedo: Vec3{0.4, 0.2, 0.1},
			},
		},
		Sphere{
			Center: Vec3{0.7, 0.6, 0.5},
			Radius: 1,
			Material: Metal{
				Fuzz:   0,
				Albedo: Vec3{0.7, 0.6, 0.5},
			},
		},
	}
	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			center := Vec3{float64(a) + 0.9 + rand.Float64(), 0.2, float64(b) + 0.9 + rand.Float64()}
			if center.Sub(Vec3{4, 0.2, 0}).Len() > 0.9 {
				w = append(w, Sphere{
					Center:   center,
					Radius:   0.2,
					Material: RandomMaterial(),
				})
			}
		}
	}
	return w
}

// Material scatters rays
type Material interface {
	Scatter(in Ray, h Hit) (out Ray, attenuation Vec3, ok bool)
}

// Lambertian is a material which diffuses light
type Lambertian struct {
	Albedo Vec3
}

// Scatter implements Material
func (l Lambertian) Scatter(in Ray, h Hit) (out Ray, attenuation Vec3, ok bool) {
	out = Ray{
		Origin: h.Pos,
		Dir:    h.Norm.Add(RandomInUnitSphere()),
	}
	return out, l.Albedo, true
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
func (m Metal) Scatter(in Ray, h Hit) (out Ray, attenuation Vec3, ok bool) {
	out = Ray{
		Origin: h.Pos,
		Dir:    m.reflect(in.Dir.Unit(), h.Norm).Add(RandomInUnitSphere().ScalarMul(m.Fuzz)),
	}
	return out, m.Albedo, out.Dir.Dot(h.Norm) > 0
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

// polynomial approximation of glass reflection/refraction
func (d Dielectric) schlick(cosine float64) float64 {
	r0 := (1 - d.RefIndex) / (1 + d.RefIndex)
	r0 = r0 * r0
	return r0 + (1 - r0) + math.Pow(1-cosine, 5)
}

// isReflect decides if the current angle should be reflected or refracted
func (d Dielectric) isReflect(cosine float64) bool {
	return d.schlick(cosine) < rand.Float64()
}

// Scatter implements material
func (d Dielectric) Scatter(in Ray, h Hit) (out Ray, attenuation Vec3, ok bool) {
	var outwardNorm Vec3
	var niOverNt float64
	var cosine float64
	if in.Dir.Dot(h.Norm) > 0 {
		outwardNorm = h.Norm.Neg()
		niOverNt = d.RefIndex
		cosine = d.RefIndex * in.Dir.Dot(h.Norm) / in.Dir.Len()
	} else {
		outwardNorm = h.Norm
		niOverNt = 1.0 / d.RefIndex
		cosine = -in.Dir.Dot(h.Norm) / in.Dir.Len()
	}
	if refracted, ok := d.refract(in.Dir, outwardNorm, niOverNt); ok && !d.isReflect(cosine) {
		out = Ray{
			Origin: h.Pos,
			Dir:    refracted,
		}
	} else {
		out = Ray{
			Origin: h.Pos,
			Dir:    d.reflect(in.Dir, h.Norm),
		}
	}
	return out, Vec3{1, 1, 1}, true
}
