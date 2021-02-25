package demos

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const damping float32 = 0.99
const force int = 1500

// Water provides a 2D ripple effect based on an algorithm explained in
// https://web.archive.org/web/20160418004149/http://freespace.virgin.net/hugo.elias/graphics/x_water.htm
type Water struct {
	prev   [Width * Height]int
	cur    [Width * Height]int
	buffer []byte
}

var waterBackground *ebiten.Image

func init() {
	// This serves as the background in which to draw the ripple effect over.
	waterBackground = ebiten.NewImage(Width, Height)
	waterBackground.Fill(color.RGBA{0x84, 0xC4, 0xEF, 0xFF})
}

// Update handles the logic to place ripples where the mouse is clicked, as well
// as generating the ripple effect.
func (w *Water) Update() error {
	// Place the start of the ripple at the mouse x,y co-ordinates when the
	// left mouse button is clicked.
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		idx := y*Width + x
		if idx > 0 && idx < (Width*Height) {
			w.prev[idx] = force
		}
	}

	// If running from a touch device (mobile, web) use the touch position.
	if t := ebiten.TouchIDs(); t != nil {
		x, y := ebiten.TouchPosition(t[0])
		idx := y*Width + x
		if idx > 0 && idx < (Width*Height) {
			w.prev[idx] = force
		}
	}

	w.buffer = make([]byte, Width*Height*4)
	for y := 1; y < Height-1; y++ {
		for x := 1; x < Width-1; x++ {
			i := y*Width + x
			w.cur[i] = (w.prev[y*Width+(x-1)]+
				w.prev[y*Width+(x+1)]+
				w.prev[(y+1)*Width+x]+
				w.prev[(y-1)*Width+x])/2.0 - w.cur[i]

			w.cur[i] = int(float32(w.cur[i]) * damping)

			// RGB
			v := uint8(w.cur[i] * 255)
			w.buffer[4*i] = v
			w.buffer[4*i+1] = v
			w.buffer[4*i+2] = v
		}
	}
	// swap the buffers.
	t := w.prev
	w.prev = w.cur
	w.cur = t
	return nil
}

// Draw will plot the ripple effect and then fill in the gaps with a blue
// background.
func (w *Water) Draw(screen *ebiten.Image) {
	screen.ReplacePixels(w.buffer)
	screen.DrawImage(waterBackground, &ebiten.DrawImageOptions{
		CompositeMode: ebiten.CompositeModeLighter,
	})
}

// Layout returns a static width and height.
func (*Water) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return Width, Height
}
