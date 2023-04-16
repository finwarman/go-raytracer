package raytracer

import "math"

type Sphere struct {
	Centre   Vector3f
	Radius   float64
	Material Material
}

// check if a given ray (originating from origin, with direction) intersects with sphere
func (s *Sphere) RayIntersect(origin Vector3f, direction Vector3f, t0 *float64) bool {
	l := s.Centre.Sub(origin)
	tca := direction.Dot(l)
	d2 := l.Dot(l) - (tca * tca)

	r2 := s.Radius * s.Radius
	if d2 > r2 {
		return false
	}
	thc := math.Sqrt(r2 - d2)

	*t0 = tca - thc // float ptr?
	t1 := tca + thc
	if *t0 < 0.0 {
		*t0 = t1
	}
	return *t0 > 0.0
}
