package demos

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lucasb-eyer/go-colorful"
)

// Plasma generates an old school plasma effect.
type Plasma struct {
	hueShift float64
	buffer   []byte
}

var plasma [Width * Height]float64

// init will pre-compute the plasma values based on the addition of sines
// as detailed in https://rosettacode.org/wiki/Plasma_effect#Java
func init() {
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			fx, fy := float64(x), float64(y)
			i := y*Width + x
			value := math.Sin(fx / 16.0)
			value += math.Sin(fy / 8.0)
			value += math.Sin((fx + fy) / 16.0)
			value += math.Sin(math.Sqrt(fx*fx+fy*fy) / 8.0)
			value += 4 // shift range from -4 .. 4 to 0 .. 8
			value /= 8 // bring range down to 0 .. 1
			plasma[i] = value
		}
	}
}

// Update loops over the canvas and calculates the RGB of the pixel based
// on a hue that shifts every tick.
func (p *Plasma) Update() error {
	p.buffer = make([]byte, Width*Height*4)
	p.hueShift = math.Mod((p.hueShift + 0.001), 1.0)
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			idx := y*Width + x
			hue := p.hueShift + math.Mod(plasma[idx], 1.0)
			r, g, b, a := colorful.Hsv(hue, 1, 1).RGBA()
			p.buffer[4*idx] = byte(r)
			p.buffer[4*idx+1] = byte(g)
			p.buffer[4*idx+2] = byte(b)
			p.buffer[4*idx+3] = byte(a)
		}
	}
	return nil
}

// Draw will take the buffer and replace the pixels with it.
func (p *Plasma) Draw(screen *ebiten.Image) {
	screen.ReplacePixels(p.buffer)
}

// Layout returns a static width and height.
func (*Plasma) Layout(_, _ int) (int, int) {
	return Width, Height
}
