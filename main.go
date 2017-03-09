package main

import (
	"image"
	"image/color"
	"image/png"
	"math"

	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/view"
)

const tileSize int = 256
const minX float64 = -2
const maxX float64 = 2
const minY float64 = -2
const maxY float64 = 2

func main() {
	app := iris.New()

	app.Adapt(iris.DevLogger(), httprouter.New())
	app.Adapt(view.HTML("./views", ".html"))
	app.Get("/", func(ctx *iris.Context) {
		ctx.Render("index.html", iris.Map{}, iris.RenderOptions{})
	})

	app.Get("/:x/:y/:z", func(ctx *iris.Context) {
		x, _ := ctx.ParamInt("x")
		y, _ := ctx.ParamInt("y")
		z, _ := ctx.ParamInt("z")

		img := renderTile(x, y, z)

		ctx.Header().Add("Content-Type", "image/png")
		png.Encode(ctx, img)
	})

	app.Listen(":6300")
}

func renderTile(x int, y int, z int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))

	numberOfTiles := math.Pow(2, float64(z))
	x1 := (maxX - minX) * (float64(x) - (numberOfTiles / 2)) / numberOfTiles
	y1 := (maxY - minY) * (float64(y) - (numberOfTiles / 2)) / numberOfTiles
	pixelSize := (maxX - minX) / (numberOfTiles * float64(tileSize))

	for dx := 0; dx < tileSize; dx++ {
		for dy := 0; dy < tileSize; dy++ {

			tx := x1 + (float64(dx) * pixelSize)
			ty := y1 + (float64(dy) * pixelSize)

			value := getColour(tx, ty)

			if value >= 0 {
				sinVal := math.Floor(255 * math.Sin(float64(value)*math.Pi/255.0))
				img.Set(dx, dy, color.RGBA{uint8(sinVal), uint8(255 - value), uint8(255 - sinVal), 255})
			} else {
				img.Set(dx, dy, color.RGBA{0, 0, 0, 255})
			}

		}
	}
	return img
}

func getColour(re float64, im float64) int {
	zRe := float64(0)
	zIm := float64(0)

	//Variables to store the squares of the real and imaginary part.
	multZre := float64(0)
	multZim := float64(0)

	//Start iterating the with the complex number to determine it's escape time (mandelValue)
	mandelValue := int(0)
	for mandelValue < 255 {
		if multZre+multZim >= 4 {
			return mandelValue
		}

		/*The new real part equals re(z)^2 - im(z)^2 + re(c), we store it in a temp variable
		  tempRe because we still need re(z) in the next calculation
		*/
		tempRe := multZre - multZim + re

		/*The new imaginary part is equal to 2*re(z)*im(z) + im(c)
		 * Instead of multiplying these by 2 I add re(z) to itself and then multiply by im(z), which
		 * means I just do 1 multiplication instead of 2.
		 */
		zRe += zRe
		zIm = zRe*zIm + im

		zRe = tempRe // We can now put the temp value in its place.

		// Do the squaring now, they will be used in the next calculation.
		multZre = zRe * zRe
		multZim = zIm * zIm

		//Increase the mandelValue by one, because the iteration is now finished.
		mandelValue++
	}
	return -1
}
