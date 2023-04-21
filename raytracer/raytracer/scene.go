package raytracer

import "image"

type Scene struct {
	EnvMap  *image.NRGBA // skybox image
	Lights  []*Light
	Spheres []*Sphere
}
