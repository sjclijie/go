// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package color

// RGBToYCbCr converts an RGB triple to a Y'CbCr triple.

// RGBToYCbCr将RGB的三重色转换为Y'CbCr模型的三重色。
func RGBToYCbCr(r, g, b uint8) (uint8, uint8, uint8) {
	// The JFIF specification says:
	//	Y' =  0.2990*R + 0.5870*G + 0.1140*B
	//	Cb = -0.1687*R - 0.3313*G + 0.5000*B + 128
	//	Cr =  0.5000*R - 0.4187*G - 0.0813*B + 128
	// http://www.w3.org/Graphics/JPEG/jfif3.pdf says Y but means Y'.
	r1 := int(r)
	g1 := int(g)
	b1 := int(b)
	yy := (19595*r1 + 38470*g1 + 7471*b1 + 1<<15) >> 16
	cb := (-11056*r1 - 21712*g1 + 32768*b1 + 257<<15) >> 16
	cr := (32768*r1 - 27440*g1 - 5328*b1 + 257<<15) >> 16
	if yy < 0 {
		yy = 0
	} else if yy > 255 {
		yy = 255
	}
	if cb < 0 {
		cb = 0
	} else if cb > 255 {
		cb = 255
	}
	if cr < 0 {
		cr = 0
	} else if cr > 255 {
		cr = 255
	}
	return uint8(yy), uint8(cb), uint8(cr)
}

// YCbCrToRGB converts a Y'CbCr triple to an RGB triple.

// YCbCrToRGB将Y'CbCr上的三重色转变成RGB的三重色。
func YCbCrToRGB(y, cb, cr uint8) (uint8, uint8, uint8) {
	// The JFIF specification says:
	//	R = Y' + 1.40200*(Cr-128)
	//	G = Y' - 0.34414*(Cb-128) - 0.71414*(Cr-128)
	//	B = Y' + 1.77200*(Cb-128)
	// http://www.w3.org/Graphics/JPEG/jfif3.pdf says Y but means Y'.
	yy1 := int(y)<<16 + 1<<15
	cb1 := int(cb) - 128
	cr1 := int(cr) - 128
	r := (yy1 + 91881*cr1) >> 16
	g := (yy1 - 22554*cb1 - 46802*cr1) >> 16
	b := (yy1 + 116130*cb1) >> 16
	if r < 0 {
		r = 0
	} else if r > 255 {
		r = 255
	}
	if g < 0 {
		g = 0
	} else if g > 255 {
		g = 255
	}
	if b < 0 {
		b = 0
	} else if b > 255 {
		b = 255
	}
	return uint8(r), uint8(g), uint8(b)
}

// YCbCr represents a fully opaque 24-bit Y'CbCr color, having 8 bits each for
// one luma and two chroma components.
//
// JPEG, VP8, the MPEG family and other codecs use this color model. Such
// codecs often use the terms YUV and Y'CbCr interchangeably, but strictly
// speaking, the term YUV applies only to analog video signals, and Y' (luma)
// is Y (luminance) after applying gamma correction.
//
// Conversion between RGB and Y'CbCr is lossy and there are multiple, slightly
// different formulae for converting between the two. This package follows
// the JFIF specification at http://www.w3.org/Graphics/JPEG/jfif3.pdf.

// YCbCr代表了完全不透明的24-bit的Y'CbCr的颜色，它的每个亮度和每两个色度分量是8位的。
//
// JPEG，VP8，MPEG家族和其他一些解码器使用这个颜色模式。每个解码器经常将YUV和Y'CbCr同等使用，
// 但是严格来说，YUV只是用于分析视频信号，Y' (luma)是Y (luminance)伽玛校正之后的结果。
//
// RGB和Y'CbCr之间的转换是有损的，并且转换的时候有许多细微的不同。这个包是遵循JFIF的说明：
// http://www.w3.org/Graphics/JPEG/jfif3.pdf。
type YCbCr struct {
	Y, Cb, Cr uint8
}

func (c YCbCr) RGBA() (uint32, uint32, uint32, uint32) {
	r, g, b := YCbCrToRGB(c.Y, c.Cb, c.Cr)
	return uint32(r) * 0x101, uint32(g) * 0x101, uint32(b) * 0x101, 0xffff
}

// YCbCrModel is the Model for Y'CbCr colors.

// YCbCrModel是Y'CbCr颜色的模型。
var YCbCrModel Model = ModelFunc(yCbCrModel)

func yCbCrModel(c Color) Color {
	if _, ok := c.(YCbCr); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	y, u, v := RGBToYCbCr(uint8(r>>8), uint8(g>>8), uint8(b>>8))
	return YCbCr{y, u, v}
}

// RGBToCMYK converts an RGB triple to a CMYK quadruple.
func RGBToCMYK(r, g, b uint8) (uint8, uint8, uint8, uint8) {
	rr := uint32(r)
	gg := uint32(g)
	bb := uint32(b)
	w := rr
	if w < gg {
		w = gg
	}
	if w < bb {
		w = bb
	}
	if w == 0 {
		return 0, 0, 0, 255
	}
	c := (w - rr) * 255 / w
	m := (w - gg) * 255 / w
	y := (w - bb) * 255 / w
	return uint8(c), uint8(m), uint8(y), uint8(255 - w)
}

// CMYKToRGB converts a CMYK quadruple to an RGB triple.
func CMYKToRGB(c, m, y, k uint8) (uint8, uint8, uint8) {
	w := uint32(255 - k)
	r := uint32(255-c) * w / 255
	g := uint32(255-m) * w / 255
	b := uint32(255-y) * w / 255
	return uint8(r), uint8(g), uint8(b)
}

// CMYK represents a fully opaque CMYK color, having 8 bits for each of cyan,
// magenta, yellow and black.
//
// It is not associated with any particular color profile.
type CMYK struct {
	C, M, Y, K uint8
}

func (c CMYK) RGBA() (uint32, uint32, uint32, uint32) {
	r, g, b := CMYKToRGB(c.C, c.M, c.Y, c.K)
	return uint32(r) * 0x101, uint32(g) * 0x101, uint32(b) * 0x101, 0xffff
}

// CMYKModel is the Model for CMYK colors.
var CMYKModel Model = ModelFunc(cmykModel)

func cmykModel(c Color) Color {
	if _, ok := c.(CMYK); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	cc, mm, yy, kk := RGBToCMYK(uint8(r>>8), uint8(g>>8), uint8(b>>8))
	return CMYK{cc, mm, yy, kk}
}
