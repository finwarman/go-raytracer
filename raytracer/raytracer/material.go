package raytracer

import "image/color"

// colour at infinity
var BackgroundColour = FloatToRGB(0.2, 0.7, 0.8)

var Ivory = Material{
	DiffuseColour:    FloatToRGB(0.4, 0.4, 0.3),
	SpecularExponent: 50.0,
	Albedo:           Vector2f{X: 0.6, Y: 0.3},
}
var RedRubber = Material{
	DiffuseColour:    FloatToRGB(0.3, 0.1, 0.1),
	SpecularExponent: 10.0,
	Albedo:           Vector2f{X: 0.9, Y: 0.1},
}

type Material struct {
	DiffuseColour    color.NRGBA
	SpecularExponent float64
	Albedo           Vector2f
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
