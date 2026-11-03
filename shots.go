package main

// Demo screens for the README (go run . -shot guide) and the app icon.

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
	"time"
)

func demoChannels() []*Channel {
	mk := func(num int, name string, items ...string) *Channel {
		c := &Channel{Num: num, Name: name, Key: name, Seed: float64(num * 977)}
		for i, t := range items {
			d := 22*60.0 + float64(i%3)*60
			if len(items) == 1 {
				d = 104 * 60
			}
			c.Items = append(c.Items, &Item{Title: t, Dur: d, W: 320, H: 240})
			c.Total += d
		}
		return c
	}
	return []*Channel{
		mk(2, "Saturday Cartoons", "Rocko's Modern Life", "Doug", "Rugrats", "Ren & Stimpy"),
		mk(3, "Back to the Future", "Back to the Future"),
		mk(4, "Home Videos", "Christmas 1994", "Beach Trip 1996", "Talent Show"),
		mk(5, "Nature Hour", "Big Buck Bunny", "Desert Life", "Ocean Deep"),
		mk(7, "Sitcom Night", "Family Matters", "Full House", "Step by Step"),
		mk(9, "Late Movie", "The Goonies"),
		mk(11, "Music Videos", "Top 20 Countdown", "Alt Rock Block"),
	}
}

func shotBackground(c *Canvas) {
	if p := os.Getenv("RABBITEARS_SHOT_BG"); p != "" {
		if f, err := os.Open(p); err == nil {
			defer f.Close()
			if img, _, err := image.Decode(f); err == nil {
				rgba := image.NewRGBA(img.Bounds())
				draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)
				c.DrawImage(rgba, 0, 0, c.W, c.H)
				return
			}
		}
	}
	DrawBars(c, newFonts(), "", "", "")
}

// overlayNRGBA scales an OSD image onto the canvas the way ffplay would.
func overlayNRGBA(c *Canvas, o *image.NRGBA) {
	b := o.Bounds()
	for y := 0; y < c.H; y++ {
		sy := y * b.Dy() / c.H
		for x := 0; x < c.W; x++ {
			sx := x * b.Dx() / c.W
			so := sy*o.Stride + sx*4
			a := int(o.Pix[so+3])
			if a == 0 {
				continue
			}
			do := y*c.Stride + x*4
			for k := 0; k < 3; k++ {
				c.Pix[do+k] = uint8((int(c.Pix[do+k])*(255-a) + int(o.Pix[so+k])*a) / 255)
			}
		}
	}
}

func renderShot(name string) error {
	c := NewCanvas(640, 480)
	fonts := newFonts()
	chans := demoChannels()
	now := time.Date(2026, 10, 18, 20, 42, 0, 0, time.Local)
	switch name {
	case "osd", "info":
		shotBackground(c)
		var info *Airing
		if name == "info" {
			a := chans[3].Upcoming(now, 1)[0]
			info = &a
		}
		overlayNRGBA(c, RenderOSD(fonts, 320, 240, chans[3], info, now))
	case "guide":
		g := &Guide{Sel: 1}
		DrawGuide(c, fonts, chans, g, now, 3)
	case "bars":
		DrawBars(c, fonts, "NO CHANNELS", "Put videos (or folders of episodes) in", "SD card: Media/TV  —  then press A")
	case "nosignal":
		DrawNoSignal(c, fonts, &Snow{s: 7}, chans[2], "Invalid data found when processing input")
	case "snow":
		(&Snow{s: 9}).Draw(c)
	case "on":
		under := NewCanvas(640, 480)
		(&Snow{s: 5}).Draw(under)
		PowerOnFrame(c, under, 0.6)
	case "off":
		last := NewCanvas(640, 480)
		shotBackground(last)
		PowerOffFrame(c, last, 0.3)
	default:
		return fmt.Errorf("unknown shot %q", name)
	}
	f, err := os.Create("shot.png")
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, c.RGBA)
}

// appIcon: a little wood-grain CRT set with rabbit-ear antennas, colour
// bars on the screen. Transparent background.
func appIcon(n int) *image.RGBA {
	c := NewCanvas(n, n)
	for i := range c.Pix {
		c.Pix[i] = 0
	}
	u := float64(n) / 74
	px := func(v float64) int { return int(math.Round(v * u)) }
	// antennas
	cx, cy := px(37), px(25)
	for _, ang := range []float64{-0.62, 0.55} {
		x1 := float64(cx) + math.Sin(ang)*float64(px(24))
		y1 := float64(cy) - math.Cos(ang)*float64(px(24))
		for t := 0.0; t <= 1; t += 0.02 {
			x := float64(cx) + (x1-float64(cx))*t
			y := float64(cy) + (y1-float64(cy))*t
			c.Fill(int(x), int(y), px(2), px(2), rgb(0xc8, 0xcc, 0xd4))
		}
		c.Circle(int(x1), int(y1), px(2.5), rgb(0xe8, 0xe8, 0xf0), true, 0)
	}
	c.Circle(cx, cy+px(1), px(4), rgb(0x30, 0x30, 0x34), true, 0)
	// cabinet
	bx, by, bw, bh := px(6), px(27), px(62), px(42)
	c.Fill(bx, by, bw, bh, rgb(0x6b, 0x3e, 0x1e))
	for y := by; y < by+bh; y += px(3) + 1 {
		c.Blend(bx, y, bw, 1, rgb(0x4a, 0x28, 0x10), 0.6)
	}
	c.Border(bx, by, bw, bh, 1, rgb(0x2a, 0x16, 0x08))
	// screen bezel + screen with colour bars
	sx, sy, sw, sh := bx+px(4), by+px(4), px(42), bh-px(8)
	c.Fill(sx-1, sy-1, sw+2, sh+2, rgb(0x14, 0x14, 0x16))
	cols := []uint32{0xc0c0c0, 0xc0c000, 0x00c0c0, 0x00c000, 0xc000c0, 0xc00000, 0x0000c0}
	for i, col := range cols {
		x0, x1 := sx+i*sw/7, sx+(i+1)*sw/7
		c.Fill(x0, sy, x1-x0, sh*3/4, rgb(uint8(col>>16), uint8(col>>8), uint8(col)))
	}
	c.Fill(sx, sy+sh*3/4, sw, sh-sh*3/4, rgb(0x10, 0x10, 0x30))
	c.Scanlines(sx, sy, sw, sh, 0.25)
	// knobs + speaker
	kx := sx + sw + px(5)
	c.Circle(kx+px(4), sy+px(6), px(3), rgb(0xd8, 0xc8, 0xa0), true, 0)
	c.Circle(kx+px(4), sy+px(15), px(3), rgb(0xd8, 0xc8, 0xa0), true, 0)
	for y := sy + px(22); y < sy+sh-px(1); y += px(2) + 1 {
		c.Fill(kx, y, px(8), 1, rgb(0x2a, 0x16, 0x08))
	}
	// feet
	c.Fill(bx+px(6), by+bh, px(6), px(3), rgb(0x2a, 0x16, 0x08))
	c.Fill(bx+bw-px(12), by+bh, px(6), px(3), rgb(0x2a, 0x16, 0x08))
	// keep alpha only where we drew
	for i := 0; i < len(c.Pix); i += 4 {
		if c.Pix[i] != 0 || c.Pix[i+1] != 0 || c.Pix[i+2] != 0 {
			c.Pix[i+3] = 255
		}
	}
	return c.RGBA
}
