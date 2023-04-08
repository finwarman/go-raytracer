package main

import (
	"image"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

func main() {
	a := app.New()
	w := a.NewWindow("wallpaper")

	c := w.Canvas()

	W, H := 800, 600
	w.Resize(fyne.NewSize(float32(W), float32(H)))

	rect := image.Rect(0, 0, W/5, H/5)
	go func() {
		for i := 0; i < 1000; i++ {
			time.Sleep(time.Millisecond * 50)

			img := createWallpaperImage(rect, float64(i), float64(i), 3.0)
			image := canvas.NewImageFromImage(img)
			image.FillMode = canvas.ImageFillContain
			image.ScaleMode = canvas.ImageScalePixels

			c.SetContent(image)
		}
	}()

	w.ShowAndRun()

}

func createWallpaperImage(rect image.Rectangle, corna float64, cornb float64, side float64) (created *image.NRGBA) {
	pix := make([]uint8, rect.Dx()*rect.Dy()*4)

	for i := 0; i < rect.Dx(); i++ {
		for j := 0; j < rect.Dy()*4; j++ {
			x := corna + float64(i)*side/100
			y := cornb + float64(j)*side/100
			c := int(math.Pow(x, 2) + math.Pow(y, 2))
			if c%2 == 0 {
				pix[rect.Dx()*j+i] = 255
			}
		}
	}

	created = &image.NRGBA{
		Pix:    pix,
		Stride: rect.Dx() * 4,
		Rect:   rect,
	}
	return
}
