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

	width, height := 800, 600
	w.Resize(fyne.NewSize(float32(width), float32(height)))

	go func() {
		for i := 100.0; i >= 0.0; i = math.Mod((i + 0.005), 0xffff) {
			size := w.Content().Size()
			width, height := int(size.Width), int(size.Height)
			rect := image.Rect(0, 0, width/4, height/4)

			img := createWallpaperImage(rect, float64(i), float64(i*2), float64(i*3))
			image := canvas.NewImageFromImage(img)
			image.FillMode = canvas.ImageFillContain
			image.ScaleMode = canvas.ImageScalePixels

			c.SetContent(image)
			time.Sleep(time.Millisecond * 50)
		}
	}()

	w.ShowAndRun()

}

func createWallpaperImage(rect image.Rectangle, corna float64, cornb float64, side float64) (created *image.NRGBA) {
	pix := make([]uint8, rect.Dx()*rect.Dy()*4)
	stride := rect.Dx() * 4

	for i := 0; i < rect.Dx()*4; i += 4 {
		for j := 0; j < rect.Dy(); j++ {
			x := corna + float64(i)*side/float64(rect.Dx())
			y := cornb + float64(j)*side/float64(rect.Dy())
			c := int(math.Pow(x, 2) + math.Pow(y, 2))
			if c%4 == 0 {
				pix[(stride*j+(i+(0)))%len(pix)] = 0xf0 // R
				pix[(stride*j+(i+(1)))%len(pix)] = 0xa0 // G
				pix[(stride*j+(i+(2)))%len(pix)] = 0x50 // B
				pix[(stride*j+(i+(3)))%len(pix)] = 0xff // A
			} else if c%2 == 0 {
				for x := 0; x < 4; x++ {
					pix[(stride*j+(i+(x)))%len(pix)] = 0xf0
				}
			} else {
				for x := 0; x < 4; x++ {
					pix[(stride*j+(i+(x)))%len(pix)] = 0x10
				}
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
