package main

import (
	"image"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"

	rt "github.com/finwarman/raytracer/raytracer"
)

func main() {
	a := app.New()
	w := a.NewWindow("wallpaper")

	c := w.Canvas()

	width, height := 1024, 768
	scale := 2.0 // pixels per pixel in image
	w.Resize(fyne.NewSize(float32(width), float32(height)))

	rect := image.Rect(0, 0, int(float64(width)/scale), int(float64(height)/scale))
	img := createImage(rect)

	image := canvas.NewImageFromImage(img)
	image.FillMode = canvas.ImageFillContain
	image.ScaleMode = canvas.ImageScalePixels
	if scale < 1.0 {
		image.ScaleMode = canvas.ImageScaleSmooth
	}

	c.SetContent(image)

	w.ShowAndRun()
}

func createImage(rect image.Rectangle) (img *image.NRGBA) {
	width, height := rect.Dx(), rect.Dy()
	fov := math.Pi / 2.0

	stride := width * 4
	pix := make([]uint8, width*stride)
	img = &image.NRGBA{
		Pix:    pix,
		Stride: stride,
		Rect:   rect,
	}

	spheres := []*rt.Sphere{
		{
			Centre:   rt.Vector3f{X: -3.0, Y: 0.0, Z: -16.0},
			Radius:   2.0,
			Material: rt.Material{Colour: rt.Ivory},
		},
		{
			Centre:   rt.Vector3f{X: -1.0, Y: -1.5, Z: -12.0},
			Radius:   2.0,
			Material: rt.Material{Colour: rt.RedRubber},
		},
		{
			Centre:   rt.Vector3f{X: 1.5, Y: -0.5, Z: -18.0},
			Radius:   3.0,
			Material: rt.Material{Colour: rt.RedRubber},
		},
		{
			Centre:   rt.Vector3f{X: 7.0, Y: 5.0, Z: -18.0},
			Radius:   4.0,
			Material: rt.Material{Colour: rt.Ivory},
		},
	}

	render(img, width, height, fov, spheres)
	return img
}

// TODO: goroutine parallelise rendering
func render(img *image.NRGBA, width, height int, fov float64, spheres []*rt.Sphere) {
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			// TODO: what are these formulae? abstract into functions?
			x := (2*(float64(i)+0.5)/float64(width) - 1) * math.Tan(float64(fov)/2.0) * float64(width) / float64(height)
			y := -1.0 * (2*(float64(j)+0.5)/float64(height) - 1) * math.Tan(float64(fov)/2.0)

			origin := rt.Vector3f{X: 0, Y: 0, Z: 0}
			direction := rt.Vector3f{X: x, Y: y, Z: -1}.Norm()

			img.Set(i, j, castRay(origin, direction, spheres))
		}
	}

}

func castRay(origin, direction rt.Vector3f, spheres []*rt.Sphere) color.NRGBA {
	var point, N rt.Vector3f
	var material rt.Material

	intersected := sceneIntersect(origin, direction, &point, &N, &material, spheres)
	if !intersected {
		return rt.BackgroundColour
	}
	return material.Colour
}

func sceneIntersect(origin, direction rt.Vector3f, hit, N *rt.Vector3f, material *rt.Material, spheres []*rt.Sphere) bool {
	spheresDist := math.MaxFloat64

	for i := 0; i < len(spheres); i++ {
		var dist_i float64
		if spheres[i].RayIntersect(origin, direction, &dist_i) && dist_i < spheresDist {
			spheresDist = dist_i
			*hit = origin.Add(direction.Multiply(dist_i))
			*N = hit.Sub(spheres[i].Centre).Norm()
			*material = spheres[i].Material
		}
	}

	return spheresDist < 1000
}
