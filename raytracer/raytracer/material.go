package raytracer

import "image/color"

// colour at infinity
var BackgroundColour = FloatToRGB(0.2, 0.7, 0.8, 1.0)

var Ivory = Material{
	DiffuseColour:    FloatToRGB(0.4, 0.4, 0.4, 1.0),
	SpecularExponent: 50.0,
	Albedo:           Vector2f{X: 0.6, Y: 0.3},
}
var RedRubber = Material{
	DiffuseColour:    FloatToRGB(0.3, 0.1, 0.1, 1.0),
	SpecularExponent: 10.0,
	Albedo:           Vector2f{X: 0.9, Y: 0.1},
}

type Material struct {
	DiffuseColour    color.NRGBA
	SpecularExponent float64
	Albedo           Vector2f
}

func FloatToRGB(r, g, b, a float64) color.NRGBA {
	_r := uint8(r * 255.0)
	_g := uint8(g * 255.0)
	_b := uint8(b * 255.0)
	_a := uint8(a * 255.0)

	return color.NRGBA{
		_r, _g, _b, _a,
	}

}
