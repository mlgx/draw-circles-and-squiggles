package main

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"os"
	"strconv"

	"github.com/fogleman/gg"
	"github.com/h8gi/canvas"
)

//nolint:gochecknoglobals // TODO: Cleanup.
var (
	frameRate       = 30  // Frames per second.
	canvasSize      = 800 // Pixels.
	numCircles      = 10
	colorBackground = color.RGBA{64, 64, 64, 255}
	colorMarker     = color.White
)

type marker struct {
	radius float64
	step   int // Marker's step size in degrees.
	angle  int // Marker's last position in degrees.
}

type circle struct {
	x, y   int // Center coordinates.
	radius float64
	marker marker
}

func (c *circle) getMarkerCoordinates() (x, y float64) {
	x = float64(c.x) + math.Sin(gg.Radians(float64(c.marker.angle)))*c.radius
	y = float64(c.y) + math.Cos(gg.Radians(float64(c.marker.angle)))*c.radius
	return x, y
}
func (c *circle) getMixedMarkerCoordinates(c2 *circle) (x, y float64) {
	x = float64(c.x) + math.Sin(gg.Radians(float64(c.marker.angle)))*c.radius
	y = float64(c2.x) + math.Cos(gg.Radians(float64(c2.marker.angle)))*c2.radius
	return x, y
}
func (c *circle) moveMarker() {
	c.marker.angle += c.marker.step
	if c.marker.angle >= 360 {
		c.marker.angle -= 360
	}
}
func (c *circle) swapedXY() *circle {
	return &circle{x: c.y, y: c.x, radius: c.radius, marker: c.marker}
}

func main() {
	if len(os.Args) == 3 {
		canvasSize, _ = strconv.Atoi(os.Args[1])
		numCircles, _ = strconv.Atoi(os.Args[2])
	}

	c := canvas.NewCanvas(&canvas.CanvasConfig{
		Width:     canvasSize,
		Height:    canvasSize,
		FrameRate: frameRate,
	})

	// Generate and configure circles.
	circles := make([]*circle, numCircles)
	for i := range circles {
		circles[i] = &circle{
			x:      canvasSize/(numCircles+1)/2 + canvasSize/(numCircles+1)*(i+1),
			y:      canvasSize / (numCircles + 1) / 2,
			radius: float64(canvasSize / (numCircles + 1) / 2 * 80 / 100), // 20% padding.
			marker: marker{
				radius: float64(canvasSize / (numCircles + 1) / 2 * 20 / 100),
				step:   360 / frameRate / 6 * (i + 1), // 6 seconds to complete a revolution (=2˚).
			},
		}
	}

	// Pre-render circles and squiggles.
	var background *image.RGBA
	c.Setup(func(ctx *canvas.Context) {
		ctx.SetColor(colorBackground)
		ctx.Clear()
		ctx.SetColor(colorMarker)

		for i1, c1 := range circles {
			// Draw circles.
			ctx.SetLineWidth(2)
			ctx.SetColor(getColor(i1, numCircles))
			ctx.DrawCircle(float64(c1.x), float64(c1.y), c1.radius)
			ctx.DrawCircle(float64(c1.y), float64(c1.x), c1.radius)
			ctx.Stroke()

			// Draw squiggles.
			ctx.SetLineWidth(1)
			for i2, c2 := range circles {
				ctx.SetColor(mixColors(getColor(i1, numCircles), getColor(i2, numCircles)))
				// Complete a full revolution of the slowest/base marker.
				// NOTE: we divide markerStep by 2 to increase drawing precision.
				for angle := 0; angle < 360+circles[0].marker.step; angle += circles[0].marker.step / 2 {
					x, y := c1.getMixedMarkerCoordinates(c2)
					ctx.LineTo(x, y)
					// Move markers.
					c1.marker.angle += c1.marker.step / 2
					c2.marker.angle += c2.marker.step / 2
				}
				c1.marker.angle = 0
				c2.marker.angle = 0
				ctx.Stroke()
			}
		}

		// Copy the image.
		src := ctx.Image()
		background = image.NewRGBA(src.Bounds())
		draw.Draw(background, src.Bounds(), src, src.Bounds().Min, draw.Src)

		// Set marker color.
		ctx.SetColor(colorMarker)
	})

	c.Draw(func(ctx *canvas.Context) {
		// Draw the pre-rendered background.
		ctx.DrawImage(image.Image(background), 0, 0)
		ctx.Stroke()

		for _, c := range circles {
			// Draw circle markers.
			x, y := c.getMarkerCoordinates()
			ctx.DrawCircle(x, y, c.marker.radius)
			x, y = c.swapedXY().getMarkerCoordinates()
			ctx.DrawCircle(x, y, c.marker.radius)

			// Draw squiggle markers.
			for _, c2 := range circles {
				x, y = c.getMixedMarkerCoordinates(c2)
				ctx.DrawCircle(x, y, c.marker.radius)
			}

			// Move markers.
			c.moveMarker()

			ctx.Fill()
			ctx.Stroke()
		}
	})
}
