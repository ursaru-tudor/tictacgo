package main

import (
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	WindowW int = 640
	WindowH int = 480
)

// Image assets
var (
	imageX *ebiten.Image
	imageO *ebiten.Image
)

func init() {
	var err error
	imageX, _, err = ebitenutil.NewImageFromFile("assets/X.png")
	if err != nil {
		log.Fatal(err)
	}

	imageO, _, err = ebitenutil.NewImageFromFile("assets/O.png")
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	screen.DrawImage(imageX, nil)

	var v ebiten.DrawImageOptions
	v.GeoM.Translate(64, 0)

	screen.DrawImage(imageO, &v)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowW, WindowH
}

func main() {
	ebiten.SetWindowSize(WindowW, WindowH)
	ebiten.SetWindowTitle("TicTacGo")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
