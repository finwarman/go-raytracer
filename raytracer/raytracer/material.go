package raytracer

import "image/color"

// colour at infinity
// var BackgroundColour = FloatToRGB(0.2, 0.7, 0.8)
var BackgroundColour = FloatToRGB(0.4, 0.4, 0.4)

var Paper = Material{
	DiffuseColour:    FloatToRGB(0.95, 0.95, 0.95),
	SpecularExponent: 10.0,
	Albedo:           [4]float64{0.5, 0.5, 0.5, 0.0},
	RefractiveIndex:  1.0,
}
var Ivory = Material{
	DiffuseColour:    FloatToRGB(0.4, 0.4, 0.3),
	SpecularExponent: 50.0,
	Albedo:           [4]float64{0.6, 0.3, 0.1, 0.0},
	RefractiveIndex:  1.0,
}
var RedRubber = Material{
	DiffuseColour:    FloatToRGB(0.3, 0.1, 0.1),
	SpecularExponent: 10.0,
	Albedo:           [4]float64{0.9, 0.1, 0.0, 0.0},
	RefractiveIndex:  1.0,
}
var Mirror = Material{
	DiffuseColour:    FloatToRGB(1.0, 1.0, 1.0),
	SpecularExponent: 1425.0,
	Albedo:           [4]float64{0.0, 10.0, 0.8, 0.0},
	RefractiveIndex:  1.0,
}
var Glass = Material{
	DiffuseColour:    FloatToRGB(0.6, 0.7, 0.8),
	SpecularExponent: 125.0,
	Albedo:           [4]float64{0.0, 0.5, 0.1, 0.8},
	RefractiveIndex:  1.5,
}

type Material struct {
	DiffuseColour    color.NRGBA
	SpecularExponent float64
	Albedo           [4]float64
	// todo: struct for albedo describing characteristics?
	RefractiveIndex float64
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
