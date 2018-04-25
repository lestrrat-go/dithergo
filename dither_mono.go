package dither

import (
	"image"
	"image/color"
)

func Monochrome(m *Matrix, input image.Image, errorMultiplier float32) image.Image {
	bounds := input.Bounds()
	img := image.NewGray(bounds)
	for x := bounds.Min.X; x < bounds.Dx(); x++ {
		for y := bounds.Min.Y; y < bounds.Dy(); y++ {
			pixel := input.At(x, y)
			img.Set(x, y, pixel)
		}
	}
	dx, dy := bounds.Dx(), bounds.Dy()

	// Prepopulate multidimensional slice
	errors := NewMatrix(dx, dy)

	ydim := m.Rows() - 1
	xdim := m.Cols() / 2
	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			pix := float32(img.GrayAt(x, y).Y)
			pix -= errors.Get(x, y) * errorMultiplier

			var quantError float32
			// Diffuse the error of each calculation to the neighboring pixels
			if pix < 128 {
				quantError = -pix
				pix = 0
			} else {
				quantError = 255 - pix
				pix = 255
			}

			img.SetGray(x, y, color.Gray{Y: uint8(pix)})

			// Diffuse error in two dimension
			for xx := 0; xx < ydim+1; xx++ {
				for yy := -xdim; yy <= xdim-1; yy++ {
					if y+yy < 0 || dy <= y+yy || x+xx < 0 || dx <= x+xx {
						continue
					}
					// Adds the error of the previous pixel to the current pixel
					prev := errors.Get(x+xx, y+yy)
					delta := quantError * m.Get(yy+ydim, xx)
					errors.Set(x+xx, y+yy, prev+delta)
				}
			}
		}
	}
	return img
}
