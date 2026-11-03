package main

// Extra framebuffer helpers for sharing the screen with ffplay: ffplay (SDL)
// page-flips between two buffers, so before drawing we pan back to page 0,
// and to animate the picture shrinking at power-off we read back whichever
// page ffplay left on screen.

import (
	"unsafe"
)

// yoffset returns the currently displayed row offset (0 or H for page 2).
func (fb *Framebuffer) yoffset() int {
	var v [40]uint32
	if ioctl(fb.f.Fd(), fbioGetVScreenInfo, unsafe.Pointer(&v[0])) != nil {
		return 0
	}
	return int(v[5])
}

// Reset pans the display back to page 0, which is where Present draws.
func (fb *Framebuffer) Reset() {
	var v [40]uint32
	if ioctl(fb.f.Fd(), fbioGetVScreenInfo, unsafe.Pointer(&v[0])) != nil {
		return
	}
	if v[5] != 0 || v[4] != 0 || v[6] != 32 {
		v[4], v[5], v[6] = 0, 0, 32
		_ = ioctl(fb.f.Fd(), fbioPutVScreenInfo, unsafe.Pointer(&v[0]))
	}
}

// Capture copies what is on screen right now (any page) into the canvas,
// undoing the panel rotation.
func (fb *Framebuffer) Capture(c *Canvas) bool {
	yo := fb.yoffset()
	if (yo+fb.H)*fb.stride > len(fb.mem) {
		yo = 0
	}
	w, h := fb.W, fb.H
	if c.W < w || c.H < h {
		return false
	}
	for y := 0; y < h; y++ {
		sy := y
		if fb.flip {
			sy = h - 1 - y
		}
		src := fb.mem[(yo+sy)*fb.stride:]
		dst := c.Pix[y*c.Stride:]
		for x := 0; x < w; x++ {
			sx := x
			if fb.flip {
				sx = w - 1 - x
			}
			s := sx * 4
			d := x * 4
			dst[d], dst[d+1], dst[d+2], dst[d+3] = src[s+2], src[s+1], src[s], 255
		}
	}
	return true
}
