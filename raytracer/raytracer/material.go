package raytracer

import "image/color"

// colour at infinity
var BackgroundColour = FloatToRGB(0.2, 0.7, 0.8, 1.0)

var Ivory = FloatToRGB(0.4, 0.4, 0.4, 1.0)
var RedRubber = FloatToRGB(0.3, 0.1, 0.1, 1.0)

type Material struct {
	Colour color.NRGBA
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
