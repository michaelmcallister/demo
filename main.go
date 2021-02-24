package main

import (
	"container/list"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/michaelmcallister/demo/demos"
)

// Runner provides a method of running several ebiten.Game instances. It itself
// satisfies this interface, but apart from providing a mechanism to swap out
// the current running ebiten.Game instance by pressing Left, or Right on the
// keyboard, it will just pass control down calling the underlying methods.
type Runner struct {
	demos    *list.List
	selected *list.Element
}

// NewRunner returns a reference to an instane of Runner. The supplied demos
// must satisfy the ebiten.Game interface.
func NewRunner(demos *list.List) *Runner {
	return &Runner{demos: demos, selected: demos.Front()}
}

// Update handles switching the demos out by pressing the Left or Right key and
// then passes through to the selected demos Update() method.
func (r *Runner) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.Key(ebiten.KeyLeft)) {
		v := r.selected.Prev()
		if v == nil {
			v = r.demos.Back()
		}
		r.selected = v
	}
	if inpututil.IsKeyJustPressed(ebiten.Key(ebiten.KeyRight)) {
		v := r.selected.Next()
		if v == nil {
			v = r.demos.Front()
		}
		r.selected = v
	}
	return r.selected.Value.(ebiten.Game).Update()
}

// Draw draws the game screen and is called every frame (typically 1/60[s] for 60Hz display).
// The current TPS (amount of times Update() is called per second) is overlayed.
func (r *Runner) Draw(screen *ebiten.Image) {
	r.selected.Value.(ebiten.Game).Draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %f\n", ebiten.CurrentTPS()))
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (*Runner) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return demos.Width, demos.Height
}

func main() {
	ebiten.SetWindowSize(demos.Width, demos.Height)
	ebiten.SetWindowTitle("Various small graphics demos")
	l := list.New()
	l.PushBack(&demos.Water{})
	if err := ebiten.RunGame(NewRunner(l)); err != nil {
		log.Fatal(err)
	}
}
