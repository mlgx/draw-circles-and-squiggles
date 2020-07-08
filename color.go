// Source: https://github.com/lucasb-eyer/go-colorful/blob/master/doc/gradientgen/gradientgen.go

package main

import (
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

// This table contains the "keypoints" of the colorgradient you want to generate.
// The position of each keypoint has to live in the range [0,1]
type gradientTable []struct {
	Col colorful.Color
	Pos float64
}

// This is the meat of the gradient computation. It returns a HCL-blend between
// the two colors around `t`.
// Note: It relies heavily on the fact that the gradient keypoints are sorted.
func (gt gradientTable) getInterpolatedColorFor(t float64) colorful.Color {
	for i := 0; i < len(gt)-1; i++ {
		c1 := gt[i]
		c2 := gt[i+1]
		if c1.Pos <= t && t <= c2.Pos {
			// We are in between c1 and c2. Go blend them!
			t := (t - c1.Pos) / (c2.Pos - c1.Pos)
			return c1.Col.BlendHcl(c2.Col, t).Clamped()
		}
	}

	// Nothing found? Means we're at (or past) the last gradient keypoint.
	return gt[len(gt)-1].Col
}

// This is a very nice thing Golang forces you to do!
// It is necessary so that we can write out the literal of the colortable below.
func mustParseHex(s string) colorful.Color {
	c, err := colorful.Hex(s)
	if err != nil {
		panic("MustParseHex: " + err.Error())
	}
	return c
}

func getColor(x, h int) color.Color {
	// The "keypoints" of the gradient.
	keypoints := gradientTable{
		{mustParseHex("#9e0142"), 0.0},
		{mustParseHex("#d53e4f"), 0.1},
		{mustParseHex("#f46d43"), 0.2},
		{mustParseHex("#fdae61"), 0.3},
		{mustParseHex("#fee090"), 0.4},
		{mustParseHex("#ffffbf"), 0.5},
		{mustParseHex("#e6f598"), 0.6},
		{mustParseHex("#abdda4"), 0.7},
		{mustParseHex("#66c2a5"), 0.8},
		{mustParseHex("#3288bd"), 0.9},
		{mustParseHex("#5e4fa2"), 1.0},
	}

	return keypoints.getInterpolatedColorFor(float64(x) / float64(h)).Clamped()
}

func mixColors(a, b color.Color) color.Color {
	if a == b {
		return a
	}
	c1, _ := colorful.MakeColor(a)
	c2, _ := colorful.MakeColor(b)
	return c1.BlendHcl(c2, 0.5).Clamped()
}
