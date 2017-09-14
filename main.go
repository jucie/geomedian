// This program is able to find out the geographic median of a region.
//
// "Geographic median is the point that divides the nation into
// equal area north-south regions and equal area east-west regions."
// See http://www.apprendre-en-ligne.net/blog/docu/centerus.pdf

package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

// medianIndex returns the index inside the array so that
// the sum of the values below the element denoted by index
// is roughly the same as the sum of the elements above it.
func medianIndex(values []int) int {
	// first pass: will calculate half the total sum of elements.
	sum := 0
	for _, val := range values {
		sum += val
	}
	half := sum / 2

	// second pass: will find the point that splits in two halves.
	sum = 0
	var index, val int
	for index, val = range values {
		sum += val
		if sum >= half { // do we have reached the middle point?
			break // se we are done.
		}
	}
	return index
}

// calcHoriz produces an array with the counters for each column.
// ref is the background color. Pixels with differing colors increment the counters.
func calcHoriz(img image.Image, ref color.Color) []int {
	refR, refG, refB, refA := ref.RGBA()

	bounds := img.Bounds()
	var horiz []int
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		sum := 0
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			if !(r == refR && g == refG && b == refB && a == refA) {
				sum++
			}
		}
		horiz = append(horiz, sum)
	}
	return horiz
}

// calcVerti produces an array with the counters for each row.
// ref is the background color. Pixels with differing colors increment the counters.
func calcVerti(img image.Image, ref color.Color) []int {
	refR, refG, refB, refA := ref.RGBA()

	bounds := img.Bounds()
	var verti []int
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		sum := 0
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if !(r == refR && g == refG && b == refB && a == refA) {
				sum++
			}
		}
		verti = append(verti, sum)
	}
	return verti
}

// medianPoint returns the geographic median of the map.
func medianPoint(img image.Image) image.Point {
	refColor := img.At(0, 0)          // consider upper left corner as the background color.
	horiz := calcHoriz(img, refColor) // produces the conters for horizontal axis
	verti := calcVerti(img, refColor) // produces the conters for vertical axis

	bounds := img.Bounds()
	x := bounds.Min.X + medianIndex(horiz) // displaces the index by the image offset
	y := bounds.Min.Y + medianIndex(verti) // displaces the index by the image offset
	return image.Point{X: x, Y: y}         // pack it into a Point
}

// drawCross draws a cross marker over the image at point.
func drawCross(img image.Image, point image.Point) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)                                    // we need a writeable image.
	draw.Draw(dst, bounds, img, image.Point{X: 0, Y: 0}, draw.Over) // copy source image onto destination.

	// calculate the cross marker size as a fraction of the image size
	hSize := (bounds.Max.X - bounds.Min.X) / 16
	vSize := (bounds.Max.Y - bounds.Min.Y) / 16

	for x := point.X - hSize; x < point.X+hSize; x++ {
		dst.Set(x, point.Y, color.Black) // draws along the horizontal axis
	}
	for y := point.Y - vSize; y < point.Y+vSize; y++ {
		dst.Set(point.X, y, color.Black) // draws along the vertical axis
	}
	return dst // returns the resulting new image.
}

func main() {
	if len(os.Args) < 3 {
		panic("Input and output filenames are missing.")
	}

	in, err := os.Open(os.Args[1]) // open input file
	if err != nil {
		panic(err)
	}
	defer in.Close()

	img, err := png.Decode(in) // extract the source image
	if err != nil {
		panic(err)
	}

	pt := medianPoint(img)   // calculate the geographic median.
	img = drawCross(img, pt) // produces a new image, with a cross.

	out, err := os.Create(os.Args[2]) // open destination for writing.
	if err != nil {
		panic(err)
	}
	defer out.Close()

	err = png.Encode(out, img) // write the new image to the output file.
	if err != nil {
		panic(err)
	}
}
