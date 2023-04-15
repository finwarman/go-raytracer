package main

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

type pixel struct {
	r uint8
	g uint8
	b uint8
	a uint8
}

func main() {
	a := app.New()
	w := a.NewWindow("wallpaper")

	c := w.Canvas()

	width, height := 1024, 768
	scale := 1 // pixels per pixel in image
	w.Resize(fyne.NewSize(float32(width), float32(height)))

	rect := image.Rect(0, 0, width/scale, height/scale)
	img := createImage(rect)

	image := canvas.NewImageFromImage(img)
	image.FillMode = canvas.ImageFillContain
	image.ScaleMode = canvas.ImageScalePixels

	c.SetContent(image)

	w.ShowAndRun()
}

func createImage(rect image.Rectangle) (created *image.NRGBA) {
	width, height := rect.Dx(), rect.Dy()
	stride := width * 4

	pix := make([]uint8, width*stride)
	created = &image.NRGBA{
		Pix:    pix,
		Stride: stride,
		Rect:   rect,
	}

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			created.Set(i, j, &color.RGBA{
				R: uint8(255 * j / height),
				G: uint8(255 * i / width),
				B: 0x00,
				A: 0xff,
			})
		}
	}

	return created
}
