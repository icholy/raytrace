package raytrace

// Ray is a ray of light
type Ray struct {
	Origin Vec3
	Dir    Vec3
}

// Pos returns the ray position at t
func (r Ray) Pos(t float64) Vec3 {
	return r.Origin.Add(r.Dir.ScalarMul(t))
}
