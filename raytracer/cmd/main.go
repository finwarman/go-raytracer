package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	rt "github.com/finwarman/raytracer/raytracer"
)

const MaxRayRecursionDepth = 4

func main() {
	a := app.New()
	w := a.NewWindow("wallpaper")

	c := w.Canvas()

	width, height := 1024, 768
	scale := 4.0 // pixels per pixel in image
	w.Resize(fyne.NewSize(float32(width), float32(height)))

	rect := image.Rect(0, 0, int(float64(width)/scale), int(float64(height)/scale))

	targetFPS := 30 // target framerate
	frametime := time.Duration(1000.0 / targetFPS)

	// set up FPS overlay
	text := widget.NewLabel("")
	text.Alignment = fyne.TextAlignLeading
	text.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	text.Move(fyne.NewPos(5.0, 2.0))
	c.Overlays().Add(text)

	// set up image
	image := canvas.NewImageFromImage(&image.NRGBA{})
	image.FillMode = canvas.ImageFillContain
	image.ScaleMode = canvas.ImageScalePixels
	if scale < 1 {
		image.ScaleMode = canvas.ImageScaleSmooth
	}
	c.SetContent(image)

	go func() {
		// rolling avg fps
		var fpsRolling float64
		durIdx, durWindow := 0, 5
		durations := make([]float64, durWindow)

		for {
			for i := 0.0; i <= 10; i += 0.05 {
				start := time.Now()

				text.SetText(fmt.Sprintf("%-4.1f fps", fpsRolling))

				image.Image = createImage(rect, i)
				image.Refresh()

				// pause if required to maintain target fps
				delay := (time.Millisecond * frametime) - time.Since(start)
				time.Sleep(delay)

				// keep rolling average of fps over a few frames

				// calculate current framerate
				fpsActual := 1.0 / time.Since(start).Seconds()

				// update circular buffer with the duration of the current frame
				durations[durIdx] = fpsActual
				durIdx = (durIdx + 1) % durWindow

				// calculate the rolling average of framerate over the last x frames
				fpsRolling = (durations[0] + durations[1] + durations[2] + durations[3] + durations[4]) / float64(durWindow)
				fpsRolling = math.Round(fpsRolling/0.2) * 0.2 // to nearest 0.2

			}
		}
	}()

	w.ShowAndRun()
}

func createImage(rect image.Rectangle, i float64) (img *image.NRGBA) {
	width, height := rect.Dx(), rect.Dy()
	fov := math.Pi / 3.0

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
			Material: rt.Ivory,
		},
		{
			Centre:   rt.Vector3f{X: -5.0 + i, Y: -1.5 + (i / 3), Z: -12.0 + (i / 2)},
			Radius:   2.0,
			Material: rt.Glass,
		},
		{
			Centre:   rt.Vector3f{X: 1.5, Y: -0.5, Z: -18.0},
			Radius:   3.0,
			Material: rt.RedRubber,
		},
		{
			Centre:   rt.Vector3f{X: 7.0, Y: 5.0, Z: -18.0},
			Radius:   4.0,
			Material: rt.Mirror,
		},
	}

	lights := []*rt.Light{
		{
			Position:  rt.Vector3f{X: -20.0, Y: 20.0, Z: 20.0},
			Intensity: 1.5,
		},
		{
			Position:  rt.Vector3f{X: 30.0, Y: 50.0, Z: -25.0},
			Intensity: 1.8,
		},
		{
			Position:  rt.Vector3f{X: 30.0, Y: 20.0, Z: 30.0},
			Intensity: 1.7,
		},
	}

	render(img, width, height, fov, lights, spheres)
	return img
}

type empty struct{}

func render(img *image.NRGBA, width, height int, fov float64, lights []*rt.Light, spheres []*rt.Sphere) {
	sem := make(chan empty, width*height) // semaphore pattern

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			// TODO: what are these formulae? abstract into functions?
			x := (2*(float64(i)+0.5)/float64(width) - 1) * math.Tan(fov/2.0) * float64(width) / float64(height)
			y := -1.0 * (2*(float64(j)+0.5)/float64(height) - 1) * math.Tan(fov/2.0)

			origin := rt.Vector3f{X: 0, Y: 0, Z: 0}
			direction := rt.Vector3f{X: x, Y: y, Z: -1}.Normalised()

			// parallelise
			go func(i, j int) {
				c := castRay(origin, direction, lights, spheres, 0)
				img.Set(i, j, c)
				sem <- empty{}
			}(i, j)
		}
	}

	// wait for goroutines to finish
	for i := 0; i < width*height; i++ {
		<-sem
	}
}

func castRay(origin, direction rt.Vector3f, lights []*rt.Light, spheres []*rt.Sphere, depth int) color.NRGBA {
	var point, normal rt.Vector3f
	var material rt.Material

	if depth > MaxRayRecursionDepth || !sceneIntersect(origin, direction, &point, &normal, &material, spheres) {
		return rt.BackgroundColour
	}

	// calculate reflections and refractions

	reflectDir := reflect(direction.Multiply(-1.0), normal).Normalised()
	refractDir := refract(direction.Multiply(1.0), normal, material.RefractiveIndex).Normalised()
	// refractDir := refract(direction.Multiply(-1.0), normal, material.RefractiveIndex).Normalised()

	// offset the original point to avoid occlusion by the object itself
	reflectOrigin := point
	if reflectDir.Dot(normal) < 0 {
		reflectOrigin = reflectOrigin.Sub(normal.Multiply(1.0 / 1000))
	} else {
		reflectOrigin = reflectOrigin.Add(normal.Multiply(1.0 / 1000))
	}

	refractOrigin := point
	if refractDir.Dot(normal) < 0 {
		refractOrigin = refractOrigin.Sub(normal.Multiply(1.0 / 1000))
	} else {
		refractOrigin = refractOrigin.Add(normal.Multiply(1.0 / 1000))
	}

	// recursively calculate reflections (up to max depth)
	reflectColour := castRay(reflectOrigin, reflectDir, lights, spheres, depth+1)
	reflectColourVec := rt.Vector3f{
		X: float64(reflectColour.R),
		Y: float64(reflectColour.G),
		Z: float64(reflectColour.B),
	}

	refractColour := castRay(refractOrigin, refractDir, lights, spheres, depth+1)
	refractColourVec := rt.Vector3f{
		X: float64(refractColour.R),
		Y: float64(refractColour.G),
		Z: float64(refractColour.B),
	}

	diffuseLightIntensity := 0.0
	specularLightIntensity := 0.0

	for i := 0; i < len(lights); i++ {
		lightDir := (lights[i].Position.Sub(point)).Normalised()
		lightDist := (lights[i].Position.Sub(point)).Norm()

		// determine shadows
		//  make sure that the segment between the current point and the light
		//  source does not intersect the objects in the scene
		//  if there is an intersection we skip the current light source
		//  (and move the point in the direction of the normal)
		shadowOrigin := point
		if lightDir.Dot(normal) < 0.0 {
			shadowOrigin = shadowOrigin.Sub(normal.Multiply(1.0 / 1000))
		} else {
			shadowOrigin = shadowOrigin.Add(normal.Multiply(1.0 / 1000))
		}
		var shadowPoint, shadowNormal rt.Vector3f
		var tmpMaterial rt.Material
		if sceneIntersect(shadowOrigin, lightDir, &shadowPoint, &shadowNormal, &tmpMaterial, spheres) &&
			(shadowPoint.Sub(shadowOrigin).Norm() < lightDist) {
			continue
		}

		// determine brightness / reflection
		diffuseLightIntensity += lights[i].Intensity * math.Max(0.0, lightDir.Dot(normal))
		specularLightIntensity += math.Pow(
			math.Max(0.0, reflect(lightDir.Multiply(-1), normal).Dot(direction)),
			material.SpecularExponent,
		) * lights[i].Intensity
	}

	// TODO: define multiply function (with limiting) for colours instead of converting to vec
	c := material.DiffuseColour
	cVec := rt.Vector3f{
		X: float64(c.R),
		Y: float64(c.G),
		Z: float64(c.B),
	}

	// phong = ambient + diffuse + specular
	cVec = cVec.Multiply(diffuseLightIntensity).Multiply(material.Albedo[0]).Add(
		rt.Vector3f{X: 0xff, Y: 0xff, Z: 0xff}.Multiply(specularLightIntensity).Multiply(material.Albedo[1]),
	).Add(reflectColourVec.Multiply(material.Albedo[2])).Add(
		refractColourVec.Multiply(material.Albedo[3]))

	// prevent brightness from exceeding maximum
	max := math.Max(float64(cVec.X), math.Max(cVec.Y, cVec.Z))
	if max > 0xff {
		cVec.X *= 0xff / max
		cVec.Y *= 0xff / max
		cVec.Z *= 0xff / max
	}

	return color.NRGBA{
		R: uint8(cVec.X),
		G: uint8(cVec.Y),
		B: uint8(cVec.Z),
		A: c.A,
	}
}

// TODO: create a scene type with spheres, lights, etc.

func sceneIntersect(origin, direction rt.Vector3f, hit, N *rt.Vector3f, material *rt.Material, spheres []*rt.Sphere) bool {
	spheresDist := math.MaxFloat64

	for i := 0; i < len(spheres); i++ {
		var dist_i float64
		if spheres[i].RayIntersect(origin, direction, &dist_i) && dist_i < spheresDist {
			spheresDist = dist_i
			*hit = origin.Add(direction.Multiply(dist_i))
			*N = hit.Sub(spheres[i].Centre).Normalised()
			*material = spheres[i].Material
		}
	}

	// checkerboard logic
	checkerboardDist := math.MaxFloat64
	if math.Abs(direction.Y) > 1.0/1000 {
		d := (-1 * (origin.Y + 4)) / direction.Y // checkerboard plane y=4
		pt := origin.Add(direction.Multiply(d))

		if d > 0 && math.Abs(pt.X) < 10 && pt.Z < -10 && pt.Z > -30 && d < spheresDist {
			checkerboardDist = d
			*hit = pt
			*N = rt.Vector3f{X: 0.0, Y: 1.0, Z: 0.0}

			*material = rt.RedRubber // copy reflective properties
			// add checkerboard pattern:
			if (int(0.5*hit.X+1000)+int(0.5*hit.Z))&1 > 0 {
				material.DiffuseColour = rt.FloatToRGB(1.0, 1.0, 1.0)
			} else {
				material.DiffuseColour = rt.FloatToRGB(1.0, 0.7, 0.3)
			}
		}
	}

	// return spheresDist < 1000
	return math.Min(checkerboardDist, spheresDist) < 1000
}

// TODO: rename args to better names
func reflect(I, N rt.Vector3f) rt.Vector3f {
	return I.Sub(N.Multiply(2.0).Cross(I.Cross(N)))
}

func refract(I, N rt.Vector3f, refractiveIndex float64) rt.Vector3f {
	// snell's law

	cosi := -1 * math.Max(-1.0, math.Min(1.0, I.Dot(N)))
	etai := 1.0
	etat := refractiveIndex
	n := N

	// if the ray is inside the object, swap the indices and invert the normal to get the correct result
	if cosi < 0 {
		cosi = -1 * cosi
		etai, etat = etat, etai
		n = N.Multiply(-1.0)
	}
	eta := etai / etat

	sin := 1.0 - (cosi * cosi)
	k := 1.0 - (eta * eta * sin)

	if k < 0 {
		zero := math.SmallestNonzeroFloat64 // stop divide by zero on normalise
		return rt.Vector3f{X: zero, Y: zero, Z: zero}
	} else {
		return I.Multiply(eta).Add(n.Multiply(eta*cosi - math.Sqrt(k)))
	}
}

// TODO: apply this fix
// https://github.com/ssloy/tinyraytracer/commit/cc608d433d37a9116eee6da2467b8ac737b0a685#diff-6e69910b828e4c7d9cb06a9b779660c6R55
