package raytracer

import "image/color"

// colour at infinity
var BackgroundColour = FloatToRGB(0.2, 0.7, 0.8)

var Ivory = Material{
	DiffuseColour:    FloatToRGB(0.4, 0.4, 0.3),
	SpecularExponent: 50.0,
	Albedo:           Vector3f{X: 0.6, Y: 0.3, Z: 0.1},
}
var RedRubber = Material{
	DiffuseColour:    FloatToRGB(0.3, 0.1, 0.1),
	SpecularExponent: 10.0,
	Albedo:           Vector3f{X: 0.9, Y: 0.1, Z: 0.0},
}
var Mirror = Material{
	DiffuseColour:    FloatToRGB(1.0, 1.0, 1.0),
	SpecularExponent: 1425.0,
	Albedo:           Vector3f{X: 0.0, Y: 10.0, Z: 0.8},
}

type Material struct {
	DiffuseColour    color.NRGBA
	SpecularExponent float64
	Albedo           Vector3f
}

func FloatToRGB(r, g, b float64) color.NRGBA {
	_r := uint8(r * 0xff)
	_g := uint8(g * 0xff)
	_b := uint8(b * 0xff)
	_a := uint8(0xff)

	return color.NRGBA{
		_r, _g, _b, _a,
	}

}
