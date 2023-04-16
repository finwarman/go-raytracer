package raytracer

import "image/color"

type Light struct {
	Position  Vector3f
	Intensity float64
	Colour    color.NRGBA // TODO: unused
}
